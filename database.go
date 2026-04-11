package main

import (
	"errors"
	"fmt"
	"log"
	"main/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
}

// -------------------- DB INIT --------------------

// initDB opens (or creates) db.sqlite3 in project root via GORM
func initDB() bool {

	// если база нужна и передана в параметре
	if len(os.Args) > 1 {
		dbname = os.Args[1]
	}

	var setupCfg *SetupConfig
	firstRun := detectFirstRun(dbname)

	if firstRun {
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
	db, err = openDatabase(dbname, gormLogger)
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

		createDefaultUser("root", "root", 9)
		createDefaultUser("admin", "admin", 2)
		createDefaultUser("user", "user", 1)

		setSetting("port", setupCfg.Listen, "Web server listen address", "1")
		setSetting("whitelist", "", "ALLOWED addresses and countries (separated by commas)", "2")
		setSetting("blacklist", "CN, 178.24.52.115", "PROHIBITED addresses and countries (separated by commas)", "3")

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
	switch typeDatabase(dsn) {
	// --- SQLite ---
	case "sqlite":
		log.Println("Using SQLite3: ", dsn)
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

	// --- PostgreSQL ---
	case "postgresql":
		log.Println("Using PostgreSQL")
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
	// --- MariaDB / MySQL ---
	case "mysql":
		log.Println("Using MySQL/MariaDB")
		return gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
	}
	return nil, errors.New("Unknown detect type database")
}

func typeDatabase(dsn string) string {
	dsn = strings.TrimSpace(dsn)
	// --- SQLite ---
	if !strings.Contains(dsn, " ") &&
		(strings.HasSuffix(dsn, ".db") ||
			strings.HasSuffix(dsn, ".sqlite") ||
			strings.HasSuffix(dsn, ".sqlite3") ||
			!strings.Contains(dsn, "@")) {
		return "sqlite"
	}
	// --- PostgreSQL ---
	if strings.Contains(dsn, " ") {
		return "postgersql"
	}
	// --- MariaDB / MySQL ---
	if strings.Contains(dsn, "@") && strings.Contains(dsn, ":") {
		return "mysql"
	}

	return "unknown"
}

func detectFirstRun(dsn string) bool {
	var err error

	switch typeDatabase(dsn) {

	case "mysql":
		_, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			cfg := parseMySQL(dsn)
			createMySQLDB(cfg)
			return true
		}

	case "postgresql":
		_, err = gorm.Open(postgres.Open(dsn))
		if err != nil {
			cfg := parsePostgres(dsn)
			createPostgresDB(cfg)
			return true
		}

	case "sqlite":
		if _, err := os.Stat(dbname); os.IsNotExist(err) {
			return true
		}
	}

	return false
}

func createPostgresDB(cfg DBConfig) error {
	cmds := []string{
		"CREATE DATABASE " + cfg.DBName + ";",
		"DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '" + cfg.User + "') THEN CREATE USER " + cfg.User + " WITH PASSWORD '" + cfg.Password + "'; END IF; END $$;",
		"GRANT ALL PRIVILEGES ON DATABASE " + cfg.DBName + " TO " + cfg.User + ";",
	}

	for _, sql := range cmds {
		cmd := exec.Command("sudo", "-u", "postgres", "psql", "-c", sql)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func createMySQLDB(cfg DBConfig) error {
	sql := `
CREATE DATABASE IF NOT EXISTS ` + cfg.DBName + ` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '` + cfg.User + `'@'localhost' IDENTIFIED BY '` + cfg.Password + `';
GRANT ALL PRIVILEGES ON ` + cfg.DBName + `.* TO '` + cfg.User + `'@'localhost';
FLUSH PRIVILEGES;
`

	cmd := exec.Command("sudo", "mysql", "-e", sql)
	return cmd.Run()
}

// MySQL: user:pass@tcp(host:port)/dbname?...
func parseMySQL(dsn string) DBConfig {
	cfg := DBConfig{}

	parts := strings.Split(dsn, "@tcp(")
	if len(parts) != 2 {
		return cfg
	}

	// user:pass
	up := strings.Split(parts[0], ":")
	cfg.User = up[0]
	if len(up) > 1 {
		cfg.Password = up[1]
	}

	// host:port)/dbname
	rest := parts[1]
	i := strings.Index(rest, ")")
	hostport := rest[:i]
	dbpart := rest[i+2:] // skip ")/"

	hp := strings.Split(hostport, ":")
	cfg.Host = hp[0]
	cfg.Port = hp[1]

	db := strings.Split(dbpart, "?")[0]
	cfg.DBName = db

	return cfg
}

// PostgreSQL: key=value
func parsePostgres(dsn string) DBConfig {
	cfg := DBConfig{}

	parts := strings.Fields(dsn)
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "user":
			cfg.User = kv[1]
		case "password":
			cfg.Password = kv[1]
		case "dbname":
			cfg.DBName = kv[1]
		case "host":
			cfg.Host = kv[1]
		case "port":
			cfg.Port = kv[1]
		}
	}

	return cfg
}
