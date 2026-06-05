package ssh

import (
	"fmt"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type Client struct {
	*gossh.Client
}

func Dial(host string, port int, username, authType, credential string) (*Client, error) {
	config := &gossh.ClientConfig{
		User:            username,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	switch authType {
	case "password":
		config.Auth = []gossh.AuthMethod{gossh.Password(credential)}
	case "key":
		signer, err := gossh.ParsePrivateKey([]byte(credential))
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		config.Auth = []gossh.AuthMethod{gossh.PublicKeys(signer)}
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", authType)
	}
	client, err := gossh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client}, nil
}
