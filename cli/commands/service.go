package commands

import (
	"io/ioutil"
	"log"
	"os"

	"encoding/json"

	"bytes"

	"github.com/mesosphere/dcos-commons/cli/client"
	"github.com/mesosphere/dcos-commons/cli/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

type DescribeHandler struct {
	DescribeName string
}

type DescribeRequest struct {
	AppID string `json:"appId"`
}

func unmarshallJSON(responseBytes []byte) map[string]interface{} {
	var responseJSON map[string]interface{}
	err := json.Unmarshal([]byte(responseBytes), &responseJSON)
	if err != nil {
		log.Printf("Failed to unmarshall JSON from response.")
		log.Printf("Original data follows:")
		outBuf := *bytes.NewBuffer(responseBytes)
		outBuf.WriteTo(os.Stdout)
		os.Exit(1)
	}
	return responseJSON
}

func parseDescribeResponse(responseBytes []byte) string {
	// TODO: what is the intended output here? Do we want to show upgradesTo/downgradesTo components?
	responseJSON := unmarshallJSON(responseBytes)
	return responseJSON["properties"].(map[string]interface{})["resolvedOptions"].(string)
}

func (cmd *DescribeHandler) DescribeConfiguration(c *kingpin.ParseContext) error {
	// TODO: add error handling
	requestContent, _ := json.Marshal(DescribeRequest{config.ServiceName})
	cosmosUrlPath := "service/describe"
	response := client.HTTPCosmosGetJSON(cosmosUrlPath, string(requestContent))
	resolvedOptions := parseDescribeResponse(client.GetResponseBytes(response))
	client.PrintText(resolvedOptions)
	return nil
}

type UpdateHandler struct {
	UpdateName     string
	OptionsFile    string
	PackageVersion string
}

type UpdateRequest struct {
	AppID          string `json:"appId"`
	PackageVersion string `json:"packageVersion"`
	OptionsJSON    string `json:"options"`
}

func checkAndReadFile(filename string) (string, error) {
	// TODO: any validation?
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseUpdateResponse(responseBytes []byte) string {
	// TODO: do something interesting with this output
	// Output should be in the same format as the `dcos marathon app update` command
	responseJSON := unmarshallJSON(responseBytes)
	return responseJSON["properties"].(map[string]interface{})["marathonDeploymentId"].(string)
}

func (cmd *UpdateHandler) UpdateConfiguration(c *kingpin.ParseContext) error {
	cosmosUrlPath := "service/update"
	request := UpdateRequest{AppID: config.ServiceName}
	if len(cmd.PackageVersion) > 0 {
		// TODO: check package version format is valid
		request.PackageVersion = cmd.PackageVersion
	}
	if len(cmd.OptionsFile) > 0 {
		optionsJSON, err := checkAndReadFile(cmd.OptionsFile)
		if err != nil {
			log.Fatalf("Failed to load specified options file %s: %s", cmd.OptionsFile, err)
		}
		request.OptionsJSON = optionsJSON
	}
	requestContent, _ := json.Marshal(request)
	response := client.HTTPCosmosPostJSON(cosmosUrlPath, string(requestContent))
	outputResponse := parseUpdateResponse(client.GetResponseBytes(response))
	client.PrintText(outputResponse)
	return nil
}

func HandleServiceSection(app *kingpin.Application) {
	pkg := app.Command("service", "Manage service package configuration")

	describeCmd := &DescribeHandler{}
	pkg.Command("describe", "View the package configuration for this DC/OS service").Action(describeCmd.DescribeConfiguration)

	updateCmd := &UpdateHandler{}
	update := pkg.Command("update", "Update the package version or configuration for this DC/OS service").Action(updateCmd.UpdateConfiguration)
	update.Flag("--options", "Path to a JSON file that contains customized package installation options").StringVar(&updateCmd.OptionsFile)
	update.Flag("--package-version", "The desired package version.").StringVar(&updateCmd.PackageVersion)
}
