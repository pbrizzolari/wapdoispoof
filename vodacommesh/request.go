package vodacommesh

type Envelope struct {
	Password              string               `xml:"Header>ServiceAuth>Password"`
	BillingInformation    BillingInformation   `xml:"Body>PerformServiceOperationRequest>BillingInformation"`
	MSISDN                string               `xml:"Body>PerformServiceOperationRequest>TSORequest>MSISDN"`
	WaspSubscriptionData  WaspSubscriptionData `xml:"Body>PerformServiceOperationRequest>TSORequest>WaspSubscriptionData"`
	Cbs                   string               `xml:"cbs,attr"`
	Soapenv               string               `xml:"soapenv,attr"`
	Username              string               `xml:"Header>ServiceAuth>Username"`
	OperationCode         string               `xml:"Body>PerformServiceOperationRequest>TSORequest>OperationCode"`
	ProvisioningPoint     string               `xml:"Body>PerformServiceOperationRequest>TSORequest>ProvisioningPoint"`
	ServiceCode           string               `xml:"Body>PerformServiceOperationRequest>TSORequest>ServiceCode"`
	ServiceDerivativeCode string               `xml:"Body>PerformServiceOperationRequest>TSORequest>ServiceDerivativeCode"`
}

type BillingInformation struct {
	CoID             string `xml:"CoID,attr"`
	CoKey            string `xml:"CoKey,attr"`
	FailedAuthReqRef string `xml:"FailedAuthReqRef,attr"`
}
type WaspSubscriptionData struct {
	BillingFrequency string `xml:"BillingFrequency,attr"`
	Channel          string `xml:"Channel,attr"`
	CustomMessage    string `xml:"CustomMessage,attr"`
	ReferenceNumber  string `xml:"ReferenceNumber,attr"`
	ServiceID        string `xml:"ServiceID,attr"`
	ServiceName      string `xml:"ServiceName,attr"`
	WaspID           string `xml:"WaspID,attr"`
	WaspName         string `xml:"WaspName,attr"`
	Price            string `xml:"Price,attr"`
}
