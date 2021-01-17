package handlers

import (
	"a15cli/models"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli/v2"
)

type MockServerConf struct {
	Pid int `yaml:"pid""`
}

type Expectation struct {
	URL              string `yaml:"url"`
	Method           string `yaml:"method"`
	ReceivedBodyPart string `yaml:"received_body_part"`
	ContentType      string `yaml:"content_type"`
	Status           int    `yaml:"status"`
	Body             string `yaml:"body"`
}

var Expectations []Expectation

func StartMockServer(c *cli.Context) error {
	scriptFilename := c.Args().Get(0)
	if scriptFilename == "" {
		return fmt.Errorf("supply expectations file")
	}

	scriptFile, err := ioutil.ReadFile(scriptFilename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(scriptFile, &Expectations)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time: ${time_rfc3339_nano} method: ${method}, uri: ${uri} status:${status} remote_ip: ${remote_ip}\n",
	}))

	e.Any("/*", handler)

	conf := MockServerConf{
		Pid: os.Getpid(),
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	confFilePath := fmt.Sprintf("%s/.swiss/mock_server.yaml", home)
	err = WriteConfig(confFilePath, conf)
	if err != nil {
		return err
	}

	fmt.Println(e.Start(models.Config.ListenPort))
	return nil
}

func handler(c echo.Context) error {
	c.Logger().Infof("Got request for URL: (%s)%s\n", c.Request().Method, c.Request().URL.String())
	body, error := ioutil.ReadAll(c.Request().Body)
	for _, expect := range Expectations {
		if expect.URL == c.Request().URL.String() &&
			expect.Method == c.Request().Method {
			if error == nil && expect.ReceivedBodyPart != "" && strings.Contains(string(body), expect.ReceivedBodyPart) {
				c.Response().Header().Set("Content-Type", expect.ContentType)
				return c.String(expect.Status, expect.Body)
			} else if expect.ReceivedBodyPart == "" {
				c.Response().Header().Set("Content-Type", expect.ContentType)
				return c.String(expect.Status, expect.Body)
			}
		}
	}
	expect := Expectations[len(Expectations)-1]
	c.Response().Header().Set("Content-Type", expect.ContentType)
	return c.String(expect.Status, expect.Body)
}

func StopMockServer(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	config := MockServerConf{}

	confFilePath := fmt.Sprintf("%s/.swiss/mock_server.yaml", home)
	err = ReadConfig(confFilePath, &config)
	if err != nil {
		return err
	}
	err = os.Remove(confFilePath)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(config.Pid)
	if err != nil {
		return err
	}
	proc.Kill()
	return nil
}
