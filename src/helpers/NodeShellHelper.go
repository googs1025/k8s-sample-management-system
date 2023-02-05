package helpers

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

var NodeShellModes  = ssh.TerminalModes{
	ssh.ECHO:           1,     // enable echoing
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}
func SSHConnect( user, password, host string, port int ) ( *ssh.Session, error ) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:               user,
		Auth:               auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback:    hostKeyCallbk,
	}
	// connet to ssh
	addr = fmt.Sprintf( "%s:%d", host, port )
	if client, err = ssh.Dial( "tcp", addr, clientConfig ); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}
