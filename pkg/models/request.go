package models

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"
)

type Request struct {
	Method  string
	Url     string
	BaseURL string
	Headers map[string]string
	Data    []byte
	Cookie  string
	Timeout time.Duration
	//Proxy []

	Client *http.Client
	Proxy  string

	RequestRaw string
}

func InitialClient(timeout int, proxy string, redirect bool) (*http.Client, error) {
	var Timeout = time.Duration(timeout) * time.Second
	var dialer = &net.Dialer{
		Timeout:   Timeout,
		KeepAlive: Timeout,
	}
	var transport = &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConnsPerHost: -1,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
	}

	if proxy != "" {
		p, err := url.Parse(proxy)
		if err != nil {
			log.Fatal(err)
		}
		transport.Proxy = http.ProxyURL(p)
	}

	var client = &http.Client{
		Transport: transport,
		Timeout:   Timeout,
	}

	// if redirect
	if redirect {
		client.CheckRedirect = nil
		return client, nil
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return client, nil
}

func (req *Request) PrepareRequest() *http.Request {
	var doReqHttp *http.Request

	// prepare method
	if req.Method == "" {
		req.Method = "GET"
	}

	req.Method = strings.ToUpper(req.Method)
	if req.Method == RequestMethodGet {
		doReqHttp, _ = http.NewRequest(RequestMethodGet, req.Url, nil)
	} else if req.Method == RequestMethodPost {
		doReqHttp, _ = http.NewRequest(RequestMethodPost, req.Url, bytes.NewBuffer(req.Data))
	}

	// prepare url

	// prepare cookie
	//doReqHttp.AddCookie(Cookies)

	// prepare body

	// prepare headers
	for k, v := range req.Headers {
		doReqHttp.Header.Add(k, v)
	}
	return doReqHttp
}

func (req *Request) DoReq() *Response {
	client, err := InitialClient(15, req.Proxy, false)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	req.Client = client

	// prepare request
	var doReq = req.PrepareRequest()
	// generate request raw data
	req.RequestRaw, _ = GetRequestRaw(doReq)

	resp, err := req.Client.Do(doReq)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return &Response{
		MetaResponse: *resp,
		StatusCode:   resp.StatusCode,
		Headers:      resp.Header,
		Body:         body,
	}
}

func GetRequestRaw(req *http.Request) (string, error) {
	var requestBody = []byte{}
	var requestRaw = bytes.Buffer{}
	requestRaw.WriteString(fmt.Sprintf("%s %s %s\r\n", req.Method, req.URL.RequestURI(), req.Proto))

	// add header
	for k, v := range req.Header {
		requestRaw.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	// write host
	requestRaw.WriteString(fmt.Sprintf("%s: %s\r\n", "Host", req.Host))

	// write connection close
	requestRaw.WriteString("Connection: close\r\n")
	if _, ok := req.Header["Content-Length"]; !ok {
		reqLength := len(requestBody)
		requestRaw.WriteString(fmt.Sprintf("Content-Length: %d\r\n", reqLength))
	}
	requestRaw.WriteString("\r\n")

	if req.Body != nil {
		requestBodyReader, err := req.GetBody()

		if err != nil {
			log.Fatal(err)
			return "", err
		}
		res, err := ioutil.ReadAll(requestBodyReader)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		requestBody = res
	}
	requestRaw.Write(requestBody)
	return requestRaw.String(), nil
}

