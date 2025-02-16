package httpserv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/titoffon/merch-store/internal/config"
	"github.com/titoffon/merch-store/internal/delivery/handlers"
	"github.com/titoffon/merch-store/internal/server"
)

type FTSendCoinResponse struct {
    code  int64
    Error *handlers.ErrorResponse
}


func TestE2ESendCoins(t *testing.T) {
    cfg := config.LoadConfig()

    go func() {
        err := server.Run(cfg)
        if err != nil {
            t.Fatal(err)
        }
    }()

    time.Sleep(1 * time.Second)

    if t.Failed() {
        t.Fatal("Server failed to start")
    }

    tClient := TestClient{
        baseURL: "http://localhost:8080/api",
    }

    t.Run("SendCoins", func(t *testing.T) {

        senderResp := tClient.Auth(t, handlers.AuthRequest{
            Username: "senderUser",
            Password: "senderPass",
        })
        if senderResp == nil || senderResp.Token == nil || senderResp.code != http.StatusOK {
            t.Fatalf("failed to create senderUser: code=%d, err=%v", senderResp.code, senderResp.Error)
        }
        senderToken := senderResp.Token.Token
        t.Logf("senderUser token: %s", senderToken)

        receiverResp := tClient.Auth(t, handlers.AuthRequest{
            Username: "receiverUser",
            Password: "receiverPass",
        })
        if receiverResp == nil || receiverResp.Token == nil || receiverResp.code != http.StatusOK {
            t.Fatalf("failed to create receiverUser: code=%d, err=%v", receiverResp.code, receiverResp.Error)
        }
        t.Logf("receiverUser is created")


        t.Run("Green scenario", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, senderToken, handlers.SendCoinRequest{
                ToUser: "receiverUser",
                Amount: 100,
            })
            if sendResp.code != http.StatusOK {
                t.Fatalf("expected 200, got %d, err=%v", sendResp.code, sendResp.Error)
            }
        })

        t.Run("No token => 401", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, "", handlers.SendCoinRequest{
                ToUser: "receiverUser",
                Amount: 50,
            })
            if sendResp.code != http.StatusUnauthorized {
                t.Fatalf("expected 401, got %d", sendResp.code)
            }
            t.Logf("Error: %+v", sendResp.Error)
        })

        t.Run("Receiver does not exist => 400", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, senderToken, handlers.SendCoinRequest{
                ToUser: "ghostUser",
                Amount: 10,
            })
            if sendResp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", sendResp.code)
            }
            if sendResp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", sendResp.Error.Error)
        })

        t.Run("Not enough coins => 400", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, senderToken, handlers.SendCoinRequest{
                ToUser: "receiverUser",
                Amount: 9999999,
            })
            if sendResp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", sendResp.code)
            }
            if sendResp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", sendResp.Error.Error)
        })

        t.Run("Non-positive amount => 400", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, senderToken, handlers.SendCoinRequest{
                ToUser: "receiverUser",
                Amount: 0,
            })
            if sendResp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", sendResp.code)
            }
            t.Logf("Error: %v", sendResp.Error)
        })

        t.Run("Empty toUser => 400", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, senderToken, handlers.SendCoinRequest{
                ToUser: "",
                Amount: 10,
            })
            if sendResp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", sendResp.code)
            }
            t.Logf("Error: %v", sendResp.Error)
        })
		
        t.Run("Invalid token => 401", func(t *testing.T) {
            sendResp := tClient.SendCoins(t, "fakeToken", handlers.SendCoinRequest{
                ToUser: "receiverUser",
                Amount: 10,
            })
            if sendResp.code != http.StatusUnauthorized {
                t.Fatalf("expected 401, got %d", sendResp.code)
            }
            t.Logf("Error: %v", sendResp.Error)
        })
    })
}

func (tc *TestClient) SendCoins(t *testing.T, token string, body handlers.SendCoinRequest) *FTSendCoinResponse {
    reqBody, err := json.Marshal(body)
    if err != nil {
        t.Fatal(err)
    }

    req, err := http.NewRequest("POST", tc.baseURL+"/sendCoin", bytes.NewReader(reqBody))
    if err != nil {
        t.Fatal("failed to create POST request:", err)
    }

    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatal("failed to do request:", err)
    }
    t.Cleanup(func() {
        resp.Body.Close()
    })

    switch resp.StatusCode {
    case http.StatusOK:

        return &FTSendCoinResponse{
            code:  200,
            Error: nil,
        }
    case http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError:
        var e handlers.ErrorResponse
        if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
            t.Fatalf("failed to decode error response: %v", err)
        }
        return &FTSendCoinResponse{
            code:  int64(resp.StatusCode),
            Error: &e,
        }
    default:

        return &FTSendCoinResponse{code: int64(resp.StatusCode)}
    }
}