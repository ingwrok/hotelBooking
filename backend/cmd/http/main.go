package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/handlers"
	"github.com/ingwrok/hotelBooking/internal/adapters/primary/web/routes"
	"github.com/ingwrok/hotelBooking/internal/adapters/secondary/postgresql"
	"github.com/ingwrok/hotelBooking/internal/common/logger"
	"github.com/ingwrok/hotelBooking/internal/core/services"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
    initTimeZone()
    initConfig()

    db := initDatabase()
    defer db.Close()

		// Repos
		roomRepo := postgresql.NewRoomRepository(db)
		amenityRepo := postgresql.NewAmenityRepository(db)
		roomTypeRepo := postgresql.NewRoomTypeRepository(db)


		// Services
		roomSvc := services.NewRoomService(roomRepo)
		amenitySvc := services.NewAmenityService(amenityRepo)
		roomTypeSvc := services.NewRoomTypeService(roomTypeRepo)

		// Handlers
		roomHandler := handlers.NewRoomHandler(roomSvc)
		amenityHandler := handlers.NewAmenityHandler(amenitySvc)
		roomTypeHandler := handlers.NewRoomTypeHandler(roomTypeSvc)


		// Server
		app := fiber.New()
		app.Use(recover.New(), fiberlogger.New())
		app.Use(cors.New(cors.Config{
			AllowOrigins: "http://localhost:3000, https://yourdomain.com",
			AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
			ExposeHeaders: "Content-Length",
			AllowCredentials: true,
		}))

		// Routes
		routes.RoomRoutes(app,roomHandler)
		routes.AmenityRoutes(app, amenityHandler)
		routes.RoomTypeRoutes(app, roomTypeHandler)

		addr := fmt.Sprintf(":%d", viper.GetInt("app.port"))
		logger.Info("Hotel service starting at port " + viper.GetString("app.port"))
		if err := app.Listen(addr); err != nil {
			logger.ErrorErr(err, "Failed to start server")
		}
}


func initConfig() {
	if err := godotenv.Load();err != nil {
        fmt.Println("No .env file found, using system env instead")
  }

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")

	viper.BindEnv("db.driver", "DB_DRIVER")
  viper.BindEnv("db.host", "DB_HOST")
  viper.BindEnv("db.port", "DB_PORT")
  viper.BindEnv("db.database", "DB_NAME")
  viper.BindEnv("db.username", "DB_USER")
  viper.BindEnv("db.password", "DB_PASSWORD")
  viper.BindEnv("db.sslmode", "DB_SSLMODE")
  viper.BindEnv("secret", "APP_SECRET")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.port", 8000)
  viper.SetDefault("db.driver", "pgx")
  viper.SetDefault("db.host", "localhost")
  viper.SetDefault("db.port", 5432)

	if err := viper.ReadInConfig(); err != nil {
		// ไม่มีไฟล์ config ก็ยังรันได้ด้วยค่า env/default
		logger.ErrorErr(err, "Failed to read config file, using env/default values")
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