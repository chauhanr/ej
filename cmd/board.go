package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var boardCmd = &cobra.Command{
	Use:   "boards",
	Short: "board command is used to list kanban or scrum boards and give numbers",
	Long:  `Board command lists the kanban/scrum boards in the Jira instance and can list the boards.`,
	Run:   boardMsg,
}

func boardMsg(cmd *cobra.Command, args []string) {
	fmt.Printf("Board command is parent command use an appropriate subcommand")
	return
}

var boardIgnoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "takes a list of board ids via flag and stores the preferences",
	Long:  `This command ignores the board ids stored from being processed.`,
	Run:   boardsIgnore,
}

func boardsIgnore(cmd *cobra.Command, args []string) {
	bIgnoreList, _ := cmd.Flags().GetString(BOARDLIST)
	l := strings.Split(bIgnoreList, " ")

	ig := IgnoreBoardConfig{}
	err := ig.cleanConfig()

	if err != nil {
		fmt.Printf("Ignore Config delete failed %s", err)
	}
	ig.BoardList = l
	err = ig.saveConfig()
	if err != nil {
		fmt.Printf("Ignore Config save failed %s", err)
	}
	return

}

var boardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the boards of a particular kind",
	Long:  `List all boards of a particular kind`,
	Run:   boardsList,
}

func boardsList(cmd *cobra.Command, args []string) {
	if !areUserCredsSaved() {
		fmt.Println("No user credentails found.")
		loginCmd.Usage()
	} else {
		bType, _ := cmd.Flags().GetString(BOARD_TYPE)
		bDB := BoardsDatabase{}
		bDB.BoardType = bType
		if bDB.configExists() {
			err := bDB.loadConfig()
			if err != nil {
				fmt.Printf("error loading data from saved config for board type: %s", bType)
				return
			}
			fmt.Printf("There are %d %s type boards in this Jira instance.", len(bDB.BoardList), bType)
			return
		}
		c := EJConfig{}
		c.loadConfig()
		boards := []Board{}
		bType = strings.ToLower(bType)
		if bType != "scrum" && bType != "kanban" {
			fmt.Printf("Kindly choose scrum/kanban board type use have entered: %s", bType)
			return
		}
		boards, err := getBoards(bType, c)
		if err != nil {
			fmt.Println(err)
			return
		}
		// calculate the last page work
		fmt.Printf("There are %d %s type boards in this Jira instance.", len(boards), bType)
		// save the file
		bDB.BoardList = boards
		err = bDB.saveConfig()
		if err != nil {
			fmt.Printf("error saving the boards to the file %s\n", err)
			return
		}
	}
}

func getBoards(bType string, c EJConfig) (boards []Board, err error) {
	b := JiraUrlBuilder{Base: c.Url}
	bURL := b.BuildBoardURL("", bType, 50, 0)
	isLastPage := false
	startAt := 0
	for !isLastPage {
		isLast, bPage, maxResults, errB := loadBoards(bURL, startAt, c)
		isLastPage = isLast
		if errB != nil {
			return nil, fmt.Errorf("error processing boards %s ", err)
		}
		startAt += maxResults
		bURL = b.BuildBoardURL("", bType, maxResults, startAt)
		boards = append(boards, bPage...)
	}
	return boards, nil
}

func loadBoards(bURL string, startIndex int, c EJConfig) (isLast bool, boards []Board, maxResults int, err error) {
	isLast = false
	boards = []Board{}
	h := &HttpClient{Client: &http.Client{}}
	rs, code := h.GET(bURL, c)
	bList := BoardList{}
	if code == http.StatusOK || code == http.StatusCreated {
		// decode the json to fields.
		reader := strings.NewReader(rs)
		d := json.NewDecoder(reader)
		err := d.Decode(&bList)
		if err != nil {
			return isLast, boards, maxResults, err
		}
	} else if code == http.StatusUnauthorized {
		return isLast, boards, maxResults, errors.New("User is not authenticated kindly check creds.")
	} else if code == http.StatusForbidden {
		return isLast, boards, maxResults, errors.New("User is not authorized to get the issue list")
	}
	maxResults = bList.MaxResults
	isLast = bList.IsLast
	boards = append(boards, bList.Values...)
	return isLast, boards, maxResults, err

}

/*BoardList is the data structure used to list the board on the current project*/
type BoardList struct {
	StartIndex int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	IsLast     bool    `json:"isLast"`
	Values     []Board `json:"values"`
}

/*Board structure represents the single board.*/
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
