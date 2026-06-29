package controllers

import (
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
)

type SystemInfo struct {
	Memory struct {
		Total uint64 `json:"total"`
		Free  uint64 `json:"free"`
		Used  uint64 `json:"used"`
	} `json:"memory"`
	Storage struct {
		Total uint64 `json:"total"`
		Free  uint64 `json:"free"`
		Used  uint64 `json:"used"`
	} `json:"storage"`
}

// GetSystemInfo mengembalikan informasi memori dan storage dari server (VPS)
func GetSystemInfo(c *gin.Context) {
	var info SystemInfo

	// Memory Info
	var sysInfo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&sysInfo); err == nil {
		// uint32 di beberapa arsitektur, sehingga di cast ke uint64
		info.Memory.Total = uint64(sysInfo.Totalram) * uint64(sysInfo.Unit)
		info.Memory.Free = uint64(sysInfo.Freeram) * uint64(sysInfo.Unit)
		info.Memory.Used = info.Memory.Total - info.Memory.Free
	}

	// Storage Info (root path "/")
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err == nil {
		info.Storage.Total = uint64(stat.Blocks) * uint64(stat.Bsize)
		info.Storage.Free = uint64(stat.Bavail) * uint64(stat.Bsize)
		info.Storage.Used = info.Storage.Total - info.Storage.Free
	}

	c.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}
