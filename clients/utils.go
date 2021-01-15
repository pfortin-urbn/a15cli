package clients

import (
	"archive/zip"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const(
	ANSI_RESET = "\u001B[0m"
	ANSI_RED = "\u001B[31m"
)

var validateVersion = regexp.MustCompile(`^\d[0-9.]*$`)
var listTerraformVersions = `^terraform_\d[0-9.]*$`

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
		if pattern != "" {
			if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
				return err
			} else if matched {
				matches = append(matches, after(info.Name(), pattern[:len(pattern)-1]))
			}
		} else {
			matches = append(matches, info.Name())
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

func getCurrentVersion(linkName, separator, binDir string) (string, error){
	filename, err := os.Readlink(fmt.Sprintf("%s/%s", binDir, linkName))
	if err != nil {
		return "", err
	}
	if separator != "" {
		return after(filename, fmt.Sprintf("%s%s", linkName, separator)), nil
	}
	return filepath.Base(filename), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func  ReadConfig(path string, receiver interface{}) error {
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)
	err = yaml.Unmarshal(yamlFile, receiver)
	if err != nil {
		return err
	}
	return nil
}

func  WriteConfig(path string, data interface{}) error {
	filename, _ := filepath.Abs(path)

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}
