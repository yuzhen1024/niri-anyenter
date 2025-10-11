package core

import (
	"log"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	launcher             = kingpin.Flag("launcher", "launcher, run your apps something, e.g. fuzzel").Default("fuzzel").String()
	mode                 = kingpin.Flag("mode", "default mode is preinput mode, in future updates include uinput mode(maybe), uinput mode is wip, now only use preinput.").Default("preinput").Enum("preinput")
	launcherPreinputFlag = kingpin.Flag("preinput-flag", "in empty workspace, type allthing start launcher, now need set preinput whatever in this. for fuzzel is --search, for rofi is -filter, you can --help your using launcher for check.").Default("--search").String()
	preinputDelay        = kingpin.Flag("preinput-delay-ms", "input anything run launcher, but how quick you type? need this check if you stop typing.").Default("300").Int64()
	// uid                  = kingpin.Flag("uid", "start launcher use uid, it decide your homedir and process uid from").Default("1000").Uint32()

	isCheckLockfile = kingpin.Flag("check-lockfile", "check lockfile for juge launcher open or close, and use --lockfile-path").Default("true").Bool()
	lockfilePath    = kingpin.Flag("lockfile-path", "use a path like /run/user/1000/fuzzel-wayland-1.lock to jude launcher is open or close, look at --check-lockfile").Default(`/run/user/1000/fuzzel-wayland-1.lock`).String()

	// excludeRule = kingpin.Flag("exclude-rule", `mean juge is has windows for exclude, whilelist, use `+"`niri msg layers`"+` see namespace to check, usage "--whitelist-rule layer=nwg-drawer" "window=Floating\ Window\ -\ Show\ Me\ The\ Key"`).
	// excludeRule = kingpin.Flag("exclude-rule", ``).
	// 	// Default().
	// 	Default("window=Floating Window - Show Me The Key", "layer=nwg-drawer"). //debug
	// 	Strings()

	excludewindows = kingpin.Flag("exclude-windows", "gnore windows if existing, the windows will not stop run launcher.").Default("Floating Window - Show Me The Key").Strings()
	excludeLayers  = kingpin.Flag("exclude-layers", "this will stop launcher run, if layers existing.").Default("nwg-drawer").Strings()

	launcherRule = kingpin.Flag("launcher-rule", "use `sleep 3 && niri msg layers` see namespace to check.\nthe fuzzel use layer, the rofi use window, so set you want.").Default("layer=launcher").Strings()
)

type RuleTypes struct {
	window []string
	layer  []string
}

var excludeRuleParsed RuleTypes
var launcherRuleParsed RuleTypes

var preinputMode = false
var uinputMode = false

func init() {
	kingpin.Parse()

	launcherRuleParsed = parseRule(*launcherRule)

	switch *mode {
	case "preinput":
		preinputMode = true
	case "uinput":
		uinputMode = true
	}
}

func parseRule(s []string) RuleTypes {
	var result RuleTypes

	result.layer = make([]string, 0)
	result.window = make([]string, 0)

	for _, v := range s {
		v = strings.Trim(v, `"`)
		// if value has `=`
		spl := strings.Split(v, "=")
		if len(spl) > 2 {
			temp := make([]string, 2)
			temp[0] = strings.Join(spl[:len(spl)-2], "")
			temp[1] = spl[len(spl)-1]
			spl = temp
		} else if len(spl) != 2 {
			log.Println("bad flag value: ", v)
			continue
		}

		if spl[0] == "layer" {
			result.layer = append(result.layer, spl[1])
		} else if spl[0] == "window" {
			result.window = append(result.window, spl[1])
		}
	}
	return result
}
