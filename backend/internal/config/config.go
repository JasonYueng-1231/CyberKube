package config

import (
    "os"
)

type Config struct {
    MySQLDSN string
    RedisAddr string
    JWTSecret string
    AESKey    string // 32字节
    SkipTLSVerify bool
}

func Load() *Config {
    c := &Config{
        MySQLDSN: getenv("MYSQL_DSN", "root:cyberkube@123@tcp(localhost:3306)/cyberkube?charset=utf8mb4&parseTime=True&loc=Local"),
        RedisAddr: getenv("REDIS_ADDR", "localhost:6379"),
        JWTSecret: getenv("JWT_SECRET", "cyberkube-secret-please-change"),
        AESKey:    getenv("K8S_KUBECONFIG_AES_KEY", "0123456789abcdef0123456789abcdef"),
        SkipTLSVerify: getenv("K8S_SKIP_TLS_VERIFY", "false") == "true" || getenv("K8S_SKIP_TLS_VERIFY", "0") == "1",
    }
    return c
}

func getenv(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}
