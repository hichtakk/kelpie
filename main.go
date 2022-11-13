package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hichtakk/kelpie/vsphere"
)

var (
	configfile string
	useSite    string
	debug      bool
	query      []string
	vsc        *vsphere.VSphereClient
)

func newCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kelpie",
		Short: "simple vSphere REST API client",
		Long:  "simple vSphere REST API client",
	}
	rootCmd.AddCommand(
		NewCmdHttpGet(&query),
		NewCmdHttpPost(&query),
		NewCmdHttpPatch(&query),
		NewCmdHttpPut(&query),
		NewCmdHttpDelete(&query),
	)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug mode")
	rootCmd.PersistentFlags().StringSliceVarP(&query, "query", "q", []string{}, "")

	return rootCmd
}

func NewCmdHttpGet(query *[]string) *cobra.Command {
	httpGetCmd := &cobra.Command{
		Use:   "get ${API-PATH}",
		Short: "call api with HTTP GET method",
		Long:  "example) kelpie get /api/vcenter/vm",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			params := map[string]string{}
			for _, q := range *query {
				qSlice := strings.Split(q, "=")
				if len(qSlice) != 2 {
					panic("invalid query parameter. it should be formatted as '<name>=<value>'.")
				}
				params[qSlice[0]] = qSlice[1]
			}
			var resp *vsphere.Response
			vsc.Login()
			resp = vsc.Request("GET", args[0], params, []byte{})
			if resp.Error != nil {
				fmt.Println(resp.Error)
			} else {
				resp.Print()
			}
			vsc.Logout()
		},
	}

	return httpGetCmd
}

func NewCmdHttpPost(query *[]string) *cobra.Command {
	fileName := ""
	data := []byte{}
	httpPostCmd := &cobra.Command{
		Use:   "post ${API-PATH}",
		Short: "call api with HTTP POST method",
		Long:  "example) kelpie post /api/vcenter/vm/{vm}/power -q action=reset",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rawData := []byte{}
			var err error
			if fileName != "" {
				rawData, err = readRequestData(fileName)
				if err != nil {
					return err
				}
			} else {
				rawData = nil
			}
			jsonObj := json.RawMessage(rawData)
			data, err = json.Marshal(jsonObj)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			params := map[string]string{}
			for _, q := range *query {
				qSlice := strings.Split(q, "=")
				if len(qSlice) != 2 {
					panic("invalid query parameter. it should be formatted as '<name>=<value>'.")
				}
				params[qSlice[0]] = qSlice[1]
			}
			var resp *vsphere.Response
			vsc.Login()
			resp = vsc.Request("POST", args[0], params, data)
			if resp.Error != nil {
				fmt.Println(resp.Error)
			} else {
				resp.Print()
			}
			vsc.Logout()
		},
	}
	httpPostCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for request data(json)")

	return httpPostCmd
}

func NewCmdHttpPatch(query *[]string) *cobra.Command {
	fileName := ""
	data := []byte{}
	httpPatchCmd := &cobra.Command{
		Use:   "patch ${API-PATH}",
		Short: "call api with HTTP PATCH method",
		Long:  "example) kelpie patch /api/vcenter/resource-pool/{resource_pool} -f ./data.json",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rawData := []byte{}
			var err error
			if fileName != "" {
				rawData, err = readRequestData(fileName)
				if err != nil {
					return err
				}
			} else {
				rawData = nil
			}
			jsonObj := json.RawMessage(rawData)
			data, err = json.Marshal(jsonObj)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			params := map[string]string{}
			for _, q := range *query {
				qSlice := strings.Split(q, "=")
				if len(qSlice) != 2 {
					panic("invalid query parameter. it should be formatted as '<name>=<value>'.")
				}
				params[qSlice[0]] = qSlice[1]
			}
			var resp *vsphere.Response
			vsc.Login()
			resp = vsc.Request("PATCH", args[0], params, data)
			if resp.Error != nil {
				fmt.Println(resp.Error)
			} else {
				resp.Print()
			}
			vsc.Logout()
		},
	}
	httpPatchCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for request data(json)")

	return httpPatchCmd
}

func NewCmdHttpPut(query *[]string) *cobra.Command {
	fileName := ""
	data := []byte{}
	httpPutCmd := &cobra.Command{
		Use:   "put ${API-PATH}",
		Short: "call api with HTTP PUT method",
		Long:  "example) kelpie put /api/vcenter/vm/{vm}/guest/customization -f ./data.json",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rawData := []byte{}
			var err error
			if fileName != "" {
				rawData, err = readRequestData(fileName)
				if err != nil {
					return err
				}
			} else {
				rawData = nil
			}
			jsonObj := json.RawMessage(rawData)
			data, err = json.Marshal(jsonObj)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			params := map[string]string{}
			for _, q := range *query {
				qSlice := strings.Split(q, "=")
				if len(qSlice) != 2 {
					panic("invalid query parameter. it should be formatted as '<name>=<value>'.")
				}
				params[qSlice[0]] = qSlice[1]
			}
			var resp *vsphere.Response
			vsc.Login()
			resp = vsc.Request("PUT", args[0], params, data)
			if resp.Error != nil {
				fmt.Println(resp.Error)
			} else {
				resp.Print()
			}
			vsc.Logout()
		},
	}
	httpPutCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for request data(json)")

	return httpPutCmd
}

func NewCmdHttpDelete(query *[]string) *cobra.Command {
	httpDeleteCmd := &cobra.Command{
		Use:   "delete ${API-PATH}",
		Short: "call api with HTTP DELETE method",
		Long:  "example) kelpie delete /api/vcenter/vm/${vm}",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			params := map[string]string{}
			for _, q := range *query {
				qSlice := strings.Split(q, "=")
				if len(qSlice) != 2 {
					panic("invalid query parameter. it should be formatted as '<name>=<value>'.")
				}
				params[qSlice[0]] = qSlice[1]
			}
			var resp *vsphere.Response
			vsc.Login()
			resp = vsc.Request("DELETE", args[0], params, []byte{})
			if resp.Error != nil {
				fmt.Println(resp.Error)
			} else {
				resp.Print()
			}
			vsc.Logout()
		},
	}

	return httpDeleteCmd
}

func readRequestData(fileName string) ([]byte, error) {
	if fileName == "-" {
		return readFromStdIn()
	} else {
		return ioutil.ReadFile(fileName)
	}
}

func readFromStdIn() ([]byte, error) {
	var body string
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		if err := stdin.Err(); err != nil {
			return []byte{}, err
		}
		body += stdin.Text()
	}
	return []byte(body), nil
}

func main() {
	vsc = vsphere.NewVSphereClient(debug)

	server := os.Getenv("KELPIE_VCENTER_SERVER")
	if server == "" {
		fmt.Printf("Environment variable 'KELPIE_VCENTER_SERVER' is not set.\n")
		os.Exit(1)
	}
	user := os.Getenv("KELPIE_VCENTER_USER")
	if user == "" {
		fmt.Printf("Environment variable 'KELPIE_VCENTER_USER' is not set.\n")
		os.Exit(1)
	}
	password := os.Getenv("KELPIE_VCENTER_PASSWORD")
	if password == "" {
		fmt.Printf("Environment variable 'KELPIE_VCENTER_PASSWORD' is not set.\n")
		os.Exit(1)
	}
	vsc.BaseUrl = server
	vsc.SetCredential(user, password)

	cmd := newCmd()
	cmd.Execute()
}
