package rundeck

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"gopkg.in/jmcvetta/napping.v2"
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

	response, err := client.Client.Send(&req)

	if err != nil {
		return err
	} else {
		if response.Status() == 404 {
			errormsg := fmt.Sprintf("No such item (%s)", response.Status)
			return errors.New(errormsg)
		}
		if response.Status() == 204 {
			return nil
		}
		if (response.Status() < 200) || (response.Status() > 299) {
			var data RundeckResult
			if err = xml.Unmarshal([]byte(response.RawText()), &data); err != nil {
				err = json.Unmarshal([]byte(response.RawText()), &data)
			}
			errormsg := fmt.Sprintf("non-2xx response (code: %d): %s", response.Status(), data.Error.Message)
			return errors.New(errormsg)
		} else {
			b := response.ResponseBody.Bytes()
			*i = append(*i, b...)
			return nil
		}
	}
}
