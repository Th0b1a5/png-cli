package main

import (
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

func getConfigPath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "png", "config.json")
}

func initFlags() (string, string) {
	var targetGroup string
	var newGroup string
	confPath := getConfigPath()

	flag.StringVar(&newGroup, "n", "", "Name of new group")
	flag.StringVar(&newGroup, "new", "", "Name of new group")
	flag.StringVar(&targetGroup, "g", "", "Group to select to ping")
	flag.StringVar(&targetGroup, "group", "", "Group to select to ping")
	flag.Parse()
	newGroup = strings.TrimSpace(newGroup)
	targetGroup = strings.TrimSpace(targetGroup)

	fmt.Println("Using path:", confPath)
	if _, err := os.Stat(confPath); err != nil {
		os.MkdirAll(filepath.Dir(confPath), 0755)
		os.WriteFile(confPath, []byte("{\n  \"groups\": []\n}"), 0644)

	} else {
		fmt.Println("File alread exists")
	}

	return newGroup, targetGroup
}

func main() {
	newGroup, targetGroup := initFlags()

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
