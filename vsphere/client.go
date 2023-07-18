package vsphere

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VSphereClient struct {
	BaseUrl    string
	Token      string
	Debug      bool
	httpClient HttpClient
	user       string
	password   string
}

func (c *VSphereClient) do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func NewVSphereClient(debug bool) *VSphereClient {
	return &VSphereClient{Debug: debug}
}

func (c *VSphereClient) SetHttpClient(hc HttpClient) {
	c.httpClient = hc
}

func (c *VSphereClient) SetVCenterServer(server string) error {
	c.BaseUrl = server
	return nil
}

func (c *VSphereClient) SetCredential(user string, password string) {
	c.user = user
	c.password = password
}

func (c *VSphereClient) SetVCenterToken(token string) {
	c.Token = token
}

func (c *VSphereClient) Login() error {
	url := c.BaseUrl + "/api/session"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	cred := c.user + ":" + c.password
	authToken := base64.StdEncoding.EncodeToString([]byte(cred))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authToken))
	res, err := c.do(req)
	if err != nil {
		// http error
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		// authentication error
		return fmt.Errorf("authentication failed")
	}
	if val, ok := res.Header["Vmware-Api-Session-Id"]; ok {
		c.Token = val[0]
	} else {
		// http response header error
		return fmt.Errorf("session id not found")
	}

	return nil
}

func (c *VSphereClient) Logout() error {
	url := c.BaseUrl + "/api/session"
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Vmware-Api-Session-Id", c.Token)
	res, err := c.do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return fmt.Errorf("logout failed")
	}

	return nil
}

func (c *VSphereClient) validatePath(path string) error {
	match := false
	for _, v := range []string{"/api/"} {
		if strings.HasPrefix(path, v) {
			match = true
		}
	}
	if !match {
		return fmt.Errorf("request path must start with \"/api/\"")
	}
	return nil
}

func (c *VSphereClient) Call(method string, path string, query map[string]string, body []byte) (*Response, error) {
	err := c.validatePath(path)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Vmware-Api-Session-Id", c.Token)
	if len(query) > 0 {
		params := req.URL.Query()
		for k, v := range query {
			params.Add(k, v)
		}
		req.URL.RawQuery = params.Encode()
	}
	res, err := c.do(req)
	if err != nil {
		log.Println(err)
		return &Response{res, ""}, err
	}
	defer res.Body.Close()
	var d any
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return &Response{res, ""}, err
	}
	b, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		log.Println(err)
		return &Response{res, ""}, err
	}
	return &Response{res, string(b)}, nil
}

type Response struct {
	*http.Response
	Body string
}

func (r *Response) Print() {
	if r.Body == "" {
		var msg string
		switch r.StatusCode {
		case 404:
			msg = "request error"
		case 500:
			msg = "server error"
		case 200, 201:
			msg = "no response body"
		}
		fmt.Printf("{\"code\": %d, \"body\": \"%v\"}\n", r.StatusCode, msg)
	} else {
		fmt.Println(r.Body)
	}
}
