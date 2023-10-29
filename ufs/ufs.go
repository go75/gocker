package ufs

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const (
	// 合并后逻辑上的完整文件系统
	mntLayerPath = "/root/gocker/mnt"
	// 工作目录
	workLayerPath = "/root/gocker/work"
	// 叠加层
	overLayerPath = "/root/gocker/over"
	// 只读的根基础文件系统
	imagePath = "ubuntu-base-23.04-base-amd64"
	oldPath   = ".old"
)

// 获取容器中work目录路径
func workerLayer(containerName string) string {
	return fmt.Sprintf("%s/%s", workLayerPath, containerName)
}

// 获取容器中逻辑根文件系统路径
func mntLayer(containerName string) string {
	return fmt.Sprintf("%s/%s", mntLayerPath, containerName)
}

// 获取容器中over目录路径
func overLayer(containerName string) string {
	return fmt.Sprintf("%s/%s", overLayerPath, containerName)
}

func mntOldLayer(containerName string) string {
	return fmt.Sprintf("%s/%s", mntLayer(containerName), oldPath)
}

func delMntNamespace(path string) error {
	_, err := exec.Command("umount", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("umount fail path=%s err=%s", path, err)
	}

	err = os.RemoveAll(path)

	if err != nil {
		return fmt.Errorf("remove dir fail path=%s err=%s", path, err)
	}
	return err
}

func DelMntNamespace(containerName string) error {
	err := delMntNamespace(mntLayer(containerName))
	if err != nil {
		return err
	}

	err = delMntNamespace(workerLayer(containerName))
	if err != nil {
		return err
	}

	err = delMntNamespace(overLayer(containerName))
	if err != nil {
		return err
	}

	return nil
}

func SetMntNamespace(containerName string) error {

	containerMntPath := mntLayer(containerName)

	containerWorkPath := workerLayer(containerName)
	
	containerOverPath := overLayer(containerName)

	containerMntOldPath := mntOldLayer(containerName)

	// 创建挂载点
	err := os.MkdirAll(containerMntPath, 0700)
	if err != nil {
		return err
	}

	// 创建工作目录
	err = os.MkdirAll(containerWorkPath, 0700)
	if err != nil {
		return err
	}

	// 创建叠加目录
	err = os.MkdirAll(containerOverPath, 0700)
	if err != nil {
		return err
	}

	// 挂载overlay的联合文件系统
	err = syscall.Mount("overlay", containerMntPath, "overlay", 0, 
						fmt.Sprintf("upperdir=%s,lowerdir=%s,workdir=%s", containerOverPath, imagePath, containerWorkPath))
	if err != nil {
		return fmt.Errorf("mount overlay fail err=%s", err)
	}

	// 设置为私有文件系统
	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("reclare rootfs private fail err=%s", err)
	}

	// 挂载
	err = syscall.Mount(containerMntPath, containerMntPath, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("mount rootfs in new mnt space fail err=%s", err)
	}

	
	err = os.MkdirAll(containerMntOldPath, 0700)
	if err != nil {
		return fmt.Errorf("mkdir mnt old layer fail err=%s", err)
	}

	// 设置私有的根文件系统
	err = syscall.PivotRoot(containerMntPath, containerMntOldPath)
	if err != nil {
		return fmt.Errorf("pivot root  fail err=%s", err)
	}
	
	return nil
}
