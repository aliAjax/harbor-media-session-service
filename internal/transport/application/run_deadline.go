package application

import "net"

func setRunDeadline(conn *net.UDPConn) error {
	conn.SetReadDeadline(nextReadDeadline())
	return nil
}
