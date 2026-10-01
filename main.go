package main

import (
	"encoding/json"
	"flag"
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

func initFlags() (string, string) {
	var targetGroup string
	var newGroup string

	// extract flags
	flag.StringVar(&newGroup, "n", "", "Name of new group")
	flag.StringVar(&newGroup, "new", "", "Name of new group")
	flag.StringVar(&targetGroup, "g", "", "Group to select to ping")
	flag.StringVar(&targetGroup, "group", "", "Group to select to ping")
	flag.Parse()
	newGroup = strings.TrimSpace(newGroup)
	targetGroup = strings.TrimSpace(targetGroup)

	return newGroup, targetGroup
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
	newGroup, targetGroup := initFlags()
	conf, err := initConfig()
	fmt.Println(conf)

	if err != nil {
		return
	}

	// Check if user is creating a group or invoking a request
	if newGroup != "" {
		// Logic for new group
		fmt.Println("Making group", newGroup)

	} else if targetGroup != "" {
		// Locic for ping
		fmt.Println("Pinging items in", targetGroup)
	} else {
		flag.Usage()
		return
	}

}
