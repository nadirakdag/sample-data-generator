package netflow

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sample-data-generator/utils"
	"time"
)

func SendNetFlowV9Message(conn net.Conn, localIP string, log bool) {
	netflowV9MsgIn := createNetFlowV9Message(localIP, "in")
	conn.Write(netflowV9MsgIn)
	if log {
		fmt.Printf("Sent NetFlow v9 IN message to %s: %x\n", conn.RemoteAddr().String(), netflowV9MsgIn)
	}

	netflowV9MsgOut := createNetFlowV9Message(localIP, "out")
	conn.Write(netflowV9MsgOut)
	if log {
		fmt.Printf("Sent NetFlow v9 OUT message to %s: %x\n", conn.RemoteAddr().String(), netflowV9MsgOut)
	}
}

func createNetFlowV9Message(localIP string, direction string) []byte {
	header := make([]byte, 20)
	binary.BigEndian.PutUint16(header[0:2], 9) // Version
	binary.BigEndian.PutUint16(header[2:4], 1) // Count of flow records
	binary.BigEndian.PutUint32(header[4:8], uint32(rand.Int31()))
	binary.BigEndian.PutUint32(header[8:12], uint32(time.Now().Unix()))
	binary.BigEndian.PutUint32(header[12:16], uint32(rand.Int31()))
	binary.BigEndian.PutUint32(header[16:20], uint32(rand.Int31()))

	templateRecord := make([]byte, 4)
	binary.BigEndian.PutUint16(templateRecord[0:2], 1) // Template ID
	binary.BigEndian.PutUint16(templateRecord[2:4], 2) // Field Count

	templateFields := make([]byte, 8)
	binary.BigEndian.PutUint16(templateFields[0:2], 8)  // Source IP address field type
	binary.BigEndian.PutUint16(templateFields[2:4], 4)  // Source IP address field length
	binary.BigEndian.PutUint16(templateFields[4:6], 18) // Destination IP address field type
	binary.BigEndian.PutUint16(templateFields[6:8], 4)  // Destination IP address field length

	dataRecord := make([]byte, 8)
	srcIP := utils.GeneratePublicIP()
	dstIP := utils.GeneratePublicIP()
	if direction == "in" {
		binary.BigEndian.PutUint32(dataRecord[0:4], utils.IpToUint32(srcIP))   // Source IP address
		binary.BigEndian.PutUint32(dataRecord[4:8], utils.IpToUint32(localIP)) // Destination IP address
	} else {
		binary.BigEndian.PutUint32(dataRecord[0:4], utils.IpToUint32(localIP)) // Source IP address
		binary.BigEndian.PutUint32(dataRecord[4:8], utils.IpToUint32(dstIP))   // Destination IP address
	}

	message := append(header, templateRecord...)
	message = append(message, templateFields...)
	message = append(message, dataRecord...)
	return message
}
