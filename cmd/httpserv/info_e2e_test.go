package httpserv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/titoffon/merch-store/internal/config"
	"github.com/titoffon/merch-store/internal/delivery/handlers"
	"github.com/titoffon/merch-store/internal/server"
)

type FTUserInfoResponse struct {
    code  int64
    Info  *handlers.InfoResponse
    Error *handlers.ErrorResponse
}

func TestE2EUserInfo(t *testing.T) {
    cfg := config.LoadConfig()
    
    go func() {
        if err := server.Run(cfg); err != nil {
            t.Fatal(err)
        }
    }()

    time.Sleep(2 * time.Second)

    if t.Failed() {
        t.Fatal("Server failed to start")
    }

    tClient := TestClient{
        baseURL: "http://localhost:8080/api",
    }

    t.Run("UserInfo", func(t *testing.T) {

        authResp := tClient.Auth(t, handlers.AuthRequest{
            Username: "infoTester",
            Password: "infoPass123",
        })
        if authResp == nil || authResp.Token == nil || authResp.code != http.StatusOK {
            t.Fatalf("Failed to create user infoTester. code=%d, err=%v", authResp.code, authResp.Error)
        }
        userToken := authResp.Token.Token
        t.Logf("Got token for user 'infoTester': %s", userToken)

        t.Run("Green scenario (check default info)", func(t *testing.T) {
            infoResp := tClient.GetUserInfo(t, userToken)
            if infoResp.code != http.StatusOK {
                t.Fatalf("expected 200, got %d (error=%v)", infoResp.code, infoResp.Error)
            }
            if infoResp.Info == nil {
                t.Fatal("expected InfoResponse, got nil")
            }

            if infoResp.Info.Coins != 1000 {
                t.Fatalf("expected 1000 coins, got %d", infoResp.Info.Coins)
            }
            if len(infoResp.Info.Inventory) != 0 {
                t.Fatalf("expected empty inventory, got %v", infoResp.Info.Inventory)
            }
            if len(infoResp.Info.CoinHistory.Received) != 0 || len(infoResp.Info.CoinHistory.Sent) != 0 {
                t.Fatalf("expected empty coinHistory, got received=%v, sent=%v",
                    infoResp.Info.CoinHistory.Received, infoResp.Info.CoinHistory.Sent)
            }
        })

        t.Run("After purchase", func(t *testing.T) {

            pResp := tClient.PurchaseMerch(t, "t-shirt", userToken)
            if pResp.code != http.StatusOK {
                t.Fatalf("Failed to buy t-shirt. code=%d, err=%v", pResp.code, pResp.Error)
            }

            infoResp := tClient.GetUserInfo(t, userToken)
            if infoResp.code != http.StatusOK {
                t.Fatalf("expected 200, got %d (error=%v)", infoResp.code, infoResp.Error)
            }
            if infoResp.Info == nil {
                t.Fatal("expected InfoResponse, got nil")
            }

            expectedCoins := int64(920)
            if infoResp.Info.Coins != expectedCoins {
                t.Fatalf("expected %d coins, got %d", expectedCoins, infoResp.Info.Coins)
            }

            if len(infoResp.Info.Inventory) != 1 {
                t.Fatalf("expected 1 inventory item, got %d", len(infoResp.Info.Inventory))
            }
            if infoResp.Info.Inventory[0].Type != "t-shirt" {
                t.Fatalf("expected 't-shirt', got '%s'", infoResp.Info.Inventory[0].Type)
            }
            if infoResp.Info.Inventory[0].Quantity != 1 {
                t.Fatalf("expected quantity=1, got %d", infoResp.Info.Inventory[0].Quantity)
            }
        })

        t.Run("No token => 401", func(t *testing.T) {
            infoResp := tClient.GetUserInfo(t, "")
            if infoResp.code != http.StatusUnauthorized {
                t.Fatalf("expected 401, got %d", infoResp.code)
            }
            if infoResp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", infoResp.Error.Error)
        })

        t.Run("Invalid token => 401", func(t *testing.T) {
            infoResp := tClient.GetUserInfo(t, "fakeTokenValue")
            if infoResp.code != http.StatusUnauthorized {
                t.Fatalf("expected 401, got %d", infoResp.code)
            }
            if infoResp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", infoResp.Error.Error)
        })
    })
}

func (tc *TestClient) GetUserInfo(t *testing.T, token string) *FTUserInfoResponse {
    url := fmt.Sprintf("%s/info", tc.baseURL)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        t.Fatal("failed to create GET request:", err)
    }

    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatal("failed to do GET request:", err)
    }
    t.Cleanup(func() { resp.Body.Close() })

    switch resp.StatusCode {
    case http.StatusOK:
        var info handlers.InfoResponse
        if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
            t.Fatalf("failed to decode info response: %v", err)
        }
        return &FTUserInfoResponse{
            code:  http.StatusOK,
            Info:  &info,
            Error: nil,
        }
    case http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError:
        var e handlers.ErrorResponse
        if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
            t.Fatalf("failed to decode error response: %v", err)
        }
        return &FTUserInfoResponse{
            code:  int64(resp.StatusCode),
            Info:  nil,
            Error: &e,
        }
    default:
        return &FTUserInfoResponse{
            code:  int64(resp.StatusCode),
            Info:  nil,
            Error: nil,
        }
    }
}