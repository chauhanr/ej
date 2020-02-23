package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var banner = `
  ______     _ 
 |  ____|   | |
 | |__      | |
 |  __| _   | |
 | |___| |__| |
 |______\____/ 
               `

var (
	username          string
	password          string
	url               string
	fieldCustomFilter bool
	fieldSystemFilter bool
	projectId         string
)

const (
	FIELD_CUSTOM_FILTER = "custom-filter"
	FIELD_SYSTEM_FILTER = "system-filter"
	PROJECT_FILTER      = "project-id"
)

var rootCmd = &cobra.Command{
	Use:   "ej",
	Short: "Jira Explorer",
	Long:  "Jira Explorer that allows for Jira schema exploration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner)
		fmt.Println("Explore Jira Cli")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	/* Manangement Commands
	 * Version - simple version command that gives the version of the application.
	 * Login - command allows for the command to login to the system
	 * Logout - this command will allow the users logout from the system.
	 */
	rootCmd.AddCommand(versionCmd)

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "username for login")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "user password login")
	loginCmd.Flags().StringVarP(&url, "url", "a", "", "url for the jira instance")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
	loginCmd.MarkFlagRequired("url")
	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(logoutCmd)

	/* Basics commands to list fields etc.*/
	fieldsCmd.Flags().BoolVarP(&fieldCustomFilter, FIELD_CUSTOM_FILTER, "c", false, "Filter custom fields only")
	fieldsCmd.Flags().BoolVarP(&fieldSystemFilter, FIELD_SYSTEM_FILTER, "s", false, "Filter system fields only")

	rootCmd.AddCommand(fieldsCmd)

	/**
	  project structure and sprints along with issues.
	  pstruct - project structure.
	*/
	ptreeCmd.Flags().StringVarP(&projectId, PROJECT_FILTER, "p", "", "project id / key to search for the issue list")
	ptreeCmd.MarkFlagRequired(PROJECT_FILTER)

	rootCmd.AddCommand(ptreeCmd)
}
