package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/ingwrok/hotelBooking/internal/common/logs"
    _ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
    initTimeZone()
    initConfig()

    db := initDatabase()
    defer db.Close()


}


func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.port", 8000)
    viper.SetDefault("db.driver", "postgres")
    viper.SetDefault("db.host", "localhost")
    viper.SetDefault("db.port", 5432)

	if err := viper.ReadInConfig(); err != nil {
		// ไม่มีไฟล์ config ก็ยังรันได้ด้วยค่า env/default
		logs.Error(err)
	}
}

func initTimeZone() {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil { panic(err) }
	time.Local = loc
}

func initDatabase() *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Bangkok",
    viper.GetString("db.host"),
    viper.GetString("db.port"),
    viper.GetString("db.username"),
    viper.GetString("db.password"),
    viper.GetString("db.database"),
    viper.GetString("db.sslmode"),
    )
	db, err := sqlx.Open(viper.GetString("db.driver"), dsn)
	if err != nil { panic(err) }

	if err := db.Ping(); err != nil { panic(err) }

	db.SetConnMaxLifetime(3 * time.Hour)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}