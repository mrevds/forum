package middlewares  

import (  
    "net/http"  
    "github.com/dgrijalva/jwt-go"  
)  

func TokenAuthMiddleware(next http.Handler) http.Handler {  
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        tokenString := r.Header.Get("Authorization")  
        claims := &jwt.StandardClaims{}  
        
        _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {  
            return []byte("secret"), nil  
        })  
        
        if err != nil {  
            http.Error(w, "Unauthorized", http.StatusUnauthorized)  
            return  
        }  
        
        next.ServeHTTP(w, r)  
    })  
}