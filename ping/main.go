package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

func BuildPingPacket(srcIP, dstIP string) []byte {
	packet := make([]byte, 46) // 20字节IP头 + 8字节ICMP头 + 18字节ICMP数据

	// IP 头
	ipHeader := packet[:20]
	ipHeader[0] = 0x45                                    // 版本和头长度
	binary.BigEndian.PutUint16(ipHeader[2:4], uint16(46)) // 总长度
	binary.BigEndian.PutUint16(ipHeader[4:6], uint16(1))  // 标识
	copy(ipHeader[6:8], []byte{0x00, 0x00})               // 分片
	ipHeader[8] = 64                                      // TTL
	ipHeader[9] = syscall.IPPROTO_ICMP                    // 协议
	copy(ipHeader[10:12], []byte{0x00, 0x00})
	copy(ipHeader[12:16], net.ParseIP(srcIP).To4())                   // 源IP
	copy(ipHeader[16:20], net.ParseIP(dstIP).To4())                   // 目的IP
	binary.BigEndian.PutUint16(ipHeader[10:12], IPChecksum(ipHeader)) // 首部校验和

	// ICMP 头 + 数据
	icmpHeader := packet[20:46]
	icmpHeader[0] = 8 // 消息类型
	icmpHeader[1] = 0 // 消息代码
	copy(icmpHeader[2:4], []byte{0x00, 0x00})
	binary.BigEndian.PutUint16(icmpHeader[4:6], uint16(1)) // 标识
	binary.BigEndian.PutUint16(icmpHeader[6:8], uint16(1)) // 序号

	// ICMP 数据
	binary.BigEndian.PutUint16(icmpHeader[8:26], uint16(0))             // 数据
	binary.BigEndian.PutUint16(icmpHeader[2:4], IPChecksum(icmpHeader)) // 校验和

	return packet
}

func IPChecksum(b []byte) uint16 {
	var sum uint32
	for i := 0; i < len(b)-1; i += 2 {
		sum += uint32(b[i])<<8 | uint32(b[i+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return uint16(^sum)
}

func main() {

	srcIP := "192.168.54.109"
	destIP := "192.168.42.173"

	// 创建socket
	var err error
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		log.Fatalf("create socket err: %v", err)
	}
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		log.Fatalf("set sock option IP_HDRINCL err: %v", err)
	}

	dstAddr := &syscall.SockaddrInet4{}
	copy(dstAddr.Addr[:], net.ParseIP(destIP).To4())

	// 封装ping包
	packet := BuildPingPacket(srcIP, destIP)

	// 发送ping请求
	err = syscall.Sendto(fd, packet, 0, dstAddr)
	if err != nil {
		log.Fatalf("send packet err: %v", err)
	}

	start := time.Now()

	// 接收ping请求
	buf := make([]byte, 46)
	_, _, err = syscall.Recvfrom(fd, buf, 0)
	if err != nil {
		log.Fatalf("Received data err: %v", err)
	}

	end := time.Since(start)
	ms := (float64(end.Milliseconds())*100 + 0.5) / 100.0

	// 自定义ping报文总长度60字节
	// Ethernet报文总长度14字节 (由于SOCK_RAW只能处理IP层报文)
	// IP报文总长度20字节
	// ICMP报文总长度26字节, 前8字节为头部, 剩余18字节为数据
	ipHeader := buf[0:20]
	ttl := ipHeader[8]
	icmp := buf[20:46]
	icmpHeader := icmp[0:8]
	Type := icmpHeader[0]
	code := icmpHeader[1]
	seq := binary.BigEndian.Uint16(icmpHeader[6:8])
	if Type != 0x00 && code != 0x00 {
		log.Fatalf("Received icmp err, type=%d code=%d", uint8(Type), uint8(code))
	}

	fmt.Printf("%v bytes from %s: icmp_seq=%d ttl=%d time=%.2f ms\n", len(icmp), destIP, seq, uint8(ttl), ms)
}
