package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"easytalk"
	"easytalk/internal/config"
	"easytalk/internal/server"
	"easytalk/internal/version"
)

func main() {
	configPath := flag.String("config", "config.json", "path to the config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	static, err := easytalk.FrontendFS()
	if err != nil {
		log.Fatalf("frontend assets: %v", err)
	}

	srv := server.New(cfg, static)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	printStartup(cfg.Server.Port, *configPath)

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func printStartup(port int, configPath string) {
	fmt.Printf("EasyTalk %s\n\n", version.Version)
	fmt.Printf("Local:\nhttp://127.0.0.1:%d\n\n", port)
	if ip := lanIP(); ip != "" {
		fmt.Printf("LAN:\nhttp://%s:%d\n\n", ip, port)
	}
	fmt.Printf("Config:\n%s\n", configPath)
}

// lanIP returns the first non-loopback IPv4 address, or an empty string.
func lanIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}
	return ""
}