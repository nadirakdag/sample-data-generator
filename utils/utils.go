package utils

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

func GetLocalIP() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addresses {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no IP address found")
}

func GeneratePublicIP() string {
	for {
		ip := rand.Uint32()
		a := byte(ip >> 24)
		b := byte(ip >> 16 & 0xFF)
		if (a >= 1 && a <= 126) || (a >= 128 && a <= 223) && !(a == 10 || a == 172 && (b >= 16 && b <= 31) || a == 192 && b == 168) {
			return fmt.Sprintf("%d.%d.%d.%d", a, b, byte(ip>>8&0xFF), byte(ip&0xFF))
		}
	}
}

func IpToUint32(ip string) uint32 {
	parts := strings.Split(ip, ".")
	b0 := parseUint8(parts[0])
	b1 := parseUint8(parts[1])
	b2 := parseUint8(parts[2])
	b3 := parseUint8(parts[3])
	return uint32(b0)<<24 | uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3)
}

func parseUint8(s string) uint8 {
	var i uint64
	fmt.Sscanf(s, "%d", &i)
	return uint8(i)
}
