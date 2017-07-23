package cachetree

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

func readCommand(conn *net.TCPConn) (byte, error) {
	cmd_data := make([]byte, 1)
	n, err := conn.Read(cmd_data)
	if err != nil {
		return 0, err
	}

	if n != 1 {
		return 0, nil
	}

	return cmd_data[0], nil
}

func readData(conn *net.TCPConn) ([]byte, error) {
	data_len_bytes := make([]byte, 4)
	n, err := conn.Read(data_len_bytes)
	if err != nil {
		return nil, err
	}

	if n != 4 {
		return nil, errors.New("Wrong API communication")
	}

	data_len := int(binary.BigEndian.Uint32(data_len_bytes))
	send_data := bytes.NewBuffer([]byte{})
	tmp_read := make([]byte, data_len)

	for {
		n, err := conn.Read(tmp_read)
		if err != nil {
			return nil, err
		}

		send_data.Write(tmp_read[:n])
		if send_data.Len() < data_len {
			tmp_read = tmp_read[n:]
			continue
		}

		return send_data.Bytes(), nil
	}
}
