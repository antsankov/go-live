package lib

import (
	"net"

	externalip "github.com/glendc/go-external-ip"
)

// GetLocalIP returns preferred local outbound ip of this machine.
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GetExternalIP returns the internet address of the machine.
func GetExternalIP() (string, error) {
	// Create the default consensus,
	// using the default configuration and no logger.
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return "", err
	}
	return ip.String(), nil // print IPv4/IPv6 in string format
}
