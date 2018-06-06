package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/basebone/wapdoispoof/cellccharge"
	"github.com/basebone/wapdoispoof/cellcsubscribe"
	"github.com/basebone/wapdoispoof/cellcunsubscribe"
	"github.com/basebone/wapdoispoof/meshresponse"
	"github.com/basebone/wapdoispoof/telkomcharge"
	"github.com/basebone/wapdoispoof/telkomsubscribe"
	"github.com/basebone/wapdoispoof/telkomunsubscribe"
	"github.com/basebone/wapdoispoof/vodacommesh"
	gxml "github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
)

type Method string

const (
	M_GET     Method = "GET"
	M_PUT            = "PUT"
	M_POST           = "POST"
	M_DELETE         = "DELETE"
	M_OPTIONS        = "OPTIONS"
	M_NONE           = ""
)

type RouteGroup struct {
	Prefix string
	Routes []Route
}

type Route struct {
	Path    string
	Method  Method
	Handler func(http.ResponseWriter, *http.Request)
}

type HttpServer struct {
	listener net.Listener
	server   *http.Server
	router   *mux.Router
	app      *App
}

type App struct {
	RPC            *rpc.Server
	FailureCounter int
	Client         *redis.Client
}

var port = flag.String("listener", ":21000", "listener")

func main() {
	flag.Parse()
	a := &App{}
	a.FailureCounter = 0
	a.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := a.Client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	jsonString, err := json.Marshal(make(map[string]bool))
	fmt.Println(err)
	err = a.Client.Set("sessiontokens", jsonString, 0).Err()

	router := mux.NewRouter()
	httpServer := &HttpServer{
		server: &http.Server{
			Addr:         *port,
			Handler:      router,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		router: router,
		app:    a,
	}

	a.RPC = rpc.NewServer()
	xmlrpcCodec := gxml.NewCodec()
	a.RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	a.RPC.RegisterService(a, "")

	httpServer.server.SetKeepAlivesEnabled(true)

	httpServer.router.HandleFunc("/", httpServer.app.Options)
	routeBase := httpServer.router.PathPrefix("/").Subrouter()
	for _, rg := range httpServer.app.GetRoutes() {
		if rg.Prefix != "" {
			fmt.Println("Add prefix", rg.Prefix)
			subRouter := routeBase.PathPrefix(rg.Prefix).Subrouter()
			for _, route := range rg.Routes {
				if route.Path == "" {
					fmt.Println("\tAdd base route", route.Path)
					r := routeBase.HandleFunc(rg.Prefix, route.Handler)
					if route.Method != "" {
						r.Methods(string(route.Method))
					}
				} else {
					fmt.Println("\tAdd route", route.Path)
					r := subRouter.HandleFunc(route.Path, route.Handler)
					if route.Method != "" {
						r.Methods(string(route.Method))
					}
				}
			}
		} else {
			for _, route := range rg.Routes {
				fmt.Println("Add route", route.Path)
				r := routeBase.HandleFunc(route.Path, route.Handler)
				if route.Method != "" {
					r.Methods(string(route.Method))
				}

			}
		}
	}
	r := routeBase.Handle("/mtnrpc", a.RPC)
	r.Methods(string(M_POST))
	httpServer.router.NotFoundHandler = http.HandlerFunc(httpServer.app.NotFound)
	httpServer.router.Schemes("http")
	go func() {
		err := httpServer.server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			// panic(err)
		}
	}()
	fmt.Scanln()
}

func (a *App) WriteNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func (a *App) WriteUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (a *App) WriteForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func (a *App) Options(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, " ", r.URL)
	// w.Header().Add("Access-Control-Allow-Origin", a.hosturl)
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type,X-XSRF-TOKEN")
	// w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
}

func (a *App) NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("NotFound")
	//w.Header().Add("Access-Control-Allow-Origin", a.hosturl)
	//w.Header().Add("Access-Control-Allow-Credentials", "true")
	http.NotFound(w, r)
}

func (a *App) GetRoutes() []RouteGroup {
	return []RouteGroup{
		RouteGroup{"/vod",
			[]Route{
				Route{"", M_GET, a.Vodacom},
			},
		},
		RouteGroup{"/vodmesh",
			[]Route{
				Route{"", M_POST, a.VodacomMesh},
			},
		},
		RouteGroup{"/mtn",
			[]Route{
				Route{"", M_GET, a.Mtn},
			},
		},
		RouteGroup{"/mtntbb",
			[]Route{
				Route{"", M_POST, a.TBBDistribute},
			},
		},
		RouteGroup{"/cellc",
			[]Route{
				Route{"", M_GET, a.CellC},
			},
		},
		RouteGroup{"/ccservice",
			[]Route{
				Route{"", M_POST, a.CellCService},
			},
		},
		RouteGroup{"/telkom",
			[]Route{
				Route{"", M_GET, a.Telkom},
			},
		},
		RouteGroup{"/telkomservices",
			[]Route{
				Route{"/RequestCharge", M_POST, a.TelkomServiceCharge},
				Route{"/RequestDoubleOptIn", M_POST, a.TelkomServiceSubscribe},
				Route{"/RequestSubscriberOptOut", M_POST, a.TelkomServiceUnsubscribe},
			},
		},
		RouteGroup{"/safcom/SubscribeManageService/services/SubscribeManage",
			[]Route{
				Route{"", M_POST, a.SafcomSubscribe},
			},
		},
	}
}

func (a *App) Vodacom(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	v := r.URL.Query()
	ServiceCode := v.Get("ServiceCode")
	ServiceDerivativeCode := v.Get("ServiceDerivativeCode")
	ReferenceNumber := v.Get("ReferenceNumber")
	OperationCode := v.Get("OperationCode")

	u, _ := url.Parse("http://localhost:8900/opvod")
	q := u.Query()
	//Fields from the request string
	q.Set("ServiceCode", ServiceCode)
	q.Set("ServiceDerivativeCode", ServiceDerivativeCode)
	q.Set("OperationCode", OperationCode)
	q.Set("ReferenceNumber", ReferenceNumber)

	//Fields that are generated on this service
	q.Set("MSISDN", "27794986880")
	q.Set("Result", "Confirmed")
	q.Set("StatusCode", "2")
	q.Set("ErrorCode", "2")
	q.Set("ErrorDescription", "Test Failure")
	if strings.ToLower(OperationCode) == "purchase" {
		q.Set("ExtAssetId", "lkjasdihgwmlkalasclkwdijadwo31")
	} else {
		q.Set("SubscriptionStartDate", time.Now().Format("2006-01-02"))
	}

	//encode url and send
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), 303)
}

func (a *App) Mtn(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	http.Redirect(w, r, "http://mtnurl", 303)
}

func (a *App) VodacomMesh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	request := &vodacommesh.Envelope{}
	xml.Unmarshal(bytes, request)
	fmt.Println("------------------------------------------------------------")
	fmt.Println(string(bytes))
	fmt.Println(request)
	fmt.Println("------------------------------------------------------------")

	response := &meshresponse.Envelope{}
	response.TSOID = fmt.Sprintf("V-%s", request.WaspSubscriptionData.ReferenceNumber)
	response.ServiceCode = request.ServiceCode
	response.ServiceDerivativeCode = request.ServiceDerivativeCode
	response.MSISDN = request.MSISDN
	response.WaspSubscriptionData.WaspID = request.WaspSubscriptionData.WaspID
	response.WaspSubscriptionData.WaspName = request.WaspSubscriptionData.WaspName
	response.WaspSubscriptionData.Channel = request.WaspSubscriptionData.Channel
	response.WaspSubscriptionData.ReferenceNumber = request.WaspSubscriptionData.ReferenceNumber
	response.WaspSubscriptionData.ServiceID = request.WaspSubscriptionData.ServiceID
	response.WaspSubscriptionData.ServiceName = request.WaspSubscriptionData.ServiceName
	response.WaspSubscriptionData.BillingFrequency = request.WaspSubscriptionData.BillingFrequency
	response.WaspSubscriptionData.Price = request.WaspSubscriptionData.Price

	response.StatusCode = 1
	response.ErrorCode = 1
	response.ErrorDescription = "SUCCESS"

	reply, err := xml.Marshal(response)

	fmt.Println(err)

	fmt.Println(string(reply))

	w.Write(reply)
}

func (a *App) CellCService(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	Body := string(bytes)
	if strings.Index(Body, "addSubscription>") != -1 {
		request := &cellcsubscribe.Envelope{}
		xml.Unmarshal(bytes, request)
		fmt.Println("------------------------------------------------------------")
		fmt.Println(Body)
		fmt.Println(string(bytes))
		fmt.Println(request)
		fmt.Println("------------------------------------------------------------")

		serviceID := "34523452345"
		result := 0
		ReplyString := fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
										<soap:Body>
											<ns2:addSubscriptionResponse xmlns:ns2="http://wasp.doi.soap.protocol.cellc.co.za">
												<return>
													<serviceID>%s</serviceID>
													<Result>%d</Result>
													<WaspReference>%s</WaspReference>
													<waspTid>%s</waspTid>
												</return>
											</ns2:addSubscriptionResponse>
										</soap:Body>
									</soap:Envelope>`, serviceID, result, request.WaspReference, request.WaspTid)

		fmt.Println(ReplyString)

		w.Write([]byte(ReplyString))
	} else if strings.Index(Body, "chargeSubscriber>") != -1 {
		request := &cellccharge.Envelope{}
		xml.Unmarshal(bytes, request)

		result := 0

		ReplyString := fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
										<soap:Body>
											<ns2:chargeSubscriberResponse xmlns:ns2="http://wasp.doi.soap.protocol.cellc.co.za">
												<return>
													<serviceID>%s</serviceID>
													<Result>%d</Result>
													<WaspReference>%s</WaspReference>
													<waspTid>%s</waspTid>
												</return>
											</ns2:chargeSubscriberResponse>
										</soap:Body>
									</soap:Envelope>`, request.ServiceID, result, request.WaspTid, request.WaspTid)

		fmt.Println(ReplyString)

		w.Write([]byte(ReplyString))

	} else if strings.Index(Body, "cancelSubscription>") != -1 {
		request := &cellcunsubscribe.Envelope{}
		xml.Unmarshal(bytes, request)

		status := "CANCELLED"

		ReplyString := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wasp="http://wasp.doi.soap.protocol.cellc.co.za">
										<soapenv:Header>
											<wasp:ServiceAuth>
												<Username>user1 </Username>
												<Password>pass1</Password>
											</wasp:ServiceAuth>
										</soapenv:Header>
										<soapenv:Body>
											<wasp:cancelSubscriptionResponse>
												<!--One or more repetitions:-->
												<services>
													<serviceID>%s</serviceID>
													<timeCreated>%s</timeCreated>
													<MSISDN>%s</MSISDN>
													<serviceName>Unknown</serviceName>
													<contentProvider>Unknown</contentProvider>
													<chargeCode>Unknown</chargeCode>
													<chargeInterval>Unknown</chargeInterval>
													<status>%s</status>
													<lastBilled>%s</lastBilled>
													<WASPReference>%s</WASPReference>
													<waspTid>%s</waspTid>
												</services>
											</wasp:cancelSubscriptionResponse>
										</soapenv:Body>
									</soapenv:Envelope>`, request.ServiceID, time.Now().Add(-48*time.Hour).Format("2006-01-02 15:04:05"), request.Msisdn, status,
			time.Now().Add(-24*time.Hour).Format("2006-01-02 15:04:05"), request.WaspTid, request.WaspTid)

		fmt.Println(ReplyString)

		w.Write([]byte(ReplyString))
	} else {
		fmt.Println("No understandable type found. Cell C Service")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *App) CellC(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	v := r.URL.Query()
	Url := v.Get("url")
	WaspTid := v.Get("wasptid")
	DoiServiceId := v.Get("doiserviceid")
	uri, _ := url.PathUnescape(Url)
	u, _ := url.Parse(uri)
	q := u.Query()

	q.Set("wasptid", WaspTid)
	q.Set("doiserviceid", DoiServiceId)
	q.Set("doistatus", "Active")
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), 303)
}

func (a *App) TelkomServiceCharge(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
	request := &telkomcharge.Envelope{}
	xml.Unmarshal(bytes, request)
	fmt.Println(request)

	TelkomId := fmt.Sprintf("Telkom:%s", request.TransactionID)

	ReplyString := fmt.Sprintf(`<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
								<SOAP-ENV:Header>
								<RelatesTo RelationshipType="{http://www.w3.org/2005/08/addressing}Reply" xmlns="http://www.w3.org/2005/08/addressing">%s</RelatesTo>
								<Action xmlns="http://www.w3.org/2005/08/addressing">http://eai.telkom.co.za/services/RequestCharge/20121001/Reply</Action>
								<SenderID wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">%s</SenderID>
								<TransactionGroupID wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">%s</TransactionGroupID>
								<ns:Metadata wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns:ns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">
								<ns:Param Key="aKey" Qualifier="aQualifier" Value="aValue"/>
								</ns:Metadata>
								</SOAP-ENV:Header>
								<SOAP-ENV:Body>
								<ns:Response xmlns:ns="http://eai.telkom.co.za/schemas/RequestCharge/20121001/DataModel/Schema.xsd" 
								xmlns:tns="http://eai.telkom.co.za/services/RequestCharge/20121001">
								<ns0:Result xmlns:ns0="http://eai.telkom.co.za/schemas/RequestCharge/20121001/DataModel/Schema.xsd">
								<ns0:ResultCode>0</ns0:ResultCode>
								</ns0:Result>
								<ns:Payload>
								<ns:TransactionID>%s</ns:TransactionID>
								<ns:TelkomInternalReferenceID>%s</ns:TelkomInternalReferenceID>
								</ns:Payload>
								</ns:Response>
								</SOAP-ENV:Body>
								</SOAP-ENV:Envelope>`, request.TransactionID, request.SenderID, request.TransactionGroupID, request.TransactionID, TelkomId)

	fmt.Println(ReplyString)

	w.Write([]byte(ReplyString))
}

func (a *App) TelkomServiceSubscribe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
	request := &telkomsubscribe.Envelope{}
	xml.Unmarshal(bytes, request)
	fmt.Println(request)
	//Body := string(bytes)
	ResultCode := 0
	ResultMsg := "CONFIRMED"
	ReplyString := fmt.Sprintf(`<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
								<SOAP-ENV:Header>
								<RelatesTo RelationshipType="{http://www.w3.org/2005/08/addressing}Reply" xmlns="http://www.w3.org/2005/08/addressing">%s</RelatesTo>
								<Action xmlns="http://www.w3.org/2005/08/addressing">http://eai.telkom.co.za/services/RequestDoubleOptIn/20121001/Reply</Action>
								<SenderID wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">%s</SenderID>
								<TransactionGroupID wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">%s</TransactionGroupID>
								<ns:Metadata wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns:ns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">
								<ns:Param Key="aKey" Qualifier="aQualifier" Value="aValue"/>
								</ns:Metadata>
								</SOAP-ENV:Header>
								<SOAP-ENV:Body>
								<ns1:Response xmlns:ns="http://eai.telkom.co.za/schemas/GetDoubleOptInHistory/20121001/DataModel/Schema.xsd" 
								xmlns:ns1="http://eai.telkom.co.za/schemas/RequestDoubleOptIn/20121001/DataModel/Schema.xsd" 
								xmlns:tns="http://eai.telkom.co.za/services/GetDoubleOptInHistory/20121001" 
								xmlns:tns1="http://eai.telkom.co.za/services/RequestDoubleOptIn/20121001">
								<ns0:Result xmlns:ns0="http://eai.telkom.co.za/schemas/RequestDoubleOptIn/20121001/DataModel/Schema.xsd">
								<ns0:ResultCode>%d</ns0:ResultCode>
								<ns0:ResultMsg>%s</ns0:ResultMsg>
								</ns0:Result>
								</ns1:Response>
								</SOAP-ENV:Body>
								</SOAP-ENV:Envelope>`, request.TransactionID, request.PartnerID, request.TransactionGroupID, ResultCode, ResultMsg)

	fmt.Println(ReplyString)

	w.Write([]byte(ReplyString))
}

func (a *App) TelkomServiceUnsubscribe(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
	request := &telkomunsubscribe.Envelope{}
	xml.Unmarshal(bytes, request)
	fmt.Println(request)

	ResultCode := 0
	ResultMsgCode := "EBO-10015"
	ResultMsg := ""
	ReplyString := fmt.Sprintf(`<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
								<SOAP-ENV:Header>
								<RelatesTo RelationshipType="{http://www.w3.org/2005/08/addressing}Reply" xmlns="http://www.w3.org/2005/08/addressing">%s</RelatesTo>
								<Action xmlns="http://www.w3.org/2005/08/addressing">http://eai.telkom.co.za/services/RequestSubscriberOptOut/20121001/Reply</Action>
								<ns:SenderID wsa:IsReferenceParameter="true" xmlns:wsa="http://www.w3.org/2005/08/addressing" 
								xmlns:ns="http://eai.telkom.co.za/EnterpriseServiceMetaData/20110801">%s</ns:SenderID>
								</SOAP-ENV:Header>
								<SOAP-ENV:Body>
								<ns0:Response xmlns:ns0="http://eai.telkom.co.za/schemas/RequestSubscriberOptOut/20121001/DataModel/Schema.xsd">
								<ns0:Result>
								<ns0:ResultCode>%d</ns0:ResultCode>
								<ns0:ResultMsgCode>%s</ns0:ResultMsgCode>
								<ns0:ResultMsg>%s</ns0:ResultMsg>
								</ns0:Result>
								</ns0:Response>
								</SOAP-ENV:Body>
								</SOAP-ENV:Envelope>`, request.MessageID, request.PartnerID, ResultCode, ResultMsgCode, ResultMsg)

	fmt.Println(ReplyString)

	w.Write([]byte(ReplyString))
}

func (a *App) TBBDistribute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	req := string(bytes)
	fmt.Println(req)
	MethodName := req[strings.Index(req, "<methodName>")+12 : strings.Index(req, "</methodName>")]
	switch MethodName {
	case "permission":
		a.TBBPermissionToLogin(w, r)
	case "login":
		a.TBBLogin(w, r, req)
	case "debitToken":
		a.TBBBill(w, r, req)
	}
}

func (a *App) TBBPermissionToLogin(w http.ResponseWriter, r *http.Request) {
	ReplyString := fmt.Sprintf(`<methodResponse>
								<params>
									<param>
									<value>
										<struct>
										<member>
											<name>key</name>
											<value><string>KEY</string> </value>
										</member>
										</struct>
									</value>
									</param>
								</params>
								</methodResponse>`)

	fmt.Println(ReplyString)

	w.Write([]byte(strings.Replace(strings.Replace(ReplyString, "\n", "", -1), "\t", "", -1)))
}

func (a *App) TBBLogin(w http.ResponseWriter, r *http.Request, req string) {
	LinkId := req[strings.Index(req, "<int>")+5 : strings.Index(req, "</int>")]
	CurrentToken := LinkId + "---" + r.Host + "-" + time.Now().Format("15:04:05")
	CTString, err := a.Client.Get("sessiontokens").Result()
	var CurrentTokens map[string]bool
	err = json.Unmarshal([]byte(CTString), &CurrentTokens)
	if err != nil {
		fmt.Println("error:", err)
	}
	CurrentTokens[CurrentToken] = true
	fmt.Println("Current Tokens", CurrentTokens)
	b, err := json.Marshal(CurrentTokens)
	if err != nil {
		fmt.Println("error:", err)
	}
	err = a.Client.Set("sessiontokens", string(b), 0).Err()
	if err != nil {
		fmt.Println("error:", err)
	}
	ReplyString := fmt.Sprintf(`<methodResponse>
									<params>
										<param>
											<value>
												<struct>
													<member>
														<name>sessionToken</name> 
														<value><string>%s</string></value>
													</member>
													<member>
														<name>allowedMethods</name>
														<value>
															<array>
																<data>
																	<value> <string>ALLOWED-METHODS</string></value>
																</data>
															</array>
														</value>
													</member>
												</struct>
											</value>
										</param>
									</params>
								</methodResponse>`, CurrentToken)

	fmt.Println(ReplyString)

	w.Write([]byte(strings.Replace(strings.Replace(ReplyString, "\n", "", -1), "\t", "", -1)))
}

func (a *App) TBBBill(w http.ResponseWriter, r *http.Request, req string) {
	StatusCode := 0
	Token := req[strings.Index(req, "<param><value><string>")+22 : strings.Index(req, "</string></value></param>")]

	CTString, err := a.Client.Get("sessiontokens").Result()
	if err != nil {
		fmt.Println("error:", err)
	}
	var CurrentTokens map[string]bool
	err = json.Unmarshal([]byte(CTString), &CurrentTokens)
	if err != nil {
		fmt.Println("error:", err)
	}
	if a.FailureCounter%5 == 0 || !CurrentTokens[Token] {
		StatusCode = 4713
	}
	a.FailureCounter += 1
	ReplyString := fmt.Sprintf(`<methodResponse>
									<params>
										<param>
										<value>
											<struct>
												<member>
													<name>msn</name> 
													<value><string>MSN</string></value>
												</member>
												<member>
													<name>rsn</name>
													<value><string>RSN</string></value> 
												</member>
												<member>
													<name>result</name>
													<value><string>RESULT</string></value> 
												</member>
												<member>
													<name>statusCode</name>
													<value><int>%d</int></value> 
												</member>
											</struct>
										</value>
										</param>
									</params>
								</methodResponse>`)

	//fmt.Println(ReplyString)
	if StatusCode == 4713 {
		fmt.Println("Replied with failure")
		fmt.Println("Received Token %s Internal Token %t", Token, CurrentTokens)
		delete(CurrentTokens, Token)
		jsonString, err := json.Marshal(CurrentTokens)
		fmt.Println(err)
		err = a.Client.Set("sessiontokens", jsonString, 0).Err()
	}

	w.Write([]byte(strings.Replace(strings.Replace(ReplyString, "\n", "", -1), "\t", "", -1)))
}

func (a *App) SafcomSubscribe(w http.ResponseWriter, r *http.Request) {
	ReplyString := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8" ?>
								<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
									<soapenv:Body>
										<ns1:subscribeProductResponse xmlns:ns1="http://www.csapi.org/schema/parlayx/subscribe/manage/v1_0/local">
											<ns1:subscribeProductRsp>
												<result>22007201</result>
											</ns1:subscribeProductRsp>
										</ns1:subscribeProductResponse>
									</soapenv:Body>
								</soapenv:Envelope>`)

	w.Write([]byte(strings.Replace(strings.Replace(ReplyString, "\n", "", -1), "\t", "", -1)))
}

func (a *App) Telkom(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	v := r.URL.Query()
	Url := v.Get("url")
	uri, _ := url.PathUnescape(Url)
	u, _ := url.Parse(uri)
	q := u.Query()

	q.Set("doistatus", "Active")
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), 303)
}

func Say(r *http.Request, args *struct{ Who string }, reply *struct{ Message string }) error {
	log.Println("Say", args.Who)
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}
