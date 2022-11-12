package vsphere

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type VSphereClient struct {
	BaseUrl    string
	Token      string
	httpClient *http.Client
	Debug      bool
	user       string
	password   string
}

func NewVSphereClient(debug bool) *VSphereClient {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   time.Duration(30) * time.Second,
	}
	vSphereClient := &VSphereClient{Token: "", httpClient: client, Debug: debug}

	return vSphereClient
}

func (c *VSphereClient) SetCredential(user string, password string) {
	c.user = user
	c.password = password
}

func (c *VSphereClient) Login() error {
	target_url := c.BaseUrl + "/api/session"
	req, _ := http.NewRequest("POST", target_url, bytes.NewBuffer([]byte{}))
	auth := c.user + ":" + c.password
	authValue := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+authValue)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		return fmt.Errorf("authentication failed")
	}
	if val, ok := res.Header["Vmware-Api-Session-Id"]; ok {
		c.Token = val[0]
	}

	return nil
}

func (c *VSphereClient) Logout() error {
	target_url := c.BaseUrl + "/api/session"
	req, _ := http.NewRequest("DELETE", target_url, nil)
	req.Header.Set("Vmware-Api-Session-Id", c.Token)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return fmt.Errorf("logout failed")
	}

	return nil
}

type Response struct {
	*http.Response
	Body  string
	Error error
}

func (r *Response) Print() {
	if r.Error != nil {
		fmt.Fprintln(os.Stderr, r.Error.Error())
		return
	}
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

func (c *VSphereClient) validatePath(path string) error {
	match := false
	for _, v := range []string{"/api/"} {
		if strings.HasPrefix(path, v) {
			match = true
		}
	}
	if match == false {
		return fmt.Errorf("path must start with \"/api/\"")
	}
	return nil
}

func (c *VSphereClient) Request(method string, path string, query_param map[string]string, req_data []byte) *Response {
	err := c.validatePath(path)
	if err != nil {
		return &Response{nil, "", err}
	}
	req, _ := http.NewRequest(method, c.BaseUrl+path, bytes.NewBuffer(req_data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Vmware-Api-Session-Id", c.Token)
	if len(query_param) > 0 {
		params := req.URL.Query()
		for k, v := range query_param {
			params.Add(k, v)
		}
		req.URL.RawQuery = params.Encode()
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return &Response{nil, "", err}
	}
	defer res.Body.Close()
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return &Response{res, "", err}
	}

	if len(res_body) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, res_body, "", "    "); err != nil {
			log.Println(err)
			return &Response{res, "", err}
		}
		return &Response{res, prettyJSON.String(), nil}
	} else {
		return &Response{res, "", nil}
	}
}
