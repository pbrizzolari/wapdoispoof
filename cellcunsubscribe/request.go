package cellcunsubscribe

type Envelope struct {
	Password  string `xml:"Header>Security>UsernameToken>Password"`
	Msisdn    string `xml:"Body>cancelSubscription>msisdn"`
	ServiceID string `xml:"Body>cancelSubscription>serviceID"`
	Wasp      string `xml:"wasp,attr"`
	Nonce     string `xml:"Header>Security>UsernameToken>Nonce"`
	Created   string `xml:"Header>Security>UsernameToken>Created"`
	Username  string `xml:"Header>Security>UsernameToken>Username"`
	WaspTid   string `xml:"Body>cancelSubscription>waspTid"`
}
