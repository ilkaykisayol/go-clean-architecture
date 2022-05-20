package sample

import "encoding/xml"

type PostSampleXmlProxyModel struct {
	XMLName    xml.Name `xml:"Sample"`
	SampleName string   `validate:"required" xml:"SampleName"`
	SampleType string   `validate:"required" xml:"SampleType"`
	SampleCode *int     `xml:"SampleCode,omitempty"`
}
