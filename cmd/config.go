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
	EJ_HOME = ".ej"
	EJ_CONF = "ejconf.json"
)

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
	b64 := base64.StdEncoding.EncodeToString([]byte(cred))
	return b64
}
