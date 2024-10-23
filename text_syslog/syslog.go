package text_syslog

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func SendSyslogMessage(conn net.Conn, appName string, log bool) {
	uniqueID := rand.Intn(1000000)
	msg := fmt.Sprintf("<134>1 %s %s - - - - [uniqueID=%d] Example syslog message", time.Now().Format(time.RFC3339), appName, uniqueID)
	conn.Write([]byte(msg))
	if log {
		fmt.Printf("Sent syslog message to %s: %s\n", conn.RemoteAddr().String(), msg)
	}
}
