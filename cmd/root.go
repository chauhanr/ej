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
	URL               string
	fieldCustomFilter bool
	fieldSystemFilter bool
	projectId         string
	boardType         string
	boardList         string
	startDate         string
	outputFormat      bool
	displayFlagUnique bool
	displayFlagType   bool
)

const (
	FIELD_CUSTOM_FILTER = "custom-filter"
	FIELD_SYSTEM_FILTER = "system-filter"
	PROJECT_FILTER      = "project-id"
	BOARD_TYPE          = "board-type"
	JSON_OUTPUT_FORMAT  = "json-format"
	BOARDLIST           = "board-list"
	UNIQUE_ISSUE        = "unique-issue"
	ISSUE_BY_TYPE       = "issue-bytype"
	START_DATE          = "start-date"
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
	loginCmd.Flags().StringVarP(&URL, "URL", "a", "", "URL for the jira instance")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
	loginCmd.MarkFlagRequired("URL")
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
	ptreeCmd.Flags().BoolVarP(&outputFormat, JSON_OUTPUT_FORMAT, "o", false, "display the project structure in json")

	rootCmd.AddCommand(ptreeCmd)

	boardListCmd.Flags().StringVarP(&boardType, BOARD_TYPE, "t", "", "board type can have two values kanban or scrum")
	boardCmd.AddCommand(boardListCmd)
	boardIgnoreCmd.Flags().StringVarP(&boardList, BOARDLIST, "l", "", "list of boards by space values while will be split")
	boardIgnoreCmd.MarkFlagRequired(BOARDLIST)
	boardCmd.AddCommand(boardIgnoreCmd)
	rootCmd.AddCommand(boardCmd)

	issueListCmd.Flags().StringVarP(&boardType, BOARD_TYPE, "t", "", "board type can have two values kanban or scrum")
	issueListCmd.Flags().StringVarP(&startDate, START_DATE, "s", "", "board type can have two values kanban or scrum")
	issueCmd.AddCommand(issueListCmd)
	issueDisplayCmd.Flags().StringVarP(&boardType, BOARD_TYPE, "t", "", "board type can have two values kanban or scrum")
	issueDisplayCmd.Flags().BoolVarP(&displayFlagUnique, UNIQUE_ISSUE, "u", false, "this will display all the unique issues across all boards.")
	issueDisplayCmd.Flags().BoolVarP(&displayFlagType, ISSUE_BY_TYPE, "i", false, "this will display all the unique issues across all boards.")
	issueCmd.AddCommand(issueDisplayCmd)
	rootCmd.AddCommand(issueCmd)

}
