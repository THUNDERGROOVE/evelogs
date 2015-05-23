package main

import "fmt"
import "time"

type Log struct {
	Channel     string
	SessionTime time.Time
	Entries     []*LogEntry
}

type LogEntry struct {
	Parent *Log
	User   string
	Text   string
	Time   time.Time
}

func (e *LogEntry) String() string {
	c := "<nil>"
	if e.Parent != nil {
		c = e.Parent.Channel
	}

	return fmt.Sprintf("[%v][%v] %v: %v", c, e.Time.Format(time.Kitchen), e.User, e.Text)
}
