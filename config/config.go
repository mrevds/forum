package config  

import (  
    "gorm.io/driver/postgres"  
    "gorm.io/gorm"  
    "forum/models" // Импортируй модели User и Post  
)  

var DB *gorm.DB  

func ConnectDatabase() {  
    dsn := "user=postgres dbname=forum password=12345678 host=localhost port=5432 sslmode=disable"  
    var err error  
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})  
    if err != nil {  
        panic("failed to connect database")  
    }  

    // Автоматическая миграция моделей  
    err = DB.AutoMigrate(&models.User{}, &models.Post{})  
    if err != nil {  
        panic("failed to migrate database")  
    }  
}