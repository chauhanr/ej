package cmd

import (
	"fmt"
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
	h := os.Getenv("HOME")

	fmt.Printf("Home Dir: %s\n", h)
	fmt.Printf("Username:  %s, Password %s, url %s\n", u.Value, p.Value, a.Value)

}
