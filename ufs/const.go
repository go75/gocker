package ufs

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