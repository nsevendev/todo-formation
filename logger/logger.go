package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

// ==================== fonction logger ==================== //

func callerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "???"
	}
	fn := runtime.FuncForPC(pc)
	funcName := filepath.Base(fn.Name())
	fileName := filepath.Base(file)
	return fmt.Sprintf("%s:%s:%d", fileName, funcName, line)
}

func logWithLocation(level string, emoji string, msg string) {
	location := callerInfo(3)
	stringFormat := fmt.Sprintf("%s [%s] [%s] %s", emoji, level, location, msg)
	log.Println(stringFormat)
}

func logfWithLocation(level string, emoji string, format string, args ...any) {
	location := callerInfo(3)
	log.Printf("%s [%s] [%s] %s", emoji, level, location, fmt.Sprintf(format, args...))
}

func Success(msg string) {
	logWithLocation("SUCCESS", "✅", msg)
}

func Info(msg string) {
	logWithLocation("INFO", "ℹ️ ", msg)
}

func Warn(msg string) {
	logWithLocation("WARN", "⚠️ ", msg)
}

func Error(msg string) {
	logWithLocation("ERROR", "❌", msg)
}

func Fatal(msg string) {
	location := callerInfo(3)
	log.Fatalf("💀 [FATAL] [%s] %s", location, msg)
}

func Successf(format string, args ...any) {
	logfWithLocation("SUCCESS", "✅", format, args...)
}

func Infof(format string, args ...any) {
	logfWithLocation("INFO", "ℹ️ ", format, args...)
}

func Warnf(format string, args ...any) {
	logfWithLocation("WARN", "⚠️ ", format, args...)
}

func Errorf(format string, args ...any) {
	logfWithLocation("ERROR", "❌", format, args...)
}

func Fatalf(format string, args ...any) {
	location := callerInfo(3)
	log.Fatalf("💀 [FATAL] [%s] %s", location, fmt.Sprintf(format, args...))
}

// ==================== init ==================== //

func InitFromEnv(env string) {
	logDir := filepath.Join("tmp", "log", env)
	logPath := filepath.Join(logDir, "log-"+time.Now().Format("2006-01-02")+".log")

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("❌ [FATAL] Impossible de créer le dossier de log : %v", err)
	}

	log.Printf("ℹ️  [INFO] Creation du fichier de log : %s", logPath)

	LogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("❌ [FATAL] Impossible d’ouvrir le fichier de log : %v", err)
	}

	log.Printf("ℹ️  [INFO] Fichier de log ouvert")

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
