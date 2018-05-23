package utils

import (
	"log"
	"os"
	"os/user"
)

func CheckErr(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CheckAndExit(err error) {
	if !CheckErr(err) {
		os.Exit(1)
	}
}

func CheckRoot() {
	u, err := user.Current()
	CheckAndExit(err)

	if u.Uid != "0" || u.Gid != "0" {
		log.Println("This command must be run as root! (sudo)")
		os.Exit(1)
	}
}
