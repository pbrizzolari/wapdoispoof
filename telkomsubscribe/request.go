package telkomsubscribe

type Envelope struct {
	SubscriberMSISDN   string `xml:"Body>Request>SubscriberMSISDN"`
	ContentDescription string `xml:"Body>Request>ContentDescription"`
	SOAPENV            string `xml:"SOAP-ENV,attr"`
	Username           string `xml:"Header>Security>UsernameToken>Username"`
	//Request                                            Request            `xml:"Body>Request"`
	ContentFrequency string `xml:"Body>Request>ContentFrequency"`
	PartnerMSISDN    string `xml:"Body>Request>PartnerMSISDN"`
	Password         string `xml:"Header>Security>UsernameToken>Password"`
	Action           string `xml:"Header>Action"`
	//From                                               From               `xml:"Header>From"`
	ContentID string `xml:"Body>Request>ContentID"`
	SenderID  string `xml:"Header>From>ReferenceParameters>SenderID"`
	//TransactionGroupIDReferenceParametersReplyToHeader string `xml:"Header>ReplyTo>ReferenceParameters>TransactionGroupID"`
	MessageID          string `xml:"Header>MessageID"`
	TransactionGroupID string `xml:"Header>From>ReferenceParameters>TransactionGroupID"`
	//ReplyTo                                            ReplyTo            `xml:"Header>ReplyTo"`
	PartnerID string `xml:"Body>Request>PartnerID"`
	//MetadataReferenceParametersReplyToHeader           Metadata           `xml:"Header>ReplyTo>ReferenceParameters>Metadata"`
	AddressReplyToHeader string `xml:"Header>ReplyTo>Address"`
	ContentType          string `xml:"Body>Request>ContentType"`
	Address              string `xml:"Header>From>Address"`
	//ParamMetadataReferenceParametersReplyToHeader      Param              `xml:"Header>ReplyTo>ReferenceParameters>Metadata>Param"`
	BearerType string `xml:"Body>Request>BearerType"`
	//Security                                           Security           `xml:"Header>Security"`
	//Metadata                                           Metadata           `xml:"Header>From>ReferenceParameters>Metadata"`
	TransactionID string `xml:"Body>Request>TransactionID"`
	//Param                                              Param              `xml:"Header>From>ReferenceParameters>Metadata>Param"`
	SenderIDReferenceParametersReplyToHeader string `xml:"Header>ReplyTo>ReferenceParameters>SenderID"`
	AmountInCents                            string `xml:"Body>Request>AmountInCents"`
}

type Action struct {
	Xsd  string `xml:"xsd,attr"`
	Xsi  string `xml:"xsi,attr"`
	Text string `xml:",chardata"`
	Add  string `xml:"add,attr"`
	Ns   string `xml:"ns,attr"`
	Ns0  string `xml:"ns0,attr"`
	Soap string `xml:"soap,attr"`
	Xs   string `xml:"xs,attr"`
}
type SenderIDReferenceParametersReplyToHeader struct {
	Xmlns string `xml:"xmlns,attr"`
	Text  string `xml:",chardata"`
}
type ParamMetadataReferenceParametersReplyToHeader struct {
	Key       string `xml:"Key,attr"`
	Qualifier string `xml:"Qualifier,attr"`
	Value     string `xml:"Value,attr"`
}
type MessageID struct {
	Soap string `xml:"soap,attr"`
	Xs   string `xml:"xs,attr"`
	Xsd  string `xml:"xsd,attr"`
	Xsi  string `xml:"xsi,attr"`
	Text string `xml:",chardata"`
	Add  string `xml:"add,attr"`
	Ns   string `xml:"ns,attr"`
	Ns0  string `xml:"ns0,attr"`
}
type SenderID struct {
	Xmlns string `xml:"xmlns,attr"`
	Text  string `xml:",chardata"`
}
type TransactionGroupID struct {
	Xmlns string `xml:"xmlns,attr"`
	Text  string `xml:",chardata"`
}
type Metadata struct {
	Ns string `xml:"ns,attr"`
}
type Param struct {
	Qualifier string `xml:"Qualifier,attr"`
	Value     string `xml:"Value,attr"`
	Key       string `xml:"Key,attr"`
}
type TransactionGroupIDReferenceParametersReplyToHeader struct {
	Xmlns string `xml:"xmlns,attr"`
	Text  string `xml:",chardata"`
}
type MetadataReferenceParametersReplyToHeader struct {
	Ns string `xml:"ns,attr"`
}
type Request struct {
	Soap string `xml:"soap,attr"`
	Tns  string `xml:"tns,attr"`
	Ns3  string `xml:"ns3,attr"`
	Xsi  string `xml:"xsi,attr"`
	Ns1  string `xml:"ns1,attr"`
	Xsd  string `xml:"xsd,attr"`
	Tns2 string `xml:"tns2,attr"`
	Ns   string `xml:"ns,attr"`
	Ns2  string `xml:"ns2,attr"`
	Tns3 string `xml:"tns3,attr"`
	Ns4  string `xml:"ns4,attr"`
	Tns1 string `xml:"tns1,attr"`
}
type Security struct {
	Xmlns string `xml:"xmlns,attr"`
}
type From struct {
	Soap string `xml:"soap,attr"`
	Xs   string `xml:"xs,attr"`
	Add  string `xml:"add,attr"`
	Xsd  string `xml:"xsd,attr"`
	Xsi  string `xml:"xsi,attr"`
	Ns   string `xml:"ns,attr"`
	Ns0  string `xml:"ns0,attr"`
}
type ReplyTo struct {
	Soap string `xml:"soap,attr"`
	Xsd  string `xml:"xsd,attr"`
	Ns   string `xml:"ns,attr"`
	Xs   string `xml:"xs,attr"`
	Xsi  string `xml:"xsi,attr"`
	Add  string `xml:"add,attr"`
	Ns0  string `xml:"ns0,attr"`
}
