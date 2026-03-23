package main

import (
	"log"
	"log/syslog"
	"os"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	llama "github.com/unidiag/go-llama"
	"gorm.io/gorm"
)

const APPNAME = "AppName"
const APPLINK = "http://github.com/unidiag/AppName"
const VERSION = "1.00"
const BUILD_DATE = "2026-03-21"
const BUILD_TIME = "09:07:43"

// pgDSN="host=localhost user=epg password=epg dbname=epg port=5432 sslmode=disable TimeZone=Europe/Minsk"
// myDSN="dvbscan:dvbscan@tcp(127.0.0.1:3306)/dvbscan?charset=utf8mb4&parseTime=True&loc=Local"
// var dbname = ""
var dbname = "db.sqlite3"

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

	if dbname != "" {
		initDB() // инициализация базы данных и wizard, если запуск впервые
	} else {
		port := ":9000"
		if len(os.Args) > 1 {
			port = os.Args[1]
		}
		setSetting("port", port)
	}
	initGeoIP()                                                        // использовать базу данных IP. можно закомментировать, тогда не используется ./geoip.mmdb
	llamaClient = llama.New(getSetting("ailink"), getSetting("aikey")) // если нужен ии

	go webserver()

	for {
		delay(1000)
	}

}
