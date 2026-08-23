package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	UDPAddr         string
	NodeID          string
	ShutdownTimeout time.Duration
	MaxRooms        int
	MaxParticipants int
}

func Load() Config {
	c := Config{HTTPAddr: env("HARBOR_HTTP_ADDR", ":8093"), UDPAddr: env("HARBOR_UDP_ADDR", ":10000"), NodeID: env("HARBOR_NODE_ID", "node-local"), ShutdownTimeout: 10 * time.Second, MaxRooms: 100, MaxParticipants: 500}
	if v := os.Getenv("HARBOR_MAX_ROOMS"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.MaxRooms = n
		}
	}
	return c
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
