package cmd

/**
  this file will contain functions which will build the rest api calls
  there will be placeholder URL that will have to be changed.
*/

import (
	"fmt"
	"net/url"
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
	BOARD_URL  = REST_AGILE + VERSION_HOLDER + "/board"

	PROJECT_HOLDER     = "{project-id}"
	PAGE_SIZE          = "{page-size}"
	ISSUE_ID_HOLDER    = "{issue-id}"
	BOARD_ID_HOLDER    = "{board-id}"
	START_DATE_HOLDER  = "start-date"
	START_AT_HOLDER    = "start-at"
	MAX_RESULTS_HOLDER = "max-results"

	PROJECT_ISSUE_JQL_URL = REST_API + VERSION_HOLDER + "/search?jql=project=" + PROJECT_HOLDER + "&fields=id,key,issuetype&maxResults=" + PAGE_SIZE

	AGILE_ISSUE_SPRINT_URL = REST_AGILE + VERSION_HOLDER + "/issue/" + ISSUE_ID_HOLDER + "?fields=parent,issuetype,sprint,closedSprints"
	BOARD_ISSUE_JQL_QUERY  = REST_AGILE + VERSION_HOLDER + "/board/" + BOARD_ID_HOLDER + "/issue?"
	JQL_PARAM              = "(updated >= '" + START_DATE_HOLDER + "' or created >= '" + START_DATE_HOLDER + "') and issuetype not in (\"Epic\",\"Program Epic\", \"Portfolio Epic\", \"Sub-task\", \"Test\", \"Story Bug\") and team is not empty"
	FIELDS_PARAM           = "issuetype,status"
	START_AT_PARAM         = "startAt=" + START_AT_HOLDER
	MAX_RESULTS_PARAM      = "maxResults=" + MAX_RESULTS_HOLDER
)

type JiraUrlBuilder struct {
	Base string
}

//BuildAuthCheckUrl generated the URL to check the creds against.
func (b *JiraUrlBuilder) BuildAuthCheckUrl(v string) string {
	if v == "" {
		v = DEFAULT_AGILE_VERSION
	}
	URL := handleTrailingSlash(b.Base)
	URL = URL + AUTH_CHECK_URL

	URL = strings.Replace(URL, VERSION_HOLDER, v, -1)
	return URL
}

//BuildBoardIssueQuery builds a query to get issues for aboard.
func (b *JiraUrlBuilder) BuildBoardIssueQuery(v string, boardID string, maxResults, startAt int, startDate string) string {
	if v == "" {
		v = DEFAULT_AGILE_VERSION
	}
	URL := handleTrailingSlash(b.Base)
	URL = URL + BOARD_ISSUE_JQL_QUERY
	URL = strings.Replace(URL, VERSION_HOLDER, v, -1)
	URL = strings.Replace(URL, BOARD_ID_HOLDER, boardID, -1)
	URL = URL + buildBoardQueryParams(startDate, startAt, maxResults)
	return URL
}

func buildBoardQueryParams(startDate string, startAt, maxResults int) string {
	p := url.Values{}
	jqlParam := strings.Replace(JQL_PARAM, START_DATE_HOLDER, startDate, -1)
	p.Set("jql", jqlParam)
	s := fmt.Sprintf("%d", startAt)
	m := fmt.Sprintf("%d", maxResults)
	p.Set("fields", FIELDS_PARAM)
	p.Set("startAt", s)
	p.Set("maxResults", m)
	return p.Encode()
}

/*BuildBoardURL create a URL to get the board either scrum, kanban or all*/
func (b *JiraUrlBuilder) BuildBoardURL(v string, bType string, maxResults int, startAt int) string {
	if v == "" {
		v = DEFAULT_AGILE_VERSION
	}
	URL := handleTrailingSlash(b.Base)
	URL = URL + BOARD_URL
	URL = strings.Replace(URL, VERSION_HOLDER, v, -1)
	switch bType {
	case "kanban":
		URL = URL + "?type=kanban"
	case "scrum":
		URL = URL + "?type=scrum"
	}
	startIndex := fmt.Sprintf("&startAt=%d", startAt)
	maxResultsValue := fmt.Sprintf("&maxResults=%d", maxResults)
	URL = URL + startIndex + maxResultsValue
	return URL
}

func (b *JiraUrlBuilder) BuildProjectIssueUrl(v string, projectId string, page string) string {
	URL := handleTrailingSlash(b.Base)
	if v == "" {
		v = DEFAULT_API_VERSION
	}
	if page == "" || page == "0" {
		page = "100"
	}
	URL = URL + PROJECT_ISSUE_JQL_URL
	URL = strings.Replace(URL, VERSION_HOLDER, v, -1)
	URL = strings.Replace(URL, PROJECT_HOLDER, projectId, -1)
	URL = strings.Replace(URL, PAGE_SIZE, page, -1)
	return URL
}

func (b *JiraUrlBuilder) BuildIssueSprintResponseUrl(v string, issueId string) string {
	URL := handleTrailingSlash(b.Base)
	if v == "" {
		v = DEFAULT_AGILE_VERSION
	}
	URL = URL + AGILE_ISSUE_SPRINT_URL
	URL = strings.Replace(URL, VERSION_HOLDER, v, -1)
	URL = strings.Replace(URL, ISSUE_ID_HOLDER, issueId, -1)
	return URL
}

func (b *JiraUrlBuilder) BuildGetFieldUrl(v string) string {
	URL := handleTrailingSlash(b.Base)
	if v == "" {
		v = DEFAULT_API_VERSION
	}
	URL = URL + FIELDS_URL
	URL = strings.Replace(URL, VERSION_HOLDER, v, -1)
	return URL
}

func handleTrailingSlash(URL string) string {
	URL = strings.TrimSuffix(URL, "/")
	return URL
}
