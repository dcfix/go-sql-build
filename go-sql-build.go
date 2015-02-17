package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	RootDirectory      string
	Command            string
	CommandArgs        []string
	DirectoryStructure []string
}

func main() {
	// see if they passed in a config.json file,
	// if not, see if there is one in the local directory.

	config := loadConfig()
	processFiles(config)
}

func processFiles(config Config) {

	//todo: remove this test code
	// fmt.Println(myConfig.DirectoryStructure) // pre, table, view, function, sproc, post

	// loop through the DirectoryStructure items and process the files
	for _, dirName := range config.DirectoryStructure {
		path := config.RootDirectory + "\\" + dirName
		// fmt.Println(path)

		err := filepath.Walk(path, visitFile)
		if err != nil {
			fmt.Println("error:", err)
		}

	}
}

func visitFile(path string, info os.FileInfo, err error) error {
	// loop through the files and pass them along to cmd.execute
	// the filename has to be the last arguement in the command line.

	if strings.HasSuffix(path, ".sql") {
		fmt.Printf(" processing %v\n", path)
		config := loadConfig()
		args := append(config.CommandArgs, path)

		runSql(config.Command, args)
		//results := runSql(config.Command, args)

		//fmt.Printf("your output is: %s\n", results)
	}

	return nil
}

func loadConfig() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Println("finished loading config")
	return config
}

func runSql(exe string, args []string) string {

	cmd := exec.Command(exe, args...)

	stdout, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(stdout))
	}
	return string(stdout)
}
