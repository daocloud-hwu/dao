package dao

import (
    "encoding/json"
    "fmt"
)

type Buildflow struct {
    ID          string `json:"buildflow_id"`
    Name        string `json:"daohub_repo_name"`
    PackageID   string `json:"package_id"`
    LatestBuild *Build `json:"latest_build"`
}

type Build struct {
    ID            int    `json:"id"`
    Status        string `json:"status"`
    Tag           string `json:"tag"`
    TriggerMethod string `json:"trigger_method"`
    CreatedAt     int64  `json:"created_at"`
}

type CiBuild struct {
    ID        int    `json:"id"`
    Status    string `json:"status"`
    CreatedAt int64  `json:"created_at"`
    Message   string `json:"message"`
}

func (c *Client) ListBuildflow() ([]*Buildflow, error) {
    status, body, _, err := c.do("GET", "/v1/buildflows?limit=-1&offset=0", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Buildflows []*Buildflow `json:"buildflows"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Buildflows, nil
}

func (c *Client) GetBuildflowByName(name string) (*Buildflow, error) {
    bs, err := c.ListBuildflow()
    if err != nil {
        return nil, err
    }

    for _, b := range bs {
        if b.Name == name {
            return b, nil
        }
    }

    return nil, nil
}

func (c *Client) GetBuildflow(id string) (*Buildflow, error) {
    url := fmt.Sprintf("/v1/buildflows/%s", id)
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
    url := fmt.Sprintf("/v1/buildflows/%s/cibuilds?limit=-1&offset=0", buildflowID)
    status, body, _, err := c.do("GET", url, nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        CiBuilds []*CiBuild `json:"ci_build_history"`
    } {}
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
    url := fmt.Sprintf("/v1/buildflows/%s/builds?limit=-1&offset=0", buildflowID)
    status, body, _, err := c.do("GET", url, nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Builds []*Build `json:"build_history"`
    } {}
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
        Type string `json:"build_type"`
        Name string `json:"name"`
    }

    bi := &BuildInfo{Type: "branch", Name: branch}
    url := fmt.Sprintf("/v1/buildflows/%s/builds", buildflowID)
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
