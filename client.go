package dao

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "strings"
)

type Client struct {
    Host string
    InternalToken string
    AuthToken string
}

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

type Runtime struct {
    ID   string `json:"runtime_id"`
    Type string `json:"app_runtime_type"`
    Name string `json:"app_runtime_name"`
}

func (c *Client) CreateUser(user, passwd string) error {
    // TBD
    return nil
}

func (c *Client) DeleteUser(user string) error {
    // TBD
    return nil
}

func (c *Client) Login(user, passwd string) error {
    // TBD
    return nil
}

func (c *Client) Logout() {
    c.AuthToken = ""
}

func (c *Client) ListRuntime() ([]*Runtime, error) {
    status, body, _, err := c.do("GET", "/v1/runtimes", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Runtimes []*Runtime `json:"runtimes"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Runtimes, nil
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

func (c *Client) do(method, url string, header map[string]string, body []byte, internal bool) (int, []byte, map[string]string, error) {
    var reader io.Reader = nil
    if body != nil {
        reader = bytes.NewBuffer(body)
    }

    client := &http.Client{}
    client.Transport = &http.Transport{DisableKeepAlives: true}
    req, err := http.NewRequest(
        method,
        fmt.Sprintf("http://%s%s", c.Host, url),
        reader,
    )
    if err != nil {
        return 0, nil, nil, err
    }

    for k, v := range header {
        req.Header.Set(k, v)
    }

    req.Header.Set("Authorization", c.AuthToken)
    if internal {
        // TBD
    }

    res, err := client.Do(req)
    if err != nil {
        return 0, nil, nil, err
    }
    defer res.Body.Close()

    outbody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return 0, nil, nil, err
    }

    var resHeader map[string]string = nil
    for k, _ := range res.Header {
        if resHeader == nil {
            resHeader = make(map[string]string)
        }
        resHeader[k] = res.Header.Get(k)
    }

    return res.StatusCode, outbody, resHeader, nil
}
