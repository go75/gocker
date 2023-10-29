package command

import (
	"fmt"
	"gocker/cgroups"
	"gocker/namespace"
	"gocker/ufs"
	"os"
	"os/exec"
	"time"
)

func Run() error {
	os.Args[1] = "init"

	cmd := exec.Command("/proc/self/exe", os.Args[1:]...)

	namespace.SetNamespace(cmd)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("cmd start fail %s", err)
	}

	containerName := os.Args[2]

	time.Sleep(2 * time.Second)

	err = cgroups.ConfigDefaultCgroups(cmd.Process.Pid, containerName)
	if err != nil {
		fmt.Printf("config cgroup fail %s\n", err)
	}

	cmd.Wait()
	cgroups.CleanCgroupsPath(containerName)
	ufs.DelMntNamespace(containerName)
	return nil
}
