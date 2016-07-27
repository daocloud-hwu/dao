package dao

import (
    "encoding/json"
    "fmt"
    "strings"
)

type Service struct {
    ID       string `json:"service_id"`
    Name     string `json:"service_name"`
    Category string `json:"category_name"`
}

type EnvVar struct {
    Name string `json:"env_var_name"`
    Value string `json:"env_var_value"`
}

type ServiceInstance struct {
    ID        string    `json:"service_instance_id"`
    Name      string    `json:"service_instance_name"`
    ServiceID string    `json:"service_id"`
    EnvVars   []*EnvVar `json:"env_vars"`
    Type      string    `json:"instance_type"`
}

func (c *Client) ListService() ([]*Service, error) {
    status, body, _, err := c.do("GET", "/v1/services", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Services []*Service `json:"services"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Services, nil
}

func (c *Client) GetService(name string) (*Service, error) {
    ss, err := c.ListService()
    if err != nil {
        return nil, err
    }

    for _, s := range ss {
        if strings.ToLower(s.Name) == strings.ToLower(name) {
            return s, nil
        }
    }

    return nil, nil
}

func (c *Client) ListServiceInstance() ([]*ServiceInstance, error) {
    status, body, _, err := c.do("GET", "/v1/service-instances", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        ServiceInstances []*ServiceInstance `json:"service_services"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.ServiceInstances, nil
}

func (c *Client) CreateServiceInstance(serviceID, name, serviceType string) (string, error) {
    type Instance struct {
        ServiceID string `json:"service_id"`
        Name      string `json:"service_instance_name"`
        Type      string `json:"service_instance_type"`
    }

    ins := new(Instance)
    ins.ServiceID = serviceID
    ins.Name = name
    ins.Type = serviceType

    inbody, err := json.Marshal(ins)
    if err != nil {
        return "", err
    }

    status, outbody, _, err := c.do("POST", "/v1/service-instances", nil, inbody, false)
    if err != nil {
        return "", err
    }
    if status/100 != 2 {
        return "", fmt.Errorf("Status code is %d, reason %s", status, outbody)
    }

    result := struct {
        InstanceID string `json:"service_instance_id"`
    } {}
    if err := json.Unmarshal(outbody, &result); err != nil {
        return "", err
    }

    return result.InstanceID, nil
}

func (c *Client) GetServiceInstance(id string) (*ServiceInstance, error) {
    instances, err := c.ListServiceInstance()
    if err != nil {
        return nil, err
    }

    for _, ins := range instances {
        if ins.ID == id {
            return ins, nil
        }
    }

    return nil, nil
}

func (c *Client) DeleteServiceInstance(id string) error {
    status, _, _, err := c.do("DELETE", fmt.Sprintf("/v1/service-instances/%s", id), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}
