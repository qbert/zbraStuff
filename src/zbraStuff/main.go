package main

import (
	"net/http"
	"syscall"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "zbraStuff/models"
)

var  (
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

func config_env() {
	port = get_env("OPENSHIFT_GO_PORT", "C9_PORT", "1323")
	ip   = get_env("OPENSHIFT_GO_IP",   "C9_IP",   "localhost")
	
	dbuser     = get_env("OPENSHIFT_MYSQL_DB_USERNAME", "C9_USER", "qbert")
	dbpassword = get_env("OPENSHIFT_MYSQL_DB_PASSWORD", "C9_DB_PASSWORD", "neinegal")
	dbhost     = get_env("OPENSHIFT_MYSQL_DB_HOST", "C9_DB_HOST", "localhost")
	dbport     = get_env("OPENSHIFT_MYSQL_DB_PORT", "C9_DB_PORT", "3306")
	dbname     = get_env("OPENSHIFT_MYSQL_DB_PASSWORD", "C9_DB_NAME", "zbrastuff_db")
}

func config_db() {
	var dbconnect string
	
	dbconnect = dbuser
	if dbpassword != "" {
		dbconnect += ":"+dbpassword
	}
	dbconnect += "@"
	
	if dbhost != "" {
		dbconnect += dbhost
		
	}
	if dbport != "" {
		dbconnect += ":"+dbport
	}
	dbconnect += "/"+dbname+"?charset=utf8&parseTime=True"
	

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
 
    if !dbm.HasTable(&models.User{}){
        dbm.CreateTable(&models.User{})
    }
    
    db = &dbm
}

func config_echo() {
	// Echo instance
	e = echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/", hello)
}

func run_echo() {
	e.Run(ip+":"+port)
}

func main() {

	config_env()
	config_db()
	config_echo()
	
	run_echo()
}

func get_env(e1, e2, d string) string {

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
