package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	var targetGroup string
	var newGroup string

	flag.StringVar(&newGroup, "n", "", "Name of new group")
	flag.StringVar(&newGroup, "new", "", "Name of new group")
	flag.StringVar(&targetGroup, "g", "", "Group to select to ping")
	flag.StringVar(&targetGroup, "group", "", "Group to select to ping")
	flag.Parse()
	newGroup = strings.TrimSpace(newGroup)
	targetGroup = strings.TrimSpace(targetGroup)

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
