package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
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
		if !strings.HasSuffix(version, "_rsa") || version == "id_rsa" {
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

	symlinkSource := fmt.Sprintf("%s/.ssh/%s_id_rsa", home, credentialsName)
	symlinkSourcePub := fmt.Sprintf("%s/.ssh/%s_id_rsa.pub", home, credentialsName)
	if !fileExists(symlinkSource) || !fileExists(symlinkSourcePub) {
		return fmt.Errorf("%s does not exist", symlinkSource)
	}
	symlinkTarget := fmt.Sprintf("%s/.ssh/id_rsa", home)
	symlinkTargetPub := fmt.Sprintf("%s/.ssh/id_rsa.pub", home)

	if _, err := os.Lstat(symlinkTarget); err == nil {
		if err := os.Remove(symlinkTarget); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}
	}
	if _, err := os.Lstat(symlinkTargetPub); err == nil {
		if err := os.Remove(symlinkTargetPub); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}
	}
	err = os.Symlink(symlinkSource, symlinkTarget)
	if err != nil {
		return fmt.Errorf("failed to link: %+v", err)
	}

	err = os.Symlink(symlinkSourcePub, symlinkTargetPub)
	if err != nil {
		return fmt.Errorf("failed to link: %+v", err)
	}

	return nil
}
