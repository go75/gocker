package main

import (
	"fmt"
	"gocker/constant"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func commitContainer(imageName string) {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("commit container failed")
	}
	containerPath := filepath.Join(pwd, constant.WorkerName, constant.ContainerName)
	imageTar := filepath.Join(pwd, imageName+".tar")
	fmt.Println("commitContainer imageTar:", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", containerPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar folder %s error %v", containerPath, err)
	}
}
