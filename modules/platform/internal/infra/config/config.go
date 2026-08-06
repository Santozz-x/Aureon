package config

import "os"

const defaultArcRPCURL = "https://rpc.testnet.arc.io"

type Config struct {
	HTTPAddr  string
	ArcRPCURL string
}

func Load() Config {
	return Config{
		HTTPAddr:  getenv("AUREON_HTTP_ADDR", ":8080"),
		ArcRPCURL: getenv("AUREON_ARC_RPC_URL", defaultArcRPCURL),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
