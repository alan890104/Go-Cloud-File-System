package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

type DiskSpace struct {
	FreeByte  uint64 `json:"free"`
	TotalByte uint64 `json:"total"`
}

type Configuration struct {
	ADMIN_GUUID string `json:"ADMIN_GUUID"`
	EMAIL       string `json:"EMAIL"`
	EMAIL_TOKEN string `json:"EMAIL_TOKEN"`
}

func Disk_Space(path string) (*DiskSpace, error) {
	if runtime.GOOS == "windows" {
		h := syscall.MustLoadDLL("kernel32.dll")
		c := h.MustFindProc("GetDiskFreeSpaceExW")
		ds := &DiskSpace{}
		utf16, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return nil, err
		}
		c.Call(
			uintptr(unsafe.Pointer(utf16)),
			uintptr(unsafe.Pointer(&ds.FreeByte)),
			uintptr(unsafe.Pointer(&ds.TotalByte)),
		)
		return ds, nil
	} else {
		// fs := syscall.Statfs_t{}
		// err := syscall.Statfs(path, &fs)
		// if err != nil {
		// 	return nil, err
		// }
		// ds := &DiskSpace{}
		// &ds.TotalByte = fs.Blocks * uint64(fs.Bsize)
		// &ds.FreeByte = fs.Bfree * uint64(fs.Bsize)
		// return ds, nil
		return nil, errors.New("this function is windows only")
	}
}

func ByteFormat(num uint64) string {
	const unit = 1024
	if num < unit {
		return fmt.Sprintf("%d Bytes", num)
	}
	div, exp := uint64(unit), 0
	for n := num / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(num)/float64(div), "KMGTPEZYBNDCX"[exp])
}

func SpaceUsed(ds *DiskSpace) uint64 {
	return uint64(ds.TotalByte) - uint64(ds.FreeByte)
}

func SpaceUsedPercent(ds *DiskSpace) float32 {
	return float32(SpaceUsed(ds)) / float32(ds.TotalByte)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func create_Folder_If_Path_Not_Exist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func create_default_config() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		log.Println("未找到config.json，將以預設值創建。")
		config := Configuration{
			ADMIN_GUUID: "administrator",
			EMAIL:       "administrator@mail",
			EMAIL_TOKEN: "administrator",
		}
		file, err := json.MarshalIndent(config, "", " ")
		if err != nil {
			log.Fatal("Error when MarshalIndent ", err)
		}
		err = ioutil.WriteFile("config.json", file, 0644)
		if err != nil {
			log.Fatal("Error when write default config.json", err)
		}
	} else {
		log.Println("成功加載config.json")
	}
}

// Use when loading global var
func InitConfig() Configuration {
	create_default_config()
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error when load config.json", err)
	}
	config := Configuration{}
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		log.Fatal("Error when Unmarshal", err)
	}
	essential_path := [...]string{"trash", "uploads"}
	for _, path := range essential_path {
		create_Folder_If_Path_Not_Exist(path)
	}
	return config
}
