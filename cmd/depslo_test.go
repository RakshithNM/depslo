// DEPSLO Command line utility tests
package main

// Importing all the required packages for our tests to work
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// Test the getFileData function
func Test_getFileData(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string   // The name of the test
		want    string   // What filepath we want our function to return.
		wantErr bool     // whether or not we want an error.
		osArgs  []string // The command arguments used for this test
	}{
		{"Default parameters", "test.json", false, []string{"depslo", "test.json"}},
		{"No parameters", "", true, []string{"depslo"}},
	}

	// Looping over the tests struct
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving the original os.Args reference
			actualOsArgs := os.Args
			// This defer function will run after the test is done
			defer func() {
				os.Args = actualOsArgs                                           // Restoring the original os.Args reference
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reseting the Flag command line. So that we can parse flags again
			}()

			os.Args = tt.osArgs             // Setting the specific command args for this test
			got, err := getFileData()       // Runing the function we want to test
			if (err != nil) != tt.wantErr { // Asserting whether or not we get the corret error value
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) { // Asserting whether or not we get the corret wanted value
				t.Errorf("getFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Testing teh checkIfValidFile function
func Test_checkIfValidFile(t *testing.T) {
	// Creating a temporal and empty JSON file
	tmpfile, err := ioutil.TempFile("", "tester.*.json")
	if err != nil {
		panic(err) // This should never happen
	}
	// Once all the tests are done. We delete the temporal file
	defer os.Remove(tmpfile.Name())
	// Defining the struct we're going to use
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{ // Defining our test cases
		{"File does exist", tmpfile.Name(), true, false},
		{"File does not exist", "nowhere/test.csv", false, true},
		{"File is not json", "test.txt", false, true},
	}
	// Iterating over our test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfValidFile(tt.filename)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Checking the returning value
			if got != tt.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Testing the processJSONFile function
func Test_processJSONFile(t *testing.T) {
	// Defining the map we're expenting to get from our function
	validTestJSONData := map[string]string{
		"HELLO": "Hello this is depslo",
		"TITLE": "The coolest developer tool",
	}

	// Validate the map data
	validJSON, marshalErr := json.Marshal(validTestJSONData)
	if marshalErr != nil {
		fmt.Printf("%v\n", marshalErr)
	}

	// Defining our test cases
	tests := []struct {
		name    string // The name of the test
		data    string // The key of our JSON file entry
		wantErr bool   // The value of our JSON file entry
	}{
		{"Valid JSON file", string(validJSON[:]), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, fileErr := ioutil.TempFile("", "valid.*.json")
			check(fileErr)

			// Removing the JSON test file before living
			defer os.Remove(tmpfile.Name())

			// Writing the content of the JSON test file
			_, writeErr := tmpfile.WriteString(fmt.Sprint(tt.data))
			if writeErr != nil {
				fmt.Printf("The write error: %v\n", writeErr)
			}

			// Persisting data on disk
			tmpfile.Sync()

			filePath := tmpfile.Name()

			// Check if the valid JSON file created is actually valid using our function
			jsonData, err := processJSONFile(filePath)

			if err == nil {
				if !reflect.DeepEqual(jsonData, validTestJSONData) {
					t.Errorf("processJSONFile() = %v, want %v", jsonData, validTestJSONData)
				}
			}
		})
	}
}
