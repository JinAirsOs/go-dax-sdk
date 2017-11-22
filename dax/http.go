package dax

import (
	"bytes"
	"errors"
	"io/ioutil"
	//"log"
	"net/http"
	"strings"

	"github.com/JinAirsOs/go-dax-sdk/aws"
)

type HttpResponse struct {
	StatusCode int
	Body       []byte
}

func (d *Dax) doRequest(method, uri string, rawbody []byte) *HttpResponse, error {
	uri = strings.TrimLeft(uri, "/")
	url := d.ApiBase + uri
	var req *http.Request
	var err error
	switch upperMethod := strings.ToUpper(method); upperMethod {
	case "GET":
		req, err = http.NewRequest("GET", url, nil)
	case "POST":
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(rawbody))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	case "DELETE":
		req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(rawbody))
	case "PUT":
		req, err = http.NewRequest("PUT", url, bytes.NewBuffer(rawbody))
	default:
		panic(errors.New("unsupported http method"))
	}

	req.Close = true

	//sign the request with aws sign s3
	aws.SignS3(req, *d.Credentials)

	//log.Println("request", req)
	resp, err := d.Client.Do(req)
	if err != nil {
		return &HttpResponse{}, err
	}
	defer resp.Body.Close()

	//log.Println("response Status:", resp.Status)
	//log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//log.Println("response Body:", string(body))

	return &HttpResponse{StatusCode: resp.StatusCode, Body: body}, nil
}
