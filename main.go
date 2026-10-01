package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Group struct {
	Name string   `json:"name"`
	URLs []string `json:"urls"`
}

type Config struct {
	Groups []Group `json:"groups"`
}

func getConfigPath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "png", "config.json")
}

func parseArgs() (string, string, string, string) {
	args := os.Args[1:]

	if len(args) < 1 {
		return "", "", "", ""
	}

	group := ""
	action := ""
	arg1 := ""
	arg2 := ""

	if len(args) > 0 {
		group = strings.TrimSpace(args[0])
	}
	if len(args) > 1 {
		action = strings.TrimSpace(args[1])
	}
	if len(args) > 2 {
		arg1 = strings.TrimSpace(args[2])
	}
	if len(args) > 3 {
		arg2 = strings.TrimSpace(args[3])
	}

	return group, action, arg1, arg2
}

func loadConfig(confPath string) (*Config, error) {

	// load raw file
	confRaw, err := os.ReadFile((confPath))

	if err != nil {
		return nil, err
	}

	// convert raw to the structs
	var cfg Config
	err = json.Unmarshal(confRaw, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getGroupInfo(config *Config, group string) *Group {

	for i, grp := range config.Groups {
		if grp.Name == group {
			return &config.Groups[i]
		}
	}
	return nil

}

func initConfig() (*Config, error) {

	// load or create configuration file
	config := &Config{}
	confPath := getConfigPath()
	fmt.Println("Using path:", confPath)
	if _, err := os.Stat(confPath); err != nil {
		os.MkdirAll(filepath.Dir(confPath), 0755)
		os.WriteFile(confPath, []byte("{\n  \"groups\": []\n}"), 0644)

	} else {
		config, _ = loadConfig(confPath)
		fmt.Println("File alread exists")
	}
	return config, nil
}

func main() {
	target_group, action, arg1, arg2 := parseArgs()
	config, err := initConfig()

	if err != nil {
		return
	}

	// Check if user is creating a group or invoking a request
	if action == "" && target_group != "" {
		// Logic for ping itself
		group := getGroupInfo(config, target_group)
		if group == nil {
			fmt.Println("Nonexistent group")
		} else {
			fmt.Println("Pinging group", group)
		}

	} else if action == "ls" {
		// Logic for listing items in groups
		group := getGroupInfo(config, target_group)
		if group == nil {
			fmt.Println("Nonexistent group")
		} else {
			fmt.Println("Listing items in", group)
		}
	} else if action == "new" {
		// Logic for creating a new group
		fmt.Println("Creating", target_group)

	} else if action == "add" {
		// Logic for adding an item to a group
		group := getGroupInfo(config, target_group)
		if group == nil {
			fmt.Println("Nonexistent group")
		} else {
			fmt.Printf("Adding %s - %s into %s \n", arg1, arg2, group)
		}

	} else if action == "rm" {
		if arg1 == "" {
			// Logic for removing a group
			group := getGroupInfo(config, target_group)
			if group == nil {
				fmt.Println("Nonexistent group")
			} else {
				fmt.Println("Deleting group", group)
			}
		} else {
			// Logic for removing a item from a group
			group := getGroupInfo(config, target_group)
			if group == nil {
				fmt.Println("Nonexistent group")
			} else {
				fmt.Printf("Deleting item %s in %s \n", arg1, group)
			}
		}

	} else if target_group == "" {
		fmt.Println("Usage:")
	}

}
