package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
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
)

func main() {
	// 初始化版本映射
	versions := make(map[string]string)

	// 从32位和64位注册表视图读取JDK
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_32, REG_JDK_KEY_1, "x86", versions)
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_32, REG_JDK_KEY_2, "x86", versions)
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_64, REG_JDK_KEY_1, "x64", versions)
	readJDKRegistryKeys(registry.LOCAL_MACHINE, REGISTRY_ACCESS_READ_64, REG_JDK_KEY_2, "x64", versions)

	// 显示指示信息
	printInstructions()

	// 获取指定前缀的环境变量
	variables := getVariableNames()

	// 获取当前的 JAVA_HOME 和 PATH
	javahome, err := getSystemEnv(ENVIRONMENT_VARIABLE_TO_SET)
	if err != nil {
		log.Printf("Error retrieving %s: %v\n", ENVIRONMENT_VARIABLE_TO_SET, err)
	}
	currentPath, err := getSystemEnv("PATH")
	if err != nil {
		log.Printf("Error retrieving PATH: %v\n", err)
	}

	// 从 PATH 中移除现有的 JAVA_HOME\bin
	newPath := removeFromPath(currentPath, filepath.Join(javahome, "bin"))

	// 打印绿色的当前 JDK 版本
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("\n%s\n", green(fmt.Sprintf("Current JDK Version: %s\n", javahome)))

	// 将环境变量添加到版本映射
	for _, varName := range variables {
		varValue, err := getSystemEnv(varName)
		if err != nil {
			continue
		}
		if _, exists := versions[varValue]; !exists {
			versions[varValue] = "(env) " + varName
		}
	}

	// 创建有序的版本列表
	versionList := make([]struct {
		Path        string
		Description string
	}, 0)

	for path, desc := range versions {
		versionList = append(versionList, struct {
			Path        string
			Description string
		}{path, desc})
	}

	if len(versionList) == 0 {
		fmt.Println("未检测到任何 JDK 版本。请确保已安装 JDK 或手动添加环境变量。")
		return
	}

	// 使用 promptui 创建选择菜单，设置 Size 为 10
	prompt := promptui.Select{
		Label: "请选择一个 JDK 版本",
		Items: formatVersionList(versionList),
		Size:  10, // 设置选择框的高度为10，根据需要调整
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U0001F336 {{ . | cyan }}",
			Inactive: "  {{ . | cyan }}",
			Selected: "\U0001F336 {{ . | green | cyan }}",
		},
		Searcher: func(input string, index int) bool {
			item := formatVersionList(versionList)[index]
			return strings.Contains(strings.ToLower(item), strings.ToLower(input))
		},
	}

	// 运行选择菜单
	index, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("选择取消: %v\n", err)
		return
	}

	if result == "Exit" {
		fmt.Println("已退出程序。")
		return
	}

	selected := versionList[index]
	fmt.Printf("使用 %s\n", selected.Path)

	// 设置 JAVA_HOME
	err = setSystemEnv(ENVIRONMENT_VARIABLE_TO_SET, selected.Path)
	if err != nil {
		log.Printf("设置 %s 时出错: %v\n", ENVIRONMENT_VARIABLE_TO_SET, err)
	} else {
		fmt.Printf("成功设置 %s 为 %s\n", ENVIRONMENT_VARIABLE_TO_SET, selected.Path)
	}

	// 更新 PATH: 在 PATH 前面添加 JAVA_HOME\bin
	newPATH := filepath.Join(selected.Path, "bin") + ";" + newPath
	err = setSystemEnv("PATH", newPATH)
	if err != nil {
		log.Printf("更新 PATH 时出错: %v\n", err)
	} else {
		fmt.Println("成功更新 PATH。")
	}

	fmt.Println("更改已成功应用。您可能需要重新启动终端或重新登录系统以使更改生效。")
}

// printInstructions 打印指示信息
func printInstructions() {
	divider := strings.Repeat("*", 50)
	fmt.Println(divider)
	color.Yellow("1. 该软件会自动搜索通过标准方式安装的 JDK 版本。")
	color.Yellow("2. 若需手动添加，请在系统环境变量中创建类似 `JAVA_SDK_{版本号}` 的变量，并指向相应的 JDK 安装目录。")
	color.Yellow("3. 在系统环境变量的 `PATH` 中，顶部添加 `%JAVA_HOME%\\bin`。")
	fmt.Println(divider)
}

// formatVersionList 格式化版本列表以供 promptui 使用
func formatVersionList(versionList []struct {
	Path        string
	Description string
}) []string {
	formatted := []string{"Exit"}
	for _, v := range versionList {
		formatted = append(formatted, fmt.Sprintf("%s\tPath: %s", v.Description, v.Path))
	}
	return formatted
}

// readJDKRegistryKeys 从指定的注册表键读取 JDK 安装信息并更新版本映射
func readJDKRegistryKeys(hive registry.Key, access uint32, subKey string, machine string, versions map[string]string) {
	k, err := registry.OpenKey(hive, subKey, access)
	if err != nil {
		// 注册表键可能不存在，记录并继续
		log.Printf("未找到注册表键: %s (%s)\n", subKey, machine)
		return
	}
	defer k.Close()

	names, err := k.ReadSubKeyNames(-1)
	if err != nil {
		log.Printf("读取子键时出错: %s (%s): %v\n", subKey, machine, err)
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

// getVariableNames 获取所有以指定前缀开头的环境变量名
func getVariableNames() []string {
	envs, err := getSystemEnvs()
	if err != nil {
		log.Printf("获取系统环境变量时出错: %v\n", err)
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

// getSystemEnv 获取指定系统环境变量的值
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

// getSystemEnvs 获取所有系统环境变量
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

// setSystemEnv 设置系统环境变量
func setSystemEnv(key, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, ENVIRONMENT_REGISTRY_PATH, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	return k.SetStringValue(key, value)
}

// removeFromPath 从 PATH 环境变量中移除指定路径
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
