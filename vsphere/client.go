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
	BaseUrl     string
	BasicAuth   bool
	Token       string
	httpClient  *http.Client
	Debug       bool
	Version     string
	FullVersion string
	user        string
	password    string
}

func NewVSphereClient(debug bool) *VSphereClient {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   time.Duration(30) * time.Second,
	}
	vSphereClient := &VSphereClient{BasicAuth: false, Token: "", httpClient: client, Debug: debug}

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

func (c *VSphereClient) ListSupervisorCluster() {

}

type Response struct {
	*http.Response
	Body  interface{}
	Error error
}

func (r *Response) BodyBytes() ([]byte, error) {
	return json.Marshal(r.Body)
}

func (r *Response) UnmarshalBody(strct interface{}) {
	bytes, _ := r.BodyBytes()
	json.Unmarshal(bytes, strct)
}

func (r *Response) Print(noPretty bool) {
	var body []byte
	if r.Error != nil {
		fmt.Fprintln(os.Stderr, r.Error.Error())
		return
	}
	if r.Body == nil {
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
		if noPretty {
			body, _ = r.BodyBytes()
		} else {
			body, _ = json.MarshalIndent(r.Body, "", "  ")
		}
		fmt.Println(string(body))
	}
}

func (c *VSphereClient) Request(method string, path string, query_param map[string]string, req_data []byte) *Response {
	err := func() error {
		var match bool
		match = false
		for _, v := range []string{"/api/"} {
			if strings.HasPrefix(path, v) {
				match = true
			}
		}
		if match == false {
			return fmt.Errorf("path must start with \"/api/\"")
		}
		return nil
	}()
	if err != nil {
		return &Response{nil, nil, err}
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
		fmt.Println(err)
		return &Response{}
	}
	defer res.Body.Close()
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return &Response{}
	}
	var data interface{}
	if len(res_body) > 0 {
		err = json.Unmarshal(res_body, &data)
		if err != nil {
			log.Println(err)
			return &Response{}
		}
		r := &Response{res, data, nil}
		return r
	} else {
		return &Response{res, nil, nil}
	}
}

func (c *VSphereClient) Exec(uri string) {
	res := c.Request("GET", uri, nil, nil)
	fmt.Println(res.Body)
}
