package core

import (
	"fmt"
	"log"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

var u *user.User
var uid uint32
var gid uint32

func init() {
	var err error
	u, err = user.Current()
	if err != nil {
		fmt.Println(err)
	}
	uidConv, err := strconv.Atoi(u.Uid)
	if err != nil {
		fmt.Println(err)
	}
	gidConv, err := strconv.Atoi(u.Gid)
	if err != nil {
		fmt.Println(err)
	}
	uid = uint32(uidConv)
	gid = uint32(gidConv)
	// fmt.Println("uid: ", uid, ", gid: ", gid)
}

func runLauncher(searchWord string) {
	log.Println("launcher start...")

	args := make([]string, 0)
	if searchWord != "" {
		args = append(args, *launcherPreinputFlag, searchWord)
	}
	cmd := exec.Command(*launcher, args...)

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
