package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	ErrNoEVELogDir         = fmt.Errorf("the EVE log directory doesn't exist")
	ErrEVELogFileIsNotADir = fmt.Errorf("the file descriptor where the EVE log dir is not a folder")
)

func FindLogPath() (logpath string, err error) {
	usr, err := user.Current()
	dir := usr.HomeDir
	dir = filepath.Join(dir, "Documents", "EVE", "logs")
	if f, err := os.Stat(dir); err == os.ErrNotExist {
		return "", ErrNoEVELogDir
	} else {
		if !f.IsDir() {
			return "", ErrEVELogFileIsNotADir
		}
	}
	return dir, nil
}

func FindLogsForChannel(channel string) (logs []*Log) {
	logs = make([]*Log, 0)
	base, err := FindLogPath()
	if err != nil {
		ReportError("[ERROR] error finding logs for specific channel [%v]\n", err.Error())
		return logs
	}

	base = filepath.Join(base, "Chatlogs")

	if files, err := ioutil.ReadDir(base); err == nil {
		for _, f := range files {
			if !f.IsDir() {
				if _, name := filepath.Split(f.Name()); strings.Contains(name, channel) {
					fpath := filepath.Join(base, f.Name())
					p := NewParserFromFile(fpath)
					logs = append(logs, p.Parse())
				}
			}
		}
	} else {
		ReportError("[ERROR] error reading logs in directory [%v]\n", err.Error())
	}
	return logs
}
