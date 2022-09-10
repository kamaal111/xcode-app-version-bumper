package bumpVersion

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DEFAULT_STRING_FLAG      = ""
	DEFAULT_INT_FLAG         = 0
	XCODE_BUILD_NUMBER_KEY   = "CURRENT_PROJECT_VERSION"
	XCODE_VERSION_NUMBER_KEY = "MARKETING_VERSION"
)

func BumpVersion(versionNumber *string, buildNumber *int, projectPath *string) {
	definitiveBuildNumber, err := getDefinitiveBuildNumber(buildNumber)
	if err != nil {
		log.Fatalln(err)
	}

	definitiveVersion, err := getDefinitiveVersion(versionNumber)
	if err != nil {
		log.Fatalln(err)
	}

	definitiveProjectPath, err := getDefinitiveProjectPath(projectPath)
	if err != nil {
		log.Fatalln(err)
	}

	hasChanged, err := editNumbers(definitiveBuildNumber, definitiveVersion, definitiveProjectPath)
	if err != nil {
		log.Fatalln(err)
	}

	if !hasChanged {
		fmt.Println("No changes where needed")
	} else {
		fmt.Println("Applied changes to xcode project")
	}
}

func editNumbers(buildNumber int, versionNumber string, projectPath string) (bool, error) {
	fullProjectConfigurationFilepath, err := findFullProjectConfigurationFilePath(projectPath)
	if err != nil {
		return false, err
	}

	configurationFileData, err := ioutil.ReadFile(fullProjectConfigurationFilepath)
	if err != nil {
		return false, err
	}

	configurationFileDataSplitByLines := strings.Split(string(configurationFileData), "\n")
	tabsMap := make(map[int]string)
	hasChanges := false
	for lineNumber, line := range configurationFileDataSplitByLines {
		isBuildNumber := strings.Contains(line, XCODE_BUILD_NUMBER_KEY)
		isVersionNumber := strings.Contains(line, XCODE_VERSION_NUMBER_KEY)

		if !isBuildNumber && !isVersionNumber {
			continue
		}

		oneTab := "	"
		amountOfTabs := strings.Count(line, oneTab)
		tabsToAdd := ""

		mappedTab, ok := tabsMap[amountOfTabs]
		if ok {
			tabsToAdd = mappedTab
		} else {
			for i := 0; i < amountOfTabs; i += 1 {
				tabsToAdd += oneTab
			}
			tabsMap[amountOfTabs] = tabsToAdd
		}

		var newLine string
		if isBuildNumber {
			newLine = fmt.Sprintf("%s%s = %d;", tabsToAdd, XCODE_BUILD_NUMBER_KEY, buildNumber)
		} else if isVersionNumber {
			newLine = fmt.Sprintf("%s%s = %s;", tabsToAdd, XCODE_VERSION_NUMBER_KEY, versionNumber)
		}

		if line == newLine {
			continue
		}

		configurationFileDataSplitByLines[lineNumber] = newLine

		if hasChanges {
			continue
		}

		hasChanges = true
	}

	if !hasChanges {
		return false, nil
	}

	newConfigurationData := strings.Join(configurationFileDataSplitByLines, "\n")
	err = ioutil.WriteFile(fullProjectConfigurationFilepath, []byte(newConfigurationData), 0644)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getDefinitiveProjectPath(flagValue *string) (string, error) {
	if flagValue != nil && *flagValue != DEFAULT_STRING_FLAG {
		return *flagValue, nil
	}

	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pathFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	var xcodeProjectFilepath string
	for _, pathFile := range pathFiles {
		if pathFile.IsDir() && strings.Contains(pathFile.Name(), ".xcodeproj") {
			xcodeProjectFilepath = filepath.Join(path, pathFile.Name())
			break
		}
	}
	if xcodeProjectFilepath == "" {
		return "", errors.New("xcode project not found at root")
	}

	return xcodeProjectFilepath, nil
}

func findFullProjectConfigurationFilePath(projectPath string) (string, error) {
	xcodeProjectFiles, err := ioutil.ReadDir(projectPath)
	if err != nil {
		return "", err
	}

	var xcodeProjectFile fs.FileInfo
	for _, pathFile := range xcodeProjectFiles {
		if pathFile.Name() == "project.pbxproj" {
			xcodeProjectFile = pathFile
			break
		}
	}

	if xcodeProjectFile == nil {
		return "", errors.New("xcode project file not found")
	}

	xcodeProjectConfigurationFilepath := filepath.Join(projectPath, xcodeProjectFile.Name())

	return xcodeProjectConfigurationFilepath, nil
}

func getDefinitiveVersion(flagValue *string) (string, error) {
	if flagValue != nil && *flagValue != DEFAULT_STRING_FLAG {
		return *flagValue, nil
	}

	var versionNumber string
	fmt.Printf("Version number: ")
	_, err := fmt.Scan(&versionNumber)
	if err != nil {
		return "", err
	}

	return versionNumber, nil
}

func getDefinitiveBuildNumber(flagValue *int) (int, error) {
	if flagValue != nil && *flagValue != DEFAULT_INT_FLAG {
		return *flagValue, nil
	}

	fmt.Printf("Build number: ")
	var buildNumberInputString string
	_, err := fmt.Scan(&buildNumberInputString)
	if err != nil {
		return 0, err
	}

	buildNumberInput, err := strconv.Atoi(buildNumberInputString)
	if err != nil {
		return 0, err
	}

	return buildNumberInput, nil
}
