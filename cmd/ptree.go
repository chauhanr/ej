package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func prepareProjectTree(pId string, config EJConfig) (JIssueList, error) {

	b := JiraUrlBuilder{Base: config.Url}
	url := b.BuildProjectIssueUrl("", pId, "")

	h := &HttpClient{Client: &http.Client{}}
	rs, code := h.GET(url, config)
	il := JIssueList{}
	if code == http.StatusOK || code == http.StatusCreated {
		// decode the json to fields.
		reader := strings.NewReader(rs)
		d := json.NewDecoder(reader)
		err := d.Decode(&il)
		if err != nil {
			return il, err
		}
		issues, _ := getAllIssueIdMap(il.Issues)
		fmt.Printf("Issue Ids %v\n", issues)
		return il, nil

	} else if code == http.StatusUnauthorized {
		return il, errors.New("User is not authenticated kindly check creds.")
	} else if code == http.StatusForbidden {
		return il, errors.New("User is not authorized to get the issue list")
	}
	errMsg := fmt.Sprintf("Error with code : %d has occurred.\n", code)
	return il, errors.New(errMsg)
}

func getProjectTreeStruct(issueIds []string, issueMap []JIssue, config EJConfig) error {
	b := JiraUrlBuilder{Base: config.Url}

	for _, id := range issueIds {
		url := b.BuildIssueSprintResponseUrl("", id)
		h := &HttpClient{Client: &http.Client{}}
		rs, code := h.GET(url, config)
		is := JIssue{}
		if code == http.StatusOK || code == http.StatusCreated {
			// decode the json to fields.
			reader := strings.NewReader(rs)
			d := json.NewDecoder(reader)
			err := d.Decode(&is)
			if err != nil {
				return err
			}

		} else if code == http.StatusUnauthorized {
			return errors.New("User is not authenticated kindly check creds.")
		} else if code == http.StatusForbidden {
			return errors.New("User is not authorized to get the issue list")
		}
		errMsg := fmt.Sprintf("Error with code : %d has occurred.\n", code)
		return errors.New(errMsg)

	}
	return nil
}

func getAllIssueIdMap(issues []JIssue) ([]string, map[string]JIssue) {
	ids := []string{}
	issueMap := map[string]JIssue{}
	for _, issue := range issues {
		ids = append(ids, issue.Id)
		key := getIssueKey(issue)
		issueMap[key] = issue
	}
	fmt.Printf("map: %v\n", issueMap)
	return ids, issueMap
}

func getIssueKey(i JIssue) string {
	key := i.IField.IssueType.Name + ":" + i.Id
	return key
}

/*
   structures we will need to be to make the structure
   1. Issue
   2. Project
   3. Sprint
*/

type JIssueList struct {
	StartIndex int      `json:"startAt"`
	MaxResults int      `json:"maxResults"`
	Total      int      `json:"total"`
	Issues     []JIssue `json:"issues"`
}

type JIssue struct {
	Id     string  `json:"id"`
	Key    string  `json:"key"`
	IField JIField `json:"fields,omitempty"`
}

type JIField struct {
	IssueType     JIssueType `json:"issuetype,omitempty"`
	Status        JStatus    `json:"status,omitempty"`
	Sprint        JSprint    `json:"sprint,omitempty"`
	ClosedSprints JSprint    `json:"closedSprints,omitempty"`
	Parent        *JIssue    `json:"parent,omitempty"`
	Priority      JPriority  `json:"priority,omitempty"`
}

type JPriority struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type JIssueType struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Subtask bool   `json:"subtask,omitempty"`
	Self    string `json:"self,omitempty"`
}

type JStatus struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Self string `json:"self,omitempty"`
}

type JSprint struct {
	Id            string `json:"id"`
	State         string `json:"state"`
	Name          string `json:"name"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	CompletedDate string `json:"completedDate,omitempty"`
	OriginBoardId int    `json:"originBoardId,omitempty"`
}
