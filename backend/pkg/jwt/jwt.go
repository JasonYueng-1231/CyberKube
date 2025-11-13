package jwt

import (
    "time"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/config"
    jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwtv5.RegisteredClaims
}

func Sign(userID uint, username, role string, ttl time.Duration) (string, error) {
    c := &Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwtv5.RegisteredClaims{
            ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(ttl)),
            IssuedAt:  jwtv5.NewNumericDate(time.Now()),
            NotBefore: jwtv5.NewNumericDate(time.Now()),
        },
    }
    token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, c)
    secret := []byte(config.Load().JWTSecret)
    return token.SignedString(secret)
}

func Parse(tokenStr string) (*Claims, error) {
    secret := []byte(config.Load().JWTSecret)
    token, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtv5.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil { return nil, err }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwtv5.ErrTokenInvalidClaims
}

