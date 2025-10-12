package core

import (
	"fmt"
	"log"
	"time"
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

		if sig.hasWin == false {
			// when overview close
			if sig.event == Overview {
				isLock = true
				go func() {
					// overview-off-animation + bin-start-time
					time.Sleep(500 * time.Millisecond)
					log.Println("overview expect animation is timeout, unlock")
					isLock = false

					v.ReceiveIPCEvent <- NiriSingle{
						event:  0,
						hasWin: HasWin(-1),
					}
				}()
				continue
			} else {
				if isLock {
					continue
				}
				v.startKeyeventMonitor()
			}

		} else {
			if sig.event == Overview {
				isLock = true
			}
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

// some time, you pressed hotkey open the app,
// niri-ipc is not send signel to event-stream
func (v *CoreVar) windowDemain() {
	nirisig := make(chan NiriSingle)
	exit := make(chan struct{})
	go ListenNiriIPC(nirisig, exit, WorkspaceChange|WindowFocusNull)
	for true {
		time.Sleep(10 * time.Millisecond)
		select {
		case <-nirisig:
			exit <- struct{}{}
			return
		default:
			if HasWin(-1) {
				v.ReceiveIPCEvent <- NiriSingle{
					event:  0,
					hasWin: true,
				}
				return
			}
		}
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
