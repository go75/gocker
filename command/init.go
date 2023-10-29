package command

import (
	"fmt"
	"syscall"
	"os"
)

func Init(containerName string, cmd string, args []string) (err error) {
	// 改变当前工作目录
	syscall.Chdir("/")

	// 设置挂载
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return fmt.Errorf("mount proc fail %s", err)
	}

	// 执行命令
	err = syscall.Exec(cmd, args, os.Environ())
	if err != nil {
		return fmt.Errorf("exec proc fail %s", err)
	}

	return nil
}
