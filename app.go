package dao

import (
    "encoding/json"
    "fmt"
)

type App struct {
    ID      string   `json:"app_id"`
    Name    string   `json:"name"`
    Runtime *Runtime `json:"runtime"`
}

func (c *Client) CreateSrApp(appName, pid, release, nodeName string, ports map[int]int) error {
    type Port struct {
        ContainerPort int `json:"container_port"`
        HostPort int `json:"host_port"`
        Protocol string `json:"protocol"`
        Published bool `json:"published"`
    }

    type Metadata struct {
        Command string `json:"command"`
        ContainerVolumes []string `json:"container_volumes"`
        Tags []map[string]string `json:"tags"`
        ContainerPorts []*Port `json:"container_ports"`
        ContainerRestart string `json:"container_restart"`
        ContainerPrivileged bool `json:"container_privileged"`
    }

    type SrApp struct {
        Name string `json:"name"`
        RuntimeID string `json:"runtime_id"`
        PackageID string `json:"package_id"`
        ReleaseName string `json:"release_name"`
        Instances int `json:"instances"`
        EnvVar map[string]string `json:"env_vars"`
        Metadata *Metadata `json:"metadata"`
    }

    m := new(Metadata)
    m.Command = ""
    m.ContainerVolumes = make([]string, 0)
    m.Tags = []map[string]string{{"name": nodeName}}
    m.ContainerPorts = make([]*Port, 0)
    for k, v := range ports {
        m.ContainerPorts = append(m.ContainerPorts, &Port{ContainerPort: k, HostPort: v, Protocol: "tcp", Published: true})
    }
    m.ContainerRestart = "always"
    m.ContainerPrivileged = false

    app := new(SrApp)
    app.Name = appName
    app.RuntimeID = "srsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsr"
    app.PackageID = pid
    app.ReleaseName = release
    app.Instances = 1
    app.EnvVar = make(map[string]string)
    app.Metadata = m

    inbody, err := json.Marshal(app)
    if err != nil {
        return err
    }

    status, _, _, err := c.do("POST", "/v1/apps", nil, inbody, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) ListApp() ([]*App, error) {
    status, body, _, err := c.do("GET", "/v1/apps", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Apps []*App `json:"apps"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Apps, nil
}

func  (c *Client) GetAppState(id string) (string, error) {
    status, body, _, err := c.do("GET", fmt.Sprintf("/v1/apps/%s/state"), nil, nil, false)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        State string `json:"state"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }

    return result.State, nil    
}

func (c *Client) StartApp(id string) error {
    status, _, _, err := c.do("POST", fmt.Sprintf("/v1/apps/%s/actions/start"), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 { 
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) StopApp(id string) error {
    status, _, _, err := c.do("POST", fmt.Sprintf("/v1/apps/%s/actions/stop"), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) DeleteApp(id string) error {
    status, _, _, err := c.do("DELETE", fmt.Sprintf("/v1/apps/%s"), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}
