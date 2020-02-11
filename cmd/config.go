package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type EJConfig struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *EJConfig) saveConfig(f string) error {
	h := os.Getenv("HOME")
	ej := filepath.Join(h, EJ_HOME)
	if _, err := os.Stat(ej); !os.IsExist(err) {
		os.MkdirAll(ej, os.ModePerm)
	}
	file, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("Error reading file: " + f)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		return errors.New("Error encoding file: " + f + " error: " + err.Error())
	}
	return nil
}

/*
  This method will check if the config file is present. if not then we have to
  ask the user to login using the login command.
*/
func (c *EJConfig) configExists() bool {
	h := os.Getenv("HOME")
	ej := filepath.Join(h, EJ_HOME, EJ_CONF)
	if _, err := os.Stat(ej); !os.IsExist(err) {
		return false
	} else {
		return true
	}
}
