package bot

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed templates/build_failed.txt
var rawTemplate string

var notifyTemplate *template.Template

type notifyContext struct {
	PipelineURL   string
	PipelineID    int
	Project       string
	Branch        string
	Commit        string
	CommitMessage string
	Author        string
	Reports       []struct {
		URL      string
		FileName string
	}
}

func init() {
	notifyTemplate = template.Must(template.New("build_failed").Parse(rawTemplate))
}

func makeNotification(ctx *notifyContext) (string, error) {
	buf := bytes.Buffer{}
	err := notifyTemplate.Execute(&buf, ctx)
	return buf.String(), err
}
