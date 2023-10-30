package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	method    Method
	header    map[string]string
	url       string
	urlParams map[string]string
	body      interface{}
}

func (c *Client) WithMethod(m Method) *Client {
	c.method = m
	return c
}

func (c *Client) WithURL(u string) *Client {
	c.url = u
	return c
}

func (c *Client) WithURLParams(params map[string]string) *Client {
	c.urlParams = params
	return c
}

func (c *Client) WithBody(body interface{}) *Client {
	c.body = body
	return c
}

func (c *Client) WithHeader(header map[string]string) *Client {
	c.header = header
	return c
}

func (c *Client) WithContentTypeJSON() *Client {
	if c.header == nil {
		c.header = make(map[string]string)
	}
	c.header["Content-Type"] = "application/json"
	return c
}

func (c *Client) Do() (resp []byte, err error) {
	if err = c.validateBeforeDO(); err != nil {
		return nil, err
	}

	params := url.Values{}
	for pK, pV := range c.urlParams {
		params.Add(pK, pV)
	}

	if paramsEncoded := params.Encode(); paramsEncoded != "" {
		c.url += "?" + paramsEncoded
	}

	body, err := json.Marshal(c.body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(c.method.Upper().String(), c.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for hK, hV := range c.header {
		req.Header.Add(hK, hV)
	}

	respRaw, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = respRaw.Body.Close()
	}()

	return io.ReadAll(respRaw.Body)
}

func (c *Client) validateBeforeDO() error {
	if !c.method.IsSupported() {
		return errors.New(fmt.Sprintf("request method %q has not been supported", c.method))
	}

	if c.url == "" {
		return errors.New("request url can not be empty")
	}

	return nil
}
