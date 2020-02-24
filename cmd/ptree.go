package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getProjectTree(pId string, config EJConfig) (*Project, error) {
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
			return nil, err
		}
		issues, imap := getAllIssueIdMap(il.Issues)
		//fmt.Printf("Issue Id: %v\n", issues)
		project, err := buildProjectTreeStruct(projectId, issues, imap, config)

		if err != nil {
			return nil, err
		}
		return project, nil

	} else if code == http.StatusUnauthorized {
		return nil, errors.New("User is not authenticated kindly check creds.")
	} else if code == http.StatusForbidden {
		return nil, errors.New("User is not authorized to get the issue list")
	}
	errMsg := fmt.Sprintf("Error with code : %d has occurred.\n", code)
	return nil, errors.New(errMsg)
}

func InitProject() *Project {
	s := []*DSprint{}
	p := Project{Id: projectId, Sprints: s}
	return &p
}

func buildProjectTreeStruct(projectId string, issueIds []string, issueMap map[string]JIssue, config EJConfig) (*Project, error) {
	b := JiraUrlBuilder{Base: config.Url}
	p := InitProject()
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
			//fmt.Printf("Issue - Sprint id %s Field %v\n", is.Id, is.IField)
			if err != nil {
				return nil, err
			}
			enrichProjectTree(p, is, issueMap)
			//fmt.Printf("Project: %v\n", p)
		} else if code == http.StatusUnauthorized {
			return nil, errors.New("User is not authenticated kindly check creds.")
		} else if code == http.StatusForbidden {
			return nil, errors.New("User is not authorized to get the issue list")
		}
	}
	return p, nil
}

func enrichProjectTree(p *Project, is JIssue, issueMap map[string]JIssue) {
	if is.IField != nil {
		parent := is.IField.Parent
		dparent := &DIssue{}
		if is.IField.Sprint != nil {
			ds := p.AddDisplaySprints(is.IField.Sprint)
			if parent != nil {
				pdi := ds.SearchIssue(parent.Id, "")
				if pdi == nil {
					k := getIssueKey(*parent)
					if i, ok := issueMap[k]; ok {
						dchild := &DIssue{}
						dchild.populateFields(&is)
						dparent.populateFields(&i)
						dparent.AddIssue(dchild)
					}
					ds.AddIssue(dparent)
				} else {
					dchild := &DIssue{}
					dchild.populateFields(&is)
					cdi := pdi.SearchChildIssues(dchild.Id, "")
					if cdi == nil {
						pdi.AddIssue(dchild)
					}
				}
			} else {
				dchild := &DIssue{}
				dchild.populateFields(&is)
				cdi := ds.SearchIssue(dchild.Id, "")
				if cdi == nil {
					ds.AddIssue(dchild)
				}
			}
		}

		// handle all the closed sprints
		closedSprints := is.IField.ClosedSprints
		if closedSprints != nil {
			for _, csp := range closedSprints {
				cds := p.AddDisplaySprints(csp)
				if parent != nil {
					dp := cds.SearchIssue(parent.Id, "")
					if dp != nil {
						dchild := dp.SearchChildIssues(is.Id, "")
						if dchild == nil {
							dchild = &DIssue{}
							dchild.populateFields(&is)
							dp.AddIssue(dchild)
						} else {
							// do nothing as child is already present
						}
					} else {
						dp = &DIssue{}
						dp.populateFields(parent)

						dchild := &DIssue{}
						dchild.populateFields(&is)

						dp.AddIssue(dchild)
						cds.AddIssue(dp)
					}
				} else {
					dchild := &DIssue{}
					dchild.populateFields(&is)
					cdi := cds.SearchIssue(dchild.Id, "")
					if cdi == nil {
						cds.AddIssue(dchild)
					}
				}
			}
		}
	} else {
		return
	}
	return
}

func getAllIssueIdMap(issues []JIssue) ([]string, map[string]JIssue) {
	ids := []string{}
	issueMap := map[string]JIssue{}
	for _, issue := range issues {
		ids = append(ids, issue.Id)
		key := getIssueKey(issue)
		issueMap[key] = issue
	}
	//fmt.Printf("map: %v\n", issueMap)
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
	Id     string   `json:"id"`
	Key    string   `json:"key"`
	IField *JIField `json:"fields,omitempty"`
}

type JIField struct {
	IssueType     *JIssueType `json:"issuetype,omitempty"`
	Status        *JStatus    `json:"status,omitempty"`
	Sprint        *JSprint    `json:"sprint,omitempty"`
	ClosedSprints []*JSprint  `json:"closedSprints,omitempty"`
	Parent        *JIssue     `json:"parent,omitempty"`
	Priority      *JPriority  `json:"priority,omitempty"`
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
	Id            int    `json:"id"`
	State         string `json:"state"`
	Name          string `json:"name"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	CompletedDate string `json:"completedDate,omitempty"`
	OriginBoardId int    `json:"originBoardId,omitempty"`
}

type Project struct {
	Id      string     `json:"projectId"`
	Name    string     `json:"name"`
	Sprints []*DSprint `json:"sprint"`
}

func (p *Project) AddDisplaySprints(sprint *JSprint) *DSprint {
	if sprint != nil {
		dsprint, _ := p.SearchSprint(strconv.Itoa(sprint.Id))
		// display spring is already present because search will return an empty sprint.
		// if it does not find a sprint.
		if dsprint != nil {
			dsprint.polulateFields(sprint)
		}
		return dsprint
	}
	return nil
}

func (p *Project) SearchSprint(id string) (*DSprint, bool) {
	if p.Sprints == nil {
		p.Sprints = []*DSprint{}
		ds := &DSprint{Id: id}
		p.Sprints = append(p.Sprints, ds)
		return ds, false
	}
	for _, s := range p.Sprints {
		if s.Id == id {
			//fmt.Printf("Found Sprint %s in Project\n", s.Id)
			return s, true
		}
	}
	// if sprint is not found create a new one add it to the collection and return ref
	ds := &DSprint{Id: id}
	p.Sprints = append(p.Sprints, ds)
	return ds, false
}

func (p *Project) AddSprint(sprint *DSprint) {
	if p.Sprints == nil {
		p.Sprints = []*DSprint{}
	}
	p.Sprints = append(p.Sprints, sprint)
}

func (p *Project) SearchIssue(id string, issueType string) *DIssue {
	if p.Sprints == nil || len(p.Sprints) == 0 {
		return nil
	}
	for _, s := range p.Sprints {
		i := s.SearchIssue(id, issueType)
		if i != nil {
			return i
		}
	}
	return nil
}

type DSprint struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	SprintState string    `json:"state"`
	StartDate   string    `json:"startDate"`
	EndDate     string    `json:"endDate"`
	Issues      []*DIssue `json:"issue"`
}

func (s *DSprint) SearchIssue(id string, issueType string) *DIssue {
	for _, i := range s.Issues {
		if id == i.Id {
			return i
		} else {
			cissue := i.SearchChildIssues(id, issueType)
			if cissue != nil {
				return cissue
			} else {
				// continue
			}
		}
	}
	return nil
}

func (d *DSprint) polulateFields(s *JSprint) {
	if s == nil {
		return
	}
	d.Id = strconv.Itoa(s.Id)
	d.Name = s.Name
	d.SprintState = s.State
	if s.StartDate != "" {
		d.StartDate = s.StartDate
	}
	if s.EndDate != "" {
		d.EndDate = s.EndDate
	}
	if d.Issues == nil {
		d.Issues = []*DIssue{}
	}
	return
}

func (s *DSprint) AddIssue(issue *DIssue) {
	if s.Issues == nil {
		s.Issues = []*DIssue{}
	}
	s.Issues = append(s.Issues, issue)
}

type DIssue struct {
	Id          string     `json:"id"`
	Key         string     `json:"key"`
	IssueType   JIssueType `json:"issuetype"`
	ChildIssues []*DIssue  `json:"child-issues"`
}

func (d *DIssue) populateFields(issue *JIssue) {
	d.Id = issue.Id
	d.Key = issue.Key
	d.IssueType = JIssueType{}
	d.IssueType.Id = issue.IField.IssueType.Id
	d.IssueType.Name = issue.IField.IssueType.Name
	d.IssueType.Subtask = issue.IField.IssueType.Subtask
	d.IssueType.Self = issue.IField.IssueType.Self

	if d.ChildIssues == nil {
		d.ChildIssues = []*DIssue{}
	}
	return
}

func (i *DIssue) AddIssue(issue *DIssue) {
	if i.ChildIssues == nil {
		i.ChildIssues = []*DIssue{}
	}
	i.ChildIssues = append(i.ChildIssues, issue)
}

func (i *DIssue) SearchChildIssues(id, issueType string) *DIssue {
	if i.ChildIssues == nil || len(i.ChildIssues) == 0 {
		return nil
	}
	for _, c := range i.ChildIssues {
		if c.Id == id {
			return c
		} else {
			cissue := c.SearchChildIssues(id, issueType)
			if cissue != nil {
				return cissue
			} else {
				// conitue
			}
		}
	}
	return nil
}
