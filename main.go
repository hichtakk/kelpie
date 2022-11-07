package main

import (
	"fmt"
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
		Short: "vSphere command-line client",
		Long:  "vSphere command-line client",
	}
	rootCmd.AddCommand(
		NewCmdHttpGet(&query),
		NewCmdHttpPost(&query),
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
			resp.Print(false)
			vsc.Logout()
		},
	}

	return httpGetCmd
}

func NewCmdHttpPost(query *[]string) *cobra.Command {
	httpGetCmd := &cobra.Command{
		Use:   "post ${API-PATH}",
		Short: "call api with HTTP POST method",
		Long:  "example) kelpie post /api/vcenter/vm/{vm}/power -q action=reset",
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
			resp = vsc.Request("POST", args[0], params, []byte{})
			resp.Print(false)
			vsc.Logout()
		},
	}

	return httpGetCmd
}

func main() {
	vsc = vsphere.NewVSphereClient(false)

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
