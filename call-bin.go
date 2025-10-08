package main

import (
	"bufio"
	_ "embed"
	"log"
	"main/utils/keycode"
	"os"
	"strings"

	"codeberg.org/msantos/embedexe/exec"
	"github.com/tidwall/gjson"
)

//go:embed bin/keyevent-monitor
var bin []byte

var pressedModifier = make([]int64, 0)
var pressedShift bool

func initPressed() {
	pressedModifier = make([]int64, 0)
	pressedShift = false
}

func keyeventMonitor(breakReceive <-chan struct{}, result chan<- string, firstPress chan<- string) {
	initPressed()
	inputs := ""

	// cmd exe, drop 200 ms
	cmd := exec.Command(bin)
	// defer func() {
	// 	fmt.Println("defer...")
	// 	cmd.Process.Signal(os.Interrupt)
	// 	cmd.Exe = nil
	// 	cmd = nil
	// }()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}

	// do not use cmd.Start
	// this will can not free money
	go func() {
		// fmt.Println("run...")
		cmd.Run()
	}()

	var callBreak = func() {
		// fmt.Println("defer...")
		/// go for fast
		go cmd.Process.Signal(os.Interrupt)
		result <- inputs
	}
	go func() {
		for range breakReceive {
			callBreak()
			break
		}
	}()
	// TODO Active call break
	// apply to debounce
	// go func() {}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		key, ignore := parseKey(strings.Trim(scanner.Text(), "\n"))
		// fmt.Println("get key: ", key, " isignore: ", ignore)
		if ignore {
			continue
		}
		// log.Println("begin send", ", inputs: ", inputs, ", len: ", len(inputs))
		if inputs == "" {
			firstPress <- key
		}
		inputs += key
	}
	// fmt.Println("break...")
	// cmd.Process.Signal(os.Interrupt)
	// cmd.Process.Kill()
}

func parseKey(json string) (letter string, ignore bool) {
	// fmt.Println("parseKey()...")
	ignore = true
	// TODO 组合键缓冲？ 300ms
	// debunce := debounce.NewDebounce(time.Millisecond * 300)
	// TODO caps lock

	key := gjson.Get(json, "key").Int()
	state := gjson.Get(json, "state").Int()
	stateBool := false
	if state == 1 {
		stateBool = true
	} else {
		stateBool = false
	}

	// lshift rshift
	if key == 42 || key == 54 {
		pressedShift = stateBool
		return
	}

	val, ok := keycode.Letters[key]
	// fmt.Println("val: ", val, " isok: ", ok, " pressedModifier: ", pressedModifier)
	if ok {
		// super + w, 完成组合键后，桌面判断 super 已松开
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
		}
		return val, false
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
