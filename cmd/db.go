package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

//IssueBoard for unique issues
type IssueBoard struct {
	IssueInBoard map[string][]string `json:"uniqueIssues"`
}

//IssueByType for unique
type IssueByType struct {
	IssueTypeMap map[string][]string `json:"issueByType"`
}

func (c *IssueBoard) saveData() error {
	ejd := getEJConfigDir()
	if _, err := os.Stat(ejd); !os.IsExist(err) {
		os.MkdirAll(ejd, os.ModePerm)
	}
	ej := getIssueBoardPath()
	file, err := os.OpenFile(ej, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("Error reading file: " + ej)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return errors.New("Error encoding file: " + ej + " error: " + err.Error())
	}
	return nil
}

func (c *IssueBoard) loadData() error {
	ej := getIssueBoardPath()
	file, err := os.Open(ej)
	defer file.Close()
	if err != nil {
		return err
	}
	d := json.NewDecoder(file)
	err = d.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *IssueByType) loadData() error {
	ej := getIssueTypePath()
	file, err := os.Open(ej)
	defer file.Close()
	if err != nil {
		return err
	}
	d := json.NewDecoder(file)
	err = d.Decode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *IssueByType) saveData() error {
	ejd := getEJConfigDir()
	if _, err := os.Stat(ejd); !os.IsExist(err) {
		os.MkdirAll(ejd, os.ModePerm)
	}
	ej := getIssueTypePath()
	file, err := os.OpenFile(ej, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("Error reading file: " + ej)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return errors.New("Error encoding file: " + ej + " error: " + err.Error())
	}
	return nil
}

func getIssueBoardPath() string {
	home, _ := os.UserHomeDir()
	ej := filepath.Join(home, EJ_HOME, ISSUE_BOARD_STORE)
	return ej
}

func getIssueTypePath() string {
	home, _ := os.UserHomeDir()
	ej := filepath.Join(home, EJ_HOME, ISSUETYPE_STORE)
	return ej
}
