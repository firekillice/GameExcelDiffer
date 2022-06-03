package main

import (
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

func diffContents(old map[string]string, new map[string]string) {
	for ko, vo := range old {
		if vn, exists := new[ko]; exists {
			if vn == vo {
				continue
			} else {
				eventInfoList = append(eventInfoList, EventInfo{eventType: ChangCell, key: ko, content: vo, newContent: vn})
			}
		} else {
			confirmDeleteEvent(ko, getNewKeys(new))
		}
	}

	for kn, _ := range new {
		if _, exists := old[kn]; !exists {
			confirmAddEvent(kn, getOldKeys(old))
		}
	}

	dumpEvents()
}

func confirmDeleteEvent(targetKey string, keys []string) {
	list := strings.Split(targetKey, ".")
	table := list[0]
	id := list[1]
	field := list[2]

	length := len(keys)
	var joinString = ""
	for i := 0; i < length; i++ {
		if strings.HasPrefix(keys[i], table) {
			joinString += keys[i]
		}
	}

	var ei EventInfo
	if len(joinString) == 0 {
		ei.eventType = DeleteTable
		ei.table = table
	} else {
		match, _ := regexp.MatchString(table+"\\."+id, joinString)
		if match {
			ei.eventType = DeleteIdentity
			ei.identity = id
		} else {
			ei.eventType = DeleteField
			ei.field = field
		}
	}
	if !checkEventExists(ei) {
		eventInfoList = append(eventInfoList, ei)
	}
}

func confirmAddEvent(targetKey string, keys []string) {
	list := strings.Split(targetKey, ".")
	table := list[0]
	id := list[1]
	field := list[2]

	length := len(keys)
	var joinString = ""
	for i := 0; i < length; i++ {
		if strings.HasPrefix(keys[i], table) {
			joinString += keys[i]
		}
	}

	var ei EventInfo
	if len(joinString) == 0 {
		ei.eventType = AddTable
		ei.table = table
	} else {
		match, _ := regexp.MatchString(table+"\\."+id, joinString)
		if match {
			ei.eventType = AddField
			ei.field = field
		} else {
			ei.eventType = AddIdentity
			ei.identity = id
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

func getOldKeys(contents map[string]string) []string {
	if len(oldKeys) > 0 {
		return oldKeys
	}
	for k := range contents {
		oldKeys = append(oldKeys, k)
	}
	return oldKeys
}

func getNewKeys(contents map[string]string) []string {
	if len(newKeys) > 0 {
		return newKeys
	}
	for k := range contents {
		newKeys = append(newKeys, k)
	}
	return newKeys
}

func dumpEvents() {
	file, err := os.OpenFile("./result.csv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")

	sort.Slice(eventInfoList, func(i, j int) bool { return eventInfoList[i].eventType > eventInfoList[j].eventType })
	for i := 0; i < len(eventInfoList); i++ {
		element := eventInfoList[i]
		switch element.eventType {
		case AddTable:
			file.WriteString("AddTable" + "\t" + element.table + "\n")
			break
		case DeleteTable:
			file.WriteString("DeleteTable" + "\t" + element.table + "\n")
			break
		case AddIdentity:
			file.WriteString("AddIdentity" + "\t" + element.table + "\t" + element.identity + "\n")
			break
		case DeleteIdentity:
			file.WriteString("DeleteIdentity" + "\t" + element.table + "\t" + element.identity + "\n")
			break
		case AddField:
			file.WriteString("AddField" + "\t" + element.table + "\t" + element.field + "\n")
			break
		case DeleteField:
			file.WriteString("DeleteField" + "\t" + element.table + "\t" + element.field + "\n")
			break
		case ChangCell:
			file.WriteString("ChangCell" + "\t" + element.key + "\t" + element.content + "\t" + element.newContent + "\n")
			break
		}
	}
}
