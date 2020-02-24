package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout and foget previous session if any",
	Long:  `Allow the user to logout from the Jira instance and forget the previous session.`,
	Run:   logout,
}

func logout(cmd *cobra.Command, args []string) {
	conf := EJConfig{}
	err := conf.cleanConfig()

	if err != nil {
		fmt.Printf("logout unsuccessful\n")
	} else {
		fmt.Printf("logout successful\n")
	}
}

var fieldsCmd = &cobra.Command{
	Use:   "field",
	Short: "lists fields in the Jira instance.",
	Long:  "The fields can be filtered on being custom or system generated one",
	Run:   getFields,
}

func getFields(cmd *cobra.Command, args []string) {
	if !areUserCredsSaved() {
		fmt.Println("No user credentails found.")
		loginCmd.Usage()
	} else {
		// call the get Fields REST call.
		c, _ := cmd.Flags().GetBool(FIELD_CUSTOM_FILTER)
		s, _ := cmd.Flags().GetBool(FIELD_SYSTEM_FILTER)
		if c && !s {
			f, err := getCustomFields()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			f.DisplayFields()
			return
		} else if !c && s {
			f, err := getSystemFields()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			f.DisplayFields()
			return
		} else {
			f, err := getAllFields()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			f.DisplayFields()
			return
		}
	}
}

var ptreeCmd = &cobra.Command{
	Use:   "ptree",
	Short: "get the project structure displayed.",
	Long:  `Project, Issues and Sprint structure in tree form`,
	Run:   pTree,
}

/*  this methos prepares the project tree structure*/
func pTree(cmd *cobra.Command, args []string) {
	if !areUserCredsSaved() {
		fmt.Println("No user credentails found.")
		loginCmd.Usage()
	} else {
		pId, _ := cmd.Flags().GetString(PROJECT_FILTER)
		c := EJConfig{}
		c.loadConfig()
		p, err := getProjectTree(pId, c)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "\t")
		enc.Encode(p)
	}
}

/*
  check the creds are present or not.
*/
func areUserCredsSaved() bool {
	conf := EJConfig{}
	err := conf.loadConfig()
	if err != nil {
		return false
	} else {
		return true
	}
}
