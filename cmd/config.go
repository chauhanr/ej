package cmd

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	EJ_HOME           = ".ej"
	EJ_CONF           = "ejconf.json"
	IGNORE_BOARD      = "board_ignore.json"
	ISSUE_BOARD_STORE = "issue-board.json"
	ISSUETYPE_STORE   = "issue-type.json"
	KANBAN_BOARDS     = "kanban_board.json"
	SCRUM_BOARDS      = "scrum_boards.json"
)

/*BoardsDatabase stores a list of boards based on scrum and kanban board*/
type BoardsDatabase struct {
	BoardType string  `json:"boardType"`
	BoardList []Board `json:"boards"`
}

func (c *BoardsDatabase) cleanConfig() error {
	if c.configExists() {
		p := getBoardDatabasePath(c.BoardType)
		err := os.Remove(p)
		if err != nil {
			fmt.Printf("Error cleaning config file %s, Error: %s\n", p, err)
			return err
		}
	} else {
		// do nothing
	}
	return nil
}

func (c *BoardsDatabase) configExists() bool {
	ej := getBoardDatabasePath(c.BoardType)
	if _, err := os.Stat(ej); err == nil {
		return true
	} else {
		return false
	}
}

func (c *BoardsDatabase) loadConfig() error {
	db := getBoardDatabasePath(c.BoardType)
	file, err := os.Open(db)
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

func getBoardDatabasePath(bType string) string {
	home, _ := os.UserHomeDir()
	var db string
	if bType == "scrum" {
		db = filepath.Join(home, EJ_HOME, SCRUM_BOARDS)
	} else {
		db = filepath.Join(home, EJ_HOME, KANBAN_BOARDS)
	}
	return db
}

func (c *BoardsDatabase) saveConfig() error {
	ejd := getEJConfigDir()
	if _, err := os.Stat(ejd); !os.IsExist(err) {
		os.MkdirAll(ejd, os.ModePerm)
	}
	var path string
	if c.BoardType == "scrum" {
		path = getBoardDatabasePath("scrum")
	} else if c.BoardType == "kanban" {
		path = getBoardDatabasePath("kanban")
	} else {
		return errors.New("unknow board type: " + c.BoardType)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("Error reading file: " + path)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return errors.New("Error encoding file: " + path + " error: " + err.Error())
	}
	return nil
}

type EJConfig struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	B64      string `json:"base64"`
}

func (c *EJConfig) saveConfig() error {
	ejd := getEJConfigDir()
	if _, err := os.Stat(ejd); !os.IsExist(err) {
		os.MkdirAll(ejd, os.ModePerm)
	}
	ej := getEJConfigPath()
	file, err := os.OpenFile(ej, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("Error reading file: " + ej)
	}
	c.Password = ""
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return errors.New("Error encoding file: " + ej + " error: " + err.Error())
	}
	return nil
}

/*
  This method will check if the config file is present. if not then we have to
  ask the user to login using the login command.
*/
func (c *EJConfig) configExists() bool {
	ej := getEJConfigPath()
	if _, err := os.Stat(ej); err == nil {
		return true
	} else {
		return false
	}
}

/*
  This method with load the ej config file to the EJConfig structure
  after this we can make use of the config to get the necessary details.
*/
func (c *EJConfig) loadConfig() error {
	ej := getEJConfigPath()
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

/*
  the clean method will remove the config.
  This method will be used in the logout functionality.
*/
func (c *EJConfig) cleanConfig() error {
	if c.configExists() {
		p := getEJConfigPath()
		err := os.Remove(p)
		if err != nil {
			fmt.Printf("Error cleaning config file %s, Error: %s\n", p, err)
			return err
		}
	} else {
		// do nothing
	}
	return nil
}

func getEJConfigPath() string {
	home, _ := os.UserHomeDir()
	ej := filepath.Join(home, EJ_HOME, EJ_CONF)
	return ej
}

func getEJConfigDir() string {
	home, _ := os.UserHomeDir()
	ej := filepath.Join(home, EJ_HOME)
	return ej
}

func EncodeCreds(username, password string) string {
	cred := username + ":" + password
	b64 := base64.URLEncoding.EncodeToString([]byte(cred))
	return b64
}

/*IgnoreBoardConfig will be used to save and load the ignore boards list*/
type IgnoreBoardConfig struct {
	BoardList []string `json:"boardIds"`
}

func (c *IgnoreBoardConfig) cleanConfig() error {
	if c.configExists() {
		p := getIgnoreBoardPath()
		err := os.Remove(p)
		if err != nil {
			fmt.Printf("Error cleaning config file %s, Error: %s\n", p, err)
			return err
		}
	} else {
		// do nothing
	}
	return nil
}

func getIgnoreBoardPath() string {
	home, _ := os.UserHomeDir()
	ej := filepath.Join(home, EJ_HOME, IGNORE_BOARD)
	return ej
}

func (c *IgnoreBoardConfig) configExists() bool {
	ej := getIgnoreBoardPath()
	if _, err := os.Stat(ej); err == nil {
		return true
	} else {
		return false
	}
}

func (c *IgnoreBoardConfig) saveConfig() error {
	ejd := getEJConfigDir()
	if _, err := os.Stat(ejd); !os.IsExist(err) {
		os.MkdirAll(ejd, os.ModePerm)
	}
	ej := getIgnoreBoardPath()
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

func (c *IgnoreBoardConfig) loadConfig() error {
	ej := getIgnoreBoardPath()
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

func (c *IgnoreBoardConfig) contains(bID string) bool {
	if !c.configExists() {
		return false
	}
	for _, b := range c.BoardList {
		if b == bID {
			return true
		}
	}
	return false

}
