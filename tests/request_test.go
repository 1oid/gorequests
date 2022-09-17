package tests

import (
	"fmt"
	"github.com/1oid/gorequests/pkg/models"
	"testing"
)

func TestReq(t *testing.T) {
	request := models.Request{
		Method: "POST",
		Url:    "https://www.baidu.com",
		Proxy:  "http://127.0.0.1:8080",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
		},
		Data: []byte("niu=1"),
	}
	response := request.DoReq()
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))

	fmt.Println("======Request Raw======")
	fmt.Println(request.RequestRaw)

}
