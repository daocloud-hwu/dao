package dao

import (
    "encoding/json"
    "fmt"
)

type Package struct {
    ID             string   `json:"package_id"`
    Name           string   `json:"package_name"`
    Namespace      string   `json:"docker_repo_namespace"`
    FullName       string   `json:"package_source_full_name"`
    IsPublic       bool     `json:"is_public"`
    UserID         string   `json:"user_id"`
    TenantID       string   `json:"tenant_id"`
    ReleaseAccount int      `json:"releases_count"`
    LatestRelease  *Release `json:"latest_release"`
}

type Release struct {
    Name string `json:"release_name"`
}

type PortInfo struct {
    Protocol string `json:"protocol"`
    Port     int    `json:"port"`
}

func (c *Client) ListPackage(ptype string) ([]*Package, error) {
    url := "/v1/packages?limit=-1"

    switch ptype {
    case "daocloud":
        url = "/v1/packages?limit=-1&is_public=true"
    }

    status, body, _, err := c.do("GET", url, nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Packages []*Package `json:"packages"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Packages, nil
}

func (c *Client) ListPackageRelease(packageID string) ([]*Release, error) {
    status, body, _, err := c.do("GET", fmt.Sprintf("/v1/packages/%s/releases", packageID), nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Releases []*Release `json:"releases"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Releases, nil
}

func (c *Client) GetPortInfo(packageID, release string) ([]*PortInfo, error) {
    status, body, _, err := c.do("GET", fmt.Sprintf("/v1/packages/%s/tags/%s/ports_info", packageID, release), nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }   
    
    result := struct {
        Ports []*PortInfo `json:"expose_ports"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }   

    return result.Ports, nil
}
