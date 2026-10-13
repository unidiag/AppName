package main

import (
	"log"
	"log/syslog"
	"main/models"
	"os"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	llama "github.com/unidiag/go-llama"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const APPNAME = "TVAPLAY"
const APPLINK = "http://github.com/unidiag/tvaplay"
const VERSION = "1.03"
const BUILD_DATE = "2026-10-13"
const BUILD_TIME = "17:27:48"

var (
	mu        sync.Mutex
	debug     = false
	err       error
	settings  map[string]string
	db        *gorm.DB
	sysLogger *syslog.Writer // для slog(..)
	runtime   time.Time
	dbip      *geoip2.Reader

	jwtSecret  = []byte("CHANGE_ME_LONG_RANDOM_SECRET_1234")
	accessTTL  = 10 * time.Minute
	refreshTTL = 14 * 24 * time.Hour

	llamaClient *llama.Client
)

// ███╗   ███╗ █████╗ ██╗███╗   ██╗
// ████╗ ████║██╔══██╗██║████╗  ██║
// ██╔████╔██║███████║██║██╔██╗ ██║
// ██║╚██╔╝██║██╔══██║██║██║╚██╗██║
// ██║ ╚═╝ ██║██║  ██║██║██║ ╚████║
// ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝

func main() {

	runtime := time.Now()
	_ = runtime

	sysLogger, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, APPNAME)
	if err != nil {
		log.Println("syslog init error:", err)
	}

	debug = isRunThroughGoRun()
	slog("Server run in DEBUG-mode", "debug")

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

	db, err = gorm.Open(
		mysql.Open(DSN),
		&gorm.Config{
			Logger: gormLogger,
		},
	)

	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err = db.AutoMigrate(models.AllModels...); err != nil {
		slog("DB migrate error: "+err.Error(), "error")
		log.Fatal(err)
	}

	initGeoIP() // использовать базу данных IP. можно закомментировать, тогда не используется ./geoip.mmdb

	go webserver()

	for {
		delay(1000)
	}

}
