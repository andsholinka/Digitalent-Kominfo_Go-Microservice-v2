package main

import (
	"fmt"

	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/service-product/config"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/service-product/database"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/service-product/handler"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Panic(err)
		return
	}

	db, err := initDB(cfg.Database)

	router := mux.NewRouter()

	authMiddleware := handler.AuthMiddleware{
		AuthService: cfg.AuthService,
	}
	menuHandler := handler.Menu{Db: db}

	router.Handle("/add-menu", authMiddleware.ValidateAuthAdmin(http.HandlerFunc(menuHandler.AddMenu)))
	router.Handle("/menu", authMiddleware.ValidateAuth(menuHandler.GetAllMenu))

	fmt.Printf("Server listen on :%s", cfg.Port)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}

func getConfig() (config.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, err
	}

	var cfg config.Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

func initDB(dbConfig config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName, dbConfig.Config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&database.Menu{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
