package telkomunsubscribe

type Envelope struct {
	AddressReplyToHeader string `xml:"Header>ReplyTo>Address"`
	PartnerID            string `xml:"Body>Request>PartnerID"`
	//Security                                 Security  `xml:"Header>Security"`
	PartnerMSISDN string `xml:"Body>Request>PartnerMSISDN"`
	TokenID       string `xml:"Body>Request>TokenID"`
	To            string `xml:"Header>To"`
	//From                                     From      `xml:"Header>From"`
	//ReplyTo                                  ReplyTo   `xml:"Header>ReplyTo"`
	Username  string `xml:"Header>Security>UsernameToken>Username"`
	Password  string `xml:"Header>Security>UsernameToken>Password"`
	Action    string `xml:"Header>Action"`
	SOAPENV   string `xml:"SOAP-ENV,attr"`
	MessageID string `xml:"Header>MessageID"`
	Address   string `xml:"Header>From>Address"`
	SenderID  string `xml:"Header>From>ReferenceParameters>SenderID"`
	//SenderIDReferenceParametersReplyToHeader SenderID  `xml:"Header>ReplyTo>ReferenceParameters>SenderID"`
	//Request                                  Request   `xml:"Body>Request"`
	SubscriberMSISDN string `xml:"Body>Request>SubscriberMSISDN"`
	Soapenv          string `xml:"soapenv,attr"`
}

type Security struct {
	Xmlns string `xml:"xmlns,attr"`
}
type Action struct {
	Add  string `xml:"add,attr"`
	Text string `xml:",chardata"`
}
type Request struct {
	Sch string `xml:"sch,attr"`
}
type MessageID struct {
	Add  string `xml:"add,attr"`
	Text string `xml:",chardata"`
}
type SenderID struct {
	Ns   string `xml:"ns,attr"`
	Text string `xml:",chardata"`
}
type From struct {
	Add string `xml:"add,attr"`
}
type SenderIDReferenceParametersReplyToHeader struct {
	Text string `xml:",chardata"`
	Ns   string `xml:"ns,attr"`
}
type To struct {
	Add  string `xml:"add,attr"`
	Text string `xml:",chardata"`
}
type ReplyTo struct {
	Add string `xml:"add,attr"`
}
