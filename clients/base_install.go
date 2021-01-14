package clients

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os/exec"
)

func BaseInstall(c *cli.Context) error {
	fmt.Println("Installing `httpie`")
	out, err := exec.Command("brew", "install", "httpie").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(out))

	fmt.Println("Installing `jq`")
	out, err = exec.Command("brew", "install", "jq").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(out))

	//brew tap rockymadden/rockymadden; brew install rockymadden/rockymadden/slack-cli
	fmt.Println("Installing `slack-cli`")
	out, err = exec.Command("brew", "tap", "rockymadden/rockymadden").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(out))
	out, err = exec.Command("brew", "install", "rockymadden/rockymadden/slack-cli").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(out))

	fmt.Println("Installing `yamllint`")
	out, err = exec.Command("brew", "install", "yamllint").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(out))

	fmt.Println("Installing `yq`")
	out, err = exec.Command("brew", "install", "yq").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(string(out))
	return nil
}

