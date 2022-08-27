package api

import (
	"bytes"
	"context"
	"errors"
	"github.com/gojektech/heimdall/v6/httpclient"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/v6/httpclient"
	"taleteller/logger"
)

func Post(
	ctx context.Context,
	url string,
	requestBody []byte,
	headers map[string]string,
	httpClient *httpclient.Client,
) (response *http.Response, err error) {
	if httpClient == nil {
		err = errors.New("no http client provided")
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	reqStartTime := time.Now().UTC()
	response, err = httpClient.Do(req)
	if err != nil {
		downstreamErrorLog(ctx, reqStartTime, req, err)
		return
	}
	downstreamInfoLog(ctx, reqStartTime, req, response)

	return
}

//PostWithCaseSensitiveHeader
//Go’s http.Client's  request.Header.Set(...) ends up calling CanonicalMIMEHeaderKey on the header key
//This converts "x-api-key" to "X-Api-Key".
//So if we want some headers with x-api-key use this method.
//Adding a new method here so other APIs continues to use the  func Post() method.
func PostWithCaseSensitiveHeader(
	ctx context.Context,
	url string,
	requestBody []byte,
	headers map[string]string,
	httpClient *httpclient.Client,
) (response *http.Response, err error) {
	if httpClient == nil {
		err = errors.New("no http client provided")
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header[k] = []string{v}
	}
	reqStartTime := time.Now().UTC()
	response, err = httpClient.Do(req)
	if err != nil {
		downstreamErrorLog(ctx, reqStartTime, req, err)
		return
	}
	downstreamInfoLog(ctx, reqStartTime, req, response)

	return
}

func Patch(
	ctx context.Context,
	url string,
	requestBody []byte,
	headers map[string]string,
	httpClient *httpclient.Client,
) (response *http.Response, err error) {
	if httpClient == nil {
		err = errors.New("no http client provided")
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	reqStartTime := time.Now().UTC()
	response, err = httpClient.Do(req)
	if err != nil {
		downstreamErrorLog(ctx, reqStartTime, req, err)
		return
	}
	downstreamInfoLog(ctx, reqStartTime, req, response)

	return
}

func Put(
	ctx context.Context,
	url string,
	requestBody []byte,
	headers map[string]string,
	httpClient *httpclient.Client,
) (response *http.Response, err error) {
	if httpClient == nil {
		err = errors.New("no http client provided")
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	reqStartTime := time.Now().UTC()
	response, err = httpClient.Do(req)
	if err != nil {
		downstreamErrorLog(ctx, reqStartTime, req, err)
		return
	}
	downstreamInfoLog(ctx, reqStartTime, req, response)

	return
}

func Get(
	ctx context.Context,
	url string,
	headers map[string]string,
	httpClient *httpclient.Client,
	requestBody []byte,
) (response *http.Response, err error) {
	if httpClient == nil {
		err = errors.New("no http client provided")
		return
	}

	var reqBody *bytes.Buffer
	if requestBody != nil {
		reqBody = bytes.NewBuffer(requestBody)
	}

	var req *http.Request
	if requestBody != nil {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, reqBody)
		if err != nil {
			return
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		if err != nil {
			return
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	reqStartTime := time.Now().UTC()
	response, err = httpClient.Do(req)
	if err != nil {
		downstreamErrorLog(ctx, reqStartTime, req, err)
		return
	}
	downstreamInfoLog(ctx, reqStartTime, req, response)

	return
}

func downstreamInfoLog(ctx context.Context, reqStartTime time.Time, req *http.Request, res *http.Response) {
	logMsg := "DOWNSTREAM_API_RESPONSE"

	logger.Infow(ctx,
		logMsg,
		"downstream_url", req.URL,
		"downstream_host", req.URL.Host,
		"downstream_path", req.URL.Path,
		"downstream_status_code", res.StatusCode,
		"downstream_response_time", time.Now().UTC().Sub(reqStartTime).Milliseconds(),
	)
}

func downstreamErrorLog(ctx context.Context, reqStartTime time.Time, req *http.Request, err error) {
	logMsg := "DOWNSTREAM_API_RESPONSE"

	logger.Errorw(ctx,
		logMsg,
		"downstream_url", req.URL,
		"downstream_host", req.URL.Host,
		"downstream_path", req.URL.Path,
		"downstream_status_code", "ERROR",
		"downstream_error", err.Error(),
		"downstream_response_time", time.Now().UTC().Sub(reqStartTime).Milliseconds(),
	)
}
