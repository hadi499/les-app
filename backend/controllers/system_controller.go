package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
	// Membaca dari /proc/meminfo untuk mendapatkan "MemAvailable" yang lebih akurat daripada "Freeram" di Sysinfo
	// karena Freeram tidak menghitung buff/cache yang sebenarnya bisa dipakai (available).
	fileBytes, err := os.ReadFile("/proc/meminfo")
	if err == nil {
		content := string(fileBytes)
		var total, available uint64

		// Mencari MemTotal
		if idx := strings.Index(content, "MemTotal:"); idx != -1 {
			fmt.Sscanf(content[idx:], "MemTotal: %d", &total)
			info.Memory.Total = total * 1024 // /proc/meminfo dalam KiB
		}

		// Mencari MemAvailable (lebih akurat untuk menunjukkan RAM sisa yang bisa dipakai)
		if idx := strings.Index(content, "MemAvailable:"); idx != -1 {
			fmt.Sscanf(content[idx:], "MemAvailable: %d", &available)
			info.Memory.Free = available * 1024
		} else if idx := strings.Index(content, "MemFree:"); idx != -1 {
			// Fallback jika MemAvailable tidak ada (kernel lama)
			fmt.Sscanf(content[idx:], "MemFree: %d", &available)
			info.Memory.Free = available * 1024
		}

		if info.Memory.Total > 0 && info.Memory.Total >= info.Memory.Free {
			info.Memory.Used = info.Memory.Total - info.Memory.Free
		}
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
