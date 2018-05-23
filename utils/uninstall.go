package utils

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func Uninstall() {

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		CheckRoot()

		log.Println("Stop gfwcheck")
		exec.Command("systemctl", "stop", "gfwcheck").Run()

		log.Println("Clean files")
		os.Remove("/usr/bin/gfwcheck")
		os.Remove("/etc/gfwcheck")
		os.Remove("/usr/lib/systemd/system/gfwcheck.service")

		log.Println("Systemd reload")
		exec.Command("systemctl", "daemon-reload").Run()
	} else {
		log.Println("Uninstall not support this platform!")
	}

}
