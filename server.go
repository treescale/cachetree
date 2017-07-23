package cachetree

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var (
	HandleServerFailure = func(err error, host string) {}
)

func startCacheServer(host string) error {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	go connectionListener(l, host)
	return nil
}

func connectionListener(l *net.TCPListener, host string) {
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			go HandleServerFailure(err, host)
			return
		}

		go handleCacheConnection(conn, fmt.Sprintf("%s-%s", host, strconv.FormatInt(time.Now().UnixNano(), 10)))
	}
}
