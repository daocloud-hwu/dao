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

func (c *Client) CreateCfApp(appName, pid, release, instanceType string, port int) (string, error) {
    type Port struct {
        ContainerPort int    `json:"container_port"`
        Protocol      string `json:"protocol"`
        PublishType   string `json:"publish_type"`
        External      string `json:"external"`
    }

    type Instance struct {
        ID    string `json:"service_instance_id"`
        Alias string `json:"service_alias"`
    }

    type Metadata struct {
        Command          string      `json:"command"`
        Volumes          []string    `json:"volumes"`
        InstanceType     string      `json:"instance_type"`
        ExposePorts      []*Port     `json:"expose_ports"`
        ServiceInstances []*Instance `json:"service_instances"`
    }

    type Options struct {
        StartAfterStage bool `json:"start_after_stage"`
    }

    type CfApp struct {
        Name         string            `json:"name"`
        RuntimeID    string            `json:"runtime_id"`
        PackageID    string            `json:"package_id"`
        ReleaseName  string            `json:"release_name"`
        Instances    int               `json:"instances"`
        EnvVar       map[string]string `json:"env_vars"`
        Metadata     *Metadata         `json:"metadata"`
        ExtraOptions *Options          `json:"extra_options"`
    }

    m := new(Metadata)
    m.Command = ""
    m.Volumes = make([]string, 0)
    m.InstanceType = instanceType
    m.ExposePorts = []*Port{&Port{ContainerPort: port, Protocol: "tcp", PublishType: "http", External: "external"}}
    m.ServiceInstances = make([]*Instance, 0)

    app := new(CfApp)
    app.Name = appName
    app.RuntimeID = "a849cdf2-c79e-4c29-83ca-50751cc388a5"
    app.PackageID = pid
    app.ReleaseName = release
    app.Instances = 1
    app.EnvVar = make(map[string]string)
    app.Metadata = m
    app.ExtraOptions = &Options{StartAfterStage: true}

    inbody, err := json.Marshal(app)
    if err != nil {
        return "", err
    }

    status, outbody, _, err := c.do("POST", "/v1/apps", nil, inbody, false)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("Status code is %d, reason %s", status, outbody)
    }

    result := struct {
        AppID string `json:"app_id"`
    } {}
    if err := json.Unmarshal(outbody, &result); err != nil {
        return "", err
    }

    return result.AppID, nil
}

func (c *Client) GetAppUrl(id string) (string, error) {
    status, body, _, err := c.do("GET", fmt.Sprintf("/v1/apps/%s/details", id), nil, nil, false)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("Status code is %d, reason %s", status, body)
    }

    result := struct {
        Url string `json:"url"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }

    return result.Url, nil
}

func (c *Client) CreateSrApp(appName, pid, release, nodeName string, ports map[int]int) (string, error) {
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
        return "", err
    }

    status, outbody, _, err := c.do("POST", "/v1/apps", nil, inbody, false)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("Status code is %d, reason %s", status, outbody)
    }

    result := struct {
        AppID string `json:"app_id"`
    } {}
    if err := json.Unmarshal(outbody, &result); err != nil {
        return "", err
    }

    return result.AppID, nil
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

func (c *Client) GetApp(id string) (*App, error) {
    apps, err := c.ListApp()
    if err != nil {
        return nil, err
    }

    for _, app := range apps {
        if app.ID == id {
            return app, nil
        }
    }

    return nil, nil
}

func  (c *Client) GetAppState(id string) (string, error) {
    status, body, _, err := c.do("GET", fmt.Sprintf("/v1/apps/%s/state", id), nil, nil, false)
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
    status, _, _, err := c.do("POST", fmt.Sprintf("/v1/apps/%s/actions/start", id), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 { 
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) StopApp(id string) error {
    status, _, _, err := c.do("POST", fmt.Sprintf("/v1/apps/%s/actions/stop", id), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) DeleteApp(id string) error {
    status, _, _, err := c.do("DELETE", fmt.Sprintf("/v1/apps/%s", id), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) RestageApp(id string, packageId string, releaseName string, startAfterStage bool) (string, error) {
    type Option struct {
        StartAfterStage bool `json:"start_after_stage"`
    }

    type RequestMetadata struct {
        PackageId string `json:"package_id"`
        ReleaseName string `json:"release_name"`
        ExtraOption *Option `json:"extra_options"`
    }

    rm := new(RequestMetadata)
    rm.PackageId = packageId
    rm.ReleaseName = releaseName
    o := new(Option)
    o.StartAfterStage = true
    rm.ExtraOption = o

    inbody, err := json.Marshal(rm)
    if err != nil {
        return "", err
    }

    status, outbody, _, err := c.do("POST", fmt.Sprintf("/v1/apps/%s/actions/restage", id), nil, inbody, false)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("Status code is %d, reason %s", status, outbody)
    }

    result := struct {
        AppID string `json:"app_id"`
    } {}
    if err := json.Unmarshal(outbody, &result); err != nil {
        return "", err
    }
    return result.AppID, nil
}


func (c *Client) UpdateAppYml(id string, yml string) error {
    type Options struct {
        Yml string `json:"compose_yml"`
        Op string `json:"operation"`
    }

    type AppT struct {
        ExtraOptions *Options `json:"extra_options"`
    }
    options := &Options{Yml: yml, Op: "update_compose_yml"}
    s := AppT{ExtraOptions: options}

    inbody, err := json.Marshal(s)
    if err != nil {
        return err
    }

    status, outbody, _, err := c.do("PATCH", fmt.Sprintf("/v1/apps/%s", id), nil, inbody, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d, reason %s", status, outbody)
    }

    return nil
}


func (c *Client) BindServiceInstance(id string, serviceInstanceId string, serviceAlias string) error {
    type Instance struct {
        ID    string `json:"service_instance_id"`
        Alias string `json:"service_alias"`
    }

    type Metadata struct {
        ServiceInstances []*Instance `json:"service_instances"`
    }

    type Options struct {
        Operation string `json:"operation"`
    }

    type RequestMetadata struct {
        Metadata *Metadata `json:"metadata"`
        ExtraOption *Options `json:"extra_options"`

    }
    i := new(Instance)
    i.ID = serviceInstanceId
    i.Alias = serviceAlias

    m := new(Metadata)
    m.ServiceInstances = append(make([]*Instance, 0), i)

    o := new(Options)
    o.Operation = "service_instances"

    rm := new(RequestMetadata)
    rm.Metadata = m
    rm.ExtraOption = o

    inbody, err := json.Marshal(rm)
    if err != nil {
        return err
    }
    //fmt.Printf("request data:\t%s\n", inbody)
    status, outbody, _,  err := c.do("PATCH", fmt.Sprintf("/v1/apps/%s", id), nil, inbody, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d, reason %s\n", status, outbody)
    }

    return nil
}
