package main

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli/v2"

	"a15cli/handlers"
	"a15cli/models"
)

var GitCommit string

func main() {
	err := envconfig.Process("", &models.Config)
	app := &cli.App{
		Name:    "a15cli",
		Usage:   "A15 tool to manage local dev stuff",
		Version: GitCommit,
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   `Install some other useful tools: (httpie, jq, slack-cli, yamllint, yq)`,
				Action:  handlers.BaseInstall,
			},
			{
				Name:    "terraform",
				Aliases: []string{"t"},
				Usage:   "terraform commands",
				Subcommands: []*cli.Command{
					{
						Name:    "install",
						Aliases: []string{"i"},
						Usage:   "install a new terraform version",
						Action:  handlers.InstallTerraformVersions,
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list available terraform versions",
						Action:  handlers.ListTerraformVersions,
					},
					{
						Name:    "switch",
						Aliases: []string{"w"},
						Usage:   "switch terraform versions",
						Action:  handlers.SwitchTerraformVersion,
					},
				},
			},
			{
				Name:    "aws",
				Aliases: []string{"a"},
				Usage:   "aws commands",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list available aws credentials",
						Action:  handlers.ListAwsCredentials,
					},
					{
						Name:    "switch",
						Aliases: []string{"w"},
						Usage:   "switch aws credentials",
						Action:  handlers.SwitchAwsCredentials,
					},
				},
			},
			{
				Name:    "gcp",
				Aliases: []string{"g"},
				Usage:   "gcp commands",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list available gcp credentials",
						Action:  handlers.ListGcpCredentials,
					},
					{
						Name:    "switch",
						Aliases: []string{"w"},
						Usage:   "switch gcp credentials",
						Action:  handlers.SwitchGcpCredentials,
					},
				},
			},
			{
				Name:    "ssh",
				Aliases: []string{"s"},
				Usage:   "ssh commands",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list available ssh identities",
						Action:  handlers.ListSshCredentials,
					},
					{
						Name:    "switch",
						Aliases: []string{"w"},
						Usage:   "switch ssh identity",
						Action:  handlers.SwitchSshCredentials,
					},
				},
			},
			{
				Name:    "api-tester",
				Aliases: []string{"at"},
				Usage:   "api tester",
				Action:  handlers.ApiTester,
			},
			{
				Name:    "static-server",
				Aliases: []string{"ss"},
				Usage:   `static files server commands`,
				Subcommands: []*cli.Command{
					{
						Name:    "start",
						Aliases: []string{"s"},
						Usage:   "serve files in the current/specified directory over http",
						Action:  handlers.StaticServer,
					},
					{
						Name:    "stop",
						Aliases: []string{"t"},
						Usage:   "kill the static server",
						Action:  handlers.StopStaticServer,
					},
				},
			},
			{
				Name:    "mock-server",
				Aliases: []string{"ms"},
				Usage:   "Webserver that matches Requests and return specific responses",
				Subcommands: []*cli.Command{
					{
						Name:    "start",
						Aliases: []string{"s"},
						Usage:   "start swiss mock server",
						Action:  handlers.StartMockServer,
					},
					{
						Name:    "t",
						Aliases: []string{"t"},
						Usage:   "stop swiss mock server",
						Action:  handlers.StopMockServer,
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("%sERROR: %s\n%s", handlers.ANSI_RED, err.Error(), handlers.ANSI_RESET)
	}
}
