package dao

import (
    "encoding/json"
    "fmt"
)

type Stack struct {
    ID   string `json:"stack_id"`
    Name string `json:"name"`
    Apps []*App `json:"apps"`
}

func (c *Client) CreateStack(stackName, nodeName, yml string) error {
    type Options struct {
        Yml string `json:"compose_yml"`
    }

    type Metadata struct {
        Tags []map[string]string `json:"tags"`
    }

    type StackT struct {
        Name string `json:"name"`
        RuntimeID string `json:"runtime_id"`
        ExtraOptions *Options `json:"extra_options"`
        Metadata *Metadata `json:"metadata"`
    }

    m := new(Metadata)
    m.Tags = []map[string]string{{"name": nodeName}}

    o := new(Options)
    o.Yml = yml

    s := new(StackT)
    s.Name = stackName
    s.RuntimeID = "srsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsrsr"
    s.Metadata = m
    s.ExtraOptions = o

    inbody, err := json.Marshal(s)
    if err != nil {
        return err
    }

    status, _, _, err := c.do("POST", "/v1/stacks", nil, inbody, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) ListStack() ([]*Stack, error) {
    status, body, _, err := c.do("GET", "/v1/stacks", nil, nil, false)
    if err != nil {
        return nil, err
    }
    if status/100 != 2 {
        return nil, fmt.Errorf("Status code is %d", status)
    }

    result := struct {
        Stacks []*Stack `json:"stacks"`
    } {}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return result.Stacks, nil
}

func  (c *Client) GetStackState(id string) (string, error) {
    status, body, _, err := c.do("GET", fmt.Sprintf("/v1/stacks/%s/state"), nil, nil, false)
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

func (c *Client) StartStack(id string) error {
    status, _, _, err := c.do("POST", fmt.Sprintf("/v1/stacks/%s/actions/start"), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 { 
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) StopStack(id string) error {
    status, _, _, err := c.do("POST", fmt.Sprintf("/v1/stacks/%s/actions/stop"), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}

func (c *Client) DeleteStack(id string) error {
    status, _, _, err := c.do("DELETE", fmt.Sprintf("/v1/stacks/%s"), nil, nil, false)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d", status)
    }

    return nil
}
