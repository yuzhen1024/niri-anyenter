package core

import (
	"bufio"
	"os"
	"os/exec"
	"time"

	"github.com/tidwall/gjson"
)

type Event uint

const (
	WorkspaceChange Event = 1 << iota
	WindowClose
	FirstStart
	WindowFocusNull // {"WindowFocusChanged":{"id":null}}
	Overview
)

type NiriSingle struct {
	event  Event
	hasWin bool
	json   string
}

func ListenNiriIPC(ch chan<- NiriSingle, exit <-chan struct{}, ev Event) {
	if ev&FirstStart != 0 {
		ch <- NiriSingle{
			event:  FirstStart,
			hasWin: HasWin(-1),
		}
	}

	cmd := exec.Command("niri", "msg", "--json", "event-stream")
	stdout, _ := cmd.StdoutPipe()

	// cmd.Start()
	go func() {
		cmd.Run()
	}()
	go func() {
		for range exit {
			cmd.Process.Signal(os.Interrupt)
		}
	}()

	scanner := bufio.NewScanner(stdout)
	pass := true
	time.AfterFunc(300*time.Millisecond, func() {
		pass = false
	})
	for scanner.Scan() {
		if pass {
			continue
		}

		json := scanner.Text()

		result := gjson.Get(json, "WorkspaceActivated")
		if ev&WorkspaceChange != 0 && result.Index != 0 {
			if result.Exists() {
				ch <- NiriSingle{
					event:  WorkspaceChange,
					hasWin: HasWin(result.Get("id").Int()),
				}
			}
		}

		result = gjson.Get(json, "WindowClosed")
		if ev&WindowClose != 0 && result.Index != 0 {
			if result.Exists() {
				ch <- NiriSingle{
					event:  WindowClose,
					hasWin: HasWin(-1),
				}
			}
		}

		result = gjson.Get(json, "WindowFocusChanged.id")
		if ev&WindowFocusNull != 0 && result.Index != 0 {
			if result.Int() == 0 {
				ch <- NiriSingle{
					event:  WindowFocusNull,
					hasWin: HasWin(-1),
				}
			}
		}

		result = gjson.Get(json, "OverviewOpenedOrClosed.is_open")
		if ev&Overview != 0 && result.Index != 0 {
			ch <- NiriSingle{
				event:  Overview,
				hasWin: HasWin(-1),
				json:   json,
			}
		}

	}
}
