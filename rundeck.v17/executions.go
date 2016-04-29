package rundeck

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"net/url"
)

type Executions struct {
	Count               int64           `xml:"count,attr"`
	Total               int64           `xml:"total,attr"`
	Max                 int64           `xml:"max,attr"`
	Offset              int64           `xml:"offset,attr"`
	Executions          []Execution     `xml:"execution"`
}

type ExecutionStep struct {
	XMLName             xml.Name        `xml:"step"`
	StepCtx             int64           `xml:"stepctx,attr"`
	ID                  int64           `xml:"id,attr"`
	StartTime           string          `xml:"startTime"`
	UpdateTime          string          `xml:"updateTime"`
	EndTime             string          `xml:"endTime"`
	ExecutionState      string          `xml:"executionState"`
	NodeStep            bool            `xml:"nodeStep"`
	NodeStates          []NodeState     `xml:"nodeStates>nodeState"`
}

type ExecutionsDeleted struct {
	XMLName             xml.Name        `xml:"deleteExecutions"`
	RequestCount        int64           `xml:"requestCount,attr"`
	AllSuccessful       bool            `xml:"allSuccessful,attr"`
	Successful struct {
		XMLName         xml.Name        `xml:"successful"`
		Count           int64           `xml:"count,attr"`
	} `xml:"successful"`
	Failed struct {
		XMLName         xml.Name                `xml:"failed"`
		Count           int64                   `xml:"count,attr"`
		Failures        []FailedExecutionDelete `xml:"execution,omitempty"`
	}                                           `xml:"failed"`
}

type FailedExecutionDelete struct {
	XMLName             xml.Name                `xml:"execution"`
	ID                  int64                   `xml:"id,attr"`
	Message             string                  `xml:"message,attr"`
}

func (c *RundeckClient) ListProjectExecutions(projectId string, params url.Values) (Executions, error) {
	var res []byte
	params.Add("project",projectId)
	var data Executions
	err := c.Get(&res, "executions", params)
	xml.Unmarshal(res, &data)
	fmt.Printf("%s\n", string(res))
	return data, err
}

func (c *RundeckClient) ListRunningExecutions(projectId string) (executions Executions, err error) {

	// init response
	var response []byte

	// set query string
	params := url.Values{}
	params.Add("project",projectId)

	// call api
	if err = c.Get(&response, "executions/running", params); err != nil {
		err = xml.Unmarshal(response, &executions)
	}

	return

}

func (c *RundeckClient) DeleteExecutions(ids []string) (executions ExecutionsDeleted, err error) {

	// init response
	var response []byte

	// set form params
	params := url.Values{}
	params.Add("ids",strings.Join(ids,","))

	// call api
	if err = c.Post(&response, "executions/delete", nil, params); err != nil {
		err = xml.Unmarshal(response, &executions)
	}

	return

}

func (c *RundeckClient) DeleteAllExecutionsForProject(project string, max int64) (dels ExecutionsDeleted, err error) {

	// init vars
	var (
		execs Executions
		delete []string
		response []byte
	)

	params := url.Values{}
	params.Add("max",strconv.FormatInt(max,10))

	// get all executions
	if execs, err = c.ListProjectExecutions(project, params); err != nil {
		return
	}

	// mount list of executions to be deleted
	for _, execution := range execs.Executions {
		delete = append(delete, strconv.Itoa(execution.ID))
	}

	// no executions found
	if len(delete) == 0 {
		return dels, errors.New("No executions found for project: " + project)
	}

	// send them to be deleted
	params.Del("max")
	params.Add("ids",strings.Join(delete,","))

	if err = c.Post(&response, "executions/delete", nil, params); err == nil {
		err = xml.Unmarshal(response,&dels)
	}

	return

}