package cellcsubscribe

type Envelope struct {
	Password        string `xml:"Header>Security>UsernameToken>Password"`
	ChargeInterval  string `xml:"Body>addSubscription>chargeInterval"`
	ServiceName     string `xml:"Body>addSubscription>serviceName"`
	ChargeCode      string `xml:"Body>addSubscription>chargeCode"`
	ContentType     string `xml:"Body>addSubscription>contentType"`
	Username        string `xml:"Header>Security>UsernameToken>Username"`
	Nonce           string `xml:"Header>Security>UsernameToken>Nonce"`
	ContentProvider string `xml:"Body>addSubscription>contentProvider"`
	Msisdn          string `xml:"Body>addSubscription>msisdn"`
	Created         string `xml:"Header>Security>UsernameToken>Created"`
	WaspReference   string `xml:"Body>addSubscription>waspReference"`
	WaspTid         string `xml:"Body>addSubscription>waspTid"`
	BearerType      string `xml:"Body>addSubscription>bearerType"`
}
