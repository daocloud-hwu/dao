package dao

import (
    "encoding/json"
    "fmt"
    "strings"
)

type Cluster struct {
    ID        string  `json:"node_cluster_id"`
    IsDefault bool    `json:"is_default"`
    Name      string  `json:"node_cluster_name"`
    Token     string  `json:"suggest_token"`
    Nodes     []*Node `json:"nodes"`
}

type Node struct {
    ID           string   `json:"node_id"`
    Addrs        []string `json:"node_addrs"`
    Hostname     string   `json:"hostname"`
    DockerStatus string   `json:"status"`
    Name         string   `json:"node_name"`
    IsConnected  bool     `json:"is_connected"`
}

type SrEnv struct {
    DaoGetUrl    string `json:"dao_get_url"`
    DaoKeeperUrl string `json:"dao_keeper_url"`
}

func (c *Client) GetUninstallCmd(os string) (string, error) {
    switch strings.ToLower(os) {
    case "ubuntu", "debian":
        return "dpkg -r daomonit", nil
    case "centos", "fedora":
        return "rpm -e daomonit", nil
    default:
        return "", fmt.Errorf("os type %s not supported", os)
    }
}

func (c *Client) GetImportCmd() (string, error) {
    env, err := c.GetSrEnv()
    if err != nil {
        return "", err
    }

    cs, err := c.ListCluster()
    if err != nil {
        return "", err
    }

    for _, c := range cs {
        if c.IsDefault {
            return fmt.Sprintf(
                "curl -sSL %s/daomonit/install.sh | sh -s %s %s",
                env.DaoGetUrl,
                c.Token,
                env.DaoKeeperUrl,
            ), nil
        }
    }

    return "", fmt.Errorf("no available cluster")
}

func (c *Client) GetSrEnv() (*SrEnv, error) {
    status, body, _, err := c.do("GET", "/v1/single_runtime/env", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := new(SrEnv)
    if err := json.Unmarshal(body, result); err != nil {
        return nil, err
    }

    return result, nil
}

func (c *Client) DeleteNode(nodeID string) error {
    status, _, _, err := c.do("DELETE", "/v1/single_runtime/nodes/" + nodeID, nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) ListCluster() ([]*Cluster, error) {
    status, body, _, err := c.do("GET", "/v1/clusters", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Clusters []*Cluster `json:"clusters"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Clusters, nil
}
