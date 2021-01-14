package clients

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func ListTerraformVersions(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	binDir := fmt.Sprintf("%s/bin/", home)

	versions, err := WalkMatch(binDir, "terraform-*")
	if err != nil {
		return err
	}
	currentVersion, err := getCurrentVersion("terraform", "-", binDir)
	if err != nil {
		return err
	}

	fmt.Println("Locally Available terraform versions:")
	for _, version := range versions {
		fmt.Printf("\t- %s", version)
		if version == currentVersion {
			fmt.Printf(" (Active)")
		}
		fmt.Printf("\n")
	}
	return nil
}

func SwitchTerraformVersion(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("version number not supplied")
	}
	version := c.Args().First()
	fmt.Printf("switching to version: %s\n", version)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	symlinkSource := fmt.Sprintf("%s/bin/terraform-%s", home, version)
	symlinkTarget := fmt.Sprintf("%s/bin/terraform", home)

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

func InstallTerraformVersions(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("Need a version number to switch to")
	}
	version := c.Args().First()
	if !validateVersion.Match([]byte(version)) {
		return fmt.Errorf("version not in proper format (#.#.#)")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	binDir := fmt.Sprintf("%s/bin", home)
	tmpDir := "/tmp"
	filename := fmt.Sprintf("terraform-%s.zip", version)
	file := fmt.Sprintf("%s/%s", tmpDir, filename)

	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_darwin_amd64.zip", version, version)

	err = DownloadFile(file, url)
	if err != nil {
		return err
	}

	unzip(file, "/tmp/")
	err = os.Rename("/tmp/terraform", fmt.Sprintf("%s/terraform-%s", binDir, version))
	if err != nil {
		return err
	}

	err = os.Remove(file)
	if err != nil {
		return err
	}
	// Save version # in versions installed

	fmt.Printf("Terraform %s has been downloaded and extracted (it is not set as active)", version)
	return nil
}


