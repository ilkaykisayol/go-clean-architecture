package sample

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"
)

type ISampleXmlProxy interface {
	PostSampleXml(ch chan *PostSampleXmlProxyResponse, model *PostSampleXmlProxyModel)
}

type SampleXmlProxy struct {
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	cachr       cacher.ICacher
	client      *http.Client
	baseUrl     *url.URL
	timeout     time.Duration
}

// NewSampleProxy
// Returns a new NewSampleXmlProxy.
func NewSampleXmlProxy(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher) ISampleXmlProxy {
	url, urlErr := url.Parse(environment.Get(env.SampleProxyUrl))
	if urlErr != nil {
		loggr.Panic("Couldn't convert SAMPLE_XML_PROXY_URL environment variable to url.URL !")
	}

	timeout, timeoutErr := strconv.Atoi(environment.Get(env.SampleProxyTimeout))
	if timeoutErr != nil {
		loggr.Panic("Couldn't convert SAMPLE_XML_PROXY_TIMEOUT environment variable to int !")
	}

	service := SampleXmlProxy{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
		client:      &http.Client{Timeout: time.Second * time.Duration(timeout)},
		baseUrl:     url,
		timeout:     time.Second * time.Duration(timeout),
	}
	return &service
}

// PostSampleXml
// Post sample xml example
func (p *SampleXmlProxy) PostSampleXml(ch chan *PostSampleXmlProxyResponse, model *PostSampleXmlProxyModel) {
	modelErr := p.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &PostSampleXmlProxyResponse{Error: modelErr}
		return
	}

	baseModelAttrs := []xml.Attr{
		{
			Name: xml.Name{
				Space: "",
				Local: "xmlns:a",
			},
			Value: "http://sample.com/",
		},
		{
			Name: xml.Name{
				Space: "",
				Local: "xmlns:b",
			},
			Value: "http://sample.net/",
		},
	}

	baseModel := NewSampleBaseXmlProxyRequestModel(baseModelAttrs)

	baseModel.Body.Model = model

	baseModelXmlOut, xmlErr := xml.Marshal(baseModel)

	if xmlErr != nil {
		ch <- &PostSampleXmlProxyResponse{Error: xmlErr}
	}

	// You can override the default timeout here.
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)

	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseUrl.String()+"", bytes.NewReader(baseModelXmlOut))

	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	request.Header.Set("SOAPAction", "http://sample.org/someAction")

	response, err := p.client.Do(request)
	if err != nil {
		ch <- &PostSampleXmlProxyResponse{Error: err}
		return
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)

	if readErr != nil {
		ch <- &PostSampleXmlProxyResponse{Error: readErr}
		return
	}

	attr := []xml.Attr{
		{
			Name: xml.Name{
				Space: "",
				Local: "xmlns:s",
			},
			Value: "http://sample.org/",
		},
	}

	responseModel := NewSampleXmlProxyResponseModel(attr)

	unmarshalErr := xml.Unmarshal(body, &responseModel)

	if unmarshalErr != nil {
		ch <- &PostSampleXmlProxyResponse{Error: unmarshalErr}
		return
	}

	ch <- &PostSampleXmlProxyResponse{
		IsSuccess: responseModel.Body.IsSuccess,
		Message:   responseModel.Body.Message,
	}
}
