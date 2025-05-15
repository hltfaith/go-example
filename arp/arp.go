package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"syscall"
)

const (
	ETH_P_ARP   = 0x0806 // ARP 协议类型
	ARP_REQUEST = 1      // ARP 请求
	ARP_REPLY   = 2      // ARP 响应
)

type EthernetHeader struct {
	DestMAC   [6]byte
	SrcMAC    [6]byte
	EtherType uint16
}

type ARPHeader struct {
	HardwareType   uint16
	ProtocolType   uint16
	HardwareLength uint8
	ProtocolLength uint8
	Operation      uint16
	SenderMAC      [6]byte
	SenderIP       [4]byte
	TargetMAC      [6]byte
	TargetIP       [4]byte
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: sudo ./arp <target_ip>")
		os.Exit(1)
	}

	targetIP := net.ParseIP(os.Args[1]).To4()
	if targetIP == nil {
		fmt.Println("Invalid target IP address")
		os.Exit(1)
	}

	// 创建原始套接字 (相当于 SOCK_PACKET)
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		fmt.Printf("Socket creation error: %v\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	// 获取接口信息
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		fmt.Printf("Interface error: %v\n", err)
		os.Exit(1)
	}

	// 构建 ARP 请求包
	requestPacket := buildARPPacket(iface, targetIP, [6]byte{0, 0, 0, 0, 0, 0})

	// 发送 ARP 请求
	if err := sendARPPacket(fd, iface, requestPacket); err != nil {
		fmt.Printf("Send ARP error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ARP request sent for %s\n", targetIP)

	// 接收 ARP 响应
	recvARPPacket(fd, iface, targetIP)

}

func buildARPPacket(iface *net.Interface, targetIP net.IP, targetMac [6]byte) []byte {
	localIP := getLocalIP(iface)
	if localIP == nil {
		fmt.Println("Could not determine local IP")
		os.Exit(1)
	}

	ethHeader := EthernetHeader{
		DestMAC:   [6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // 广播地址
		SrcMAC:    [6]byte(iface.HardwareAddr),
		EtherType: ETH_P_ARP,
	}
	arpHeader := ARPHeader{
		HardwareType:   1, // 以太网
		ProtocolType:   0x800,
		HardwareLength: 6,           // MAC 地址长度
		ProtocolLength: 4,           // IP 地址长度
		Operation:      ARP_REQUEST, // ARP 请求
		SenderMAC:      [6]byte(iface.HardwareAddr),
		SenderIP:       [4]byte(localIP.To4()),
		TargetMAC:      targetMac,
		TargetIP:       [4]byte(targetIP),
	}

	packet := make([]byte, 0)
	packet = append(packet, toBytes(ethHeader)...)
	packet = append(packet, toBytes(arpHeader)...)
	return packet
}

func sendARPPacket(fd int, iface *net.Interface, packet []byte) error {
	// 构建 sockaddr_ll
	sa := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  iface.Index,
		Hatype:   1, // ARPHRD_ETHER
		Pkttype:  0, // PACKET_HOST
		Halen:    6, // MAC 地址长度
	}

	if err := syscall.Bind(fd, &sa); err != nil {
		fmt.Printf("Bind error: %v\n", err)
		os.Exit(1)
	}

	// 发送数据包
	return syscall.Sendto(fd, packet, 0, &sa)
}

func recvARPPacket(fd int, iface *net.Interface, targetIP net.IP) {
	buffer := make([]byte, 1024)

	// 设置超时
	timeout := syscall.Timeval{Sec: 5, Usec: 0}
	syscall.SetsockoptTimeval(fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &timeout)

	for {
		_, _, err := robustRecvfrom(fd, buffer)
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				fmt.Println("Timeout waiting for ARP reply")
			} else {
				fmt.Printf("Receive error: %v\n", err)
			}
			return
		}
		// 解析以太网头
		ethHeader := EthernetHeader{}
		fromBytes(buffer[:14], &ethHeader)
		if ethHeader.EtherType != ETH_P_ARP {
			continue
		}
		// 解析 ARP 头
		arpHeader := ARPHeader{}
		fromBytes(buffer[14:42], &arpHeader)
		if arpHeader.Operation == ARP_REPLY &&
			net.IP(arpHeader.SenderIP[:]).Equal(targetIP) {
			fmt.Printf("ARP reply from %s: %02x:%02x:%02x:%02x:%02x:%02x\n",
				targetIP,
				arpHeader.SenderMAC[0], arpHeader.SenderMAC[1], arpHeader.SenderMAC[2],
				arpHeader.SenderMAC[3], arpHeader.SenderMAC[4], arpHeader.SenderMAC[5])
			return
		}
	}
}

func robustRecvfrom(fd int, buf []byte) (int, syscall.Sockaddr, error) {
	for {
		n, addr, err := syscall.Recvfrom(fd, buf, 0)
		if err == syscall.EINTR {
			continue // 被中断，重试
		}
		return n, addr, err
	}
}

func getLocalIP(iface *net.Interface) net.IP {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return nil
}

func toBytes(data interface{}) []byte {
	size := binary.Size(data)
	buf := make([]byte, size)
	binary.LittleEndian.PutUint64(buf, 0) // 清零
	binary.BigEndian.PutUint16(buf, 0)    // 清零
	binary.BigEndian.PutUint32(buf, 0)    // 清零

	switch v := data.(type) {
	case EthernetHeader:
		copy(buf[0:6], v.DestMAC[:])
		copy(buf[6:12], v.SrcMAC[:])
		binary.BigEndian.PutUint16(buf[12:14], v.EtherType)
	case ARPHeader:
		binary.BigEndian.PutUint16(buf[0:2], v.HardwareType)
		binary.BigEndian.PutUint16(buf[2:4], v.ProtocolType)
		buf[4] = v.HardwareLength
		buf[5] = v.ProtocolLength
		binary.BigEndian.PutUint16(buf[6:8], v.Operation)
		copy(buf[8:14], v.SenderMAC[:])
		copy(buf[14:18], v.SenderIP[:])
		copy(buf[18:24], v.TargetMAC[:])
		copy(buf[24:28], v.TargetIP[:])
	}
	return buf
}

func fromBytes(data []byte, result interface{}) {
	switch v := result.(type) {
	case *EthernetHeader:
		copy(v.DestMAC[:], data[0:6])
		copy(v.SrcMAC[:], data[6:12])
		v.EtherType = binary.BigEndian.Uint16(data[12:14])
	case *ARPHeader:
		v.HardwareType = binary.BigEndian.Uint16(data[0:2])
		v.ProtocolType = binary.BigEndian.Uint16(data[2:4])
		v.HardwareLength = data[4]
		v.ProtocolLength = data[5]
		v.Operation = binary.BigEndian.Uint16(data[6:8])
		copy(v.SenderMAC[:], data[8:14])
		copy(v.SenderIP[:], data[14:18])
		copy(v.TargetMAC[:], data[18:24])
		copy(v.TargetIP[:], data[24:28])
	}
}

// 主机字节序转网络字节序
func htons(i uint16) uint16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return binary.LittleEndian.Uint16(b)
}
