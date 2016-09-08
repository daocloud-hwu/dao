package dao

import (
	"encoding/json"
	"fmt"
)

type Project struct {
	ID        string `json:"buildflow_id"`
	Name      string `json:"name"`
	PackageID string `package_id`
}

type Buildflow struct {
	Project     *Project `json:"project"`
	LatestBuild *Build   `json:"last_build"`
}

type Build struct {
	ID            int    `json:"id"`
	Status        string `json:"status"`
	Sha           string `json:"sha"`
	Branch        string `json:"branch"`
	Tag           string `json:"tag"`
	TriggerMethod string `json:"trigger_method"`
	CreatedAt     int64  `json:"created_at"`
}

type CiBuild struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Sha       string `json:"sha"`
	Branch    string `json:"branch"`
	CreatedAt int64  `json:"created_at"`
	Message   string `json:"message"`
}

func (c *Client) ListBuildflow() ([]*Buildflow, error) {
	status, body, _, err := c.do("GET", "/v1/ship/projects?size=-1&offset=0", nil, nil, false)
	if err != nil {
		return nil, err
	}
	if status/100 != 2 {
		return nil, fmt.Errorf("Status code is %d", status)
	}

	result := struct {
		Projects   []*Buildflow `json:"projects"`
		TotalCount int          `json:"total_count"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Projects, nil
}

func (c *Client) GetBuildflowByName(name string) (*Buildflow, error) {
	bs, err := c.ListBuildflow()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		if b.Project.Name == name {
			return b, nil
		}
	}

	return nil, nil
}

func (c *Client) GetBuildflow(id string) (*Buildflow, error) {
	url := fmt.Sprintf("/v1/ship/project/%s", id)
	status, body, _, err := c.do("GET", url, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if status/100 != 2 {
		return nil, fmt.Errorf("Status code is %d", status)
	}

	result := new(Buildflow)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) ListCiBuild(buildflowID string) ([]*CiBuild, error) {
	url := fmt.Sprintf("/v1/ship/project/%s/ci_build?size=-1&offset=0", buildflowID)
	status, body, _, err := c.do("GET", url, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if status/100 != 2 {
		return nil, fmt.Errorf("Status code is %d", status)
	}

	result := struct {
		CiBuilds []*CiBuild `json:"builds"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.CiBuilds, nil
}

func (c *Client) GetCiBuildByMessage(buildflowID, message string) (*CiBuild, error) {
	cibuilds, err := c.ListCiBuild(buildflowID)
	if err != nil {
		return nil, err
	}

	for _, b := range cibuilds {
		if b.Message == message {
			return b, nil
		}
	}

	return nil, nil
}

func (c *Client) ListBuild(buildflowID string) ([]*Build, error) {
	url := fmt.Sprintf("/v1/ship/project/%s/image_build?size=-1&offset=0", buildflowID)
	status, body, _, err := c.do("GET", url, nil, nil, false)
	if err != nil {
		return nil, err
	}
	if status/100 != 2 {
		return nil, fmt.Errorf("Status code is %d", status)
	}

	result := struct {
		Builds []*Build `json:"builds"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Builds, nil
}

func (c *Client) GetBuild(buildflowID string, id int) (*Build, error) {
	builds, err := c.ListBuild(buildflowID)
	if err != nil {
		return nil, err
	}

	for _, b := range builds {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, nil
}

func (c *Client) GetBuildByTag(buildflowID, tag string) (*Build, error) {
	builds, err := c.ListBuild(buildflowID)
	if err != nil {
		return nil, err
	}

	for _, b := range builds {
		if b.Tag == tag {
			return b, nil
		}
	}

	return nil, nil
}

func (c *Client) PostManualBuild(buildflowID, branch string) (int, error) {
	type BuildInfo struct {
		Branch string `json:"branch"`
	}

	bi := &BuildInfo{Branch: branch}
	url := fmt.Sprintf("/v1/ship/project/%s/image_build", buildflowID)
	inbody, err := json.Marshal(bi)
	if err != nil {
		return 0, err
	}

	status, outbody, _, err := c.do("POST", url, nil, inbody, false)
	if err != nil {
		return 0, err
	}
	if status/100 != 2 {
		return 0, fmt.Errorf("Status code is %d, reason %s", status, outbody)
	}

	result := new(Build)
	if err := json.Unmarshal(outbody, result); err != nil {
		return 0, err
	}

	return result.ID, nil
}
