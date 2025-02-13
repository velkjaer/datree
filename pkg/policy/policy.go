package policy

import (
	_ "embed"
	"fmt"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/datreeio/datree/pkg/fileReader"
	"github.com/datreeio/datree/pkg/jsonSchemaValidator"
	"github.com/ghodss/yaml"
)

//go:embed defaultRules.yaml
var defaultRulesYamlContent string

//go:embed policiesSchema.json
var policiesSchemaContent string

type DefaultRulesDefinitions struct {
	ApiVersion string                   `yaml:"apiVersion"`
	Rules      []*DefaultRuleDefinition `yaml:"rules"`
}

type DefaultRuleDefinition struct {
	ID               int                    `yaml:"id"`
	Name             string                 `yaml:"name"`
	UniqueName       string                 `yaml:"uniqueName"`
	EnabledByDefault bool                   `yaml:"enabledByDefault"`
	DocumentationUrl string                 `yaml:"documentationUrl"`
	MessageOnFailure string                 `yaml:"messageOnFailure"`
	Category         string                 `yaml:"category"`
	Schema           map[string]interface{} `yaml:"schema"`
}

func GetDefaultRules() (*DefaultRulesDefinitions, error) {
	defaultRulesDefinitions, err := yamlToStruct(defaultRulesYamlContent)
	return defaultRulesDefinitions, err
}

func GetPoliciesFileFromPath(path string) (*cliClient.EvaluationPrerunPolicies, error) {
	fileReader := fileReader.CreateFileReader(nil)
	policiesStr, err := fileReader.ReadFileContent(path)
	if err != nil {
		return nil, err
	}

	err = validatePoliciesYaml(policiesStr, path)
	if err != nil {
		return nil, err
	}

	var policies *cliClient.EvaluationPrerunPolicies
	policiesBytes, err := yaml.YAMLToJSON([]byte(policiesStr))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(policiesBytes, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func validatePoliciesYaml(content string, policyYamlPath string) error {
	jsonSchemaValidator := jsonSchemaValidator.New()
	result, err := jsonSchemaValidator.Validate(policiesSchemaContent, content)

	if err != nil {
		return err
	}

	if !result.Valid() {
		validationErrors := fmt.Errorf("Found errors in policies file %s:\n", policyYamlPath)

		for _, validationError := range result.Errors() {
			validationErrors = fmt.Errorf("%s\n%s", validationErrors, validationError)
		}

		return validationErrors
	}

	return nil
}

func yamlToStruct(content string) (*DefaultRulesDefinitions, error) {
	var defaultRulesDefinitions DefaultRulesDefinitions
	err := yaml.Unmarshal([]byte(content), &defaultRulesDefinitions)
	if err != nil {
		return nil, err
	}
	return &defaultRulesDefinitions, err
}
