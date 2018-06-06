package mtnrpc

type MethodCall struct {
	XMLName    xml.Name `xml:"MethodCall"`
	MethodName string   `xml:"methodName"`
	Params     []param  `xml:"params"`
}

type Param struct {
	XMLName xml.Name `xml:"param"`
	Value   BaseType `xml:"value"`
}

type BaseType struct {
	Value string `xml:",chardata"`
}

type StructType struct {
	BaseType
	XMLName xml.Name `xml:"struct"`
	Members []Member ``
}

type IntType struct {
	BaseType
	XMLName xml.Name `xml:"int"`
}

type StringType struct {
	BaseType
	XMLName xml.Name `xml:"string"`
}

type Member struct {
	XMLName xml.Name `xml:"member"`
	Name    string   `xml:"name"`
	Value   BaseType `xml:"value"`
}
