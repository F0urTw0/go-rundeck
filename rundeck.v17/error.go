package rundeck

import "encoding/xml"

type RundeckError struct {
	XMLName         xml.Name    `xml:"result"           json:"-"`
	Error           bool        `xml:"error,attr"       json:"error"`
	ErrorCode       string      `xml:"errorCode,attr"   json:"errorCode"`
	ApiVersion      int         `xml:"apiversion,attr"  json:"apiversion"`
	Message         string      `xml:"error>message"    json:"message"`
}