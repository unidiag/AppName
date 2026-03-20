package main

import (
	"fmt"
	"log"
	"main/models"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// -------------------- DB INIT --------------------

// initDB opens (or creates) db.sqlite3 in project root via GORM
func initDB() bool {

	if len(os.Args) > 1 {
		dbname = os.Args[1]
	}

	firstRun := false
	var setupCfg *SetupConfig

	if _, err := os.Stat(dbname); os.IsNotExist(err) {
		firstRun = true

		var ok bool
		setupCfg, ok = runSetupWizard()
		if !ok {
			log.Fatal("Setup cancelled")
		}
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
			LogLevel: func() logger.LogLevel {
				if debug {
					return logger.Warn
				}
				return logger.Error
			}(),
		},
	)

	var err error
	db, err = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		slog("DB open error: "+err.Error(), "error")
		log.Fatal(err)
	}

	if err = db.AutoMigrate(models.AllModels...); err != nil {
		slog("DB migrate error: "+err.Error(), "error")
		log.Fatal(err)
	}

	// save setup settings
	if firstRun && setupCfg != nil {

		createDefaultUser("admin", "admin", 9)
		//createDefaultUser("user", "user", 2)

		setSetting("port", setupCfg.Listen, "Web server listen address")
		setSetting("whitelist", "", "ALLOWED addresses and countries (separated by commas)")
		setSetting("blacklist", "CN, 178.24.52.115", "PROHIBITED addresses and countries (separated by commas)")

		if setupCfg.InstallSvc {
			// вычисляем полный путь
			exe, _ := os.Executable()
			exe, _ = filepath.EvalSymlinks(exe)
			svc, err := installSystemdService(exe, os.Args[1:])
			if err != nil {
				fmt.Println("Systemd install error:", err)
			} else {
				fmt.Println()
				fmt.Println("Created and started: " + svc)
				fmt.Printf("Open WebUI: http://%s:%s (admin/admin)\n", getLocalIP(), setupCfg.Listen)
				fmt.Println("Application will now continue running in background.")
				os.Exit(0)
			}
		}
	}

	loadSettings()

	return firstRun
}

func openDatabase(dsn string, gormLogger logger.Interface) (*gorm.DB, error) {
	dsn = strings.TrimSpace(dsn)
	// --- SQLite ---
	if !strings.Contains(dsn, " ") &&
		(strings.HasSuffix(dsn, ".db") ||
			strings.HasSuffix(dsn, ".sqlite") ||
			strings.HasSuffix(dsn, ".sqlite3") ||
			!strings.Contains(dsn, "@")) {
		log.Println("Using SQLite: ", dsn)
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
	}
	// --- PostgreSQL ---
	if strings.Contains(dsn, " ") {
		log.Println("Using PostgreSQL")
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
	}
	// --- MariaDB / MySQL ---
	log.Println("Using MySQL/MariaDB")
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
}
