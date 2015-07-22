package locking

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

func InstanceLocked(lockfileName string) bool {
	_, err := os.Stat(lockfileName)
	if err == nil {
		var pidBytes []byte
		pidBytes, err = ioutil.ReadFile(lockfileName)
		pid, _ := strconv.Atoi(string(pidBytes))
		err = syscall.Kill(pid, 0)
	}

	return err == nil
}

func LockInstance(lockfileName string) {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(lockfileName, []byte(pid), 0644); err != nil {
		panic(err)
	}
}

func UnlockInstance(lockfileName string) {
	syscall.Unlink(lockfileName)
}
