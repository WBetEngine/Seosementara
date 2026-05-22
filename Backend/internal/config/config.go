package config

import (
	"os"
	"strconv"
)

type Config struct {
	Addr                string
	DatabaseURL         string
	PixelEncryptionKey  string
	MasterEncryptionKey string
	AdminTemplatesDir   string
	AdminPartialsDir    string
	StaticDir           string
	DispatchIntervalSec int
	SuperAdminToken     string
}

func Load() Config {
	return Config{
		Addr:                getenv("ADDR", ":8080"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		PixelEncryptionKey:  os.Getenv("PIXEL_ENCRYPTION_KEY"),
		MasterEncryptionKey: os.Getenv("MASTER_ENCRYPTION_KEY"),
		AdminTemplatesDir:   getenv("ADMIN_TEMPLATES_DIR", "../Frontend-admin/templates"),
		AdminPartialsDir:    getenv("ADMIN_PARTIALS_DIR", "../Frontend-admin/public/admin/_partials"),
		StaticDir:           getenv("STATIC_DIR", "../Frontend-admin/static"),
		DispatchIntervalSec: getenvInt("PIXEL_DISPATCH_INTERVAL_SEC", 10),
		SuperAdminToken:     os.Getenv("SUPER_ADMIN_TOKEN"),
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
