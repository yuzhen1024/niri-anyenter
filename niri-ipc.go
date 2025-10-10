package main

import (
	"bufio"
	"os/exec"

	"github.com/tidwall/gjson"
)

type Event uint

// TODO map
const (
	WorkspaceChange Event = 1 << iota
	WindowClose
	FirstStart
	// {"WindowFocusChanged":{"id":null}}
	WindowFocusNull
	Overview
)

var EventMapFeild = map[Event]string{ // TODO typeof 更快
	FirstStart:      "firstStart",
	WorkspaceChange: "workspace-change",
	WindowClose:     "window-close",
	WindowFocusNull: "window-focus-null",
	Overview:        "overview",
}

type NiriSingle struct {
	event  Event
	hasWin bool
}

func ListenNiriIPC(ch chan<- NiriSingle, ev Event) {
	if ev&FirstStart != 0 {
		ch <- NiriSingle{
			event:  FirstStart,
			hasWin: hasWin(-1),
		}
	}

	cmd := exec.Command("niri", "msg", "--json", "event-stream")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		json := scanner.Text()
		// fmt.Println(json)

		result := gjson.Get(json, "WorkspaceActivated")
		if ev&WorkspaceChange != 0 && result.Index != 0 {
			if result.Exists() {
				ch <- NiriSingle{
					event:  WorkspaceChange,
					hasWin: hasWin(result.Get("id").Int()),
				}
			}
		}

		result = gjson.Get(json, "WindowClosed")
		if ev&WindowClose != 0 && result.Index != 0 {
			if result.Exists() {
				ch <- NiriSingle{
					event:  WindowClose,
					hasWin: hasWin(-1),
				}
			}
		}

		result = gjson.Get(json, "WindowFocusChanged.id")
		if ev&WindowFocusNull != 0 && result.Index != 0 {
			if result.Int() == 0 {
				ch <- NiriSingle{
					event:  WindowFocusNull,
					hasWin: hasWin(-1),
				}
			}
		}

		result = gjson.Get(json, "OverviewOpenedOrClosed.is_open")
		if ev&Overview != 0 && result.Index != 0 {
			hasWinResult := false
			if result.Bool() { // open
				hasWinResult = true
			} else {
				hasWinResult = hasWin(-1)
			}
			ch <- NiriSingle{
				event:  Overview,
				hasWin: hasWinResult,
			}
		}

	}
}

// -1 == null
func hasWin(wkspcId int64) (result bool) {
	if wkspcId < 0 {
		output, _ := exec.Command("niri", "msg", "--json", "workspaces").CombinedOutput()
		gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {
			if gjson.Get(value.Raw, "is_active").Bool() {
				id := gjson.Get(value.Raw, "id")
				wkspcId = id.Int()
				return false
			}
			return true
		})
	}
	output, _ := exec.Command("niri", "msg", "--json", "windows").CombinedOutput()
	gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {
		id := gjson.Get(value.Raw, "workspace_id").Int()
		if id == wkspcId {
			result = true
			return false
		}
		return true
	})
	return
}
