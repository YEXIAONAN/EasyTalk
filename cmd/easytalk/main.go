package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"easytalk"
	"easytalk/internal/config"
	"easytalk/internal/server"
	"easytalk/internal/version"
)

func main() {
	env := config.LoadEnv()

	configPath := flag.String("config", env.ConfigPath, "path to the config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			printMissingConfig(*configPath)
			os.Exit(1)
		}
		log.Fatalf("load config: %v", err)
	}

	static, err := easytalk.FrontendFS()
	if err != nil {
		log.Fatalf("frontend assets: %v", err)
	}

	srv := server.New(cfg, *configPath, static)

	addr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	printStartup(env.Port, *configPath)

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

func printMissingConfig(path string) {
	fmt.Printf("EasyTalk %s\n\n", version.Version)
	fmt.Printf("Config file not found:\n%s\n\n", path)
	fmt.Printf("Please copy:\nconfig.example.json\n\n")
	fmt.Printf("to:\nconfig.json\n\n")
	fmt.Printf("Then configure your AI providers.\n")
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