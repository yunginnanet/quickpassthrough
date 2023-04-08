package configs

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/klauspost/cpuid/v2"
)

type Path struct {
	MODPROBE   string
	INITRAMFS  string
	ETCMODULES string
	DEFAULT    string
	QUICKEMU   string
	DRACUT     string
	MKINITCPIO string
}

type Config struct {
	bootloader string
	cpuvendor  string
	path       *Path
}

func GetConfigPaths() *Path {
	Paths := &Path{
		MODPROBE:   "config/etc/modprobe.d",
		INITRAMFS:  "config/etc/initramfs-tools",
		ETCMODULES: "config/etc/modules",
		DEFAULT:    "config/etc/default",
		QUICKEMU:   "config/quickemu",
		DRACUT:     "config/etc/dracut.conf.d",
		MKINITCPIO: "config/etc/mkinitcpio.conf",
	}

	return Paths
}

func GetConfig() *Config {
	config := &Config{}
	config.path = GetConfigPaths()

	// Set default value for bootloader
	config.bootloader = "unknown"

	// Detect the bootloader we are using
	getBootloader(config)

	// Detect the cpu vendor
	config.cpuvendor = cpuid.CPU.VendorID.String()

	return config
}

func InitConfigs() {
	config := GetConfig()

	dirs := []string{
		config.path.MODPROBE,
		config.path.INITRAMFS,
		config.path.DEFAULT,
		config.path.DRACUT,
	}

	// Remove old config
	os.RemoveAll("config")

	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)

	// For each directory
	for _, confpath := range dirs {
		// Get the system path
		syspath := syspath_re.ReplaceAllString(confpath, "")

		// If the path exists
		if _, err := os.Stat(syspath); !errors.Is(err, os.ErrNotExist) {
			// Create the directories for our configs
			err := os.MkdirAll(confpath, os.ModePerm)
			errorcheck.ErrorCheck(err)
		}
	}

	files := []string{
		config.path.ETCMODULES,
		config.path.MKINITCPIO,
		fmt.Sprintf("%s/modules", config.path.INITRAMFS),
		fmt.Sprintf("%s/grub", config.path.DEFAULT),
	}

	for _, conffile := range files {
		// Get the system file path
		sysfile := syspath_re.ReplaceAllString(conffile, "")

		// If the file exists
		if _, err := os.Stat(sysfile); !errors.Is(err, os.ErrNotExist) {
			// Create the directories for our configs
			file, err := os.Create(conffile)
			errorcheck.ErrorCheck(err)
			defer file.Close()
		}
	}
}
