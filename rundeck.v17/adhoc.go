package rundeck

import (
	"encoding/xml"
	"strings"
	"net/url"
)

type ExecutionId struct {
	ID string `xml:"id,attr"`
}

func (c *RundeckClient) RunAdhoc(projectId string, exec string, node_filter string) (ExecutionId, error) {

	params := url.Values{}
	params.Add("project",projectId)
	params.Add("exec",exec)

	n := strings.Split(node_filter, " ")
	for _, i := range n {
		f := strings.Split(i, "=")
		k, v := f[0], f[1]
		params.Add(k,v)
	}

	var res []byte
	var data ExecutionId
	err := c.Get(&res, "run/command", params)
	xml.Unmarshal(res, &data)
	return data, err
}
