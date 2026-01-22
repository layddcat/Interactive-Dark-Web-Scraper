package services

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"
)

func StartRealMonitoring() {
	waitForProxy("tor:9150")

	for {
		targets := loadTargets()
		
		if len(targets) == 0 {
			fmt.Println("[KRİTİK HATA] targets.yaml içinden hiçbir URL okunamadı.")
		} else {
			fmt.Printf("[LOG] %d adet hedef Tor ağı üzerinden taranıyor...\n", len(targets))
			for _, targetURL := range targets {
				content, isActive, err := fetchOnionContent(targetURL)
				
				if err != nil {
					fmt.Printf("[SYSTEM] %s erisilemedi: %v\n", targetURL, err)
					ProcessRawData("Interactive-Crawler", targetURL, "", time.Now().Format("2006-01-02 15:04"), false)
					continue 
				}

				ProcessRawData("Interactive-Crawler", targetURL, content, time.Now().Format("2006-01-02 15:04"), isActive)
				time.Sleep(2 * time.Second) 
			}
		}
		time.Sleep(1 * time.Hour)
	}
}

func waitForProxy(addr string) {
	fmt.Printf("[SYSTEM] Tor Proxy (%s) bekleniyor...\n", addr)
	for {
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			conn.Close()
			fmt.Println("[SUCCESS] Tor Proxy hazır. İstihbarat toplama başlatılıyor.")
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func fetchOnionContent(onionURL string) (string, bool, error) {
	proxyAddr := "tor:9150"
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return "", false, fmt.Errorf("proxy baglantisi kurulamadi")
	}

	transport := &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}}

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	resp, err := client.Get(onionURL)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("HTTP %d hatası", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}

	return string(body), true, nil
}

func loadTargets() []string {
	data, err := os.ReadFile("targets.yaml")
	if err != nil {
		return nil
	}

	var targets []string
	if err := yaml.Unmarshal(data, &targets); err == nil && len(targets) > 0 {
		return targets
	}
	var structTargets struct { Targets []string `yaml:"targets"` }
	if err := yaml.Unmarshal(data, &structTargets); err == nil {
		return structTargets.Targets
	}
	return nil
}