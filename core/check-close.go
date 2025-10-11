package core

import (
	"os"
)

func checkLockFile() bool {
	_, err := os.Stat(*lockfilePath)
	// if os.IsNotExist(err) { return false }
	if err != nil {
		return false
	}
	return true
}
