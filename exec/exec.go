package exec

import (
	"os/exec"

	"log"

	"os"

	"strings"

	"sync"

	"github.com/mritd/gfwcheck/alarm"
	"github.com/mritd/gfwcheck/proxy"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var mutex sync.RWMutex

func (server *Server) RemoteExec() bool {
	client, err := server.Connection()
	if err != nil {
		log.Printf("Connect to server [%s] failed!\n", server.Host)
		log.Println(err.Error())
		return false
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		log.Println("Session create failed:", err)
		log.Println(err.Error())
		return false
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(server.RemoteCmd)
	if err != nil {
		log.Printf("Server %s remote command [%s] exec failed!\n", server.Name, server.RemoteCmd)
		log.Println(err.Error())
		return false
	} else {
		log.Printf("Server %s remote command [%s] exec success!\n", server.Name, server.RemoteCmd)
		return true
	}
}

func (server *Server) LocalExec() bool {
	var cmd *exec.Cmd
	localCmd := strings.Fields(server.LocalCmd)
	if len(localCmd) < 1 {
		log.Printf("Local command missing,Server %s\n", server.Name)
		return false
	} else if len(localCmd) == 1 {
		cmd = exec.Command(localCmd[0])
	} else {
		cmd = exec.Command(localCmd[0], localCmd[1:]...)
	}

	err := cmd.Run()
	if err != nil {
		log.Printf("Server %s local command [%s] exec failed!\n", server.Name, server.LocalCmd)
		log.Println(err.Error())
		return false
	} else {
		log.Printf("Server %s local command [%s] exec success!\n", server.Name, server.LocalCmd)
		return true
	}
}

func (server *Server) CheckGFWAndExec() {
	log.Printf("%s checking...\n", server.Name)
	if !proxy.Check(server.Proxy) {
		mutex.Lock()
		server.failedCount++
		mutex.Unlock()
		if server.failedCount >= server.MaxFailed {
			alarm.Alarm(server.Name)
		}
		server.RemoteExec()
		server.LocalExec()
	} else {
		mutex.Lock()
		server.failedCount = 0
		mutex.Unlock()
	}
}

func Start() {
	var servers []Server
	err := viper.UnmarshalKey("servers", &servers)
	if err != nil {
		log.Println("Can't parse server config!")
		return
	}
	c := cron.New()
	for i, _ := range servers {
		x := i
		c.AddFunc(servers[x].Cron, func() {
			servers[x].CheckGFWAndExec()
		})
	}
	c.Start()
	select {}
}
