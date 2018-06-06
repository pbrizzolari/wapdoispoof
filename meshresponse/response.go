package meshresponse

type Envelope struct {
	TSOID                 string               `xml:"Body>PerformServiceOperationResponse>TSOID"`
	CustomerSMS           string               `xml:"Body>PerformServiceOperationResponse>TSOResponse>Service>CustomerSMS"`
	Soapenv               string               `xml:"soapenv,attr"`
	MSISDN                string               `xml:"Body>PerformServiceOperationResponse>TSOResponse>Service>MSISDN"`
	Header                Header               `xml:"Header"`
	SubscriptionCode      string               `xml:"Body>PerformServiceOperationResponse>TSOResponse>Service>SubscriptionCode"`
	WaspSubscriptionData  WaspSubscriptionData `xml:"Body>PerformServiceOperationResponse>TSOResponse>WaspSubscriptionData"`
	StatusCode            int                  `xml:"Body>PerformServiceOperationResponse>TSOResult>StatusCode"`
	ErrorCode             int                  `xml:"Body>PerformServiceOperationResponse>TSOResult>ErrorCode"`
	ErrorDescription      string               `xml:"Body>PerformServiceOperationResponse>TSOResult>ErrorDescription"`
	ServiceCode           string               `xml:"Body>PerformServiceOperationResponse>TSOResponse>Service>ServiceCode"`
	ServiceDerivativeCode string               `xml:"Body>PerformServiceOperationResponse>TSOResponse>Service>ServiceDerivativeCode"`
}

type Header struct {
	Cbs string `xml:"cbs,attr"`
}
type PerformServiceOperationResponse struct {
	Cbs string `xml:"cbs,attr"`
}
type WaspSubscriptionData struct {
	WaspName         string `xml:"WaspName,attr"`
	ServiceName      string `xml:"ServiceName,attr"`
	URL              string `xml:"URL,attr"`
	WaspID           string `xml:"WaspID,attr"`
	BillingFrequency string `xml:"BillingFrequency,attr"`
	Channel          string `xml:"Channel,attr"`
	Price            string `xml:"Price,attr"`
	ReferenceNumber  string `xml:"ReferenceNumber,attr"`
	ServiceID        string `xml:"ServiceID,attr"`
}
