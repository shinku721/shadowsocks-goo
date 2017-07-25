package shadowsocks

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

// IsIPv6 checks whether an address is IPv6.
// It does not guarantee that the address is valid,
// so you should only check it on an IP address or
// a hostname.
func IsIPv6(host string) bool {
	if r := strings.Index(host, ":"); r != -1 {
		return true
	}
	return false
}

func WrapIPv6(host string) string {
	if IsIPv6(host) {
		return "[" + host + "]"
	}
	return host
}

func WrapAddr(host string, port uint16) string {
	return WrapIPv6(host) + ":" + strconv.Itoa(int(port))
}

// ParseAddress parses an address buffer into string.
// It returns resulting address, length of bytes required,
// and error if exists.
// You must check that len(buf) >= n, otherwise addr is not
// a valid address.
func ParseAddress(buf []byte) (addr string, n int, err error) {
	if len(buf) < 2 {
		n = 2
		return
	}
	typ := buf[0]
	if typ == 0x1 { // ipv4
		n = 7
		if len(buf) < n {
			return
		}
		host := net.IPv4(buf[1], buf[2], buf[3], buf[4]).String()
		var p16 uint16
		binary.Read(bytes.NewBuffer(buf[5:7]), binary.BigEndian, &p16)
		port := strconv.Itoa(int(p16))
		addr = "[" + host + "]:" + port
	} else if typ == 0x3 { // host
		alen := int(buf[1])
		n = 1 + 1 + alen + 2
		if len(buf) < n {
			return
		}
		host := string(buf[2 : 2+alen])
		if IsIPv6(host) {
			host = "[" + host + "]"
		}
		var p16 uint16
		binary.Read(bytes.NewBuffer(buf[2+alen:n]), binary.BigEndian, &p16)
		port := strconv.Itoa(int(p16))
		addr = host + ":" + port
	} else if typ == 0x4 { // ipv6
		n = 19
		if len(buf) < n {
			return
		}
		host := net.IP(buf[1:17]).String()
		var p16 uint16
		binary.Read(bytes.NewBuffer(buf[17:19]), binary.BigEndian, &p16)
		port := strconv.Itoa(int(p16))
		addr = "[" + host + "]:" + port
	} else { // error
		err = ERR_INVALID_ADDR_TYPE
	}
	return
}
