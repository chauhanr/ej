package cmd

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func getAllFields() ([]Field, error) {
	f := []Field{}
	c, err := checkEJConfigAndLoad()
	if err != nil {
		return f, err
	}

	b := JiraUrlBuilder{Base: c.Url}
	url := b.BuildGetFieldUrl("")
	h := &HttpClient{Client: &http.Client{}}

	rs, code := h.GET(url, c)
	if code == http.StatusOK || code == http.StatusCreated {
		// decode the json to fields.
		reader := strings.NewReader(rs)
		d := json.NewDecoder(reader)
		err := d.Decode(&f)
		if err != nil {
			return f, err
		}

	} else if code == http.StatusUnauthorized {
		return f, errors.New("User is unauthorized")
	} else if code == http.StatusForbidden {
		return f, errors.New("User is forbidden to access this resource")
	}
	return f, errors.New("Unknow Error with http status code: " + string(code))
}

func getSystemFields() ([]Field, error) {
	f := []Field{}
	_, err := checkEJConfigAndLoad()
	if err != nil {
		return f, err
	}
	f, err = getAllFields()
	if err != nil {
		return f, err
	}
	filtered := filterFields(f, false)
	return filtered, nil
}

func filterFields(fields []Field, needCustom bool) []Field {
	fs := []Field{}
	for _, f := range fields {
		if f.Custom == needCustom {
			fs = append(fs, f)
		}
	}
	return fs
}

func getCustomFields() ([]Field, error) {
	f := []Field{}
	_, err := checkEJConfigAndLoad()
	if err != nil {
		return f, err
	}
	f, err = getAllFields()
	filtered := filterFields(f, true)
	return filtered, nil
}

func checkEJConfigAndLoad() (EJConfig, error) {
	c := EJConfig{}
	err := c.loadConfig()
	if err != nil {
		return EJConfig{}, errors.New("User credentials could not be found")
	}
	return c, nil
}

/* Structure for Fields*/

type Field struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Custom bool   `json:"custom"`
}
