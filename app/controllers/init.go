package controllers

import (
	"database/sql"
	"fmt"
	"revel101/app/models"
	"strings"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql" // MySQL Dialect
	"github.com/revel/revel"
)

func init(){
    revel.OnAppStart(InitDb)
    revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
    revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
    revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

func getParamString(param string, defaultValue string) string {
    p, found := revel.Config.String(param)
    if !found {
        if defaultValue == "" {
            revel.ERROR.Fatal("Could not find parameter: " + param)
        } else {
            return defaultValue
        }
    }
    return p
}

func getConnectionString() string {
    host := getParamString("db.host", "")
    port := getParamString("db.port", "3306")
    user := getParamString("db.user", "")
    pass := getParamString("db.password", "")
    dbname := getParamString("db.name", "golangdb")
    protocol := getParamString("db.protocol", "tcp")
    dbargs := getParamString("dbargs", " ")

    if strings.Trim(dbargs, " ") != "" {
        dbargs = "?" + dbargs
    } else {
        dbargs = ""
    }
    return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", 
        user, pass, protocol, host, port, dbname, dbargs)
}

func defineUserTable(dbm *gorp.DbMap) {
  // set "id" as primary key and autoincrement
  t := dbm.AddTable(models.User{}).SetKeys(true, "id")
  // e.g. VARCHAR(25)
  t.ColMap("first_name").SetMaxSize(25)
  t.ColMap("last_name").SetMaxSize(25)
}

var InitDb func() = func() {
  connectionString := getConnectionString()
  if db, err := sql.Open("mysql", connectionString); err != nil {
    revel.ERROR.Fatal(err)
  } else {
    Dbm = &gorp.DbMap{
      Db: db,
      Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
  }
  // Defines the table for use by GORP
  // This is a function we will create soon.
  defineUserTable(Dbm)
  if err := Dbm.CreateTablesIfNotExists(); err != nil {
    revel.ERROR.Fatal(err)
  }
}
