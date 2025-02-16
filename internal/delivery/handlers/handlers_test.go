package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


func TestCheckPassword(t *testing.T) {
    hashed, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
    if err != nil {
        t.Fatalf("failed to generate hash: %v", err)
    }

    ok, err := CheckPassword(string(hashed), "test123")
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if !ok {
        t.Errorf("expected password to match, got false")
    }

    ok, err = CheckPassword(string(hashed), "wrongPass")
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if ok {
        t.Errorf("expected password mismatch (false), got true")
    }
}

func TestResponseError(t *testing.T) {
    rr := httptest.NewRecorder()

    ResponseError(rr, http.StatusBadRequest, "some error message")

    if rr.Code != http.StatusBadRequest {
        t.Errorf("expected status 400, got %d", rr.Code)
    }

    var resp ErrorResponse
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal body: %v", err)
    }
    if resp.Error != "some error message" {
        t.Errorf("expected 'some error message', got '%s'", resp.Error)
    }

    ct := rr.Header().Get("Content-Type")
    if ct != "application/json" {
        t.Errorf("expected Content-Type=application/json, got %s", ct)
    }
}

func TestResponseJWT(t *testing.T) {
    rr := httptest.NewRecorder()

    ResponseJWT(rr, "fake-jwt-token")

    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rr.Code)
    }

    // Проверяем JSON-ответ
    var resp AuthResponse
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal body: %v", err)
    }
    if resp.Token != "fake-jwt-token" {
        t.Errorf("expected 'fake-jwt-token', got '%s'", resp.Token)
    }

    if rr.Header().Get("Content-Type") != "application/json" {
        t.Errorf("expected Content-Type=application/json")
    }
}

func TestHashedPass(t *testing.T) {
    password := "mySecretPassword"

    hashed, err := HashedPass(password)
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if len(hashed) == 0 {
        t.Error("expected non-empty hashed password")
    }

    if err := bcrypt.CompareHashAndPassword(hashed, []byte(password)); err != nil {
        t.Errorf("bcrypt.CompareHashAndPassword should succeed for the correct password, got err=%v", err)
    }

    if err := bcrypt.CompareHashAndPassword(hashed, []byte("wrong")); err == nil {
        t.Error("expected error for wrong password, got nil")
    }
}

func TestGenerateJWTToken(t *testing.T) {
    secret := []byte("testSecretKey")
    username := "testUser"

    tokenStr, err := generateJWTToken(username, secret)
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if tokenStr == "" {
        t.Error("expected non-empty token string")
    }


    parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secret, nil
    })
    if err != nil {
        t.Fatalf("failed to parse token: %v", err)
    }
    if !parsedToken.Valid {
        t.Error("parsed token is invalid")
    }

    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok {
        t.Fatal("could not cast claims to jwt.MapClaims")
    }
    if claims["sub"] != username {
        t.Errorf("expected sub=%s, got %v", username, claims["sub"])
    }
}

func TestExtractJWT(t *testing.T) {

    slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

    os.Setenv("JWT_SECRET", "testSecretKey")


    validToken := createTestJWTToken(t, "validUser", []byte("testSecretKey"))

    wrongSignatureToken := createTestJWTToken(t, "userWrongSignature", []byte("otherSecretKey"))

    emptySubToken := createTestJWTToken(t, "", []byte("testSecretKey"))

    tests := []struct {
        name            string
        authHeaderValue string 
        wantStatus      int    
        wantErrMessage  string 
        wantUsername    string 
    }{
        {
            name:           "No Authorization header => 401",
            authHeaderValue: "",
            wantStatus:     http.StatusUnauthorized,
            wantErrMessage: "Authorization token is required",
        },
        {
            name:           "Wrong signature => 401",
            authHeaderValue: "Bearer " + wrongSignatureToken,
            wantStatus:     http.StatusUnauthorized,
            wantErrMessage: "Invalid token",
        },
        {
            name:           "Empty sub => 401",
            authHeaderValue: "Bearer " + emptySubToken,
            wantStatus:     http.StatusUnauthorized,
            wantErrMessage: "Empty Username Plaload",
        },
        {
            name:           "Valid token => 200, returns username",
            authHeaderValue: "Bearer " + validToken,
            wantStatus:     http.StatusOK,
            wantUsername:   "validUser",
        },
        {
            name:           "Malformed token => 401",
            authHeaderValue: "Bearer some.garbage.token",
            wantStatus:     http.StatusUnauthorized,
            wantErrMessage: "Invalid token",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {

            req := httptest.NewRequest(http.MethodGet, "/someEndpoint", nil)
            if tc.authHeaderValue != "" {
                req.Header.Set("Authorization", tc.authHeaderValue)
            }

            rr := httptest.NewRecorder()

            username, err := ExtractJWT(rr, req)

            if rr.Code != tc.wantStatus {
                t.Errorf("expected status %d, got %d", tc.wantStatus, rr.Code)
            }

            if tc.wantStatus == http.StatusOK {

                if err != nil {
                    t.Errorf("expected no error, got %v", err)
                }
                if username != tc.wantUsername {
                    t.Errorf("expected username=%q, got %q", tc.wantUsername, username)
                }
            } else {

                if err == nil {
                    t.Error("expected an error, got nil")
                }
                if !strings.Contains(rr.Body.String(), tc.wantErrMessage) && tc.wantErrMessage != "" {
                    t.Errorf("expected error message %q in response body, got %q",
                        tc.wantErrMessage, rr.Body.String())
                }
            }
        })
    }
}

func createTestJWTToken(t *testing.T, username string, secret []byte) string {
    t.Helper()
    claims := jwt.MapClaims{
        "sub": username,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(secret)
    if err != nil {
        t.Fatalf("failed to sign token: %v", err)
    }
    return signed
}