package network

import "net"

type TcpClient struct {
	*TcpConn
}

func (c *TcpClient) Connect(addr string) error {
	conn, err := tcpDial(addr)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}


func tcpDial(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
