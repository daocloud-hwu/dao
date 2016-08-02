package dao

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
)

type Client struct {
    Host string
    AuthToken string
    InternalHost string
    InternalToken string
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
    type Info struct {
        Username string `json:"email_or_mobile"`
        Password string `json:"password"`
    }

    info := new(Info)
    info.Username = user
    info.Password = passwd

    inbody, err := json.Marshal(info)
    if err != nil {
        return err
    }

    status, outbody, _, err := c.do("POST", "/internal/access-token", nil, inbody, true)
    if err != nil {
        return err
    }
    if status/100 != 2 {
        return fmt.Errorf("Status code is %d, reason %s", status, outbody)
    }

    result := struct {
        Token string `json:"access_token"`
    } {}
    if err := json.Unmarshal(outbody, &result); err != nil {
        return err
    }

    c.AuthToken = result.Token
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

func (c *Client) do(method, url string, header map[string]string, body []byte, internal bool) (int, []byte, map[string]string, error) {
    var reader io.Reader = nil
    if body != nil {
        reader = bytes.NewBuffer(body)
    }

    client := &http.Client{}
    client.Transport = &http.Transport{DisableKeepAlives: true}

    link := fmt.Sprintf("http://%s%s", c.Host, url)
    if internal {
        link = fmt.Sprintf("http://%s%s", c.InternalHost, url)
    }

    req, err := http.NewRequest(method, link, reader)
    if err != nil {
        return 0, nil, nil, err
    }

    for k, v := range header {
        req.Header.Set(k, v)
    }

    req.Header.Set("Authorization", c.AuthToken)
    req.Header.Set("Content-Type", "application/json")
    if internal {
        req.Header.Set("X-DAO-INTERNAL-TOKEN", c.InternalToken)
    } else {
        req.Header.Set("Authorization", c.AuthToken)
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
