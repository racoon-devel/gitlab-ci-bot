package bot

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates
var templatesFS embed.FS

var templates *template.Template

type notifyContext struct {
	PipelineURL   string
	PipelineID    int
	Project       string
	Branch        string
	Commit        string
	CommitMessage string
	Author        string
	ChangelogURL  string
	Reports       []struct {
		URL      string
		FileName string
	}
}

type notifyType = int

const (
	kBuildFailed notifyType = iota
	kVersionReleased
)

func init() {
	templates = template.Must(template.ParseFS(templatesFS, "templates/*.txt"))
}

func makeNotification(ctx *notifyContext, nType notifyType) (string, error) {
	buf := bytes.Buffer{}
	tmpl := ""

	switch nType {
	case kBuildFailed:
		tmpl = "build_failed.txt"
	case kVersionReleased:
		tmpl = "version_released.txt"
	default:
		panic("unexpected notification type")
	}

	err := templates.ExecuteTemplate(&buf, tmpl, ctx)
	return buf.String(), err
}
