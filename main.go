package main

import (
	"flag"
	"fmt"
	"net"
	"sample-data-generator/netflow"
	"sample-data-generator/text_syslog"
	"sample-data-generator/utils"
	"strings"
	"time"
)

func main() {
	ipListPtr := flag.String("ip", "127.0.0.1", "Comma-separated list of destination IP addresses")
	portListPtr := flag.String("port", "514", "Comma-separated list of destination UDP ports")
	ratePtr := flag.Int("rate", 10, "Messages per second")
	appNamePtr := flag.String("app", "myApp", "Application name for syslog messages")
	formatPtr := flag.String("format", "syslog", "Message format: syslog, netflow-v5, or netflow-v9")
	logPtr := flag.Bool("log", false, "Enable logging")

	flag.Parse()

	localIP, err := utils.GetLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP address:", err)
		return
	}

	connections, done := createConnections(ipListPtr, portListPtr)
	if done {
		return
	}

	ticker := time.NewTicker(time.Second / time.Duration(*ratePtr))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			switch *formatPtr {
			case "syslog":
				for _, conn := range connections {
					text_syslog.SendSyslogMessage(conn, *appNamePtr, *logPtr)
				}
			case "netflow-v5":
				for _, conn := range connections {
					netflow.SendNetFlowV5Message(conn, localIP, *logPtr)
				}
			case "netflow-v9":
				for _, conn := range connections {
					netflow.SendNetFlowV9Message(conn, localIP, *logPtr)
				}
			default:
				fmt.Println("Unknown format:", *formatPtr)
			}
		}
	}
}

func createConnections(ipListPtr *string, portListPtr *string) ([]net.Conn, bool) {
	ips := strings.Split(*ipListPtr, ",")
	ports := strings.Split(*portListPtr, ",")

	connections := make([]net.Conn, 0)
	for _, ip := range ips {
		for _, port := range ports {
			address := fmt.Sprintf("%s:%s", strings.TrimSpace(ip), strings.TrimSpace(port))
			conn, err := net.Dial("udp", address)
			if err != nil {
				fmt.Println("Error creating UDP connection:", err)
				return nil, true
			}
			connections = append(connections, conn)
		}
	}
	defer func() {
		for _, conn := range connections {
			if err := conn.Close(); err != nil {
				panic(err)
			}
		}
	}()
	return connections, false
}
