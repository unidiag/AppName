package main

import (
	"bufio"
	"fmt"
	"log"
	"main/models"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

type SetupConfig struct {
	Listen       string
	ClusterNodes string
	SecretKey    string
	InstallSvc   bool
}

func runSetupWizard() (*SetupConfig, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(APPNAME + " first run setup")

	cfg := &SetupConfig{}

	// PORT
	for {
		addr := ask(reader, "Listen address", ":9000")

		addr, err := normalizeListenAddr(addr)
		if err != nil {
			fmt.Println("Invalid address:", err)
			continue
		}

		if !checkPort(addr) {
			fmt.Printf("Address %s is already in use\n", addr)
			continue
		}

		cfg.Listen = addr
		break
	}

	// install systemd unit
	ans := strings.ToLower(ask(reader, "Install `"+strings.ToLower(APPNAME)+".service` to systemd? (y/N)", "n"))
	cfg.InstallSvc = ans == "y" || ans == "yes"

	fmt.Println("SETUP COMPLETED!")

	return cfg, true
}

func ask(reader *bufio.Reader, question string, def string) string {
	fmt.Printf("%s [%s]: ", question, def)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return def
	}
	return text
}

// /////////// HELPERS
func checkPort(addr string) bool {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return false // порт занят
	}
	l.Close()
	return true // порт свободен
}

// systemd
func installSystemdService(binPath string, args []string) (string, error) {

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	serviceName := strings.ToLower(APPNAME) + ".service"

	// формируем ExecStart
	execCmd := binPath
	if len(args) > 0 {
		execCmd += " " + strings.Join(args, " ")
	}

	unit := fmt.Sprintf(`[Unit]
Description=`+APPNAME+`
After=network.target

[Service]
ExecStart=%s
Restart=always
RestartSec=3
WorkingDirectory=%s

[Install]
WantedBy=default.target
`, execCmd, filepath.Dir(binPath))

	var servicePath string
	var cmdPrefix []string

	if u.Uid == "0" {
		servicePath = "/etc/systemd/system/" + serviceName
		cmdPrefix = []string{"systemctl"}
	} else {
		dir := filepath.Join(u.HomeDir, ".config/systemd/user")
		os.MkdirAll(dir, 0755)
		servicePath = filepath.Join(dir, serviceName)
		cmdPrefix = []string{"systemctl", "--user"}
	}

	if err := os.WriteFile(servicePath, []byte(unit), 0644); err != nil {
		return "", err
	}

	exec.Command(cmdPrefix[0], append(cmdPrefix[1:], "daemon-reload")...).Run()
	exec.Command(cmdPrefix[0], append(cmdPrefix[1:], "enable", serviceName)...).Run()
	exec.Command(cmdPrefix[0], append(cmdPrefix[1:], "start", serviceName)...).Run()

	return servicePath, nil
}

func createDefaultUser(login, pass string, status uint) {
	var count int64
	db.Model(&models.User{}).Where("login = ?", login).Count(&count)
	if count > 0 {
		return
	}
	hashpassword, _ := hashPassword(pass)
	u := models.User{
		Login:    login,
		Password: hashpassword,
		Status:   status,
	}
	if err := db.Create(&u).Error; err != nil {
		log.Println("Create user error:", err)
		return
	}
}

func normalizeListenAddr(input string) (string, error) {

	input = strings.TrimSpace(input)

	if input == "" {
		return ":8888", nil
	}

	// если введён только порт
	if !strings.Contains(input, ":") {
		port := cast.ToInt(input)
		if port <= 0 || port > 65535 {
			return "", fmt.Errorf("invalid port")
		}
		return fmt.Sprintf(":%d", port), nil
	}

	// если формат ":8888"
	if strings.HasPrefix(input, ":") {
		port := cast.ToInt(strings.TrimPrefix(input, ":"))
		if port <= 0 || port > 65535 {
			return "", fmt.Errorf("invalid port")
		}
		return fmt.Sprintf(":%d", port), nil
	}

	// формат host:port
	host, portStr, err := net.SplitHostPort(input)
	if err != nil {
		return "", fmt.Errorf("invalid address format")
	}

	port := cast.ToInt(portStr)
	if port <= 0 || port > 65535 {
		return "", fmt.Errorf("invalid port")
	}

	return net.JoinHostPort(host, fmt.Sprintf("%d", port)), nil
}
