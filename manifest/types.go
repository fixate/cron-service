package manifest

type PubSubDef struct {
	Topic      string            `yaml:topic`
	Message    string            `yaml:message`
	Attributes map[string]string `yaml:attributes`
}

type Header struct {
	Name  string `yaml:name`
	Value string `yaml:value`
}

type RequestDef struct {
	Url     string   `yaml:url`
	Method  string   `yaml:method`
	Headers []Header `yaml:headers`
	Body    string   `yaml:body`
}

type CronTaskDef struct {
	Description string      `yaml:description`
	Enabled     bool        `yaml:enabled`
	FireOnStart bool        `yaml:fireonstart`
	Schedule    string      `yaml:schedule`
	PubSub      *PubSubDef  `yaml:subsub`
	Request     *RequestDef `yaml:request`
}

type CronManifest []CronTaskDef

func (c *CronTaskDef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type localAlias CronTaskDef
	obj := localAlias{Enabled: true}
	if err := unmarshal(&obj); err != nil {
		return err
	}

	*c = CronTaskDef(obj)
	return nil
}
