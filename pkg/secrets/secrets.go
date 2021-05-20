package secrets

// CAInfo represents the detailed information about a CA
// swagger:model
type CACrt struct {
	// PEM ca certificate
	// required: false
	// example: ----BEGIN CERTIFICATE-----\nMIID2TCCAsGgAwIBAgIUcYimUsFDI6395PM2WbAvPEtbfjowDQYJKoZIhvcNAQEL\nBQAwczELMAkGA1UEBhMCRVMxETAPBgNVBAgTCEdpcHV6a29hMREwDwYDVQQHEwhB\ncnJhc2F0ZTEhMA4GA1UEChMHUy4gQ29vcDAPBgNVBAoTCExLUyBOZXh0MRswGQYD\nVQQDExJMS1MgTmV4dCBSb290IENBIDIwIBcNMjEwNTE4MTEzNzM2WhgPMjA1MTA1\nMTExMTM4MDZaMHMxCzAJBgNVBAYTAkVTMREwDwYDVQQIEwhHaXB1emtvYTERMA8G\nA1UEBxMIQXJyYXNhdGUxITAOBgNVBAoTB1MuIENvb3AwDwYDVQQKEwhMS1MgTmV4\ndDEbMBkGA1UEAxMSTEtTIE5leHQgUm9vdCBDQSAyMIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEA2ePwTAHaGPd3H/I3mRkLqL0GxgcZw/VlSHfT0I6clIvQ\n1Ulc7kL0NZRTYPOsBQIjWuu61PwSwPgop/N+slMYpG/NOJwKzH9JHAjNKISuNasS\n66Q3pLBK/QMHIZsaRkPOCfVlQeV75YFhehtabxM10CLdJq9HE5iKY/B1SEdCcAz4\nGbzVy/DzdqAtHrdwyjlS2DM+hYWEvUwbZIzSAWlOtIMHCYypd5wvYTN3tfsYtjft\nTwT3gIdoQTz4eOF/HGmE3NglO3qJspze7sgMDmcfBrgo51C+XOfmZ5zYk1cJSkjM\nsT3tcmwJlBP6va2AGTuTtCQDhbGbnXM33uIlh7L9JQIDAQABo2MwYTAOBgNVHQ8B\nAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUMfBBHj4BJqvG8FRH\n2nMrdF/JZo8wHwYDVR0jBBgwFoAUMfBBHj4BJqvG8FRH2nMrdF/JZo8wDQYJKoZI\nhvcNAQELBQADggEBANWH9n6Ezh0hmozhtu8HGKybIxAVTmxiXirY3sYwIsMyB1Ns\nljbGaah4qqmzQKCAqfaeQbd1YMER+C98OnA7S/xV0Vxucu5g/obFekXyJf1U9SLW\nfh5tuCtsfgkSNPLk21hWMFfZR3hJKfcK6GuoTOW6cBUf+VbWLO6tsO011xWF4tYj\nfppbk7wHT6LIFY3wsKl5ti16U0gd/s9XfqYR84y9bZWZ+SGzNC3n9OWxvYnOrX/B\nNO/ucnBKon7kpHX91kkj9kWRNONAf2lWTeg0WcUm2e1sim6fEekux7cg1PCqz3Li\n2zRuHYvLO1cBeXQ+8olyCpBQDWaXMWkoNW49xbY=\n-----END CERTIFICATE-----
	CRT string `json:"crt"`

	PublicKey string `json:"pub_key"`
}

// CA represents a registered CA minimum information
// swagger:model
type CA struct {
	// The serial number of the CA
	// required: true
	// example: 7e:36:13:a5:31:9f:4a:76:10:64:2e:9b:0a:11:07:b7:e6:3e:cf:94
	SerialNumber string `json:"serial_number,omitempty"`

	// The serial number of the CA
	// required: true
	// example: 7e:36:13:a5:31:9f:4a:76:10:64:2e:9b:0a:11:07:b7:e6:3e:cf:94
	CaName string `json:"ca_name,omitempty"`

	// Common name of the CA certificate
	// required: true
	// example: Lamassu-Root-CA1-RSA4096
	CN string `json:"common_name"`

	// Algorithm used to create CA key
	// required: true
	// example: RSA
	KeyType string `json:"key_type"`

	// Length used to create CA key
	// required: true
	// example: 4096
	KeyBits int `json:"key_bits"`

	// Organization of the CA certificate
	// required: true
	// example: Lamassu IoT
	O string `json:"organization"`

	// Organization Unit of the CA certificate
	// required: true
	// example: Lamassu IoT department 1
	OU string `json:"organization_unit"`

	// Country Name of the CA certificate
	// required: true
	// example: ES
	C string `json:"country"`

	// State of the CA certificate
	// required: true
	// example: Guipuzcoa
	ST string `json:"province"`

	// Locality of the CA certificate
	// required: true
	// example: Arrasate
	L string `json:"locality"`

	// Expiration period of the new emmited CA
	// required: true
	// example: 262800h
	TTL string `json:"ttl,omitempty"`
}

// CAs represents a list of CAs with minimum information
// swagger:model
type CAs struct {
	CAs []CA
}

type Secrets interface {
	GetCAs() (CAs, error)
	GetCACrt(caName string) (CACrt, error)
	CreateCA(caName string, ca CA) error
	DeleteCA(caName string) error
}
