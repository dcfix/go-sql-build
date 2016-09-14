package main

import (
	"encoding/json"
        "flag"
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
        configPath := flag.String("configPath", "gobolt.config.json", "The config file to process.")
        section := flag.String("section", "all", "The section to run.")
        flag.Parse()

	config := loadConfig(*configPath)
        if *section != "all" {
          newDirStructure := make([]string,1)
          newDirStructure[0] = *section
          config.DirectoryStructure = newDirStructure
        }
	processFiles(config)
}

func processFiles(config Config) {

	// loop through the DirectoryStructure items and process the files
	for _, dirName := range config.DirectoryStructure {
		path := config.RootDirectory + "/" + dirName
		//fmt.Println(path)

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
                  // fmt.Println(path)
                  if strings.HasSuffix(path, ".sql") {
                    fmt.Printf("  processing %v\n", path)
                    args := append(config.CommandArgs, path)
                    runSql(config.Command, args)
                }
                return nil
            })
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}

func loadConfig(configPath string) Config {
	file, _ := os.Open(configPath)
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
