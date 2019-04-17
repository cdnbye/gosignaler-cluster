package util

import (
	"fmt"
	"net"
	"os"
)

var (
	LOCAL_IP = InitIntranetIp()
)

//get local ipAddr
func InitIntranetIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := len(addrs) - 1; i >= 0; i-- {
		address := addrs[i]
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
			//	log.Info( ipnet.IP.String())
				fmt.Println("ip:", ipnet.IP.String())
				return ipnet.IP.String()
			}

		}

	}
	return ""
}
