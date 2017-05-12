package commands

import (
	"fmt"
	"log"
	"os"

	"encoding/json"

	"bytes"

	"github.com/mesosphere/dcos-commons/cli/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

type DescribeHandler struct {
	DescribeName string
}

func parseResolvedOptions(responseBytes []byte) string {
	var responseJSON map[string]interface{}
	err := json.Unmarshal([]byte(responseBytes), &responseJSON)
	if err != nil {
		log.Printf("Failed to unmarshall JSON from response.")
		log.Printf("Original data follows:")
		outBuf := *bytes.NewBuffer(responseBytes)
		outBuf.WriteTo(os.Stdout)
		os.Exit(1)
	}
	var resolvedOptions = responseJSON["properties"].(map[string]interface{})["resolvedOptions"].(string)
	return resolvedOptions
}

func (cmd *DescribeHandler) DescribeConfiguration(c *kingpin.ParseContext) error {
	// Call out to <dcos_url>/cosmos/service/describe w/ auth header
	//      Accept:'application/vnd.dcos.service.describe-response+json;charset=utf-8;version=v1' \
	//      Content-Type:'application/vnd.dcos.service.describe-request+json;charset=utf-8;version=v1' \
	//      appId="$1"
	// TODO: what is the intended output here? Do we want to show upgradesTo/downgradesTo components?
	// TODO: set arguments
	urlPath := "cosmos/service/describe"
	response := client.HTTPGet(urlPath)
	resolvedOptions := parseResolvedOptions(client.GetResponseBytes(response))
	client.PrintText(resolvedOptions)
	return nil
}

type UpdateHandler struct {
	UpdateName     string
	File           *os.File
	PackageVersion string
}

func (cmd *UpdateHandler) UpdateConfiguration(c *kingpin.ParseContext) error {
	// Call out to <dcos_url>/cosmos/service/update w/ auth header
	// Accept:'application/vnd.dcos.service.update-response+json;charset=utf-8;version=v1' \
	// Content-Type:'application/vnd.dcos.service.update-request+json;charset=utf-8;version=v1' \
	// TODO: read file in
	fmt.Printf("Updatin")
	return nil
}

func HandleServiceSection(app *kingpin.Application) {
	pkg := app.Command("service", "Manage service package configuration")

	describeCmd := &DescribeHandler{}
	pkg.Command("describe", "View the package configuration for this DC/OS service").Action(describeCmd.DescribeConfiguration)

	updateCmd := &UpdateHandler{}
	update := pkg.Command("update", "Update the package version or configuration for this DC/OS service").Action(updateCmd.UpdateConfiguration)
	update.Flag("--options", "Path to a JSON file that contains customized package installation options").FileVar(&updateCmd.File)
	update.Flag("--package-version", "The desired package version.").StringVar(&updateCmd.PackageVersion)
}
