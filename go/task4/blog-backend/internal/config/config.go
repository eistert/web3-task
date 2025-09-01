package config

import "os"

// 读取配置（端口、DSN、JWT 密钥等）

type Config struct {
	// DB：优先 MySQL（MYSQL_DSN），否则用 SQLite（SQLitePath）
	MySQLDSN   string
	SQLitePath string
	JWTSecret  string
	Port       string
	TokenTTL   int64 // seconds
}

func Load() Config {
	c := Config{
		MySQLDSN:   os.Getenv("MYSQL_DSN"),
		SQLitePath: getEnv("SQLITE_PATH", "blog.db"),
		JWTSecret:  getEnv("JWT_SECRET", "123456"),
		Port:       getEnv("PORT", "8080"),
		TokenTTL:   24 * 3600,
	}
	return c
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
