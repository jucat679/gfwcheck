package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"

	"fmt"

	cfg "github.com/mritd/gfwcheck/exec"
	"github.com/spf13/viper"
)

func Install() {

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		Uninstall()

		log.Println("Create config dir /etc/gfwcheck")
		os.MkdirAll("/etc/gfwcheck", 0755)

		log.Println("Copy file to /usr/bin")
		currentPath, err := exec.LookPath(os.Args[0])
		CheckAndExit(err)

		currentFile, err := os.Open(currentPath)
		defer currentFile.Close()
		CheckAndExit(err)

		installFile, err := os.OpenFile("/usr/bin/gfwcheck", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		defer installFile.Close()
		CheckAndExit(err)
		_, err = io.Copy(installFile, currentFile)
		CheckAndExit(err)

		log.Println("Create config file /etc/gfwcheck/config.yaml")
		configFile, err := os.Create("/etc/gfwcheck/config.yaml")
		defer configFile.Close()
		CheckAndExit(err)
		viper.SetConfigFile("/etc/gfwcheck/config.yaml")
		viper.Set("Servers", cfg.ExampleConfig())
		CheckAndExit(viper.WriteConfig())

		log.Println("Create systemd config file /lib/systemd/system/gfwcheck.service")
		systemdServiceFile, err := os.Create("/lib/systemd/system/gfwcheck.service")
		defer systemdServiceFile.Close()
		CheckAndExit(err)
		fmt.Fprint(systemdServiceFile, cfg.SystemdConfig)

	} else {
		log.Println("Install not support this platform!")
	}
}
