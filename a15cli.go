package main

import (
	"a15cli/clients"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "a15cli",
		Usage: "A15 tool to manage local dev stuff",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   `Install some other useful tools: (httpie, jq, slack-cli, yamllint, yq)`,
				Action:  clients.BaseInstall,
			},
			{
				Name:    "terra-install",
				Aliases: []string{"ti"},
				Usage:   "install/overwrite a new terraform versions",
				Action:  clients.InstallTerraformVersions,
			},
			{
				Name:    "terra-list",
				Aliases: []string{"tl"},
				Usage:   "list all available terraform versions",
				Action:  clients.ListTerraformVersions,
			},
			{
				Name:    "terra-switch",
				Aliases: []string{"tw"},
				Usage:   "switch active terraform versions",
				Action:  clients.SwitchTerraformVersion,
			},
			{
				Name:    "aws-list",
				Aliases: []string{"al"},
				Usage:   "list available AWS credentials",
				Action:  clients.ListAwsCredentials,
			},
			{
				Name:    "aws-switch",
				Aliases: []string{"aw"},
				Usage:   "switch the active AWS credentials",
				Action:  clients.SwitchAwsCredentials,
			},
			{
				Name:    "gcp-list",
				Aliases: []string{"gl"},
				Usage:   "list available AWS credentials",
				Action:  clients.ListGcpCredentials,
			},
			{
				Name:    "gcp-switch",
				Aliases: []string{"gw"},
				Usage:   "switch the active AWS credentials",
				Action:  clients.SwitchGcpCredentials,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
