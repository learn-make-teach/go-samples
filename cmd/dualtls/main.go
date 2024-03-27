package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) != 2 {
		panic("Expecting certificate dir as first and only argument")
	}
	certDir := os.Args[1]
	err := os.Chdir(certDir)
	if err != nil {
		panic(fmt.Errorf("specified directory is not valid: %v", err))
	}
	rsaCert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		panic(fmt.Errorf("can't load rsa key pair: %v", err))
	}
	ecCert, err := tls.LoadX509KeyPair("server-p256.crt", "server-p256.key")
	if err != nil {
		panic(fmt.Errorf("can't load ecdsa key pair: %v", err))
	}
	config := tls.Config{
		Certificates: []tls.Certificate{ecCert, rsaCert},
	}
	l, err := tls.Listen("tcp", "127.0.0.1:8443", &config)
	if err != nil {
		panic(fmt.Errorf("can't listen: %v", err))
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Printf("accept failed: %v\n", err)
				break
			}
			fmt.Printf("Connection from %s\n", conn.RemoteAddr())
			err = conn.(*tls.Conn).Handshake()
			if err != nil {
				fmt.Printf("handshake failed: %v\n", err)
			}
			conn.Close()
		}
	}()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Listening on %s, press Ctrl+C to stop\n", l.Addr())
	<-done
	l.Close()
}
