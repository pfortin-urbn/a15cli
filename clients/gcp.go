package clients

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func ListGcpCredentials(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	gcpDir := fmt.Sprintf("%s/.gcp/", home)

	versions, err := WalkMatch(gcpDir, "")
	if err != nil {
		return err
	}
	currentVersion, err := getCurrentVersion("credentials", "", gcpDir)
	if err != nil {
		return err
	}

	fmt.Println("Locally Available gcp credentials:")
	for _, version := range versions {
		if version == "credentials"{
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

func SwitchGcpCredentials(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("GCP credentials credentialsName not supplied")
	}
	credentialsName := c.Args().First()
	fmt.Printf("switching to GCP credentials: %s\n", credentialsName)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	symlinkSource := fmt.Sprintf("%s/.gcp/%s", home, credentialsName)
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
