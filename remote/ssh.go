package remote

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"fmt"

	"golang.org/x/crypto/ssh"
)

type Server struct {
	Name     string `yml:"name"`
	Host     string `yml:"ip"`
	Port     int    `yml:"port"`
	User     string `yml:"user"`
	Password string `yml:"password"`
	Method   string `yml:"method"`
	Key      string `yml:"key"`
}

func (server *Server) Connection() (*ssh.Session, error) {
	authMethods, err := parseAuthMethods(server)

	if err != nil {
		fmt.Println("parse auth methods:", err)
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: server.User,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := server.Host + ":" + strconv.Itoa(server.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("ssh connect error:", err)
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("session create failed:", err)
		return nil, err
	}

	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	return session, nil

}

func parseAuthMethods(server *Server) ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	switch server.Method {
	case "password":
		authMethods = append(authMethods, ssh.Password(server.Password))
		break

	case "pem":
		method, err := pemKey(server)
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, method)
		break

	default:
		return nil, errors.New("invalid auth method: " + server.Method)
	}

	return authMethods, nil
}

func pemKey(server *Server) (ssh.AuthMethod, error) {
	pemBytes, err := ioutil.ReadFile(server.Key)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if server.Password == "" {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(server.Password))
	}

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

func ExampleConfig() []*Server {
	return []*Server{
		{
			Name:     "test1",
			Host:     "test1.com",
			Port:     22,
			User:     "root",
			Password: "test123",
			Key:      "",
			Method:   "password",
		},
		{
			Name:     "test2",
			Host:     "test2.com",
			Port:     22,
			User:     "root",
			Password: "",
			Key:      "~/.ssh/id_rsa",
			Method:   "pem",
		},
	}
}
