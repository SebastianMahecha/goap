package goap

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const userAgent = "goap/v1.0"

// Client handle soap
type Client struct {
	// HTTPClient used for send/recieve soap request
	HTTPClient     *http.Client
	RequestBuilder func(method string, url string, body io.Reader) (*http.Request, error)
}

var (
	// DebugMode log request/response if set true
	DebugMode = false

	// DefaultRequestBuilder create simple post method
	DefaultRequestBuilder = func(method string, url string, body io.Reader) (*http.Request, error) {
		req, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
		req.Header.Set("User-Agent", userAgent)
		return req, err
	}

	// DefaultClient used by default
	DefaultClient = &Client{
		HTTPClient:     http.DefaultClient,
		RequestBuilder: DefaultRequestBuilder,
	}
)

// Call soap action
func (c *Client) Call(url string, action string, requestHeaders, request, responseHeaders, response interface{}) error {
	envelope := SOAPEnvelope{
		Header: SOAPHeader{Header: requestHeaders},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	if DebugMode {
		rawreq, _ := ioutil.ReadAll(buffer)
		log.Println("soap request:", string(rawreq))
		buffer = bytes.NewBuffer(rawreq)
	}

	req, err := c.RequestBuilder(action, url, buffer)
	if err != nil {
		return err
	}

	if action != "" {
		req.Header.Add("SOAPAction", action)
	}

	req.Close = true

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if DebugMode {
		log.Println("soap response:", string(rawbody))
	}

	if len(rawbody) == 0 {
		return nil
	}

	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Header = SOAPHeader{Header: responseHeaders}
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}

// Call soap action
func Call(url string, action string, requestHeaders, request, responseHeaders, response interface{}) error {
	return DefaultClient.Call(url, action, requestHeaders, request, responseHeaders, response)
}
