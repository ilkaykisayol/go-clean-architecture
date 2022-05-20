package sample

import (
	"encoding/xml"
)

type SampleBaseXmlProxyRequestModel struct {
	XMLName    xml.Name                        `xml:"Soapenv:Envelope"`
	Body       SampleBaseProxyRequestBodyModel `xml:"Soapenv:Body"`
	Attributes []xml.Attr                      `xml:",any,attr"`
}

type SampleBaseProxyRequestBodyModel struct {
	XMLName  xml.Name    `xml:"Soapenv:Body"`
	Username string      `xml:"Username"`
	Password string      `xml:"Password"`
	Model    interface{} `xml:"Model"`
}

func NewSampleBaseXmlProxyRequestModel(attributes []xml.Attr) SampleBaseXmlProxyRequestModel {
	return SampleBaseXmlProxyRequestModel{
		Attributes: attributes,
	}
}
