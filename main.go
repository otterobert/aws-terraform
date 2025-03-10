package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"io/ioutil"
	"log"
	"os"
)

type TerraformState struct {
	Resources []struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		Instances []struct {
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

func parseTfstateFile() {
	roleIdToArn := make(map[string]string)
	roleIdToPolicies := make(map[string][]string)
	policyArnToName := make(map[string]string)

	// Open the .tfstate file
	stateFile, err := os.Open("terraform.tfstate")
	if err != nil {
		log.Fatal("Error opening state file:", err)
	}
	defer stateFile.Close()

	// Read the entire content of the .tfstate file
	stateBytes, err := ioutil.ReadAll(stateFile)
	if err != nil {
		log.Fatal("Error reading state file:", err)
	}

	// Unmarshal the JSON content into the TerraformState struct
	var state TerraformState
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		log.Fatal("Error unmarshalling state file:", err)
	}

	// Iterate through resources in the state file
	for _, resource := range state.Resources {
		if resource.Type == "aws_iam_role" {
			for _, instance := range resource.Instances {
				id := ""
				arn := ""
				for key, value := range instance.Attributes {
					if key == "arn" {
						arn = value.(string)
					}
					if key == "id" {
						id = value.(string)
					}
				}
				roleIdToArn[id] = arn
			}
		}

		if resource.Type == "aws_iam_role_policy_attachment" {
			for _, instance := range resource.Instances {
				roleId := ""
				policyArn := ""
				for key, value := range instance.Attributes {
					if key == "role" {
						roleId = value.(string)
					}
					if key == "policy_arn" {
						policyArn = value.(string)
					}
				}

				_, ok := roleIdToPolicies[roleId]
				if ok {
					roleIdToPolicies[roleId] = append(roleIdToPolicies[roleId], policyArn)
				} else {
					roleIdToPolicies[roleId] = []string{policyArn}
				}
			}
		}

		if resource.Type == "aws_iam_policy" {
			name := resource.Name
			for _, instance := range resource.Instances {
				policyArn := ""
				for key, value := range instance.Attributes {
					if key == "arn" {
						policyArn = value.(string)
					}
				}

				policyArnToName[policyArn] = name
			}
		}
	}

	for id, policies := range roleIdToPolicies {
		fmt.Printf("Role ID: %s\n", id)
		if arn, ok := roleIdToArn[id]; ok {
			fmt.Printf("Role ARN: %s\n", arn)
		}

		for _, policyArn := range policies {
			fmt.Printf("Policy ARN: %s\n", policyArn)
			if name, ok := policyArnToName[policyArn]; ok {
				fmt.Printf("Policy terraform name: %s\n", name)
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	parseTfstateFile()

	module, diags := tfconfig.LoadModule(".")
	if diags.HasErrors() {
		fmt.Printf("has errors")
	}

	for _, resource := range module.DataResources {
		fmt.Println(resource.Type)
	}
}
