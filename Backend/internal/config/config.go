package config

import (
	"os"
	"strconv"
)

type Config struct {
	Addr               string
	DatabaseURL        string
	PixelEncryptionKey string
	AdminTemplatesDir  string
	StaticDir          string
	DispatchIntervalSec int
}

func Load() Config {
	return Config{
		Addr:                getenv("ADDR", ":8080"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		PixelEncryptionKey:  os.Getenv("PIXEL_ENCRYPTION_KEY"),
		AdminTemplatesDir:   getenv("ADMIN_TEMPLATES_DIR", "../Frontend-admin/templates"),
		StaticDir:           getenv("STATIC_DIR", "../Frontend-admin/static"),
		DispatchIntervalSec: getenvInt("PIXEL_DISPATCH_INTERVAL_SEC", 10),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
