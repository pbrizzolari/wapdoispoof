package telkomcharge

type Envelope struct {
	Address string `xml:"Header>From>Address"`
	//MetadataReferenceParametersReplyToHeader           Metadata           `xml:"Header>ReplyTo>ReferenceParameters>Metadata"`
	//ParamMetadataReferenceParametersReplyToHeader      Param              `xml:"Header>ReplyTo>ReferenceParameters>Metadata>Param"`
	ContentDescription string `xml:"Body>Request>ContentDescription"`
	PartnerMSISDN      string `xml:"Body>Request>PartnerMSISDN"`
	ContentType        string `xml:"Body>Request>ContentType"`
	SenderID           string `xml:"Header>From>ReferenceParameters>SenderID"`
	TransactionGroupID string `xml:"Header>From>ReferenceParameters>TransactionGroupID"`
	//Param                                              Param              `xml:"Header>From>ReferenceParameters>Metadata>Param"`
	//SenderIDReferenceParametersReplyToHeader           SenderID           `xml:"Header>ReplyTo>ReferenceParameters>SenderID"`
	AddressReplyToHeader string `xml:"Header>ReplyTo>Address"`
	PartnerID            string `xml:"Body>Request>PartnerID"`
	MessageID            string `xml:"Header>MessageID"`
	//From                                               From               `xml:"Header>From"`
	//TransactionGroupIDReferenceParametersReplyToHeader TransactionGroupID `xml:"Header>ReplyTo>ReferenceParameters>TransactionGroupID"`
	ContentID     string `xml:"Body>Request>ContentID"`
	TransactionID string `xml:"Body>Request>TransactionID"`
	//Request                                            Request            `xml:"Body>Request"`
	SOAPENV string `xml:"SOAP-ENV,attr"`
	Action  string `xml:"Header>Action"`
	//Metadata         Metadata `xml:"Header>From>ReferenceParameters>Metadata"`
	//ReplyTo          ReplyTo `xml:"Header>ReplyTo"`
	TokenID          string `xml:"Body>Request>TokenID"`
	SubscriberMSISDN string `xml:"Body>Request>SubscriberMSISDN"`
	AmountInCents    string `xml:"Body>Request>AmountInCents"`
}

type From struct {
	Xsd  string `xml:"xsd,attr"`
	Xsi  string `xml:"xsi,attr"`
	Ns0  string `xml:"ns0,attr"`
	Xs   string `xml:"xs,attr"`
	Soap string `xml:"soap,attr"`
	Add  string `xml:"add,attr"`
	Ns   string `xml:"ns,attr"`
}
type SenderID struct {
	Xmlns string `xml:"xmlns,attr"`
	Text  string `xml:",chardata"`
}
type ReplyTo struct {
	Soap string `xml:"soap,attr"`
	Xs   string `xml:"xs,attr"`
	Xsi  string `xml:"xsi,attr"`
	Ns   string `xml:"ns,attr"`
	Xsd  string `xml:"xsd,attr"`
	Add  string `xml:"add,attr"`
	Ns0  string `xml:"ns0,attr"`
}
type TransactionGroupIDReferenceParametersReplyToHeader struct {
	Xmlns string `xml:"xmlns,attr"`
	Text  string `xml:",chardata"`
}
type ParamMetadataReferenceParametersReplyToHeader struct {
	Qualifier string `xml:"Qualifier,attr"`
	Value     string `xml:"Value,attr"`
	Key       string `xml:"Key,attr"`
}
type Request struct {
	Tns  string `xml:"tns,attr"`
	Soap string `xml:"soap,attr"`
	Xsd  string `xml:"xsd,attr"`
	Ns   string `xml:"ns,attr"`
	Xsi  string `xml:"xsi,attr"`
}
type MessageID struct {
	Ns   string `xml:"ns,attr"`
	Ns0  string `xml:"ns0,attr"`
	Soap string `xml:"soap,attr"`
	Xs   string `xml:"xs,attr"`
	Xsd  string `xml:"xsd,attr"`
	Xsi  string `xml:"xsi,attr"`
	Text string `xml:",chardata"`
	Add  string `xml:"add,attr"`
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
type Param struct {
	Key       string `xml:"Key,attr"`
	Qualifier string `xml:"Qualifier,attr"`
	Value     string `xml:"Value,attr"`
}
type MetadataReferenceParametersReplyToHeader struct {
	Ns string `xml:"ns,attr"`
}
type SenderIDReferenceParametersReplyToHeader struct {
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
