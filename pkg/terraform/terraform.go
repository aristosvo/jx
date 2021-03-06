package terraform

import (
	"fmt"
	"github.com/jenkins-x/jx/pkg/util"
	"io/ioutil"
	"os"
	"strings"
)

func Init(terraformDir string, serviceAccountPath string) error {
	os.Setenv("GOOGLE_CREDENTIALS", serviceAccountPath)
	err := util.RunCommand("", "terraform", "init", terraformDir)
	if err != nil {
		return err
	}
	return nil
}

func Plan(terraformDir string, terraformVars string, serviceAccountPath string) error {
	err := util.RunCommand("", "terraform", "plan",
		fmt.Sprintf("-var-file=%s", terraformVars),
		"-var",
		fmt.Sprintf("credentials=%s", serviceAccountPath),
		terraformDir)
	if err != nil {
		return err
	}
	return nil
}

func Apply(terraformDir string, terraformVars string, serviceAccountPath string) error {
	err := util.RunCommand("", "terraform", "apply", "-auto-approve",
		fmt.Sprintf("-var-file=%s", terraformVars),
		"-var",
		fmt.Sprintf("credentials=%s", serviceAccountPath),
		terraformDir)
	if err != nil {
		return err
	}
	return nil
}

func WriteKeyValueToFileIfNotExists(path string, key string, value string) error {
	// file exists
	if _, err := os.Stat(path); err == nil {
		buffer, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		contents := string(buffer)

		if strings.Contains(contents, key) {
			return nil
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintf("%s = \"%s\"\n", key, value)

	_, err = file.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}

func ReadValueFromFile(path string, key string) (string, error) {
	if _, err := os.Stat(path); err == nil {
		buffer, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		contents := string(buffer)
		lines := strings.Split(contents, "\n")
		for _, line := range lines {
			if strings.Contains(line, key) {
				tokens := strings.Split(line, "=")
				trimmedValue := strings.Trim(strings.TrimSpace(tokens[1]), "\"")
				return trimmedValue, nil
			}
		}

	}
	return "", nil
}
