package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	heartbeat "github.com/qbert/heartbeat-golang"
	"github.com/qbert/zbraStuff/models"
	"net/http"
	"syscall"
)

var (
	port       string
	ip         string
	e          *echo.Echo
	dbuser     string
	dbpassword string
	dbhost     string
	dbport     string
	dbname     string
	db         *gorm.DB
)

func configEnv() {
	port = getEnv("OPENSHIFT_GO_PORT", "C9_PORT", "1323")
	ip = getEnv("OPENSHIFT_GO_IP", "C9_IP", "localhost")

	dbuser = getEnv("OPENSHIFT_MYSQL_DB_USERNAME", "C9_USER", "qbert")
	dbpassword = getEnv("OPENSHIFT_MYSQL_DB_PASSWORD", "C9_DB_PASSWORD", "neinegal")
	dbhost = getEnv("OPENSHIFT_MYSQL_DB_HOST", "C9_DB_HOST", "localhost")
	dbport = getEnv("OPENSHIFT_MYSQL_DB_PORT", "C9_DB_PORT", "3306")
	dbname = getEnv("OPENSHIFT_MYSQL_DB_NAME", "C9_DB_NAME", "zbrastuff_db")
}

func configDb() {
	var dbconnect string

	dbconnect = dbuser
	if dbpassword != "" {
		dbconnect += ":" + dbpassword
	}
	dbconnect += "@"

	if dbhost != "" {
		dbconnect += "tcp(" + dbhost
		if dbport != "" {
			dbconnect += ":" + dbport
		}
		dbconnect += ")"
	}
	dbconnect += "/" + dbname + "?charset=utf8&parseTime=True"

	log.Info("Attempting to connect to db with:%s", dbconnect)

	dbm, err := gorm.Open("mysql", dbconnect)
	if err != nil {
		panic("Unable to connect to the database")
	}

	dbm.DB().Ping()
	dbm.DB().SetMaxIdleConns(10)
	dbm.DB().SetMaxOpenConns(100)
	dbm.LogMode(true)
	dbm.SingularTable(true)
	dbm.Set("gorm:table_options", "ENGINE=InnoDB")
	// and fix db on startup
	dbm.AutoMigrate(&models.User{})

	db = &dbm
}

func configEcho() {
	// Echo instance
	e = echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/", hello)
}

func runHeartbeat() {
	log.Info("Running Heartbeat on:%s:%s/heartbeat", ip, port)
	e.Get("/heartbeat", heartbeat.Handler)
}

func runEcho() {
	log.Info("Running Echo Server on:%s:%s", ip, port)
	e.Run(ip + ":" + port)
}

func main() {

	configEnv()
	configDb()
	configEcho()

	runHeartbeat()
	runEcho()
}

func getEnv(e1, e2, d string) string {

	if r, found := syscall.Getenv(e1); found {
		return r
	}

	if r, found := syscall.Getenv(e2); found {
		return r
	}

	return d
}

// Handler
func hello(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! I was here\n")
}
