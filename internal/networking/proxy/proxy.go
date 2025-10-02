package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ProxyInterface interface {
	NewRequest(ctx context.Context, method Method, url string) Request
	SetAuthorization(request Request, key, value string)
	Do(request Request) Response
}

type Proxy struct {
	client ClientInterface
}

func NewProxyFromRaw(
	client ClientInterface,
) ProxyInterface {
	return Proxy{
		client: client,
	}
}

func NewProxy() ProxyInterface {
	transport := NewTransport()
	client := &http.Client{
		Timeout:   0,
		Transport: transport,
	}
	return Proxy{
		client: client,
	}
}

func (p Proxy) NewRequest(ctx context.Context, method Method, url string) (request Request) {

	req, err := http.NewRequestWithContext(ctx, method.String(), url, nil)
	request.Req = req

	if err != nil {
		request.Err = fmt.Errorf("proxy: new request: %w", err)
		return request
	}
	return request
}

func (p Proxy) SetAuthorization(request Request, key, value string) {
	request.Req.Header.Set(key, value)
}

func (p Proxy) Do(request Request) (response Response) {
	timer := response.Timer

	timer.Start()
	resp, err := p.client.Do(request.Req)
	timer.Stop()

	if resp == nil {
		response.Err = errors.New("nil response")
		return response
	}

	response.Status = resp.StatusCode
	if err != nil {
		response.Err = fmt.Errorf("proxy: do: %w", err)
		return response
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	response.Body = body

	if err != nil {
		response.Err = err
		return response
	}
	return response
}
