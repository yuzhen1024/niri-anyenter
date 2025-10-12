package core

import (
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	launcher             = kingpin.Flag("launcher", "launcher or runner, run your apps, e.g. fuzzel").Default("fuzzel").String()
	mode                 = kingpin.Flag("mode", "default mode is preinput mode, in future updates include uinput mode(maybe), uinput mode is wip, now only use preinput.").Default("preinput").Enum("preinput")
	launcherPreinputFlag = kingpin.Flag("preinput-flag", "in empty workspace, type allthing start launcher, now need set for preinput whatever in this. for fuzzel is --search, for rofi is -filter, you can `--help` your using launcher for check.").Default("--search").String()
	preinputDelay        = kingpin.Flag("preinput-delay-ms", "input anything run launcher, but how quick can you type? need the this value to check if you stopped typing.").Default("400").Int64()
	// uid                  = kingpin.Flag("uid", "start launcher use uid, it decide your homedir and process uid from").Default("1000").Uint32()

	// excludewindows = kingpin.Flag("exclude-windows", `gnore windows if existing, the windows will not stop run launcher. example: --exclude-windows "app1" "app2" ...`).Default("Floating Window - Show Me The Key").Strings()
	excludewindows = kingpin.Flag("exclude-windows", `gnore windows if existing, the windows will not stop run launcher. example: --exclude-windows "app1" "app2" ...`).Default("").Strings()
	excludeLayers  = kingpin.Flag("exclude-layers", `this will stop launcher run, if layers existing. example: --exclude-layers "app1" "app2" ...`).Default("nwg-drawer").Strings()
	launcherRule   = kingpin.Flag("launcher-rule", `use `+"`sleep 3 && niri msg layers`"+`see namespace to check. the fuzzel use layer, the rofi use window, so, set you want. default value is `+"`layer=launcher`"+`, fit for fuzzel.`).Default("layer=launcher").String()
)

type RuleTypes struct {
	window []string
	layer  []string
}

var launcherRuleParsed RuleTypes

var preinputMode = false
var uinputMode = false

func init() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("v0.2.3")
	// kingpin.CommandLine.Help = ""
	kingpin.Parse()

	launcherRuleParsed = parseRule(*launcherRule)

	switch *mode {
	case "preinput":
		preinputMode = true
	case "uinput":
		uinputMode = true
	}
}

func parseRule(v string) RuleTypes {
	var result RuleTypes

	result.layer = make([]string, 0)
	result.window = make([]string, 0)

	v = strings.Trim(v, `"`)
	// if value has `=`
	spl := strings.Split(v, "=")
	if len(spl) > 2 {
		temp := make([]string, 2)
		temp[0] = strings.Join(spl[:len(spl)-2], "=")
		temp[1] = spl[len(spl)-1]
		spl = temp
	} else if len(spl) != 2 {
		log.Println("bad flag value: ", v)
		os.Exit(1)
	}

	if spl[0] == "layer" {
		result.layer = append(result.layer, spl[1])
	} else if spl[0] == "window" {
		result.window = append(result.window, spl[1])
	}
	return result
}
