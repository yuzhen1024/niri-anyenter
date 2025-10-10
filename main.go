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

var u *user.User
var gid uint32

var preinputMode = false
var uinputMode = false

func init() {
	var err error
	u, err = user.LookupId(fmt.Sprint(*uid))
	if err != nil {
		fmt.Println("user lookup error: ", err)
	}
	getGid, err := strconv.Atoi(u.Gid)
	if err != nil {
		log.Panicln("get gui error")
	}
	gid = uint32(getGid)
	// fmt.Println("uid: ", uid, ", gid: ", gid)

	switch *mode {
	case "preinput":
		preinputMode = true
	case "uinput":
		uinputMode = true
	}
}

func main() {
	receiveIPCEvent := make(chan NiriSingle)
	go ListenNiriIPC(receiveIPCEvent, FirstStart|WorkspaceChange|WindowClose|Overview)

	breakSend := make(chan struct{})
	inputRec := make(chan string)
	inputDelay := time.Millisecond * time.Duration(*preinputDelay)

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
			go keyeventMonitor(breakSend, inputRec, inputDelay)
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
					hasWin: HasWin(-1), // or true
				}
				break
			}
		}
	}
	go func() {
		for v := range inputRec {
			if preinputMode {
				closeKeyeventMonitor()
				go runLauncher(v)
				checkLockFileAndSendEvent()
			} else {
				log.Panicln("bad flag --mode")
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
					// overview-off-animation + bin-start-time
					time.Sleep(500 * time.Millisecond)
					isLock = false

					if HasWin(-1) == false {
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
	_, err := os.Stat(*lockfilePath)
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
		args = append(args, *launcherPreinputFlag, searchWord)
	}
	cmd := exec.Command(*launcher, args...)

	cmd.Dir = u.HomeDir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
		Credential: &syscall.Credential{
			Uid: *uid,
			Gid: gid,
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

// func typeInput(input string) { }
