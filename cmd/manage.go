package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	EJ_HOME = ".ej"
	EJ_CONF = "conf.json"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number for the ej app",
	Long:  `All software has versions. This is ej's`,
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println(banner)
	fmt.Println("EJ - Explore Jira version 1.0")
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login the user to a Jira instance",
	Long:  `Allow the user to login to the Jira instance by asking the credentials.`,
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {
	u := cmd.Flag("username")
	p := cmd.Flag("password")
	a := cmd.Flag("url")
	h := os.Getenv("HOME")

	ep := filepath.Join(h, EJ_HOME, EJ_CONF)

	un := u.Value.String()
	pa := p.Value.String()
	url := a.Value.String()

	c := EJConfig{Username: un, Password: pa, Url: url}
	err := c.saveConfig(ep)
	if err != nil {
		fmt.Printf("Error saving creds: %s\n", err)
	}
}
