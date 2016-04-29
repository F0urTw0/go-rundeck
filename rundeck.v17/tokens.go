package rundeck

import (
	"encoding/xml"
	"net/url"
)

type Tokens struct {
	XMLName     xml.Name    `xml:"tokens"`
	Count       int64       `xml:"count,attr"`
	AllUsers    *bool       `xml:"allusers,omitempty"`
	User        *string     `xml:"user,attr"`
	Tokens      []*Token    `xml:"token"`
}

type Token struct {
	XMLName     xml.Name    `xml:"token"`
	ID          string      `xml:"id,attr"`
	User        string      `xml:"user,attr"`
}

func (c *RundeckClient) GetTokens() (tokens Tokens, err error) {

	// init response
	var response []byte

	// call api
	if err = c.Get(&response, "tokens", url.Values{}); err == nil {
		err = xml.Unmarshal(response, &tokens)
	}

	return

}

func (c *RundeckClient) GetUserTokens(user string) (tokens Tokens, err error) {

	// init response
	var response []byte

	// call api
	if err = c.Get(&response, "tokens/"+user, url.Values{}); err == nil {
		err = xml.Unmarshal(response, &tokens)
	}

	return

}

func (c *RundeckClient) GetToken(id string) (token Token, err error) {

	// init response
	var response []byte

	// call api
	if err = c.Get(&response, "token/"+id, url.Values{}); err == nil {
		err = xml.Unmarshal(response, &token)
	}

	return

}
