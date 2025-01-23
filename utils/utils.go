package utils  

import (  
    "golang.org/x/crypto/bcrypt"  
)  

// HashPassword хеширует пароль  
func HashPassword(password string) string {  
    bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)  
    return string(bytes)  
}  

// CheckPassword проверяет совпадение пароля  
func CheckPassword(hashedPassword, password string) bool {  
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))  
    return err == nil  
}