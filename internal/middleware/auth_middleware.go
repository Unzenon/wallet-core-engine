package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"belajar-go-docker/internal/handler"

	"github.com/golang-jwt/jwt/v5"
)

// Membuat tipe data custom untuk kunci context agar tidak bentrok
type contextKey string
const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1. Ambil isi Header "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token tidak ditemukan, silakan login"})
			return
		}

		// 2. Format token biasanya: "Bearer <token_asli>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Format token salah (Wajib 'Bearer <token>')" })
			return
		}

		tokenString := tokenParts[1]

		// 3. Validasi token JWT
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return handler.JwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token kedaluwarsa atau tidak sah"})
			return
		}

		// 4. Ambil user_id dari dalam token, lalu titipkan ke dalam "Context" Request
		userIDFl, ok := (*claims)["user_id"].(float64) // JWT membaca angka sebagai float64
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Klaim user_id tidak valid"})
			return
		}

		userID := int(userIDFl)
		ctx := context.WithValue(r.Context(), "userID", userID)

		// 5. Lolos pemeriksaan! Teruskan ke handler asli yang dituju
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}