package main

import (
	"net"
	"net/http"
	"strings"
)

func serviceIP() (string, error) {
	ip := "127.0.0.1"
	ifaces, err := net.Interfaces()
	if err != nil {
		return ip, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ip, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return ip, nil
}

// GetInfoFromRequestHeader 함수는 리퀘스트 헤더를 불러와서 Device, OS, Browser 정보를 가지고 온다. 알 수 없다면 unknown 으로 출력한다.
func GetInfoFromRequestHeader(r *http.Request) (string, string, string) {
	device := "desktop/notebook"
	osname := "unknown"
	browser := "unknown"
	devices := []string{"iPad", "iPhone", "Android", "Kindle"}
	os := []string{"Linux", "Mac OS", "Windows"}
	browsers := []string{"Opera", "Whale", "Firefox", "Edg", "Safari", "Explorer", "Chrome"} // 사용률 통계의 역순으로 리스트를 넣는다. 요즘 브라우저는 여러 소스코드의 조합으로 이루어진다.
	ua := r.UserAgent()
	for _, d := range devices {
		if !strings.Contains(ua, d) {
			continue
		}
		device = strings.ToLower(d)
	}
	for _, o := range os {
		if !strings.Contains(ua, o) {
			continue
		}
		osname = strings.Replace(strings.ToLower(o), " ", "", -1)
	}
	for _, b := range browsers {
		if strings.Contains(ua, b) {
			browser = strings.ToLower(b)
			break
		}
	}
	return device, osname, browser
}
