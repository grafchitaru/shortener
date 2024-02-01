package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/grafchitaru/shortener/internal/config"
	"net/http"
)

func GenerateToken(userID uuid.UUID, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithUserCookie(ctx config.HandlerContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				userID, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				token, _ := GenerateToken(userID, ctx.Config.SecretKey)

				http.SetCookie(w, &http.Cookie{
					Name:  "token",
					Value: token,
					Path:  "/",
				})
			} else {
				_, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
					return []byte(ctx.Config.SecretKey), nil
				})

				if err != nil {
					userID := uuid.New()
					token, _ := GenerateToken(userID, ctx.Config.SecretKey)

					//nolint:exhaustruct
					cook := &http.Cookie{
						Name:  "token",
						Value: token,
						Path:  "/",
					}
					w.Header().Add("Authorization", "Bearer "+token)
					r.Header.Add("Authorization", "Bearer "+token)
					http.SetCookie(w, cook)
					r.AddCookie(cook)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(req *http.Request, secretKey string) (string, error) {
	cookie, err := req.Cookie("token")
	if err != nil {
		return "", err
	}
	tokenString := cookie.Value

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	return claims["user_id"].(string), nil
}
