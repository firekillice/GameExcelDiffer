package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type EventType int32

const (
	AddTable       EventType = 1
	DeleteTable    EventType = 2
	AddIdentity    EventType = 3
	DeleteIdentity EventType = 4
	AddField       EventType = 5
	DeleteField    EventType = 6
	ChangCell      EventType = 7
)

type EventInfo struct {
	eventType  EventType
	table      string
	identity   string
	field      string
	key        string
	content    string
	newContent string
}

var eventInfoList []EventInfo
var oldKeys []string
var newKeys []string
var allNewKeysJoinString string
var allOldKeysJoinString string

func diffContents(old map[string]string, new map[string]string) {
	for ko, vo := range old {
		if vn, exists := new[ko]; exists {
			if vn == vo {
				continue
			} else {
				eventInfoList = append(eventInfoList, EventInfo{eventType: ChangCell, key: ko, content: vo, newContent: vn})
			}
		} else {
			getNewKeys(new)

			confirmDeleteEvent(ko, allNewKeysJoinString)
		}
	}

	for kn, _ := range new {
		if _, exists := old[kn]; !exists {
			getOldKeys(old)
			confirmAddEvent(kn, allOldKeysJoinString)
		}
	}

	dumpEvents()
}

func confirmDeleteEvent(targetKey string, newJoinString string) {
	list := strings.Split(targetKey, ".")
	table := list[0]
	id := list[1]
	field := list[2]

	var ei EventInfo
	match, _ := regexp.MatchString(table, newJoinString)
	if !match {
		ei.eventType = DeleteTable
		ei.table = table
	} else {
		match, _ := regexp.MatchString(table+"\\."+id, newJoinString)
		if match {
			ei.eventType = DeleteIdentity
			ei.identity = id
			ei.table = table
		} else {
			ei.eventType = DeleteField
			ei.field = field
			ei.table = table
		}
	}

	if !checkEventExists(ei) {
		eventInfoList = append(eventInfoList, ei)
	}
}

func confirmAddEvent(targetKey string, oldJoinString string) {
	list := strings.Split(targetKey, ".")
	table := list[0]
	id := list[1]
	field := list[2]

	var ei EventInfo
	match, _ := regexp.MatchString(table, oldJoinString)
	if !match {
		ei.eventType = AddTable
		ei.table = table
	} else {
		match, _ := regexp.MatchString(table+"\\."+id, oldJoinString)
		if match {
			ei.eventType = AddField
			ei.field = field
			ei.table = table
		} else {
			ei.eventType = AddIdentity
			ei.identity = id
			ei.table = table
		}
	}
	if !checkEventExists(ei) {
		eventInfoList = append(eventInfoList, ei)
	}
}

func checkEventExists(ei EventInfo) bool {
	et := ei.eventType
	for i := 0; i < len(eventInfoList); i++ {
		element := eventInfoList[i]
		if element.eventType == et {
			switch et {
			case AddTable:
				if element.table == ei.table {
					return true
				}
				break
			case DeleteTable:
				if element.table == ei.table {
					return true
				}
				break
			case AddIdentity:
				if element.table == ei.table && element.identity == ei.identity {
					return true
				}
				break
			case DeleteIdentity:
				if element.table == ei.table && element.identity == ei.identity {
					return true
				}
				break
			case AddField:
				if element.table == ei.table && element.field == ei.field {
					return true
				}
				break
			case DeleteField:
				if element.table == ei.table && element.field == ei.field {
					return true
				}
				break
			}
		}
	}
	return false
}

func getOldKeys(contents map[string]string) {
	if len(oldKeys) > 0 {
		return
	}
	for k := range contents {
		oldKeys = append(oldKeys, k)
	}
	allOldKeysJoinString = strings.Join(oldKeys, "")
}

func getNewKeys(contents map[string]string) {
	if len(newKeys) > 0 {
		return
	}
	for k := range contents {
		newKeys = append(newKeys, k)
	}
	allNewKeysJoinString = strings.Join(newKeys, "")
}

func dumpEvents() {
	file, err := os.OpenFile("./result.csv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")

	var comparor = func(i, j int) bool {
		if eventInfoList[i].eventType > eventInfoList[j].eventType {
			return true
		} else if eventInfoList[i].eventType < eventInfoList[j].eventType {
			return false
		} else {
			return eventInfoList[i].table != eventInfoList[j].table
		}
	}

	sort.Slice(eventInfoList, comparor)

	for i := 0; i < len(eventInfoList); i++ {
		var content = ""
		element := eventInfoList[i]
		switch element.eventType {
		case AddTable:
			content = "AddTable" + "\t" + element.table + "\n"
			break
		case DeleteTable:
			content = "DeleteTable" + "\t" + element.table + "\n"
			break
		case AddIdentity:
			content = "AddIdentity" + "\t" + element.table + "\t" + element.identity + "\n"
			break
		case DeleteIdentity:
			content = "DeleteIdentity" + "\t" + element.table + "\t" + element.identity + "\n"
			break
		case AddField:
			content = "AddField" + "\t" + element.table + "\t" + element.field + "\n"
			break
		case DeleteField:
			content = "DeleteField" + "\t" + element.table + "\t" + element.field + "\n"
			break
		case ChangCell:
			content = "ChangCell" + "\t" + element.key + "\t" + element.content + "\t" + element.newContent + "\n"
			break
		}

		file.WriteString(content)
		fmt.Print(content)
	}
}
