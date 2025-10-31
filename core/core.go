package core

import (
	"fmt"
	"log"
	"time"

	"github.com/tidwall/gjson"
)

type CoreVar struct {
	ReceiveIPCEvent      chan NiriSingle
	BreakSend            chan struct{}
	InputRec             chan string
	InputDelay           time.Duration
	KeyeventMonitorState bool
}

func Create() *CoreVar {
	return &CoreVar{
		ReceiveIPCEvent:      make(chan NiriSingle),
		BreakSend:            make(chan struct{}),
		InputRec:             make(chan string),
		InputDelay:           time.Millisecond * time.Duration(*preinputDelay),
		KeyeventMonitorState: false,
	}
}

func (v *CoreVar) Run() {
	go ListenNiriIPC(v.ReceiveIPCEvent, make(<-chan struct{}), FirstStart|WorkspaceChange|WindowClose|Overview)

	go func() {
		for input := range v.InputRec {
			if preinputMode {
				log.Println("preinputMode begin start launcher...")
				v.closeKeyeventMonitor()
				matching := MatchExcludeLayer()
				if matching == "" {
					go runLauncher(input)
				}
				v.listenStateAndSendEvent(matching)
			} else {
				log.Panicln("bad flag --mode")
			}
		}
	}()

	isLock := false

	for sig := range v.ReceiveIPCEvent {
		debugIPCEventPrint(sig, isLock)

		if sig.event == Overview {
			isLock = true
			result := gjson.Get(sig.json, "OverviewOpenedOrClosed.is_open").Bool()
			// when overview close, hotkey is break, so it can effect a-z key
			if result == false {
				go func() {
					// overview-off-animation + bin-start-time
					time.Sleep(200 * time.Millisecond)
					log.Println("overview expect animation is timeout, unlock")
					isLock = false

					v.ReceiveIPCEvent <- NiriSingle{
						event:  0,
						hasWin: HasWin(-1),
					}
				}()
			} else {
				v.closeKeyeventMonitor()
			}
			continue
		}

		if sig.hasWin == false {
			if isLock {
				continue
			}
			v.startKeyeventMonitor()
		} else {
			v.closeKeyeventMonitor()
		}
	}
}

func (v *CoreVar) closeKeyeventMonitor() {
	fmt.Println("clear...", " state: ", v.KeyeventMonitorState)
	if v.KeyeventMonitorState {
		v.KeyeventMonitorState = false
		v.BreakSend <- struct{}{}
	}
}

func (v *CoreVar) startKeyeventMonitor() {
	fmt.Println("start...", " state: ", v.KeyeventMonitorState)
	if v.KeyeventMonitorState == false {
		v.KeyeventMonitorState = true
		go KeyeventMonitor(v.BreakSend, v.InputRec, v.InputDelay)
		go v.windowDemain()
	}
}

func (v *CoreVar) windowDemain() {
	nirisig := make(chan NiriSingle)
	exit := make(chan struct{})
	go ListenNiriIPC(nirisig, exit, WorkspaceChange|WindowFocusNull|Overview)

	// some time, you pressed hotkey open the app,
	// niri-ipc is not send signel to event-stream
	done := make(chan struct{})
	go func() {
		for true {
			time.Sleep(200 * time.Millisecond)
			if HasWin(-1) {
				break
			}
		}
		close(done)
	}()

	select {
	case <-nirisig:
		close(exit)

	case <-done:
		v.ReceiveIPCEvent <- NiriSingle{
			event:  0,
			hasWin: true,
		}
		close(exit)
	}
}

func (v *CoreVar) listenStateAndSendEvent(matching string) {
	isStart := false
	for matching == "" {
		time.Sleep(100 * time.Millisecond)
		if isStart == false && MatchLauncherRule() == true {
			isStart = true
		} else if isStart == true && MatchLauncherRule() == false {
			log.Println("launcher closed...")
			v.closeKeyeventMonitor()
			v.ReceiveIPCEvent <- NiriSingle{
				event:  WindowClose,
				hasWin: HasWin(-1),
			}
			return
		}
	}
	for true {
		time.Sleep(100 * time.Millisecond)
		if MatchExcludeLayer() == "" {
			log.Println("exclude layer closed...")
			v.closeKeyeventMonitor()
			v.ReceiveIPCEvent <- NiriSingle{
				event:  WindowClose,
				hasWin: HasWin(-1),
			}
			return
		}
	}
}

// Window focus changed: None
// Window focus changed: Some(52)
// func callNiriIPCWinFocusNone() {
// 	// go ListenNiriIPC()
// }
// func typeInput(input string) { }

func debugIPCEventPrint(sig NiriSingle, isLock bool) {
	eventname := ""
	if sig.event == Overview {
		eventname = "overview"
	} else if sig.event == WorkspaceChange {
		eventname = "WorkspaceChange"
	} else if sig.event == WindowClose {
		eventname = "WindowClose"
	} else if sig.event == FirstStart {
		eventname = "FirstStart"
	} else if sig.event == WindowFocusNull {
		eventname = "WindowFocusNull"
	} else {
		eventname = "Other"
	}
	log.Println("ev: ", eventname, "hasWin: ", sig.hasWin, ", islock: ", isLock)
}
