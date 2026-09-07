package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"path/filepath"

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
	Uploads struct {
		Size  uint64 `json:"size"`
		Error string `json:"error,omitempty"`
	} `json:"uploads"`
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
		
		// stat.Bavail = ruang kosong yang benar-benar bisa Anda pakai (sama dengan kolom Avail di df -h)
		info.Storage.Free = uint64(stat.Bavail) * uint64(stat.Bsize)
		
		// stat.Bfree = total ruang kosong SEBELUM dikurangi cadangan 5% untuk sistem root Linux
		// Agar perhitungan 'Terpakai' sama dengan perintah 'df -h', kita kurangi Total dengan Bfree.
		freeTermasukRoot := uint64(stat.Bfree) * uint64(stat.Bsize)
		if info.Storage.Total >= freeTermasukRoot {
			info.Storage.Used = info.Storage.Total - freeTermasukRoot
		}
	}

	// Uploads Folder Info
	var dirSize uint64
	var walkErr error
	
	err = filepath.Walk("./uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkErr = err
			return nil // lanjutkan ke file lain meskipun ada error
		}
		if !info.IsDir() {
			dirSize += uint64(info.Size())
		}
		return nil
	})
	
	info.Uploads.Size = dirSize
	if err != nil {
		info.Uploads.Error = err.Error()
	} else if walkErr != nil {
		info.Uploads.Error = walkErr.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}
