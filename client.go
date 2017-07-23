package cachetree

import (
	"net"
	"sync"
	"time"
)

var (
	HandleMemberConnectionFail = func(err error, host string) {}
	clientConnections          = make(map[string]*net.TCPConn)
	connectionsLocker          = sync.Mutex{}
)

func memberConnector(targets ...string) {
	for {
		for _, target := range targets {

			if _, ok := clientConnections[target]; ok {
				continue
			}

			connectToMember(target)
		}

		time.Sleep(time.Second * 5)
	}
}

func connectToMember(target string) {
	addr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		go HandleMemberConnectionFail(err, target)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		go HandleMemberConnectionFail(err, target)
		return
	}

	connectionsLocker.Lock()
	clientConnections[target] = conn
	connectionsLocker.Unlock()
}
