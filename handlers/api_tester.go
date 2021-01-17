package handlers

import (
	"a15cli/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type TestRun struct {
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
}
type ScriptLine struct {
	Env        string                       `yaml:"env"`
	Service    string                       `yaml:"service"`
	Site       string                       `yaml:"site"`
	Headers    map[string]string            `yaml:"headers"`
	Method     string                       `yaml:"method"`
	URLPattern string                       `yaml:"url_pattern"`
	Body       string                       `yaml:"body"`
	RespValues map[string]map[string]string `yaml:"resp_values"`
}

type Script struct {
	ScriptValues map[string]string `yaml:"vars"`
	ScriptLines  []ScriptLine      `yaml:"actions"`
}

func ApiTester(c *cli.Context) error {
	scriptFilename := c.Args().Get(0)

	scriptFile, err := ioutil.ReadFile(scriptFilename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	var script Script
	err = yaml.Unmarshal(scriptFile, &script)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	for _, line := range script.ScriptLines {
		vals, errs := runLine(line, script.ScriptValues)
		if len(errs) == 0 {
			for k, v := range vals {
				script.ScriptValues[k] = v
			}
		}
	}

	fmt.Printf("%+v\n", script.ScriptValues)
	return nil
}

func runLine(inLine ScriptLine, values map[string]string) (map[string]string, []error) {
	fmt.Printf("- %+v\n\t", inLine)
	line := applyTemplate(&inLine, values)
	service := fmt.Sprintf("https:%s", models.Config.ApiServices[line.Service])
	urlPattern := line.URLPattern
	body := &line.Body
	statusCode, respBoby, err := MakeHttpRequest(line.Method, fmt.Sprintf("%s/%s", service, urlPattern), line.Site, line.Headers, body)

	errs := make([]error, 0)
	if err != nil {
		errs = append(errs, err)
	}
	if statusCode > 299 {
		errs = append(errs, errors.New(fmt.Sprintf("http error: %d", statusCode)))
	}

	scriptValues := make(map[string]string)
	if values, ok := line.RespValues["json"]; ok {
		for name, jsonPosition := range values {
			scriptValues[name] = gjson.Get(respBoby, jsonPosition).String()
		}
	}

	return scriptValues, errs
}

func applyTemplate(lineInfo *ScriptLine, data map[string]string) *ScriptLine {
	bites, err := json.Marshal(lineInfo)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("test").Parse(string(bites))
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		panic(err)
	}

	var output ScriptLine
	err = json.Unmarshal(tpl.Bytes(), &output)
	if err != nil {
		panic(err)
	}

	return &output
}
