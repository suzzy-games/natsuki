package kaho

import (
	"context"
	"log"
	"natsuki/db"
	"os"
	"time"

	"github.com/mediocregopher/radix/v4"
)

type Severity string

const (
	DEFAULT   Severity = "DEFAULT"
	DEBUG     Severity = "DEBUG"
	INFO      Severity = "INFO"
	NOTICE    Severity = "NOTICE"
	WARNING   Severity = "WARNING"
	ERROR     Severity = "ERROR"
	CRITICAL  Severity = "CRITICAL"
	ALERT     Severity = "ALERT"
	EMERGENCY Severity = "EMERGENCY"
	FATAL     Severity = "FATAL"
)

type KahoLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Severity  Severity  `json:"severity"`
	Service   string    `json:"service"`
	Message   string    `json:"message"`
	Payload   any       `json:"payload"`
}

var LogToConsole bool = (os.Getenv("NATSUKI_KAHO_PRINT") != "")
var StoreLogRedis bool = (os.Getenv("NATSUKI_KAHO_ENABLE") != "")

func Log(severity Severity, service string, message string, payload any) {
	KahoLogRaw(KahoLogEntry{
		Timestamp: time.Now(),
		Severity:  severity,
		Message:   message,
		Service:   service,
		Payload:   payload,
	})
}

func KahoLogRaw(entry KahoLogEntry) error {

	// Log to Console (if specified)
	if LogToConsole {
		log.Printf("[%s][%s] %s", entry.Service, entry.Severity, entry.Message)
	}

	// Do Not Store in Redis (if specified)
	if StoreLogRedis {
		return nil
	}

	// Write Log to Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Redis.Do(ctx, radix.FlatCmd(nil, "LPUSH", "kaho:entries", entry)); err != nil {
		return err
	}

	// Trim Entries to recent 100,000
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Redis.Do(ctx, radix.Cmd(nil, "LTRIM", "kaho:entries", "0", "100000")); err != nil {
		return err
	}

	// Publish new entry event (we dont publish the entry data because that would be extra wasteful if nobody is tailing)
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Redis.Do(ctx, radix.Cmd(nil, "PUBLISH", "kaho:entries", "")); err != nil {
		return err
	}

	// Exit if Fatal Error
	if entry.Severity == Severity(FATAL) {
		os.Exit(1)
	}

	return nil
}
