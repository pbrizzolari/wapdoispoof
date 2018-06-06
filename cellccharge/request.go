package cellccharge

type Envelope struct {
	Created   string `xml:"Header>Security>UsernameToken>Created"`
	Username  string `xml:"Header>Security>UsernameToken>Username"`
	Password  string `xml:"Header>Security>UsernameToken>Password"`
	WaspTid   string `xml:"Body>chargeSubscriber>waspTid"`
	ServiceID string `xml:"Body>chargeSubscriber>serviceID"`
	Soapenv   string `xml:"soapenv,attr"`
	Wasp      string `xml:"wasp,attr"`
	Nonce     string `xml:"Header>Security>UsernameToken>Nonce"`
	Msisdn    string `xml:"Body>chargeSubscriber>msisdn"`
}
