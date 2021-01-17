package handlers

import (
	"a15cli/models"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
)

type StaticServerConf struct {
	Pid int `yaml:"pid""`
}

func StaticServer(c *cli.Context) error {
	dir := models.Config.ServerDirectory
	addr := models.Config.ListenPort

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Static("/", dir)
	conf := StaticServerConf{
		Pid: os.Getpid(),
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	confFilePath := fmt.Sprintf("%s/.swiss/static_server.yaml", home)
	err = WriteConfig(confFilePath, conf)
	if err != nil {
		return err
	}

	fmt.Println(e.Start(addr))
	return nil
}

func StopStaticServer(ctx *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	config := StaticServerConf{}

	confFilePath := fmt.Sprintf("%s/.swiss/static_server.yaml", home)
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
