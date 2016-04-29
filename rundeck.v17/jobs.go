package rundeck

import (
	"encoding/xml"
	"encoding/json"
	"errors"
	"net/url"
	"fmt"
	"io"
)

type JobList struct {
	XMLName             xml.Name        `xml:"joblist"`
	Job                 []Job           `xml:"job"`
}

type Job struct {
	XMLName             xml.Name        `xml:"job"`
	UUID                string          `xml:"uuid"`
	Name                string          `xml:"name"`
	Group               string          `xml:"group,omitempty"`
	Description         string          `xml:"description,omitempty"`
	MultipleExecutions  bool            `xml:"multipleExecutions,omitempty"`
	Sequence            *Sequence       `xml:"sequence"`
	Schedule            *Schedule       `xml:"schedule"`
	Context             *Context        `xml:"context"`
	Dispatch            *Dispatch       `xml:"dispatch"`
	Notification        *Notification   `xml:"notification"`
	NodeFilters         *NodeFilters    `xml:"nodefilters"`
	LogLevel            string          `xml:"loglevel,omitempty"`
	Logging             *Logging        `xml:"logging"`
}

type Sequence struct {
	XMLName             xml.Name        `xml:"sequence"`
	KeepGoing	        bool            `xml:"keepgoing,omitempty"`
	Strategy            string          `xml:"strategy,omitempty"`
	Command             *Command        `xml:"command"`
}

type Command struct {
	XMLName             xml.Name        `xml:"command"`
	Exec                string          `xml:"exec,omitempty"`
	Script              string          `xml:"script,omitempty"`
	ScriptFile          string          `xml:"scriptfile,omitempty"`
	ScriptURL           string          `xml:"scripturl,omitempty"`
	ScriptArgs          string          `xml:"scriptargs,omitempty"`
	ScriptInterpreter   string          `xml:"scriptinterpreter,omitempty"`
}


type Schedule struct {
	XMLName             xml.Name        `xml:"schedule"`
	Crontab             string          `xml:"crontab,attr,omitempty"`
	Time                *Time           `xml:"time"`
	Weekday             *Weekday        `xml:"weekday"`
	Month               *Month          `xml:"month"`
	Year                *Year           `xml:"year"`
}

type Time struct {
	Hour                string          `xml:"hour,attr,omitempty"`
	Minute              string          `xml:"minute,attr,omitempty"`
	Seconds             string          `xml:"seconds,attr,omitempty"`
}

type Weekday struct {
	Day                 string          `xml:"day,attr,omitempty"`
}

type Month struct {
	Month               string          `xml:"month,attr,omitempty"`
}

type Year struct {
	Year                string          `xml:"year,attr,omitempty"`
}

type Context struct {
	XMLName             xml.Name        `xml:"context"`
	Project             string          `xml:"project"`
	Options             *[]Option       `xml:"options>option,omitempty"`
}

type Options struct {
	XMLName             xml.Name        `xml:"options"`
	Options             []Option        `xml:"option"`
}

type Option struct {
	XMLName             xml.Name        `xml:"option"`
	Name                string          `xml:"name,attr"`
	Value               string          `xml:"value,attr,omitempty"`
}

type Notification struct {
	XMLName             xml.Name        `xml:"notification"`
	OnFailure           *OnStatus       `xml:"onfailure"`
	OnSuccess           *OnStatus       `xml:"onsuccess"`
	OnStart             *OnStatus       `xml:"onstart"`
}

type OnStatus struct {
	Email               *Email          `xml:"email"`
	WebHook             *WebHook        `xml:"webhook"`
	Plugin              *Plugin         `xml:"plugin"`
}

type Email struct {
	Recipients          string          `xml:"recipients,attr,omitempty"`
}

type WebHook struct {
	URLs                string          `xml:"urls,attr,omitempty"`
}

type Plugin struct {
	XMLName             xml.Name        `xml:"plugin"`
	Type                string          `xml:"type,attr,omitempty"`
	Configuration       *Configuration  `xml:"configuration"`
}

type Configuration struct {
	XMLName             xml.Name        `xml:"configuration"`
	Entry               *[]Entry        `xml:"entry"`
}

type Entry struct {
	XMLName             xml.Name        `xml:"entry"`
	Key                 string          `xml:"key,attr,omitempty"`
	Value               string          `xml:"value,attr,omitempty"`
}

type Dispatch struct {
	XMLName             xml.Name        `xml:"dispatch"`
	ThreadCount         int64           `xml:"threadcount,omitempty"`
	KeepGoing           bool            `xml:"keepgoing,omitempty"`
	RankAttribute       string          `xml:"rankAttribute,omitempty"`
	RankOrder           string          `xml:"rankOrder,omitempty"`
}

type NodeFilters struct {
	XMLName             xml.Name        `xml:"nodefilters"`
	ExcludePrecedence   bool            `xml:"excludeprecedence,attr,omitempty"`
	Filter              string          `xml:"filter"`
}

type Logging struct {
	Limit               string          `xml:"limit,attr,omitempty"`
	LimitAction         string          `xml:"limitAction,attr,omitempty"`
	Status              string          `xml:"status,attr,omitempty"`
}

//type Job struct {
//	XMLName     xml.Name `xml:"job"`
//	ID          string   `xml:"id,attr"`
//	Name        string   `xml:"name"`
//	Group       string   `xml:"group"`
//	Project     string   `xml:"project"`
//	Description string   `xml:"description,omitempty"`
//	// These two come from Execution output
//	AverageDuration int64   `xml:"averageDuration,attr,omitempty"`
//	Options         Options `xml:"options,omitempty"`
//	// These four come from Import output (depending on success,error,skipped)
//	Index int    `xml:"index,attr,omitempty"`
//	Href  string `xml:"href,attr,omitempty"`
//	Error string `xml:"error,omitempty"`
//	Url   string `xml:"url,omitempty"`
//}

type JobImportResultJob struct {
	XMLName             xml.Name        `xml:"job"`
	ID                  string          `xml:"id,omitempty"`
	Name                string          `xml:"name"`
	Group               string          `xml:"group"`
	Project             string          `xml:"project"`
	Index               int             `xml:"index,attr,omitempty"`
	Href                string          `xml:"href,attr,omitempty"`
	Error               string          `xml:"error,omitempty"`
	Url                 string          `xml:"url,omitempty"`
}
type JobImportResult struct {
	XMLName             xml.Name        `xml:"result"`
	Success             bool            `xml:"success,attr,omitempty"`
	Error               bool            `xml:"error,attr,omitempty"`
	APIVersion          int64           `xml:"apiversion,attr"`
	Succeeded           struct {
		XMLName xml.Name                `xml:"succeeded"`
		Count           int64           `xml:"count,attr"`
		Jobs            []JobImportResultJob `xml:"job,omitempty"`
	}                                   `xml:"succeeded,omitempty"`
	Failed              struct {
		XMLName         xml.Name        `xml:"failed"`
		Count           int64           `xml:"count,attr"`
		Jobs            []JobImportResultJob `xml:"job,omitempty"`
	}                                   `xml:"failed,omitempty"`
	Skipped struct {
		XMLName         xml.Name        `xml:"skipped"`
		Count           int64           `xml:"count,attr"`
		Jobs            []JobImportResultJob `xml:"job,omitempty"`
	}                                   `xml:"skipped,omitempty"`
}

//type Options struct {
//	XMLName xml.Name
//	Options []Option `xml:"option"`
//}

//type Option struct {
//	XMLName xml.Name `xml:"option"`
//	Name    string   `xml:"name,attr"`
//	Value   string   `xml:"value,attr,omitempty"`
//}

// used for listing jobs
type Jobs struct {
	XMLName             xml.Name
	Count               int64           `xml:"count,attr"`
	Jobs                []Job           `xml:"job"`
}

type RunOptions struct {
	Filter              string          `qp:"filter,omitempty"`
	LogLevel            string          `qp:"loglevel,omitempty"`
	RunAs               string          `qp:"runAs,omitempty"`
	Arguments           string          `qp:"argString,omitempty"`
}

//func (ro *RunOptions) toQueryParams() (u map[string]string) {
//	q := make(map[string]string)
//	f := reflect.TypeOf(ro).Elem()
//	for i := 0; i < f.NumField(); i++ {
//		field := f.Field(i)
//		tag := field.Tag
//		mytag := tag.Get("qp")
//		tokens := strings.Split(mytag, ",")
//		if len(tokens) == 1 {
//			switch tokens[0] {
//			case "-":
//				//skip
//			default:
//				k := tokens[0]
//				v := reflect.ValueOf(*ro).Field(i).String()
//				q[k] = v
//			}
//		} else {
//			switch tokens[1] {
//			case "omitempty":
//				if tokens[0] == "" {
//					// skip
//				} else {
//					k := tokens[0]
//					v := reflect.ValueOf(*ro).Field(i).String()
//					q[k] = v
//				}
//			default:
//				//skip
//			}
//		}
//	}
//	return q
//}

//type JobList struct {
//	XMLName xml.Name   `xml:"joblist"`
//	Job     JobDetails `xml:"job"`
//}

//type JobDetails struct {
//	ID                string          `xml:"id"`
//	Name              string          `xml:"name"`
//	LogLevel          string          `xml:"loglevel"`
//	Description       string          `xml:"description,omitempty"`
//	UUID              string          `xml:"uuid"`
//	Group             string          `xml:"group"`
//	Context           JobContext      `xml:"context"`
//	Notification      JobNotification `xml:"notification"`
//	MultipleExections bool            `xml:"multipleExecutions"`
//	Dispatch          JobDispatch     `xml:"dispatch"`
//	NodeFilters       struct {
//		Filter []string `xml:"filter"`
//	} `xml:"nodefilters"`
//	Sequence JobSequence `xml:"sequence"`
//}

type JobSequence struct {
	XMLName   xml.Name
	KeepGoing bool           `xml:"keepgoing,attr"`
	Strategy  string         `xml:"strategy,attr"`
	Steps     []SequenceStep `xml:"command"`
}

type SequenceStep struct {
	XMLName        xml.Name
	Description    string      `xml:"description,omitempty"`
	JobRef         *JobRefStep `xml:"jobref,omitempty"`
	NodeStepPlugin *PluginStep `xml:"node-step-plugin,omitempty"`
	StepPlugin     *PluginStep `xml:"step-plugin,omitempty"`
	Exec           *string     `xml:"exec,omitempty"`
	*ScriptStep    `xml:",omitempty"`
}

type ExecStep struct {
	XMLName xml.Name
	string  `xml:"exec,omitempty"`
}

type ScriptStep struct {
	XMLName           xml.Name
	Script            *string `xml:"script,omitempty"`
	ScriptArgs        *string `xml:"scriptargs,omitempty"`
	ScriptFile        *string `xml:"scriptfile,omitempty"`
	ScriptUrl         *string `xml:"scripturl,omitempty"`
	ScriptInterpreter *string `xml:"scriptinterpreter,omitempty"`
}

type PluginStep struct {
	XMLName       xml.Name
	Type          string `xml:"type,attr"`
	Configuration []struct {
		XMLName xml.Name `xml:"entry"`
		Key     string   `xml:"key,attr"`
		Value   string   `xml:"value,attr"`
	} `xml:"configuration>entry,omitempty"`
}

type JobRefStep struct {
	XMLName  xml.Name
	Name     string `xml:"name,attr,omitempty"`
	Group    string `xml:"group,attr,omitempty"`
	NodeStep bool   `xml:"nodeStep,attr,omitempty"`
}

type JobContext struct {
	XMLName xml.Name     `xml:"context"`
	Project string       `xml:"project"`
	Options *[]JobOption `xml:"options>option,omitempty"`
}

type JobOptions struct {
	XMLName xml.Name
	Options []JobOption `xml:"option"`
}

type JobOption struct {
	XMLName      xml.Name `xml:"option"`
	Name         string   `xml:"name,attr"`
	Required     bool     `xml:"required,attr,omitempty"`
	Secure       bool     `xml:"secure,attr,omitempty"`
	ValueExposed bool     `xml:"valueExposed,attr,omitempty"`
	DefaultValue string   `xml:"value,attr,omitempty"`
	Description  string   `xml:"description,omitempty"`
}

type JobNotifications struct {
	Notifications []JobNotification `xml:"notification,omitempty"`
}

type JobNotification struct {
	XMLName   xml.Name   `xml:"notification"`
	OnStart   JobPlugins `xml:"onstart,omitempty"`
	OnSuccess JobPlugins `xml:"onsuccess,omitempty"`
	OnFailure JobPlugins `xml:"onfailure,omitempty"`
}

type JobPlugins struct {
	Plugins []JobPlugin `xml:"plugin,omitempty"`
}

type JobPlugin struct {
	XMLName       xml.Name               `xml:"plugin"`
	PluginType    string                 `xml:"type,attr"`
	Configuration JobPluginConfiguration `xml:"configuration,omitempty"`
}

type JobPluginConfiguration struct {
	XMLName xml.Name                      `xml:"configuration"`
	Entries []JobPluginConfigurationEntry `xml:"entry,omitempty"`
}

type JobPluginConfigurationEntry struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr,omitempty"`
}

type JobDispatch struct {
	XMLName           xml.Name `xml:"dispatch"`
	ThreadCount       int64    `xml:"threadcount"`
	KeepGoing         bool     `xml:"keepgoing"`
	ExcludePrecedence bool     `xml:"excludePrecendence"`
	RankOrder         string   `xml:"rankOrder"`
}

type ImportParams struct {
	Format   string
	Dupe     string
	Uuid     string
	Project  string
}

func (c *RundeckClient) GetJob(id string) (list JobList, err error) {

	var response []byte

	if err = c.Get(&response, "job/"+id, url.Values{}); err == nil {
		err = xml.Unmarshal(response, &list)
	}

	return

}

func (c *RundeckClient) DeleteJob(id string) error {
	return c.Delete("job/"+id, nil)
}

func (c *RundeckClient) ExportJob(id string) (list JobList, err error) {

	// init response
	var response []byte

	// set query string
	params := url.Values{}
	params.Add("format","xml")

	// call api
	if err = c.Get(&response, "job/"+id, params); err == nil {
		err = xml.Unmarshal(response, &list)
	}

	return

}

func (c *RundeckClient) ImportJob(jobList *JobList, params *ImportParams) (result JobImportResult, err error) {

	// init response
	var response, jobs []byte

	// set form parameters
	form := url.Values{}
	form.Add("project",params.Project)
	form.Add("format",params.Format)
	form.Add("uuidOption",params.Uuid)
	form.Add("dupeOption",params.Dupe)

	// marshal job to xml
	if jobs, err = xml.Marshal(&jobList); err != nil {
		return
	}

	// post data to api
	if err = c.Post(&response, "jobs/import", jobs, form); err == nil {
		err = xml.Unmarshal(response, &result)
	}

	return

}

func (c *RundeckClient) RunJob(id string, options RunOptions) (execution Execution, err error) {

	// init response
	var response []byte

	// set form parameters
	params := url.Values{}
	params.Add("filter", options.Filter)
	params.Add("loglevel", options.LogLevel)
	params.Add("runAs", options.RunAs)
	params.Add("argString", options.Arguments)

	// call api
	if err = c.Post(&response, fmt.Sprintf("job/%s/run",id), nil, params); err == nil {
		if err = xml.Unmarshal(response, &execution); err == io.EOF {
			err = json.Unmarshal(response, &execution)
		}
	}

	return

}

func (c *RundeckClient) FindJobByName(name string, project string) (job *Job, err error) {

	var (
		jobs Jobs
		list JobList
	)

	// get all jobs
	if jobs, err = c.ListJobs(project); err == nil {

		if len(jobs.Jobs) > 0 {

			for _, d := range jobs.Jobs {

				if d.Name == name {

					if list, err = c.GetJob(d.UUID); err == nil {
						job = &list.Job[0]
					}

				}

			}

		} else {

			err = errors.New("Job not found")

		}

	}

	return
}

func (c *RundeckClient) ListJobs(projectId string) (jobs Jobs, err error) {

	// init response
	var response []byte

	// set query string
	params := url.Values{}
	params.Add("project",projectId)

	// call api
	if err = c.Get(&response, "jobs", params); err == nil {
		err = xml.Unmarshal(response, &jobs)
	}

	return

}