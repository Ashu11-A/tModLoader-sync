package network

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

var ipServices = []string{
	"https://ipinfo.io/json",
	"https://ifconfig.me/all.json",
	"https://icanhazip.com/",
	"https://api.ipify.org/?format=json",
}

func GetPublicIP() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, service := range ipServices {
		resp, err := client.Get(service)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var ip string
		if strings.Contains(service, "ipinfo.io") {
			var data struct {
				IP string `json:"ip"`
			}
			if err := json.Unmarshal(body, &data); err == nil {
				ip = data.IP
			}
		} else if strings.Contains(service, "ifconfig.me") {
			var data struct {
				IP string `json:"ip_addr"`
			}
			if err := json.Unmarshal(body, &data); err == nil {
				ip = data.IP
			}
		} else if strings.Contains(service, "icanhazip.com") {
			ip = strings.TrimSpace(string(body))
		} else if strings.Contains(service, "ipify.org") {
			var data struct {
				IP string `json:"ip"`
			}
			if err := json.Unmarshal(body, &data); err == nil {
				ip = data.IP
			}
		}

		if ip != "" {
			return ip, nil
		}
	}

	return "", io.EOF
}
