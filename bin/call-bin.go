package bin

import (
	"bufio"
	_ "embed"
	"log"
	"niri-anyenter/utils/keycode"
	"os"
	"strings"
	"time"

	"codeberg.org/msantos/embedexe/exec"
	"github.com/tidwall/gjson"
)

//go:embed keyevent-monitor
var bin []byte

var pressedModifier = make([]int64, 0)
var pressedShift = false

func clearPressed() {
	pressedModifier = make([]int64, 0)
	pressedShift = false
}

func KeyeventMonitor(
	returnReceiver <-chan struct{},
	inputSender chan<- string,
	inputSendDelay time.Duration,
) {
	clearPressed()

	cmd := exec.Command(bin)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}

	// do not use cmd.Start
	// They will can not free money
	go func() {
		cmd.Run()
	}()

	inputs := ""
	isBreak := false

	var callBreak = func() {
		go cmd.Process.Signal(os.Interrupt)
		log.Println("closed keyevent-monitor...")
		// go cmd.Process.Signal(os.Kill)
	}
	go func() {
		for range returnReceiver {
			log.Println("begin close keyevent-monitor...")
			isBreak = true
			callBreak()
			break
		}
	}()

	isSendFirstInput := false
	isSending := false
	flashSendFirstInput := make(chan struct{})
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		key, ignore := parseKey(strings.Trim(scanner.Text(), "\n"))
		if ignore {
			if (len(pressedModifier) == 0 && pressedShift == false) == false {
				flashSendFirstInput <- struct{}{}
			}
			continue
		}
		if isSendFirstInput == false {
			if isSending {
				flashSendFirstInput <- struct{}{}
			}
			isSending = true
			go func() {
				select {
				case <-flashSendFirstInput:
					return
				case <-time.After(inputSendDelay):
					if isBreak == false {
						isSendFirstInput = false
						inputSender <- inputs
					}
				}
			}()
		}
		inputs += key
	}
}

func parseKey(json string) (letter string, ignore bool) {
	ignore = true

	key := gjson.Get(json, "key").Int()
	state := gjson.Get(json, "state").Int()

	stateBool := false
	if state == 1 {
		stateBool = true
	} else {
		stateBool = false
	}

	// lshift || rshift
	if key == 42 || key == 54 {
		pressedShift = stateBool
		return
	}

	// TODO caps lock
	val, ok := keycode.Letters[key]
	// fmt.Println("val: ", val, " isok: ", ok, " pressedModifier: ", pressedModifier)
	if ok {
		for _, v := range pressedModifier {
			if v != 0 {
				return
			}
		}
		if stateBool == false {
			return
		}

		if pressedShift {
			return strings.ToUpper(val), false
		} else {
			return val, false
		}
	} else {
		if stateBool {
			pressedModifier = append(pressedModifier, key)
		} else {
			temp := make([]int64, len(pressedModifier))
			for _, v := range pressedModifier {
				if v == key {
					continue
				}
				temp = append(temp, v)
			}
			pressedModifier = temp
		}
	}
	return
}
