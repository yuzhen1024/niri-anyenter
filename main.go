package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

var launcher = "fuzzel"
var launcherLockFile = `/run/user/1000/fuzzel-wayland-1.lock`
var searchMode = true
var firstInputPressTimming = 300 * time.Millisecond
var waitNiriStart = true

func main() {
	// TODO waitingNiriStart

	receiveIPCEvent := make(chan NiriSingle)
	go ListenNiriIPC(receiveIPCEvent, FirstStart|WorkspaceChange|WindowClose|Overview)

	breakSend := make(chan struct{})
	inputResultRec := make(chan string)
	firstInput := make(chan string)
	keyeventMonitorState := false

	var closeKeyeventMonitor = func() {
		// fmt.Println("clear...")
		if keyeventMonitorState {
			breakSend <- struct{}{}
			keyeventMonitorState = false
		}
	}
	var startKeyeventMonitor = func() {
		// fmt.Println("start...")
		if keyeventMonitorState == false {
			go keyeventMonitor(breakSend, inputResultRec, firstInput, firstInputPressTimming)
			keyeventMonitorState = true
		}
	}

	var checkLockFileAndSendEvent = func() {
		isExisted := false
		for true {
			time.Sleep(100 * time.Millisecond)

			if isExisted == false && checkLockFile() == true {

				isExisted = true
			} else if isExisted == true && checkLockFile() == false {
				// log.Println("lock file removed...")
				closeKeyeventMonitor()
				receiveIPCEvent <- NiriSingle{
					event:  WindowClose,
					hasWin: hasWin(-1), // or true
				}
				break
			}
		}
	}
	// searchMode, when press first key
	go func() {
		for v := range firstInput {
			// log.Println("firstInput...")
			if searchMode {
				closeKeyeventMonitor()
				go runLauncher(v)
				checkLockFileAndSendEvent()
			} else {
				runLauncher("")
			}
		}
	}()
	go func() {
		for range <-inputResultRec {
		}
	}()

	isLock := false
	for v := range receiveIPCEvent {
		// fmt.Println("ev: ", EventMapFeild[v.event], "hasWin: ", v.hasWin)

		if v.hasWin == false {
			// when overview close
			if v.event == Overview {
				isLock = true
				go func() {
					// keyevnet start need 200 ms, 200 + 100 = 300ms
					time.Sleep(100 * time.Millisecond)
					isLock = false

					if hasWin(-1) == false {
						startKeyeventMonitor()
					}
				}()
				continue
			} else {
				if isLock {
					continue
				}
				startKeyeventMonitor()
			}

		} else {
			if v.event == Overview {
				isLock = true
			}
			closeKeyeventMonitor()
		}
	}
}

func checkLockFile() bool {
	_, err := os.Stat(launcherLockFile)
	// if os.IsNotExist(err) { return false }
	if err != nil {
		return false
	}
	return true
}

// Window focus changed: None
// Window focus changed: Some(52)
func callNiriIPCWinFocusNone() {
	// go ListenNiriIPC()
}

func runLauncher(searchWord string) {
	args := make([]string, 0)
	if searchWord != "" {
		args = append(args, "--search", searchWord)
	}
	if err := exec.Command(launcher, args...).Run(); err != nil {
		fmt.Println(err)
	}
}

func typeInput(input string) {

}
