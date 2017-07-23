package cachetree

import (
	"net"
)

var (
	HandleServerFailure           = func(err error, host string) {}
	HandleClientConnectionFailure = func(err error, host string) {}
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

		go handleServerConnection(conn)
	}
}

func handleServerConnection(conn *net.TCPConn) {
	defer conn.Close()

	for {
		command, err := readCommand(conn)
		if err != nil {
			go HandleClientConnectionFailure(err, conn.RemoteAddr().String())
			return
		}

		switch command {
		case CMD_REQUEST_FILE:
			filename_bytes, err := readData(conn)
			if err != nil {
				go HandleClientConnectionFailure(err, conn.RemoteAddr().String())
				return
			}

			err = sendFileIfExists(conn, string(filename_bytes))
			if err != nil {
				go HandleClientConnectionFailure(err, conn.RemoteAddr().String())
				return
			}
		}
	}
}
