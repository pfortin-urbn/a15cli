package clients

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func ListAwsCredentials(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	awsDir := fmt.Sprintf("%s/.aws/", home)

	versions, err := WalkMatch(awsDir, "credentials.*")
	if err != nil {
		return err
	}
	currentVersion, err := getCurrentVersion("credentials", ".", awsDir)
	if err != nil {
		return err
	}

	fmt.Println("Locally Available aws credentials:")
	for _, version := range versions {
		fmt.Printf("\t- %s", version)
		if version == currentVersion {
			fmt.Printf(" (Active)")
		}
		fmt.Printf("\n")
	}
	return nil
}

func SwitchAwsCredentials(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("AWS credentials credentialsName not supplied")
	}
	credentialsName := c.Args().First()
	fmt.Printf("switching to AWS credentials: %s\n", credentialsName)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	symlinkSource := fmt.Sprintf("%s/.aws/credentials.%s", home, credentialsName)
	symlinkTarget := fmt.Sprintf("%s/.aws/credentials", home)

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
