package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var LogFile *os.File

const (
    RED    = "\033[31m"
    GREEN  = "\033[32m"
    YELLOW = "\033[33m"
    CYAN   = "\033[36m"
    RESET  = "\033[0m"
)

// ==================== fonction log juste avec message ==================== //

func Success(msg string) {
	log.Println("✅ [SUCCESS] " + msg)
}

func Info(msg string) {
	log.Println("ℹ️  [INFO] " + msg)
}

func Warn(msg string) {
	log.Println("⚠️ [WARN] " + msg)
}

func Error(msg string) {
	log.Panic("❌ [ERROR] " + msg) 
}

func Fatal(msg string) {
	log.Fatal("💀 [FATAL] " + msg)
}

// =========== fonction log avec message et passage de variables ============ //

func Successf(format string, a ...any) {
	log.Printf("✅ [SUCCESS] "+format, a...)
}

func Infof(format string, a ...any) {
	log.Printf("ℹ️  [INFO] "+format, a...)
}

func Warnf(format string, a ...any) {
	log.Printf("⚠️ [WARN] "+format, a...)
}

func Errorf(format string, a ...any) {
	log.Panicf("❌ [ERROR] "+format, a...)
}

func Fatalf(format string, a ...any) {
	log.Fatalf("💀 [FATAL] "+format, a...)
} 

func InitFromEnv(env string) {
	logDir := filepath.Join("runtime", "logs", env)
	logPath := filepath.Join(logDir, "todof-"+time.Now().Format("2006-01-02")+".log")

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("❌ [FATAL] Impossible de créer le dossier de logs : %v", err)
	}

	log.Printf("ℹ️  [INFO] Creation du fichier de log : %s", logPath)

	LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("❌ [FATAL] Impossible d’ouvrir le fichier de log : %v", err)
	}

	log.Printf("✅ [SUCCESS] Fichier de log ouvert avec succès")

	multi := io.MultiWriter(os.Stdout, LogFile)
	log.SetOutput(multi)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("✅ [SUCCESS] Logger initialisé avec succès")
}

func Init() {
	env := os.Getenv("APP_ENV")
	
	if env == "" { 
		env = "dev"
	}
	
	InitFromEnv(env)
}

func Close() {
	if LogFile != nil {
		LogFile.Close()
	}
}
