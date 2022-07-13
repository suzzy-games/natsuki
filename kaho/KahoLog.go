package kaho

import (
	"context"
	"encoding/json"
	"fmt"
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
	Data      any       `json:"data"`
}

var LogToConsole bool = (os.Getenv("NATSUKI_KAHO_PRINT") != "")
var BroadcastLog bool = (os.Getenv("NATSUKI_KAHO_BROADCAST") != "")
var StoreLogInDB bool = ShouldStoreInDb()
var InitializationQuery = `
-- Create Kaho Schema
CREATE SCHEMA IF NOT EXISTS kaho;

-- Create Kaho Entries Table
CREATE TABLE IF NOT EXISTS kaho.entries (
	id SERIAL NOT NULL PRIMARY KEY,
	"timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
	severity text,
	service text,
	message text,
	data text
) WITH (autovacuum_enabled='true');`

func ShouldStoreInDb() bool {
	enabled := (os.Getenv("NATSUKI_KAHO_STORE") != "")

	if enabled {
		// Create kaho.entries Table
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		_, err := db.Postgres.Exec(ctx, InitializationQuery)

		// Unable to Initialize kaho.entries Table
		if err != nil {
			log.Printf("[SQL][ERROR] unable to create table 'kaho.entries': %s", err.Error())
			return false
		}
	}

	return enabled
}

func Log(severity Severity, service string, message string, data any) {
	KahoLogRaw(KahoLogEntry{
		Timestamp: time.Now(),
		Severity:  severity,
		Message:   message,
		Service:   service,
		Data:      data,
	})
}

func KahoLogRaw(entry KahoLogEntry) error {
	// Log to Console (if Enabled)
	if LogToConsole {
		if entry.Data != nil {
			log.Printf("[%s][%s] %s • %v", entry.Service, entry.Severity, entry.Message, entry.Data)
		} else {
			log.Printf("[%s][%s] %s ", entry.Service, entry.Severity, entry.Message)
		}
	}

	// Write to Database (if Enabled)
	if StoreLogInDB {
		// Convert Data to String
		jsonString, err := json.Marshal(entry.Data)
		if err != nil {
			return err
		}

		// Insert Entry into Database
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		results, err := db.Postgres.Query(
			ctx, "INSERT INTO kaho.entries (severity, service, message, data) VALUES ($1, $2, $3, $4) RETURNING id;",
			string(entry.Severity), entry.Service, entry.Message, string(jsonString),
		)
		if err != nil {
			results.Close()
			return err
		}

		// Get Insert Id
		var insertId int64
		results.Next()
		results.Scan(&insertId)
		results.Close()

		// Publish Insert Id to channel 'kaho:create'
		if BroadcastLog {
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if err := db.Redis.Do(ctx, radix.Cmd(nil, "PUBLISH", "kaho:create", fmt.Sprintf("%v", insertId))); err != nil {
				return err
			}
		}
	}

	// Exit if Fatal Error
	if entry.Severity == Severity(FATAL) {
		os.Exit(1)
	}

	return nil
}
