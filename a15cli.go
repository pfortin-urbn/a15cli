package main

import (
	"a15cli/clients"
	"fmt"
	"github.com/urfave/cli/v2"
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
				Name:        "terraform",
				Aliases:     []string{"t"},
				Usage:       "terraform commands",
				Subcommands: []*cli.Command{
					{
						Name:  "install",
						Aliases:     []string{"i"},
						Usage: "install a new terraform version",
						Action: clients.InstallTerraformVersions,
					},
					{
						Name:  "list",
						Aliases:     []string{"l"},
						Usage: "list available terraform versions",
						Action: clients.ListTerraformVersions,
					},
					{
						Name:  "switch",
						Aliases:     []string{"s"},
						Usage: "switch terraform versions",
						Action: clients.SwitchTerraformVersion,
					},
				},
			},
			{
				Name:        "aws",
				Aliases:     []string{"a"},
				Usage:       "aws commands",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Aliases:     []string{"l"},
						Usage: "list available aws credentials",
						Action: clients.ListAwsCredentials,
					},
					{
						Name:  "switch",
						Aliases:     []string{"s"},
						Usage: "switch aws credentials",
						Action: clients.SwitchAwsCredentials,
					},
				},
			},
			{
				Name:        "gcp",
				Aliases:     []string{"g"},
				Usage:       "gcp commands",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Aliases:     []string{"l"},
						Usage: "list available gcp credentials",
						Action: clients.ListGcpCredentials,
					},
					{
						Name:  "switch",
						Aliases:     []string{"s"},
						Usage: "switch gcp credentials",
						Action: clients.SwitchGcpCredentials,
					},
				},
			},
			{
				Name:        "ssh",
				Aliases:     []string{"s"},
				Usage:       "ssh commands",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Aliases:     []string{"l"},
						Usage: "list available ssh identities",
						Action: clients.ListSshCredentials,
					},
					{
						Name:  "switch",
						Aliases:     []string{"s"},
						Usage: "switch ssh identity",
						Action: clients.SwitchSshCredentials,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%sERROR: %s\n%s", clients.ANSI_RED, err.Error(), clients.ANSI_RESET)
	}
}
