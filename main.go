package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
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

	action := ""
	group := ""
	arg1 := ""
	arg2 := ""

	if len(args) > 0 {
		action = strings.TrimSpace(args[0])
	}
	if len(args) > 1 {
		group = strings.TrimSpace(args[1])
	}
	if len(args) > 2 {
		arg1 = strings.TrimSpace(args[2])
	}
	if len(args) > 3 {
		arg2 = strings.TrimSpace(args[3])
	}

	return action, group, arg1, arg2
}

func loadConfig(confPath string) (*Config, error) {

	// load raw file
	confRaw, err := os.ReadFile(confPath)
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

func writeConfig(config *Config) (err error) {
	configRaw, err := json.MarshalIndent(config, "", "    ")
	if err == nil {
		path := getConfigPath()
		err := os.WriteFile(path, configRaw, 0644)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return

}

func getGroupInfo(config *Config, group string) *Group {

	for i, grp := range config.Groups {
		if grp.Name == group {
			return &config.Groups[i]
		}
	}
	return nil

}

// Yes absolutely stolen from Stack Overflow
func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) || string(c) == "." {
			return false
		}
	}
	return true
}

func pingURL(url string) error {
	// Figure out the type of endpoint to ping, move to a seperate function later
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		fmt.Println(url, "isnt a rest endpoint!")
		return errors.New("Not REST")
	}
	client := http.Client{Timeout: 3 * time.Second}
	start := time.Now()

	// Send request
	response, err := client.Get(url)

	if err != nil {
		fmt.Println(url, "-", "Unreachable")
		return err
	}
	// Wait and measure
	if response.StatusCode == 200 {
		fmt.Println(url, "-", time.Since(start))
		response.Body.Close()
		return nil
	} else {
		fmt.Println(url, "-", "Offline")
		return nil
	}

}

func initConfig() (*Config, error) {

	// load or create configuration file
	config := &Config{}
	confPath := getConfigPath()
	// fmt.Println("Using path:", confPath)
	if _, err := os.Stat(confPath); err != nil {
		os.MkdirAll(filepath.Dir(confPath), 0755)
		os.WriteFile(confPath, []byte("{\n  \"groups\": []\n}"), 0644)

	} else {
		config, _ = loadConfig(confPath)
		// fmt.Println("File alread exists")
	}
	return config, nil
}

func main() {
	action, target_group, arg1, _ := parseArgs()
	config, err := initConfig()

	if err != nil {
		return
	}

	// Check if user is creating a group or invoking a request
	if action != "" && target_group == "" {
		// Logic for ping itself
		// ACTION NOW HOLDS THE GROUP IN THIS CASE
		group := getGroupInfo(config, action)
		if group == nil {
			fmt.Println("Nonexistent group")
		} else {
			fmt.Println("Pinging group", group.Name)
			for _, url := range group.URLs {
				pingURL(url)
			}
		}

	} else if action == "ls" {
		// Logic for listing items in groups
		group := getGroupInfo(config, target_group)
		if group == nil {
			fmt.Println("Nonexistent group")
		} else {

			// Goes throught every group.URLs in specified group, then listing them.
			for i, url := range group.URLs {
				fmt.Println(i, "-", url)
			}
		}

	} else if action == "new" {
		// First check if the group name is a known command to prevent delete conflicts later
		commands := []string{"ping", "ls", "new", "add", "rm", ""}
		for _, cmd := range commands {
			if cmd == target_group {
				fmt.Println("Cant use name of", target_group, "because its reseved for commands")
				return
			}
		}
		if isInt(target_group) == true {
			fmt.Println("New group cant be a plain number")
			return
		}

		// Logic for creating a new group
		group := getGroupInfo(config, target_group)
		if group == nil {

			// Create a new empty struct
			newGroup := Group{Name: target_group, URLs: []string{}}

			// Take existing config, take Groups group and append a new item to the slice
			config.Groups = append(config.Groups, newGroup)
			err := writeConfig(config)

			if err == nil {
				fmt.Println("Successfully created new group", string(target_group))
			} else {
				fmt.Println("Error writing to configuration file")
			}
		} else {
			fmt.Println("Group already exists")
		}

	} else if action == "add" {
		if arg1 == "" {
			fmt.Println("Missing URL to add to group")
			return
		}
		// Check if its a IP or a URL, this is to prevent deleting conflicts later
		if isInt(arg1) == true {
			fmt.Println("New item in a group cant be a plain number")
			return
		}

		// Logic for adding an item to a group
		group := getGroupInfo(config, target_group)
		if group == nil {
			fmt.Println("Nonexistent group")
		} else {

			// Append new URL to target group
			// I dont need to reappend stuff to "config" because group is a pointer to a Group inside config
			// Trippy but i have to get used to it
			group.URLs = append(group.URLs, arg1)
			err := writeConfig(config)

			if err == nil {
				fmt.Println("Successfully added", string(arg1), "to", string(target_group))
			} else {
				fmt.Println("Error writing to configuration file")
			}

		}

	} else if action == "rm" {
		if arg1 == "" {
			// Allows deleting a group by specifying its order number, no need to search for group in config
			if isInt(target_group) == true {
				i, _ := strconv.Atoi(target_group)
				if len(config.Groups) == 0 {
					fmt.Println("Configuration file has no groups!")
					return
				}
				if len(config.Groups) < i {
					fmt.Println("Index", i, "is out of range, maximum is", len(config.Groups))
					return
				}
				fmt.Println(len(config.Groups), " ", i)
				config.Groups = slices.Delete(config.Groups, i, i+1)

				// Write the final config
				err := writeConfig(config)

				if err == nil {
					fmt.Println("Successfully removed target group")
				} else {
					fmt.Println("Error writing to configuration file")
				}
				return
			}

			// Logic for removing a group
			group := getGroupInfo(config, target_group)
			if group == nil {
				fmt.Println("Nonexistent group")
			} else {
				// Find which position the target_group is in config, then delete it from config
				for i, grp := range config.Groups {
					if grp.Name == group.Name {
						config.Groups = slices.Delete(config.Groups, i, i+1)
					}
				}
				// Write the final config
				err := writeConfig(config)

				if err == nil {
					fmt.Println("Successfully removed", string(group.Name))
				} else {
					fmt.Println("Error writing to configuration file")
				}
			}
		} else {
			// Logic for removing a item from a group
			group := getGroupInfo(config, target_group)
			if group == nil {
				fmt.Println("Nonexistent group")
			} else {
				fmt.Printf("Deleting item %s in %s \n", arg1, group.Name)

				// Allows deleting a URL by specifying its order number, no need to search for URL in Group
				if isInt(arg1) == true {
					i, _ := strconv.Atoi(arg1)
					if len(config.Groups) == 0 {
						fmt.Println("Group has no URLs!")
						return
					}
					if len(group.URLs) < i {
						fmt.Println("Index", i, "is out of range, maximum is", len(group.URLs))
						return
					}
					group.URLs = slices.Delete(group.URLs, i, i+1)

					// Write the final config
					err := writeConfig(config)

					if err == nil {
						fmt.Println("Successfully removed target group")
					} else {
						fmt.Println("Error writing to configuration file")
					}
					return
				}

				// Find which position the URL is in the target_group, then delete it from group
				for i, url := range group.URLs {
					if url == arg1 {
						group.URLs = slices.Delete(group.URLs, i, i+1)
					}
				}
				// Write the final config
				err := writeConfig(config)

				if err == nil {
					fmt.Println("Successfully removed", string(arg1), "from", string(group.Name))
				} else {
					fmt.Println("Error writing to configuration file")
				}

			}
		}

	} else if action != "" {
		fmt.Println("Unknown action", action)

	} else if target_group == "" {
		// In case the command is without args
		fmt.Println("Listing all groups")
		// Goes thru the entire config listing only the position of the group and its name
		for i, grp := range config.Groups {
			fmt.Println(i, "-", grp.Name)
		}
	}

}
