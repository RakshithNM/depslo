// DEPSLO command line utility source code
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	core "github.com/RakshithNM/depslo/core"
)

// getFileData can be used to get the meta data of the file passed as the parameter to the cli
func getFileData() (string, error) {
	// Validate the number of arguments passed
	if len(os.Args) < 2 {
		return "", errors.New("A filepath argument is required")
	}

	flag.Parse()

	// Get the location of the file that is supposed to be psuedo localized
	fileLocation := flag.Arg(0)

	return fileLocation, nil
}

// CheckIfValidFile can be used to check if the file is JSON file and that the file exists
func checkIfValidFile(filename string) (bool, error) {
	// Checking if entered file is CSV by using the filepath package from the standard library
	if fileExtension := filepath.Ext(filename); fileExtension != ".json" {
		return false, fmt.Errorf("File %s is not JSON", filename)
	}

	// Checking if filepath entered belongs to an existing file. We use the Stat method from the os package (standard library)
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("File %s does not exist", filename)
	}
	// If we get to this point, it means this is a valid file
	return true, nil
}

// exitGracefully can be used to exit gracefully on error
func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

// check can be used if the previous function call resulted in an error, this inturn calls exitGracefully
func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
}

// processJSONFile reads the content of the file path and checks if there are syntax errors
func processJSONFile(filePath string) (map[string]string, error) {
	// Read the file
	file, err := ioutil.ReadFile(filePath)
	check(err)

	var data map[string]string

	// Unmarshall to a map to see if there are any errors in JSON content
	jErr := json.Unmarshal(file, &data)

	if jErr != nil {
		switch jsErr := jErr.(type) {
		case *json.SyntaxError:
			fmt.Printf("Error in input syntax at byte %d: %s\n", jsErr.Offset, jsErr.Error())
			break
		default:
			fmt.Printf("Other error decoding JSON: %s\n", jsErr.Error())
		}
	}
	return data, nil
}

// writeToJSONFile can be used to create indented JSON file with the passed in content
func writeToJSONFile(contentToWrite map[string]string) (string, error) {
	// Make indented JSON value out of the content passed in
	jsonValue, marshalErr := json.MarshalIndent(contentToWrite, "", "\t")
	if marshalErr != nil {
		fmt.Printf("%v\n", marshalErr)
	}

	outputFile := "depslo.json"

	// Write the JSON value to a file
	fileWrtingErr := ioutil.WriteFile(outputFile, jsonValue, 0644)

	if fileWrtingErr != nil {
		return "", errors.New("File path is empty")
	}
	return outputFile, nil
}

func main() {
	// Get the file name and handle errors in cli cmd input
	fileData, errGettingFile := getFileData()
	// Check and exit gracefully if there was an error
	check(errGettingFile)

	if len(fileData) > 0 {
		// Check if the file was a valid JSON file (extension and if file exists)
		isValidJSONFile, errValidFile := checkIfValidFile(fileData)
		// Check and exit gracefully if there was an error
		check(errValidFile)

		if isValidJSONFile == true {
			// Check if the content of the file is a valid JSON
			fileContent, errSyntaxCheck := processJSONFile(fileData)
			// Check and exit gracefully if there was an error
			check(errSyntaxCheck)

			// Call the PsuedoLocalize function from the depslo core package(~root/core)
			psuedoLocalContent := core.PsuedoLocalize(fileContent)

			// Write the psuedo localized JSON value to a file and return the file path for printing
			writtenFilePath, err := writeToJSONFile(psuedoLocalContent)
			// Check and exit gracefully if there was an error
			check(err)

			// Prompt the user that the process was complete and let the user know of the file name
			fmt.Printf("The psuedo local converted JSON file is at %s", writtenFilePath)
		}
	}
}
