package simpletrace

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net"
	"strconv"
)

func randomID(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println(err)
		return ""
	}
	return hex.EncodeToString(bytes)
}

func (s *Service) separateAddresses(address string) (err error) {
	// split address to ip/port
	host, p, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	// convert separated port to an int
	s.Port, err = strconv.Atoi(p)
	if err != nil {
		return err
	}
	// check if IP is an valid IP address
	ip := net.ParseIP(host)
	if ip == nil {
		return err
	}

	// check if IP is an IPv4 address; assign to correct field
	switch ip.To4() {
	case nil:
		s.IPv6 = ip
	default:
		s.IPv4 = ip
	}
	return nil
}
