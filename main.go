package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	ENVIRONMENT_VARIABLE_PREFIX = "JAVA_SDK_"
	ENVIRONMENT_VARIABLE_TO_SET = "JAVA_HOME"
	REG_JDK_KEY_1               = `SOFTWARE\JavaSoft\JDK`
	REG_JDK_KEY_2               = `SOFTWARE\JavaSoft\Java Development Kit`
	CURRENT_VERSION_KEY         = "CurrentVersion"
	ENVIRONMENT_REGISTRY_PATH   = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	REGISTRY_ACCESS_READ        = registry.READ
	REGISTRY_ACCESS_READ_32     = registry.READ | registry.WOW64_32KEY
	REGISTRY_ACCESS_READ_64     = registry.READ | registry.WOW64_64KEY
	REGISTRY_ACCESS_WRITE       = registry.SET_VALUE
	REGISTRY_ACCESS_WRITE_32    = registry.SET_VALUE | registry.WOW64_32KEY
	REGISTRY_ACCESS_WRITE_64    = registry.SET_VALUE | registry.WOW64_64KEY
	Green                       = "\033[32m"
	Reset                       = "\033[0m"
)

// 颜色常量

func main() {
	// Initialize versions map
	versions := make(map[string]string)

	// Read JDKs from 32-bit and 64-bit registry views
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_32, REG_JDK_KEY_1, "x86", versions)
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_32, REG_JDK_KEY_2, "x86", versions)
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_64, REG_JDK_KEY_1, "x64", versions)
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_64, REG_JDK_KEY_2, "x64", versions)

	// Display instructions
	fmt.Println("**************************************************")
	fmt.Println("1. The software automatically searches for JDK versions installed via standard methods.")
	fmt.Println("2. To add manually, create system environment variables for all JDK versions with names like: JAVA_SDK_{version}")
	fmt.Println("3. Add `%JAVA_HOME%\\bin` at the top of the `PATH` system environment variable.")
	fmt.Println("**************************************************")

	// Retrieve environment variables with the specified prefix
	variables := getVariableNames()

	// Retrieve current JAVA_HOME and PATH
	javahome, err := getSystemEnv(ENVIRONMENT_VARIABLE_TO_SET)
	if err != nil {
		log.Printf("Error retrieving %s: %v\n", ENVIRONMENT_VARIABLE_TO_SET, err)
	}
	currentPath, err := getSystemEnv("PATH")
	if err != nil {
		log.Printf("Error retrieving PATH: %v\n", err)
	}

	// Remove existing JAVA_HOME\bin from PATH
	newPath := removeFromPath(currentPath, filepath.Join(javahome, "bin"))

	fmt.Printf(Green+"Current JDK Version: %s\n"+Reset, javahome)

	// Display selection menu
	fmt.Println("Please select a JDK version:")
	fmt.Println("0\t: Exit")
	index := 1

	// Add environment variables to versions map
	for _, varName := range variables {
		varValue, err := getSystemEnv(varName)
		if err != nil {
			continue
		}
		if _, exists := versions[varValue]; !exists {
			versions[varValue] = "(env) " + varName
		}
	}

	// Create a slice to maintain order
	versionList := make([]struct {
		Path        string
		Description string
	}, 0)

	for path, desc := range versions {
		fmt.Printf("%d:\t\t %s \t\t %s\n", index, desc, path)
		versionList = append(versionList, struct {
			Path        string
			Description string
		}{path, desc})
		index++
	}

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	input = strings.TrimSpace(input)
	intInput, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Please enter a number.")
		return
	}

	if intInput == 0 {
		fmt.Println("Exiting the application as per user request.")
		return
	}

	if intInput < 1 || intInput > len(versionList) {
		fmt.Println("Invalid selection. Please run the program again and choose a valid option.")
		return
	}

	selected := versionList[intInput-1]
	fmt.Printf("Using %s\n", selected.Path)

	// Set JAVA_HOME
	err = setSystemEnv(ENVIRONMENT_VARIABLE_TO_SET, selected.Path)
	if err != nil {
		log.Printf("Error setting %s: %v\n", ENVIRONMENT_VARIABLE_TO_SET, err)
	} else {
		fmt.Printf("Set %s to %s successfully.\n", ENVIRONMENT_VARIABLE_TO_SET, selected.Path)
	}

	// Update PATH: Prepend JAVA_HOME\bin
	newPATH := filepath.Join(selected.Path, "bin") + ";" + newPath
	err = setSystemEnv("PATH", newPATH)
	if err != nil {
		log.Printf("Error updating PATH: %v\n", err)
	} else {
		fmt.Println("Updated PATH successfully.")
	}

	fmt.Println("Changes applied successfully. You may need to restart your terminal or log out and log back in for changes to take effect.")
}

// readJDKRegistryKeys reads JDK installations from the specified registry key and updates the versions map.
func readJDKRegistryKeys(hive registry.Key, access uint32, subKey string, machine string, versions map[string]string) {
	k, err := registry.OpenKey(hive, subKey, access)
	if err != nil {
		// It's possible that the key doesn't exist; log and continue
		log.Printf("Registry key not found: %s (%s)\n", subKey, machine)
		return
	}
	defer k.Close()

	names, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Printf("Error reading subkeys from %s (%s): %v\n", subKey, machine, err)
		return
	}

	for _, name := range names {
		versionKey, err := registry.OpenKey(k, name, registry.QUERY_VALUE)
		if err != nil {
			continue
		}

		javaHome, _, err := versionKey.GetStringValue("JavaHome")
		versionKey.Close()
		if err != nil || javaHome == "" {
			continue
		}

		if _, exists := versions[javaHome]; !exists {
			versions[javaHome] = fmt.Sprintf("(%s) %s", machine, name)
		}
	}
}

// getVariableNames retrieves environment variable names that start with the specified prefix.
func getVariableNames() []string {
	envs, err := getSystemEnvs()
	if err != nil {
		log.Printf("Error retrieving system environment variables: %v\n", err)
		return []string{}
	}

	var variables []string
	for key := range envs {
		if strings.HasPrefix(key, ENVIRONMENT_VARIABLE_PREFIX) {
			varName := strings.TrimPrefix(key, ENVIRONMENT_VARIABLE_PREFIX)
			variables = append(variables, varName)
		}
	}
	return variables
}

// getSystemEnv retrieves the value of a system environment variable from the registry.
func getSystemEnv(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, ENVIRONMENT_REGISTRY_PATH, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	val, _, err := k.GetStringValue(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

// getSystemEnvs retrieves all system environment variables from the registry.
func getSystemEnvs() (map[string]string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, ENVIRONMENT_REGISTRY_PATH, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer k.Close()

	names, err := k.ReadValueNames(-1)
	if err != nil {
		return nil, err
	}

	envs := make(map[string]string)
	for _, name := range names {
		val, _, err := k.GetStringValue(name)
		if err == nil {
			envs[name] = val
		}
	}
	return envs, nil
}

// setSystemEnv sets a system environment variable in the registry.
func setSystemEnv(key, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, ENVIRONMENT_REGISTRY_PATH, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	return k.SetStringValue(key, value)
}

// removeFromPath removes a specific path from the PATH environment variable.
func removeFromPath(path, toRemove string) string {
	parts := strings.Split(path, ";")
	newParts := []string{}
	toRemove = strings.ToLower(toRemove)
	for _, p := range parts {
		pTrimmed := strings.TrimSpace(p)
		if strings.ToLower(pTrimmed) == strings.ToLower(toRemove) {
			continue
		}
		newParts = append(newParts, pTrimmed)
	}
	return strings.Join(newParts, ";")
}
