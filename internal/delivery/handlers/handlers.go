package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/titoffon/merch-store/internal/db"
	"golang.org/x/crypto/bcrypt"
)

const WelcomCoins = 1000

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct{
	Error string `json:"error"`
}

type Handlers struct {
	Dal *db.DB
}

func (h *Handlers) Auth(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
			ResponseError(w, http.StatusBadRequest, "Invalid request body")
			slog.Warn("Username and password must not be empty.")
            return
        }

		if req.Username == "" || req.Password == "" {
			slog.Warn("Validation failed")
			ResponseError(w, 400, "Validation failed")
			return

		}
		
		user, err := h.Dal.GetUserByName(r.Context(), req.Username)
		if err != nil {
			ResponseError(w, 500, "internal error")
			slog.Error("Database error", slog.String("error", err.Error()))
			return
		}

		if user == nil {
			hashPassword, err := HashedPass( req.Password )
			if err != nil{
				ResponseError(w, 400, err.Error())
				slog.Error("Failed to hash pass")
				return
			}
			
			user, err = h.Dal.CreateUser(r.Context(), db.User{
				Username: req.Username,
				HashedPassword: string(hashPassword),
				Balance: WelcomCoins,
			})
			if err != nil {
				ResponseError(w, 500, "internal error")
				slog.Error("Failed to create user", slog.String("error", err.Error()))
				return
			}
			

			token, err := generateJWTToken(user.Username, []byte(os.Getenv("JWT_SECRET")))
			if err != nil {
				ResponseError(w, 500, "internal error")
				slog.Error("Failed to generate token", slog.String("error", err.Error()))
				return
			}
			
			ResponseJWT(w, token)
			return 

		}

		valid, err := CheckPassword(user.HashedPassword, req.Password)
		if err != nil || !valid {
			ResponseError(w, http.StatusUnauthorized, "Invalid password")
			slog.Warn("Invalid password")
			return
		}

		token, err := generateJWTToken(user.Username, []byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			ResponseError(w, 500, "internal error")
			slog.Error("Failed to generate token", slog.String("error", err.Error()))
			return
		}		
		ResponseJWT(w, token)
		
	}



func CheckPassword(hashPassword, password string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		slog.Info("The password is uncorrect")
		return false, nil
	}
	slog.Info("The password is correct")
	return true, nil
}

func ResponseError(w http.ResponseWriter, code int, message string){
			res, err := json.Marshal(ErrorResponse{
				Error: message,
			})
			if err != nil{
				slog.Error("failed Unmarshall")
				http.Error(w, "Internal error", http.StatusInternalServerError  )
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			w.Write(res)
}


func ResponseJWT(w http.ResponseWriter, token string){
			res, err := json.Marshal(AuthResponse{
				Token: token,
			})
			if err != nil{
				slog.Error("failed Unmarshall")
				http.Error(w, "Internal error", http.StatusInternalServerError  )
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(res)
}

	func HashedPass(password string) ([]byte, error){
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
		return hashedPassword, nil
	}


func generateJWTToken(username string, secretKey []byte) (string, error) {

	claims := jwt.MapClaims{
		"sub": username,
		//"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GetUserInfo(dal *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserPurchases(dal *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func GetUserTransactions(dal *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}

func TransferCoins(dal *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Использование pool для работы с БД
	}
}