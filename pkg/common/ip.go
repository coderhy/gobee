package common

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//IP 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作（公网）
func IP(r *http.Request) string {

	ip := ""
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if ip == "::1" { //localhost
			ip = "127.0.0.1"
		}
		return ip
	}

	return ""
}

//IP2long IP2long
func IP2long(ipstr string) (ip int64) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += int64(ip1 * 0x1000000)
	ip += int64(ip2 * 0x10000)
	ip += int64(ip3 * 0x100)
	ip += int64(ip4)

	return
}

// LocalIP 本地IP
func LocalIP() string {

	ip := ""
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return ip
	}

	for i := 0; i < len(netInterfaces); i++ {
		//net.FlagUp标志进行确认，排除掉无用的网卡
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						// fmt.Println(ipnet.IP.String())
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ip
}

//Long2IP Long2IP
func Long2IP(properAddress uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, properAddress)
	ip := net.IP(ipByte)
	return ip.String()
}
