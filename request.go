package cachetree

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

func AskFileFromMembers(filename string) (data []byte) {
	conns := clientConnections[:]
	var err error
	for host, conn := range conns {
		err = requestFile(conn, filename)
		if err != nil {
			closeConn(host)
			go HandleMemberConnectionFail(err, host)
			continue
		}

		data, err = readFileIfExists(conn, filename)
		if err != nil {
			closeConn(host)
			go HandleMemberConnectionFail(err, host)
			continue
		}

		if data != nil {
			break
		}
	}

	return data
}

func DeleteFileFromMembers(filename string) {
	conns := clientConnections[:]
	for host, conn := range conns {
		err := requestFileDelete(conn, filename)
		if err != nil {
			closeConn(host)
			go HandleMemberConnectionFail(err, host)
			continue
		}

		file_name_data, err := readData(conn)
		if err != nil ||
			// Got wrong API
			string(file_name_data) != filename {
			closeConn(host)
			go HandleMemberConnectionFail(err, host)
			continue
		}
	}
}

func requestFile(conn *net.TCPConn, filename string) error {
	filename_len_bytes := make([]byte, 4)
	filename_bytes := []byte(filename)
	binary.BigEndian.PutUint32(filename_len_bytes, uint32(len(filename_bytes)))

	send_bytes := bytes.NewBuffer([]byte{})
	send_bytes.Write([]byte{CMD_REQUEST_FILE})
	send_bytes.Write(filename_len_bytes)
	send_bytes.Write(filename_bytes)

	_, err := io.Copy(conn, send_bytes)
	return err
}

func requestFileDelete(conn *net.TCPConn, filename string) error {
	filename_len_bytes := make([]byte, 4)
	filename_bytes := []byte(filename)
	binary.BigEndian.PutUint32(filename_len_bytes, uint32(len(filename_bytes)))

	send_bytes := bytes.NewBuffer([]byte{})
	send_bytes.Write([]byte{CMD_DELETE_FILE})
	send_bytes.Write(filename_len_bytes)
	send_bytes.Write(filename_bytes)

	_, err := io.Copy(conn, send_bytes)
	return err
}

func readFileIfExists(conn *net.TCPConn, filename string) ([]byte, error) {
	data, err := readData(conn)
	if err != nil {
		return nil, err
	}

	// if data length is 1, then this means connection don't have requested cache file
	if len(data) == 1 {
		return nil, nil
	}

	return data, nil
}
