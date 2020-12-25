package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "issue command will allow for issues to be listed or purged",
	Long:  `issue command will allow for issues to be listed or purged from local database`,
	Run:   issueMsg,
}

func issueMsg(cmd *cobra.Command, args []string) {
	fmt.Printf("Issue command to list issues in various ways")
	return
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "list sub command with issue allows for issues to be listed by either board or issue type",
	Run:   issueList,
}

func issueList(cmd *cobra.Command, args []string) {
	if !areUserCredsSaved() {
		fmt.Println("No user credentails found.")
		loginCmd.Usage()
	} else {
		bType, _ := cmd.Flags().GetString(BOARD_TYPE)
		startDate, _ := cmd.Flags().GetString(START_DATE)
		fmt.Printf("StartDate: %s\n", startDate)
		bDB := BoardsDatabase{}
		bIgnore := IgnoreBoardConfig{}
		bIgnore.loadConfig()
		bDB.BoardType = bType
		issueBoardMap := make(map[string][]string)
		issueTypeMap := make(map[string][]string)
		if bDB.configExists() {
			err := bDB.loadConfig()
			if err != nil {
				fmt.Printf("error loading data from saved config for board type: %s", bType)
				return
			}
			boards := bDB.BoardList
			// now iterate through the board
			for _, b := range boards {
				boardID := fmt.Sprintf("%d", b.ID)
				if bIgnore.contains(boardID) {
					continue
				}
				fmt.Printf("Processing issues for board id: %s, name: %s\n", boardID, b.Name)
				// process the board for issues.
				ilist, err := getIssueForBoard(boardID, startDate)
				fmt.Printf("Issues in board %s is %d\n", boardID, len(ilist))
				if err != nil {
					fmt.Printf("error querying issue using the board %s, error: %s\n", boardID, err)
				}
				for _, i := range ilist {
					iKey := i.Key
					if v, ok := issueBoardMap[iKey]; ok {
						v = append(v, boardID)
						issueBoardMap[iKey] = v
						// there is not need to add to the issuetype as we add it only the first time.
					} else {
						// key does not exist add key to teh issueBoardMap and also the issueType Map
						boards := []string{boardID}
						issueBoardMap[iKey] = boards
						issueTypeKey := i.Fields.IssueType.Name
						if il, ok := issueTypeMap[issueTypeKey]; ok {
							il = append(il, iKey)
							issueTypeMap[issueTypeKey] = il
						} else {
							v := []string{iKey}
							issueTypeMap[issueTypeKey] = v
						}
					}
				}
			}
			fmt.Printf("Unique Issue in %s type boards %d", bType, len(issueBoardMap))
			ib := IssueBoard{}
			ib.IssueInBoard = issueBoardMap
			ib.saveData()

			it := IssueByType{}
			it.IssueTypeMap = issueTypeMap
			it.saveData()
		} else {
			// because the bDB does not exist give the command help for boards
			boardCmd.Usage()
			return
		}
	}
}

func getIssueForBoard(boardID, startDate string) (issues []Issue, err error) {
	issues = []Issue{}
	ej := EJConfig{}
	ej.loadConfig()
	b := JiraUrlBuilder{Base: ej.Url}
	qURL := b.BuildBoardIssueQuery("", boardID, 50, 0, startDate)
	startAt := 0
	isLastPage := false
	for !isLastPage {
		bPage, maxResults, total, errB := queryIssues(qURL, startAt, ej)
		if total <= startAt {
			isLastPage = true
		}
		if errB != nil {
			return nil, fmt.Errorf("error processing boards %s ", err)
		}
		startAt += maxResults
		qURL = b.BuildBoardIssueQuery("", boardID, maxResults, startAt, startDate)
		issues = append(issues, bPage...)
	}
	return issues, nil
}

func queryIssues(query string, startAt int, c EJConfig) (issues []Issue, maxResults, total int, err error) {
	total = 0
	issues = []Issue{}
	h := &HttpClient{Client: &http.Client{}}
	rs, code := h.GET(query, c)
	iList := IssueList{}
	if code == http.StatusOK || code == http.StatusCreated {
		// decode the json to fields.
		reader := strings.NewReader(rs)
		d := json.NewDecoder(reader)
		err := d.Decode(&iList)
		if err != nil {
			return issues, maxResults, total, err
		}
	} else if code == http.StatusUnauthorized {
		return issues, maxResults, total, errors.New("User is not authenticated kindly check creds.")
	} else if code == http.StatusForbidden {
		return issues, maxResults, total, errors.New("User is not authorized to get the issue list")
	}
	maxResults = iList.MaxResults
	total = iList.Total
	issues = append(issues, iList.Issues...)
	return issues, maxResults, total, nil
}

var issueDisplayCmd = &cobra.Command{
	Use:   "display",
	Short: "display commands displays issues based on flags for data collected under the list command.",
	Run:   displayIssue,
}

func displayIssue(cmd *cobra.Command, args []string) {
	bType, _ := cmd.Flags().GetString(BOARD_TYPE)
	unique, _ := cmd.Flags().GetBool(UNIQUE_ISSUE)
	iType, _ := cmd.Flags().GetBool(ISSUE_BY_TYPE)

	if bType == "kanban" {
		if !unique && !iType {
			fmt.Println("choose either of the flags unique or issue by type to display the stats")
			return
		} else if unique && iType {
			fmt.Println("you can choose just either of issue by type or unique flags but not both")
			return
		} else if unique {
			ibDB := IssueBoard{}
			err := ibDB.loadData()
			if err != nil {
				fmt.Println("currently the kanban board database is not loaded")
				issueListCmd.Usage()
				return
			}
			fmt.Printf("Total number of unique issues in the kanban boards is: %d\n", len(ibDB.IssueInBoard))
			return
		} else if iType {
			ibyType := IssueByType{}
			err := ibyType.loadData()
			if err != nil {
				fmt.Printf("currently the kanban board database in not loaded")
				issueListCmd.Usage()
				return
			}
			fmt.Printf("%40s | %12s\n", "Issue Type", "Issue Count")
			fmt.Printf("--------------------------------------------------------------------------\n")
			for k, v := range ibyType.IssueTypeMap {
				fmt.Printf("%40s | %12d\n", k, len(v))
			}
		}

	} else {
		fmt.Printf("Board type %s is currently not suppported.", bType)
	}
}

//IssueList has the list that is give on a response
type IssueList struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

//Issue is used as the struct to capture the issue.
type Issue struct {
	Key    string `json:"key"`
	Fields JField `json:"fields"`
}

//JField lists the fields are are interested in.
type JField struct {
	IssueType IssueType `json:"issuetype"`
	Status    Status    `json:"status"`
	//Sprint    Sprint    `json:"sprint,omitempty"`
}

//Status stores the status of the issue
type Status struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//IssueType is a structure that wil be used for storing the issue type.
type IssueType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//Sprint just captures basic sprint details and will be used to determine if issue has a sprint.
type Sprint struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SprintState string `json:"state"`
}
