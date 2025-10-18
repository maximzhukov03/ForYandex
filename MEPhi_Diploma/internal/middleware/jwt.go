package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func JWTAuth(secret string) gin.HandlerFunc{
    return func(c *gin.Context){
        authHeader := c.GetHeader("Authorization")
        if authHeader == ""{
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
            return
        }
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer"{
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
            return
        }
        tokenStr := parts[1]

        secret := os.Getenv("JWT_SECRET")
        if secret == ""{
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
            return
        }
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error){
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(secret), nil
        })
        if err != nil || !token.Valid{
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            return
        }
        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
            return
        }

        c.Set("userID", int64(userIDFloat))
        c.Next()
    }
}