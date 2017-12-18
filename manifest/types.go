package manifest

type PubSubDef struct {
	Topic string `yaml:topic`
}

type Header struct {
	name  string `yaml:name`
	value string `yaml:value`
}

type RequestDef struct {
	Url     string   `yaml:url`
	Method  string   `yaml:method`
	Headers []Header `yaml:headers`
	Body    string   `yaml:body`
}

type CronTaskDef struct {
	Description string      `yaml:description`
	Schedule    string      `yaml:schedule`
	PubSub      *PubSubDef  `yaml:subsub`
	Request     *RequestDef `yaml:request`
}

type CronManifest []CronTaskDef
