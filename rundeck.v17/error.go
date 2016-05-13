package rundeck

import "encoding/xml"

type RundeckResult struct {
	XMLName         xml.Name        `xml:"result"           json:"-"`
	HasError        bool            `xml:"error,attr"       json:"error"`
	ApiVersion      int             `xml:"apiversion,attr"  json:"apiversion"`
	Error           RundeckError    `xml:"error"            json:"-"`
}

type RundeckError struct {
	XMLName         xml.Name        `xml:"error"            json:"-"`
	Code            string          `xml:"code,attr"        json:"errorCode"`
	Message         string          `xml:"message"          json:"message"`
}