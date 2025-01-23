package config

import (
	"fmt"
	"log"
	"os"

	"forum/models"

	"gopkg.in/yaml.v3" // Используем yaml.v3
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Database struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		Dbname         string `yaml:"dbname"`
		Sslmode        string `yaml:"sslmode"`
		ClientEncoding string `yaml:"client_encoding"`
	} `yaml:"database"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

var DB *gorm.DB
var Cfg Config

func LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла конфигурации: %w", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		return fmt.Errorf("ошибка декодирования YAML: %w", err)
	}
	return nil
}

func ConnectDatabase() {
	err := LoadConfig("config/config.yaml") // Загружаем конфиг
	if err != nil {
		log.Fatal(err) // Обрабатываем ошибку загрузки конфига
	}

	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=%s client_encoding=%s",
		Cfg.Database.User, Cfg.Database.Dbname, Cfg.Database.Password, Cfg.Database.Host, Cfg.Database.Port, Cfg.Database.Sslmode, Cfg.Database.ClientEncoding)

	var errDb error
	DB, errDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDb != nil {
		panic(fmt.Sprintf("failed to connect database: %s", errDb)) // More informative error message
	}

	errDb = DB.AutoMigrate(&models.User{}, &models.Post{})
	if errDb != nil {
		panic(fmt.Sprintf("failed to migrate database: %s", errDb)) // More informative error message
	}

	fmt.Println("Database connected")
}
