package netflow

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sample-data-generator/utils"
	"time"
)

func SendNetFlowV5Message(conn net.Conn, localIP string, log bool) {
	netflowV5MsgIn := createNetFlowV5Message(localIP, "in")
	conn.Write(netflowV5MsgIn)
	if log {
		fmt.Printf("Sent NetFlow v5 IN message to %s: %x\n", conn.RemoteAddr().String(), netflowV5MsgIn)
	}

	netflowV5MsgOut := createNetFlowV5Message(localIP, "out")
	conn.Write(netflowV5MsgOut)
	if log {
		fmt.Printf("Sent NetFlow v5 OUT message to %s: %x\n", conn.RemoteAddr().String(), netflowV5MsgOut)
	}
}

func createNetFlowV5Message(localIP string, direction string) []byte {
	header := make([]byte, 24)
	binary.BigEndian.PutUint16(header[0:2], 5) // Version
	binary.BigEndian.PutUint16(header[2:4], 1) // Count of flow records
	binary.BigEndian.PutUint32(header[4:8], uint32(rand.Int31()))
	binary.BigEndian.PutUint32(header[8:12], uint32(time.Now().Unix()))
	binary.BigEndian.PutUint32(header[12:16], uint32(rand.Int31()))
	binary.BigEndian.PutUint32(header[16:20], uint32(rand.Int31()))
	binary.BigEndian.PutUint32(header[20:24], uint32(rand.Int31()))

	record := make([]byte, 48)
	srcIP := utils.GeneratePublicIP()
	dstIP := utils.GeneratePublicIP()
	if direction == "in" {
		binary.BigEndian.PutUint32(record[0:4], utils.IpToUint32(srcIP))
		binary.BigEndian.PutUint32(record[4:8], utils.IpToUint32(localIP))
	} else {
		binary.BigEndian.PutUint32(record[0:4], utils.IpToUint32(localIP))
		binary.BigEndian.PutUint32(record[4:8], utils.IpToUint32(dstIP))
	}
	binary.BigEndian.PutUint32(record[8:12], rand.Uint32())
	binary.BigEndian.PutUint16(record[12:14], uint16(rand.Intn(65535)))
	binary.BigEndian.PutUint16(record[14:16], uint16(rand.Intn(65535)))
	binary.BigEndian.PutUint32(record[16:20], rand.Uint32())
	binary.BigEndian.PutUint32(record[20:24], rand.Uint32())
	binary.BigEndian.PutUint32(record[24:28], uint32(rand.Int31()))
	binary.BigEndian.PutUint32(record[28:32], uint32(rand.Int31()))
	binary.BigEndian.PutUint16(record[32:34], uint16(rand.Intn(65535)))
	binary.BigEndian.PutUint16(record[34:36], uint16(rand.Intn(65535)))
	record[38] = byte(rand.Intn(256)) // Protocol
	record[39] = byte(rand.Intn(256)) // TCP flags
	record[40] = byte(rand.Intn(256)) // Type of service
	binary.BigEndian.PutUint16(record[41:43], uint16(rand.Intn(65535)))
	binary.BigEndian.PutUint16(record[43:45], uint16(rand.Intn(65535)))
	if direction == "in" {
		binary.BigEndian.PutUint32(record[44:48], uint32(rand.Int31())) // Traffic In
	} else {
		binary.BigEndian.PutUint32(record[44:48], uint32(rand.Int31())) // Traffic Out
	}

	return append(header, record...)
}
