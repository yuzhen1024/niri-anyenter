package main

import "github.com/alecthomas/kingpin/v2"

var (
	launcher             = kingpin.Flag("launcher", "launcher, run your apps something, e.g. fuzzel").Default("fuzzel").String()
	mode                 = kingpin.Flag("mode", "default mode is preinput mode, in future updates include uinput mode(maybe), uinput mode is wip, now only use preinput.").Default("preinput").Enum("preinput")
	launcherPreinputFlag = kingpin.Flag("preinput-flag", "in empty workspace, type allthing start launcher, now need set preinput whatever in this. for fuzzel is --search, for rofi is -filter, you can --help your using launcher for check.").Default("--search").String()
	preinputDelay        = kingpin.Flag("preinput-delay-ms", "input anything run launcher, but how quick you type? need this check if you stop typing.").Default("300").Int64()
	isCheckLockfile      = kingpin.Flag("check-lockfile", "check lockfile for juge launcher open or close, and use --lockfile-path").Default("true").Bool()
	lockfilePath         = kingpin.Flag("lockfile-path", "use a path like /run/user/1000/fuzzel-wayland-1.lock to jude launcher is open or close, look at --check-lockfile").Default(`/run/user/1000/fuzzel-wayland-1.lock`).String()
	// uid                  = kingpin.Flag("uid", "start launcher use uid, it decide your homedir and process uid from").Default("1000").Uint32()
)

func init() {
	kingpin.Parse()
}
