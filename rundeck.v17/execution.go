package rundeck

import (
	"encoding/xml"
	"net/url"
)

type Execution struct {
	XMLName                 xml.Name            `xml:"execution"`
	ID                      int                 `xml:"id,attr"                  json:"id"`
	HRef                    string              `xml:"href,attr"                json:"href"`
	Permalink               string              `xml:"permalink"                json:"permalink"`
	Status                  string              `xml:"status,attr"              json:"status"`
	Project                 string              `xml:"project,attr"             json:"project"`
	User                    string              `xml:"user"                     json:"user"`
	Description             string              `xml:"description,omitempty"    json:"description"`
	ArgString               string              `xml:"argstring,omitempty"      json:"argstring"`
	DateStarted             *ExecutionDateTime  `xml:"date-started"             json:"date-started"`
	DateEnded               *ExecutionDateTime  `xml:"date-ended"               json:"date-ended"`
	Job                     *Job                `xml:"job"                      json:"job"`
	SuccessfulNodes         Nodes               `xml:"successfulNodes,omitempty"`
	FailedNodes             Nodes               `xml:"failedNodes,omitempty"`
}

type ExecutionDateTime struct {
	UnixTime                int64               `xml:"unixtime,attr,omitempty"`
	Date                    string              `xml:"date,attr,omitempty"`
}

type ExecutionOutput struct {
	XMLName                 xml.Name            `xml:"output"`
	ID                      int64               `xml:"id"`
	Offset                  int64               `xml:"offset"`
	Completed               bool                `xml:"completed"`
	ExecCompleted           bool                `xml:"execCompleted"`
	HasFailedNodes          bool                `xml:"hasFailedNodes"`
	ExecState               string              `xml:"execState"`
	LastModified            ExecutionDateTime   `xml:"lastModified"`
	ExecDuration            int64               `xml:"execDuration"`
	TotalSize               int64               `xml:"totalSize"`
	Entries                 OutputEntries       `xml:"entries"`
}

type OutputEntries struct {
	Entry                   []OutputEntry       `xml:"entry"`
}

type OutputEntry struct {
	XMLName                 xml.Name
	Time                    string              `xml:"time,attr"`
	AbsoluteTime            string              `xml:"absolute_time,attr"`
	Log                     string              `xml:"log,attr"`
	Level                   string              `xml:"level,attr"`
	User                    string              `xml:"user,attr"`
	Command                 string              `xml:"command,attr"`
	Stepctx                 string              `xml:"stepctx,attr"`
	Node                    string              `xml:"node,attr"`
}

type ExecutionState struct {
	XMLName                 xml.Name            `xml:"result"`
	Success                 bool                `xml:"success,attr"`
	ApiVersion              int64               `xml:"apiversion,attr"`
	StartTime               string              `xml:"executionState>startTime"`
	StepCount               int64               `xml:"executionState>stepCount"`
	AllNodes                []Node              `xml:"executionState>allNodes>nodes>node,omitempty"`
	TargetNodes             []Node              `xml:"executionState>targetNodes>nodes>node,omitempty"`
	ExecutionID             int64               `xml:"executionState>executionId"`
	Completed               bool                `xml:"executionState>completed"`
	UpdateTime              string              `xml:"executionState>updateTime,omitempty"`
	Steps                   []ExecutionStep     `xml:"executionState>steps>step,omitempty"`
	Nodes                   []NodeWithSteps     `xml:"executionState>nodes>node"`
}

func (c *RundeckClient) GetExecution(executionId string) (exec Execution, err error) {
	var res []byte
	var execs Executions
	err = c.Get(&res, "execution/"+executionId, nil)
	xmlerr := xml.Unmarshal(res, &execs)
	if xmlerr != nil {
		return exec, xmlerr
	}
	return execs.Executions[0], err
}

func (c *RundeckClient) GetExecutionState(executionId string) (ExecutionState, error) {
	var res []byte
	var data ExecutionState
	err := c.Get(&res, "execution/"+executionId+"/state", url.Values{})
	err = xml.Unmarshal(res, &data)
	return data, err
}

func (c *RundeckClient) GetExecutionOutput(executionId string) (ExecutionOutput, error) {
	var res []byte
	var data ExecutionOutput
	err := c.Get(&res, "execution/"+executionId+"/output", url.Values{})
	xml.Unmarshal(res, &data)
	return data, err
}

func (c *RundeckClient) DeleteExecution(id string) error {
	return c.Delete("execution/"+id, nil)
}