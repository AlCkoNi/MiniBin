package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	AutoCleanMinutes int
	MaxFillSizeMB    int64
}

func defaultConfig() Config {
	return Config{
		AutoCleanMinutes: 0,
		MaxFillSizeMB:    1024,
	}
}

func configPath() string {
	return filepath.Join(executableDir(), "minibin.ini")
}

func loadConfig() Config {
	cfg := defaultConfig()
	f, err := os.Open(configPath())
	if err != nil {
		return cfg
	}
	defer f.Close()

	section := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		switch section + "." + key {
		case "configure.autocleanminutes":
			if n, err := strconv.Atoi(value); err == nil && validAutoCleanMinutes(n) {
				cfg.AutoCleanMinutes = n
			}
		case "display.maxfillsizemb":
			if n, err := strconv.ParseInt(value, 10, 64); err == nil && n > 0 {
				cfg.MaxFillSizeMB = n
			}
		}
	}
	return cfg
}

func saveConfig(cfg Config) error {
	data := fmt.Sprintf("; MiniBin %s configuration\r\n[Configure]\r\nAutoCleanMinutes=%d\r\n\r\n[Display]\r\nMaxFillSizeMB=%d\r\n", AppVersion, cfg.AutoCleanMinutes, cfg.MaxFillSizeMB)
	return os.WriteFile(configPath(), []byte(data), 0644)
}

func validAutoCleanMinutes(n int) bool {
	return n == 0 || n == 30 || n == 60 || n == 180
}
