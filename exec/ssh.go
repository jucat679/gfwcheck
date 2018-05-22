package exec

import (
	"errors"
	"io/ioutil"
	"strconv"

	"log"

	"golang.org/x/crypto/ssh"
)

func (server *ServerConfig) Connection() (*ssh.Client, error) {
	authMethods, err := parseAuthMethods(server)

	if err != nil {
		log.Println("Parse auth methods:", err)
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            server.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         server.Timeout,
	}

	addr := server.Host + ":" + strconv.Itoa(server.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Println("SSH connect error:", err)
		return nil, err
	}
	return client, nil

}

func parseAuthMethods(server *ServerConfig) ([]ssh.AuthMethod, error) {
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

func pemKey(server *ServerConfig) (ssh.AuthMethod, error) {
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
