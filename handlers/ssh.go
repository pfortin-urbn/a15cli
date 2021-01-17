package handlers

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func ListSshCredentials(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshDir := fmt.Sprintf("%s/.ssh/", home)

	versions, err := WalkMatch(sshDir, "")
	if err != nil {
		return err
	}
	currentVersion, err := getCurrentVersion("id_rsa", "", sshDir)
	if err != nil {
		return err
	}

	fmt.Println("Locally Available ssh credentials:")
	for _, version := range versions {
		if !strings.HasSuffix(version, "_rsa") || version == "id_rsa"{
			continue
		}
		fmt.Printf("\t- %s", version)
		if version == currentVersion {
			fmt.Printf(" (Active)")
		}
		fmt.Printf("\n")
	}
	return nil
}

func SwitchSshCredentials(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("SSH credentials credentialsName not supplied")
	}
	credentialsName := c.Args().First()
	fmt.Printf("switching to SSH credentials: %s\n", credentialsName)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	symlinkSource := fmt.Sprintf("%s/.ssh/%s", home, credentialsName)
	if !fileExists(symlinkSource) {
		return fmt.Errorf("%s does not exist", symlinkSource)
	}
	symlinkTarget := fmt.Sprintf("%s/.gcp/credentials", home)

	if _, err := os.Lstat(symlinkTarget); err == nil {
		if err := os.Remove(symlinkTarget); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}
	}
	err = os.Symlink(symlinkSource, symlinkTarget)
	if err != nil {
		return fmt.Errorf("failed to link: %+v", err)
	}

	return nil
}
