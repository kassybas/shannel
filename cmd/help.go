package cmd

import (
	"html/template"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kassybas/shannel/internal/snlapi"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Help command prints the help for the user
// TODO
func Help(c *cli.Context) error {
	return nil
}

func getDefText(sp *string) string {
	if sp == nil {
		return ""
	}
	return "[default:" + *sp + "]"
}

// PrintHelpText templates and prints the help text for the current shannel file
func PrintHelpText(sf snlapi.SnlFile) {
	funcMap := template.FuncMap{
		"GetDefText": getDefText,
		"ToUpper":    strings.ToUpper,
		"Yellow":     color.New(color.FgYellow).SprintFunc(),
		"Blue":       color.New(color.FgCyan).SprintFunc(),
		"BoldWhite":  color.New(color.FgHiWhite, color.Bold).SprintFunc(),
		"White":      color.New(color.FgHiWhite).SprintFunc(),
		"Faint":      color.New(color.Faint).SprintFunc(),
		"Bold":       color.New(color.Bold).SprintFunc(),
	}
	helptext := `
{{BoldWhite "Usage"}}
      {{White "snl make"}} {{Blue "[target]"}}{{range  $arg := .Args}}{{if eq $arg.Type "pos"}} {{$arg.Name | ToUpper | Yellow }}{{end}}{{end}} [options]`

	helptext += "\n\n{{BoldWhite \"Targets\"}}"
	helptext += `{{range $name, $trg := .Target}}
      {{Blue $name}}{{if ne $trg.Usage ""}}	{{$length := len $name}}{{if gt 8 $length}}	{{end}}{{White $trg.Usage}}{{end}}{{end}}
`

	hasPosArgs := false
	hasNamedArgs := false
	for _, f := range sf.Args {
		if f.Type != "pos" {
			hasNamedArgs = true
		} else {
			hasPosArgs = true
		}
		if hasNamedArgs && hasPosArgs {
			break
		}
	}
	if hasPosArgs {
		helptext += "\n{{BoldWhite \"Positional arguments\"}}"
		helptext += `{{range  $arg := .Args}}{{if eq $arg.Type "pos"}}
      {{$arg.Name | ToUpper | Yellow}}{{if ne $arg.Usage ""}}{{$length := len $arg.Name}}{{if lt 7 $length}} {{end}}		{{$arg.Usage}}{{end}}{{end}}{{end}}`
	}
	if hasNamedArgs {
		helptext += "\n\n{{BoldWhite \"Options\"}}"
		helptext += `{{range  $arg := .Args}}{{if ne $arg.Type "pos"}}
{{if ne $arg.Alias ""}}  -{{$arg.Alias}}, {{else}}      {{end}}{{Bold "--"}}{{Bold $arg.Name}}{{if ne $arg.Type "bool"}} {{$arg.Name | ToUpper}}{{else}}	{{end}}{{if ne $arg.Usage ""}}{{$length := len $arg.Name}}{{if lt 7 $length}}
		{{end}}	{{$arg.Usage}}{{$def := GetDefText $arg.Default }}{{if ne $def ""}} {{Faint $def}}{{end}}{{end}}{{end}}{{end}}`
	}

	helptext += "\n"

	logrus.Trace("printing help text")
	t, err := template.New("helptext").Funcs(funcMap).Parse(helptext)
	if err != nil {
		logrus.WithField("error", err).Fatal("internal error: could not template usage text")
	}
	err = t.Execute(os.Stdout, sf)
	if err != nil {
		logrus.WithField("error", err).Fatal("internal error: could not print usage text")
	}
}
