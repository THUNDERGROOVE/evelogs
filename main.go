package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

var (
	List     = flag.Bool("l", false, "List information about your chat logs")
	Channel  = flag.String("c", "", "the channel to look for chat logs in")
	User     = flag.String("u", "", "the user to filter chat logs from")
	FilterRE = flag.String("fr", "", "A regex string to filter messages that match the given regex string")
	Filter   = flag.String("f", "", "A string that filters log entries that contain the given string")
)

func main() {
	flag.Parse()

	var logs []*Log
	var filteredEntries = make([]*LogEntry, 0)

	logs = FindLogsForChannel(*Channel)

	if logs == nil || len(logs) == 0 {
		ReportError("[ERROR] no logs were returned\n")
	}

	if *List {
		logc := len(logs)
		var ec int
		var pc int
		users := make([]string, 0)

		for _, v := range logs {
			for _, e := range v.Entries {
				if !contains(e.User, users) {
					users = append(users, e.User)
					pc += 1
				}
				ec += 1
			}
		}
		fmt.Printf("You have:\n")
		fmt.Printf("=> %v log files containing\n", logc)
		fmt.Printf("=> %v individual lines of chat\n", ec)
		fmt.Printf("=> written by %v individual characters\n", pc)
		return
	}

	fmt.Printf("Got %v logs\n", len(logs))

	if *FilterRE != "" {
		if re, err := regexp.Compile(*FilterRE); err == nil {
			for _, v := range logs {
				for _, l := range v.Entries {
					if re.Match([]byte(l.Text)) {
						if strings.Contains(l.User, *User) {
							filteredEntries = append(filteredEntries, l)
						}
					}
				}
			}
		} else {
			ReportError("[ERROR] couldn't compile regex [%v]\n", err.Error())
			return
		}
	} else {
		for _, v := range logs {
			for _, l := range v.Entries {
				if strings.Contains(l.Text, *Filter) {

					if *User == "" || strings.Contains(l.User, *User) {
						filteredEntries = append(filteredEntries, l)
					}
				}
			}
		}
	}

	if filteredEntries == nil || len(filteredEntries) == 0 {
		ReportError("[ERROR] no filtered logs were returned\n")
	}
	fmt.Printf("Found %v entries\n", len(filteredEntries))
	for _, v := range filteredEntries {
		fmt.Println(v)
	}
}

func contains(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
