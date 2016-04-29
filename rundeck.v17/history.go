package rundeck

import (
	"encoding/xml"
	"net/url"
)

type Events struct {
	XMLName xml.Name `xml:"events"`
	Count   int64    `xml:"count,attr"`
	Total   int64    `xml:"total,attr"`
	Max     int64    `xml:"max,attr"`
	Offset  int64    `xml:"offset,attr"`
	Events  []Event  `xml:"event"`
}

type Event struct {
	XMLName     xml.Name `xml:"event"`
	StartTime   string   `xml:"starttime,attr"`
	EndTime     string   `xml:"endtime,attr"`
	Title       string   `xml:"title"`
	Status      string   `xml:"status"`
	Summary     string   `xml:"summary"`
	NodeSummary struct {
		XMLName   xml.Name
		Succeeded int64 `xml:"succeeded,attr"`
		Failed    int64 `xml:"failed,attr"`
		Total     int64 `xml:"total,attr"`
	} `xml:"node-summary"`
	User        string `xml:"user"`
	Project     string `xml:"project"`
	DateStarted string `xml:"date-started"`
	DateEnded   string `xml:"date-ended"`
	AbortedBy   string `xml:"abortedby,omitempty"`
	Job         *struct {
		XMLName xml.Name
		ID      string `xml:"id,attr"`
	} `xml:"job,omitempty"`
	Execution struct {
		XMLName xml.Name
		ID      int64 `xml:"id,attr"`
	} `xml:"execution,omitempty"`
}

func (c *RundeckClient) GetHistory(project string) (events Events, err error) {

	// init vars
	var (
		params url.Values
		response []byte
	)

	// set query string
	params.Add("project",project)

	// call api
	if err = c.Get(&response, "history", params); err == nil {
		err = xml.Unmarshal(response,&events)
	}

	return

}