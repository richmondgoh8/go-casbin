package custommiddleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/richmondgoh8/go-casbin/pkg/middleware/models"
	"github.com/richmondgoh8/go-casbin/pkg/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var SecretKey []byte

const (
	ENV_SECRET string = "SECRET"
	NO_AUTH    string = "unauthorized access"
	ClaimsKey  string = "jwt_claims"
)

// Strict Auth, does not allow invalid tokens
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		setSecret()
		msg := map[string]interface{}{}
		if len(SecretKey) <= 0 {
			utils.JSONError(ctx, w, "secret has not been configured", http.StatusUnauthorized)
			return
		}

		if auth := r.Header.Get("Authorization"); auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			utils.JSONError(ctx, w, "malformed or missing token", http.StatusUnauthorized)
			return
		}

		jwtToken := authHeader[1]
		token, err := jwt.ParseWithClaims(jwtToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			// verify the signing method
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// return secret key for validating the token
			return SecretKey, nil
		})
		if err != nil {
			utils.JSONError(ctx, w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok || !token.Valid {
			msg["message"] = NO_AUTH
			utils.JSONError(ctx, w, NO_AUTH, http.StatusUnauthorized)
			return
		}

		// set to context map
		ctx = context.WithValue(ctx, ClaimsKey, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GenerateToken(jwtPayload models.JWTPayload, customExpTimeInMinute int) (string, error) {
	setSecret()
	if customExpTimeInMinute <= 0 {
		return "", errors.New("Expiry Time cannot be 0")
	}

	// Create a map to set & store our token claims
	claims := jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Minute * time.Duration(customExpTimeInMinute)).Unix(),
		"role": jwtPayload.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}

	return tokenString, nil
}

func GetClaimsFromJWTTokenCtx(ctx context.Context) (*models.Claims, error) {
	claims, ok := ctx.Value(ClaimsKey).(*models.Claims)
	if !ok {
		return nil, errors.New("no claims")
	}
	return claims, nil
}

// Set Secret If Secret is Not Initialized
func setSecret() {
	var m sync.Mutex
	if len(SecretKey) <= 0 {
		m.Lock()
		defer m.Unlock()
		SecretKey = []byte(os.Getenv(ENV_SECRET))
	}
}
