package sample

import (
	"encoding/xml"
)

type PostSampleXmlProxyResponse struct {
	Error     error `json:"-"`
	IsSuccess bool
	Message   string
}

type SampleXmlProxyResponseModel struct {
	XMLName    xml.Name `xml:"Envelope"`
	Body       SampleXmlProxyResponseBodyModel
	Attributes []xml.Attr `xml:",any,attr"`
}

type SampleXmlProxyResponseBodyModel struct {
	XMLName   xml.Name `xml:"Body"`
	IsSuccess bool     `xml:"isSuccess"`
	Message   string   `xml:"message"`
}

func NewSampleXmlProxyResponseModel(attributes []xml.Attr) SampleXmlProxyResponseModel {
	return SampleXmlProxyResponseModel{
		Body:       SampleXmlProxyResponseBodyModel{},
		Attributes: attributes,
	}
}
