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

var filePath string

func getFileData() (string, error) {
	// We need to validate that we're getting the correct number of arguments
	if len(os.Args) < 2 {
		return "", errors.New("A filepath argument is required")
	}

	flag.Parse() // This will parse all the arguments from the terminal

	fileLocation := flag.Arg(0) // The only argument (that is not a flag option) is the file location (JSON file)

	return fileLocation, nil
}

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

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
}

func processJSONFile(filePath string) (map[string]string, error) {
	file, err := ioutil.ReadFile(filePath)
	check(err)

	var data map[string]string

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

func writeToJSONFile(contentToWrite map[string]string) (string, error) {
	jsonValue, marshalErr := json.MarshalIndent(contentToWrite, "", "\t")
	if marshalErr != nil {
		fmt.Printf("ERRR")
		fmt.Printf("%v\n", marshalErr)
	}

	fileWrtingErr := ioutil.WriteFile("depslo.json", jsonValue, 0644)

	if fileWrtingErr != nil {
		return "", errors.New("File path is empty")
	}
	return "depslo.json", nil
}

func main() {
	fileData, errGettingFile := getFileData()
	check(errGettingFile)

	if len(fileData) > 0 {
		isValidJSONFile, errValidFile := checkIfValidFile(fileData)
		check(errValidFile)

		if isValidJSONFile == true {
			fileContent, errSyntaxCheck := processJSONFile(fileData)
			check(errSyntaxCheck)

			psuedoLocalContent := core.PsuedoLocalize(fileContent)

			writtenFilePath, err := writeToJSONFile(psuedoLocalContent)
			check(err)

			fmt.Printf("The psuedo local converted JSON file is at %s", writtenFilePath)
		}
	}
}
