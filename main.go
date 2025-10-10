package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

var launcher = "fuzzel"
var launcherLockFile = `/run/user/1000/fuzzel-wayland-1.lock`
var launcherPreInputFlag = "--search"
var searchMode = true
var firstInputPressTimming = 300 * time.Millisecond
var waitNiriStart = true
var uid uint32 = 1000
var gid uint32 = 100

var u, userLookupErr = user.LookupId(fmt.Sprint(uid))

func init() {
	if userLookupErr != nil {
		fmt.Println("user lookup error: ", userLookupErr)
	}
	getGid, err := strconv.Atoi(u.Gid)
	if err != nil {
		log.Panicln("get gui error")
	}
	gid = uint32(getGid)
	// fmt.Println("uid: ", uid, ", gid: ", gid)
}

func main() {
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
	// for uinput mode, wip...
	go func() {
		for range <-inputResultRec {
		}
	}()

	// filter init recive
	func() {
		for true {
			select {
			case <-time.After(300 * time.Millisecond):
				return
			case <-receiveIPCEvent:
			}
		}
	}()

	isLock := false
	for v := range receiveIPCEvent {
		// log.Println("ev: ", EventMapFeild[v.event], "hasWin: ", v.hasWin, ", islock: ", isLock)

		if v.hasWin == false {
			// when overview close
			if v.event == Overview {
				isLock = true
				go func() {
					// time-be-will = overview-off + bin-start
					time.Sleep(500 * time.Millisecond)
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
		args = append(args, launcherPreInputFlag, searchWord)
	}
	cmd := exec.Command(launcher, args...)

	cmd.Dir = u.HomeDir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
		Credential: &syscall.Credential{
			Uid: uid,
			Gid: gid,
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func typeInput(input string) {

}
