package cmd

/**
  this file will contain functions which will build the rest api calls
  there will be placeholder URL that will have to be changed.
*/

import (
	"strings"
)

const (
	DEFAULT_AGILE_VERSION = "1.0"
	DEFAULT_API_VERSION   = "2"

	REST_API   = "/rest/api/"
	REST_AGILE = "/rest/agile/"

	AUTH_CHECK_URL = REST_AGILE + VERSION_HOLDER + "/board"
	VERSION_HOLDER = "{version}"

	FIELDS_URL = REST_API + VERSION_HOLDER + "/field"

	PROJECT_HOLDER        = "{project-id}"
	PAGE_SIZE             = "{page-size}"
	ISSUE_ID_HOLDER       = "{issue-id}"
	PROJECT_ISSUE_JQL_URL = REST_API + VERSION_HOLDER + "/search?jql=project=" + PROJECT_HOLDER + "&fields=id,key,issuetype&maxResults=" + PAGE_SIZE

	AGILE_ISSUE_SPRINT_URL = REST_AGILE + VERSION_HOLDER + "/issue/" + ISSUE_ID_HOLDER + "?fields=parent,sprint,closedSprints"
)

type JiraUrlBuilder struct {
	Base string
}

func (b *JiraUrlBuilder) BuildAuthCheckUrl(v string) string {
	if v == "" {
		v = DEFAULT_AGILE_VERSION
	}
	url := handleTrailingSlash(b.Base)
	url = url + AUTH_CHECK_URL

	url = strings.Replace(url, VERSION_HOLDER, v, -1)
	return url
}

func (b *JiraUrlBuilder) BuildProjectIssueUrl(v string, projectId string, page string) string {
	url := handleTrailingSlash(b.Base)
	if v == "" {
		v = DEFAULT_API_VERSION
	}
	if page == "" || page == "0" {
		page = "100"
	}
	url = url + PROJECT_ISSUE_JQL_URL
	url = strings.Replace(url, VERSION_HOLDER, v, -1)
	url = strings.Replace(url, PROJECT_HOLDER, projectId, -1)
	url = strings.Replace(url, PAGE_SIZE, page, -1)
	return url
}

func (b *JiraUrlBuilder) BuildIssueSprintResponseUrl(v string, issueId string) string {
	url := handleTrailingSlash(b.Base)
	if v == "" {
		v = DEFAULT_AGILE_VERSION
	}
	url = url + PROJECT_ISSUE_JQL_URL
	url = strings.Replace(url, VERSION_HOLDER, v, -1)
	url = strings.Replace(url, ISSUE_ID_HOLDER, issueId, -1)
	return url
}

func (b *JiraUrlBuilder) BuildGetFieldUrl(v string) string {
	url := handleTrailingSlash(b.Base)
	if v == "" {
		v = DEFAULT_API_VERSION
	}
	url = url + FIELDS_URL
	url = strings.Replace(url, VERSION_HOLDER, v, -1)
	return url
}

func handleTrailingSlash(url string) string {
	url = strings.TrimSuffix(url, "/")
	return url
}
