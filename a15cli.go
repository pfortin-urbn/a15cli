package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

var validateVersion = regexp.MustCompile(`^\d[0-9.]*$`)
var listTerraformVersions = `^terraform_\d[0-9.]*$`


func main() {
	app := &cli.App{
		Name: "a15cli",
		Usage: "A15 tool to manage local dev stuff",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   `Install some other useful tools: (httpie, jq, slack-cli, yamllint, yq)`,
				Action:  func(c *cli.Context) error {
					return nil
				},
			},
			{
				// https://releases.hashicorp.com/terraform/0.14.4/terraform_0.14.4_darwin_amd64.zip
				Name:    "terra-install",
				Aliases: []string{"ti"},
				Usage:   "install/overwrite a new terraform versions",
				Action:  func(c *cli.Context) error {
					if c.Args().Len() != 1 {
						return errors.New("Need a version number to switch to")
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
					filename := fmt.Sprintf("terraform_%s.zip", version)
					file := fmt.Sprintf("%s/%s", tmpDir, filename)

					url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_darwin_amd64.zip", version, version)

					err = DownloadFile(file, url)
					if err != nil {
						return err
					}

					unzip(file, "/tmp/")
					err = os.Rename("/tmp/terraform", fmt.Sprintf("%s/terraform_%s", binDir, version))
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
				},
			},
			{
				Name:    "terra-show",
				Aliases: []string{"ts"},
				Usage:   "list all available terraform versions",
				Action:  func(c *cli.Context) error {
					home, err := os.UserHomeDir()
					if err != nil {
						return err
					}

					binDir := fmt.Sprintf("%s/bin/", home)

					versions, err := WalkMatch(binDir, "terraform_*")
					if err != nil {
						return err
					}
					currentVersion, err := getCurrentVersion(binDir)
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
				},
			},
			{
				Name:    "terra-switch",
				Aliases: []string{"tw"},
				Usage:   "switch active terraform versions",
				Action:  func(c *cli.Context) error {
					if c.Args().Len() != 1 {
						return fmt.Errorf("version number not supplied")
					}
					version := c.Args().First()
					fmt.Printf("switching to version: %s\n", version)

					home, err := os.UserHomeDir()
					if err != nil {
						return err
					}

					symlinkSource := fmt.Sprintf("%s/bin/terraform_%s", home, version)
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
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, destDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath,string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, after(filepath.Base(path), "terraform_"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func getCurrentVersion(binDir string) (string, error){
	filename, err := os.Readlink(fmt.Sprintf("%s/terraform", binDir))
	if err != nil {
		return "", err
	}
	return after(filename, "terraform_"), nil
}