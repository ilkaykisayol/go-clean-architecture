package sample

import (
	"context"
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

type ISampleProxy interface {
	GetGoogle(ch chan *GetSampleProxyResponse, model *GetSampleProxyModel)
}

type SampleProxy struct {
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	cachr       cacher.ICacher
	client      *http.Client
	baseUrl     *url.URL
	timeout     time.Duration
}

// NewSampleProxy
// Returns a new SampleProxy.
func NewSampleProxy(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher) ISampleProxy {
	url, urlErr := url.Parse(environment.Get(env.SampleProxyUrl))
	if urlErr != nil {
		loggr.Panic("Couldn't convert SAMPLE_PROXY_URL environment variable to url.URL !")
	}

	timeout, timeoutErr := strconv.Atoi(environment.Get(env.SampleProxyTimeout))
	if timeoutErr != nil {
		loggr.Panic("Couldn't convert SAMPLE_PROXY_TIMEOUT environment variable to int !")
	}

	service := SampleProxy{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
		client:      &http.Client{Timeout: time.Second * time.Duration(5)},
		baseUrl:     url,
		timeout:     time.Second * time.Duration(timeout),
	}
	return &service
}

// GetGoogle
// Gets sample response from google.com.
func (p *SampleProxy) GetGoogle(ch chan *GetSampleProxyResponse, model *GetSampleProxyModel) {
	modelErr := p.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetSampleProxyResponse{Error: modelErr}
		return
	}

	// You can override the default timeout here.
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseUrl.String()+"", nil)

	response, err := p.client.Do(request)
	if err != nil {
		ch <- &GetSampleProxyResponse{Error: err}
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		ch <- &GetSampleProxyResponse{Error: err}
		return
	}

	ch <- &GetSampleProxyResponse{
		Id:         model.Id,
		SampleName: model.SampleName + string(body[0]),
	}
}
