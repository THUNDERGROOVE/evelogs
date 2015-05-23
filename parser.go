package main

import (
	"io/ioutil"
	"strings"
	"time"
)

const (
	metaDataSpliter        = "---------------------------------------------------------------"
	metaDataChannelID      = "Channel ID:" // Not parsing for now
	metaDataChannelName    = "Channel Name:"
	metaDataSessionStarted = "Session started:"

	timeFmt = "2006.01.02 15:04:05"
)

type Parser struct {
	file     string
	curLine  int
	prevLine string
}

func NewParserFromFile(file string) *Parser {
	return &Parser{
		file: file,
	}
}

func (p *Parser) Parse() *Log {
	out := &Log{}
	p.curLine = 0
	var (
		hasMetadata        bool
		lookingForMetadata bool
	)

	out.Entries = make([]*LogEntry, 0)

	data := string(readFile(p.file))
	data = strings.Replace(data, "\xFF\xFE", "", -1) // Trim this weird bullshit CCP adds before each lien
	data = strings.Replace(data, "\x00", "", -1)     // Trim remaining NULL characters
	data = strings.Replace(data, "\x0D", "", -1)     // Trim carriage return
	lines := strings.Split(data, "\n")
	for _, v := range lines {
		p.curLine += 1
		if len(v) == 0 { // Ignore if line is blank
			continue
		}
		if v == metaDataSpliter && !hasMetadata && !lookingForMetadata {
			lookingForMetadata = true
			continue
		}

		if v == metaDataSpliter && lookingForMetadata {
			lookingForMetadata = false
			hasMetadata = true
			continue
		}

		if lookingForMetadata && !hasMetadata {
			if strings.Contains(v, metaDataChannelName) {
				out.Channel = p.parseValueFromMetadataString(v, metaDataChannelName)
			} else if strings.Contains(v, metaDataSessionStarted) {
				var err error
				out.SessionTime, err = time.Parse(timeFmt, p.parseValueFromMetadataString(v, metaDataSessionStarted))
				if err != nil {
					ReportError("[ERROR] couldn't parse time from Session metadata at %v:%v [%v]\n", p.file, p.curLine, err.Error())
				}
			}
		}
		if hasMetadata {
			if v[0] != '[' {
				continue
			}
			out.Entries = append(out.Entries, p.parseLineFromLog(v, out))
			p.prevLine = v
		}
	}
	return out
}

func (p *Parser) parseLineFromLog(v string, parent *Log) *LogEntry {
	defer func() {
		if r := recover(); r != nil {
			ReportError("[ERROR] paniced while trying to parse the following line\nFile:%v:%v\nString:%v\nPrevString:%v\nPrevBytes:%v\n[%]\n", p.file, p.curLine, v, p.prevLine, []byte(p.prevLine), r)
		}
	}()

	out := new(LogEntry)
	out.Parent = parent

	t, err := time.Parse(timeFmt, v[2:][:19])
	if err != nil {
		ReportError("[ERROR] couldn't parse time from log entry \n%v:%v\n%v[%v]\n", p.file, p.curLine, v, err.Error())
		return nil
	}
	out.Time = t
	nameandtext := v[24:]

	vv := strings.Split(nameandtext, " > ")
	name := vv[0]
	text := vv[1]

	out.User = name
	out.Text = text

	return out
}

func (p *Parser) parseValueFromMetadataString(v string, metadataident string) string {
	s := strings.Join(strings.Split(v, metadataident), "")
	s = strings.Trim(s, " ")
	return s

}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		ReportError("[Error] Couldn't open the file: %v [%v]", err.Error())
	}
	return data
}
