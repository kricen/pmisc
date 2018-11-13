package lib

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	HttpClient *http.Client
)

func init() {
	initClient()
}

//ResolveHostIp : a function to resolve localhost ip
func ResolveHostIP() string {
	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIp, ok := netInterfaceAddress.(*net.IPNet)
		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {
			ip := networkIp.IP.String()
			return ip
		}
	}
	return ""
}

// Decimal : remain three significant digits lg: 3.1415926 to 3.14
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value
}

// init common http client
func initClient() {
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
	}

	HttpClient = &http.Client{
		Transport: trans,
		Timeout:   3 * time.Second,
	}
}

// ValidateEmailAddress : validata whether email address is illedge
func ValidateEmailAddress(addr string) bool {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", addr); !m {
		return false
	}
	return true
}
