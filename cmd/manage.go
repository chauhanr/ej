package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
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

	un := u.Value.String()
	pa := p.Value.String()
	url := a.Value.String()

	b := JiraUrlBuilder{Base: url}

	c := EJConfig{Username: un, Password: pa, Url: url, B64: EncodeCreds(un, pa)}
	hc := HttpClient{Client: &http.Client{}}

	auth := isUserAuthCorrect(b.BuildAuthCheckUrl(""), c, hc)
	if auth == http.StatusOK {
		err := c.saveConfig()
		if err != nil {
			fmt.Printf("Error saving creds: %s\n", err)
		} else {
			fmt.Printf("Success: you have successfully logged into Jira instance\n")
		}
	} else {
		fmt.Println("Explorer could not authenticate you against the Jira instance.")
	}
}
