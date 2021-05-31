package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lamassuiot/lamassu-ca/pkg/auth"
	"github.com/lamassuiot/lamassu-ca/pkg/secrets"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/gorilla/mux"
)

type errorer interface {
	error() error
}

var (
	errCAName = errors.New("ca name not provided")
)

var claims = &auth.KeycloakClaims{}

func MakeHTTPHandler(s Service, logger log.Logger, auth auth.Auth, otTracer stdopentracing.Tracer) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s, otTracer)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	r.Methods("GET").Path("/v1/health").Handler(httptransport.NewServer(
		e.HealthEndpoint,
		decodeHealthRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Health", logger)))...,
	))

	r.Methods("GET").Path("/v1/cas").Handler(httptransport.NewServer(
		jwt.NewParser(auth.Kf, stdjwt.SigningMethodRS256, auth.KeycloakClaimsFactory)(e.GetCAsEndpoint),
		decodeGetCAsRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "GetCAs", logger)))...,
	))

	r.Methods("GET").Path("/v1/cas/{ca}").Handler(httptransport.NewServer(
		jwt.NewParser(auth.Kf, stdjwt.SigningMethodRS256, auth.KeycloakClaimsFactory)(e.GetCACrtEndpoint),
		decodeGetCACrtRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "GetCACrt", logger)))...,
	))

	r.Methods("POST").Path("/v1/cas/{ca}").Handler(httptransport.NewServer(
		jwt.NewParser(auth.Kf, stdjwt.SigningMethodRS256, auth.KeycloakClaimsFactory)(e.CreateCAEndpoint),
		decodeCreateCARequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "CreateCA", logger)))...,
	))

	r.Methods("POST").Path("/v1/cas/import/{ca}").Handler(httptransport.NewServer(
		jwt.NewParser(auth.Kf, stdjwt.SigningMethodRS256, auth.KeycloakClaimsFactory)(e.ImportCAEndpoint),
		decodeImportCARequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "ImportCA", logger)))...,
	))

	r.Methods("DELETE").Path("/v1/cas/{ca}").Handler(httptransport.NewServer(
		jwt.NewParser(auth.Kf, stdjwt.SigningMethodRS256, auth.KeycloakClaimsFactory)(e.DeleteCAEndpoint),
		decodeDeleteCARequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "DeleteCA", logger)))...,
	))

	return r
}

func decodeHealthRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req healthRequest
	return req, nil
}

func decodeGetCAsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req getCAsRequest
	return req, nil
}

func decodeGetCACrtRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	CA, ok := vars["ca"]
	if !ok {
		return nil, errCAName
	}
	return getCACrtRequest{CA: CA}, nil
}

func decodeCreateCARequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	var caRequestInfo secrets.CA
	json.NewDecoder(r.Body).Decode((&caRequestInfo))
	if err != nil {
		return nil, errors.New("Cannot decode JSON request")
	}

	caName, ok := vars["ca"]
	if !ok {
		return nil, errCAName
	}
	return createCARequest{CAName: caName, CA: caRequestInfo}, nil
}

func decodeImportCARequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	var importCaRequest secrets.CAImport
	json.NewDecoder(r.Body).Decode((&importCaRequest))
	if err != nil {
		return nil, errors.New("Cannot decode JSON request")
	}

	caName, ok := vars["ca"]
	if !ok {
		return nil, errCAName
	}
	return importCARequest{CAName: caName, CAImport: importCaRequest}, nil
}

func decodeDeleteCARequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	CA, ok := vars["ca"]
	if !ok {
		return nil, errCAName
	}
	return deleteCARequest{CA: CA}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.

		// https://medium.com/@ozdemir.zynl/rest-api-error-handling-in-go-behavioral-type-assertion-509d93636afd
		//
		encodeError(ctx, e.error(), w)

		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	http.Error(w, err.Error(), codeFrom(err))
}

func codeFrom(err error) int {
	switch err {
	case errCAName:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
