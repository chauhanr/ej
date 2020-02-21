package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func getAllFields() (Fields, error) {
	f := []Field{}
	c, err := checkEJConfigAndLoad()
	if err != nil {
		return Fields{Field: f}, err
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
			return Fields{Field: f}, err
		}
		return Fields{Field: f}, nil

	} else if code == http.StatusUnauthorized {
		return Fields{Field: f}, errors.New("User is unauthorized")
	} else if code == http.StatusForbidden {
		return Fields{Field: f}, errors.New("User is forbidden to access this resource")
	}
	e := fmt.Sprintf("Unknow Error with http status code: %d\n", code)
	return Fields{Field: f}, errors.New(e)
}

func getSystemFields() (Fields, error) {
	_, err := checkEJConfigAndLoad()
	if err != nil {
		return Fields{}, err
	}
	f, err := getAllFields()
	if err != nil {
		return f, err
	}
	filtered := filterFields(f.Field, false)
	return Fields{Field: filtered}, nil
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

func getCustomFields() (Fields, error) {
	_, err := checkEJConfigAndLoad()
	if err != nil {
		return Fields{}, err
	}
	f, err := getAllFields()
	filtered := filterFields(f.Field, true)
	return Fields{Field: filtered}, nil
}

func checkEJConfigAndLoad() (EJConfig, error) {
	c := EJConfig{}
	err := c.loadConfig()
	if err != nil {
		return EJConfig{}, errors.New("User credentials could not be found")
	}
	return c, nil
}

type Fields struct {
	Field []Field
}

func (f *Fields) DisplayFields() {
	if len(f.Field) != 0 {
		fmt.Printf("------------------------------------------------------------------------\n")
		fmt.Printf("%20.20s | %30.30s | %8.8s\n", "Field Id", "Field Name", "Custom")
		fmt.Printf("------------------------------------------------------------------------\n")
		for _, f := range f.Field {
			fmt.Printf("%20.20s | %30.30s | %5.5t\n", f.Id, f.Name, f.Custom)
		}
	}
}

/* Structure for Fields*/

type Field struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Custom bool   `json:"custom"`
}
