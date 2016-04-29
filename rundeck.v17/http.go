package rundeck

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"gopkg.in/jmcvetta/napping.v2"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"github.com/davecgh/go-spew/spew"
)

func (rc *RundeckClient) Get(i *[]byte, path string, params url.Values) error {
	return rc.makeRequest(i, nil, "GET", path, params)
}

func (rc *RundeckClient) Delete(path string, params url.Values) error {
	var b []byte
	return rc.makeRequest(&b, nil, "DELETE", path, params)
}

func (rc *RundeckClient) Post(i *[]byte, path string, data []byte, params url.Values) error {
	return rc.makeRequest(i, data, "POST", path, params)
}

func (rc *RundeckClient) Put(i *[]byte, path string, data []byte, params url.Values) error {
	return rc.makeRequest(i, data, "PUT", path, params)
}

func (client *RundeckClient) makeRequest(i *[]byte, payload []byte, method string, path string, params url.Values) error {

	// set endpoint
	endpoint := client.Config.BaseURL + "/api/17/" + path

	// add headers
	headers := http.Header{}
	headers.Add("Content-Type","application/xml")

	jar, _ := cookiejar.New(nil)
	client.Client.Client.Jar = jar
	if client.Config.AuthMethod == "basic" {
		authQs := url.Values{}
		authQs.Add("j_username", client.Config.Username)
		authQs.Add("j_password", client.Config.Password)
		authPayload := bytes.NewBuffer(nil)
		base_auth_url := client.Config.BaseURL + "/j_security_check"
		cookieReq := napping.Request{
			Url:        base_auth_url,
			Params:     &authQs,
			Method:     "POST",
			RawPayload: true,
			Payload:    authPayload,
		}
		r, err := client.Client.Send(&cookieReq)
		if err != nil {
			return err
		}
		if r.Status() != 200 {
			return errors.New(r.RawText())
		}
	} else {
		headers.Add("X-Rundeck-Auth-Token", client.Config.Token)
	}

	req := napping.Request{
		Url:                 endpoint,
		Header:              &headers,
		Params:              &params,
		Method:              method,
		RawPayload:          true,
		Payload:             bytes.NewBuffer(payload),
		CaptureResponseBody: true,
	}

	r, err := client.Client.Send(&req)
	if err != nil {
		return err
	} else {
		if r.Status() == 404 {
			errormsg := fmt.Sprintf("No such item (%s)", r.Status)
			return errors.New(errormsg)
		}
		if r.Status() == 204 {
			return nil
		}
		if (r.Status() < 200) || (r.Status() > 299) {
			var data RundeckError



			xml.Unmarshal([]byte(r.RawText()), &data)
			err = json.Unmarshal([]byte(r.RawText()), &data)

			spew.Dump(err, r.RawText(), data)

			errormsg := fmt.Sprintf("non-2xx response (code: %d): %s", r.Status(), data.Message)
			return errors.New(errormsg)
		} else {
			b := r.ResponseBody.Bytes()
			*i = append(*i, b...)
			return nil
		}
	}
}
