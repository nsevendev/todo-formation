package config

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"sync"
	"todof/mod/projectroot"
)

var once sync.Once

func loadEnvOnce() {
	envPath := filepath.Join(projectroot.Root(), ".env")
	_ = godotenv.Load(envPath)
}

func Get(key string) string {
	once.Do(loadEnvOnce)
	return os.Getenv(key)
}
