package request

import (
	"io/ioutil"
	"log"
	"net/http"

	mfst "github.com/fixate/cron-server/manifest"
	"github.com/urfave/cli"
)

type HttpRequestProvider struct {
	cli  *cli.Context
	Task *mfst.CronTaskDef

	client *http.Client
}

func NewProvider(cli *cli.Context, task *mfst.CronTaskDef) *HttpRequestProvider {
	return &HttpRequestProvider{
		cli:  cli,
		Task: task,
	}
}

func (p *HttpRequestProvider) Setup() error {
	p.client = &http.Client{}
	return nil
}

func (p *HttpRequestProvider) Handler() func() {
	var task mfst.CronTaskDef = *p.Task
	req := task.Request
	return func() {
		log.Printf("[REQUEST] Task start: '%s'\n", task.Description)
		log.Printf("[REQUEST] Publishing to url: '%s'\n", req.Url)
		status, _, err := p.makeRequest(req)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[REQUEST] Got status from '%s' of %d\n", req.Url, status)
	}
}

func (p *HttpRequestProvider) makeRequest(def *mfst.RequestDef) (int, string, error) {
	req, err := http.NewRequest(def.Method, def.Url, nil)
	if err != nil {
		return 0, "", err
	}

	for _, header := range def.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	resp, err := p.client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body), nil
}
