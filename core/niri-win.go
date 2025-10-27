package core

import (
	"fmt"
	"os/exec"
	"slices"

	"github.com/tidwall/gjson"
)

// TODO *exclude... parts of code, move to other place

// -1 == null
func HasWin(wkspcId int64) (result bool) {

	if wkspcId < 0 {
		wkspcId = getCurrentWorkspaceID()
	}

	output, _ := exec.Command("niri", "msg", "--json", "windows").CombinedOutput()

	gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {

		id := gjson.Get(value.Raw, "workspace_id").Int()

		if id == wkspcId {
			if slices.Contains(
				*excludewindows,
				gjson.Get(value.Raw, "title").String(),
			) == false {
				result = true
				return false
			}
		}
		return true
	})

	return
}

func getCurrentWorkspaceID() (wkspcId int64) {
	// for multipie monitor, and more fast get value
	output, _ := exec.Command("niri", "msg", "--json", "focused-window").CombinedOutput()
	focusedWorkspace := gjson.GetBytes(output, "workspace_id")
	if focusedWorkspace.Index != 0 {
		return focusedWorkspace.Int()
	}

	output, _ = exec.Command("niri", "msg", "--json", "workspaces").CombinedOutput()
	gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {
		if gjson.Get(value.Raw, "is_active").Bool() {
			id := gjson.Get(value.Raw, "id")
			// null == 0
			wkspcId = id.Int()
			return false
		}
		return true
	})
	return wkspcId
}

func MatchLauncherRule() (result bool) {
	r := launcherRuleParsed
	id := getCurrentWorkspaceID()
	output, _ := exec.Command("niri", "msg", "--json", "windows").CombinedOutput()
	gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {
		if gjson.Get(value.Raw, "workspace_id").Int() == id {
			for _, v := range r.window {
				if v == gjson.Get(value.Raw, "title").String() {
					result = true
					return false
				}
			}
		}
		return true
	})
	if result {
		return
	}

	output, _ = exec.Command("niri", "msg", "--json", "layers").CombinedOutput()
	gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {
		for _, v := range r.layer {
			if v == gjson.Get(value.Raw, "namespace").String() {
				result = true
				return false
			}
		}
		return true
	})

	return
}

func MatchExcludeLayer() string {
	result := ""
	output, _ := exec.Command("niri", "msg", "--json", "layers").CombinedOutput()
	gjson.ParseBytes(output).ForEach(func(key, value gjson.Result) bool {
		name := gjson.Get(value.Raw, "namespace").String()
		fmt.Println("debug: name: ", name)
		if slices.Contains(*excludeLayers, name) {
			result = name
			return false
		}
		return true
	})
	return result
}
