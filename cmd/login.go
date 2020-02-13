package cmd

import (
	"fmt"
	"net/http"
)

func isUserAuthCorrect(url string, c EJConfig, hc HttpClient) int {
	code, err := hc.HEAD(url, c)
	if code == 200 {
		// means authentication is successful. Need to save the config
		return http.StatusOK
	} else if code == 401 {
		fmt.Printf("Error: User %s does not have access to: %s check your creds\n", c.Username, c.Url)
		return http.StatusUnauthorized
	} else if code == 403 {
		fmt.Printf("Error: User %s has the correct creds but no access \n", c.Username)
		return http.StatusForbidden
	} else {
		fmt.Printf("Error: %s\n", err)
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
