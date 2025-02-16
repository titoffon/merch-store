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

func TestE2EPurchaseMerch(t *testing.T) {
    cfg := config.LoadConfig()

    go func() {
        err := server.Run(cfg)
        if err != nil {
            t.Fatal(err)
        }
    }()
    time.Sleep(1 * time.Second)

    if t.Failed() {
        t.Fatal("fail to start")
    }

    tClient := TestClient{
        baseURL: "http://localhost:8080/api",
    }

    resp := tClient.Auth(t, handlers.AuthRequest{
        Username: "merchBuyer2",
        Password: "merchPass",
    })
    if resp == nil || resp.Token == nil || resp.code != 200 {
        t.Fatalf("failed to create user for PurchaseMerch tests: code=%d, err=%v", resp.code, resp.Error)
    }
    userToken := resp.Token.Token
    t.Logf("Got token for user 'merchBuyer': %s", userToken)

    t.Run("PurchaseMerch", func(t *testing.T) {

        t.Run("Green", func(t *testing.T) {
            buyResp := tClient.PurchaseMerch(t, "t-shirt", userToken)
            if buyResp.code != http.StatusOK {
                t.Fatalf("expected 200, got %d (err=%v)", buyResp.code, buyResp.Error)
            }
        })

        t.Run("No token => 401", func(t *testing.T) {
            buyResp := tClient.PurchaseMerch(t, "t-shirt", "") // передаём пустой токен
            if buyResp.code != http.StatusUnauthorized {
                t.Fatalf("expected 401, got %d", buyResp.code)
            }
            if buyResp.Error == nil {
                t.Fatal("expected error response, got nil")
            }
            t.Logf("Error: %s", buyResp.Error.Error)
        })

        t.Run("Wrong item => 400", func(t *testing.T) {
            buyResp := tClient.PurchaseMerch(t, "fake-item", userToken)
            if buyResp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", buyResp.code)
            }
            if buyResp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", buyResp.Error.Error)
        })

        t.Run("Not enough coins => 400", func(t *testing.T) {
			var buyResp *FTPurchaseMerchResp
            for i := 1; i <= 2; i++ {
                buyResp = tClient.PurchaseMerch(t, "pink-hoody", userToken)
            }
            if  buyResp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", buyResp.Error.Error)
        })

	t.Run("Wrong item => 400", func(t *testing.T) {
        buyResp := tClient.PurchaseMerch(t, "thisItemDoesNotExist", userToken)
        if buyResp.code != http.StatusBadRequest {
            t.Fatalf("expected 400, got %d", buyResp.code)
        }
        if buyResp.Error == nil {
            t.Fatal("expected error message, got nil")
        }
        t.Logf("Error: %s", buyResp.Error.Error)
    })

	    t.Run("Invalid token => 401", func(t *testing.T) {
        buyResp := tClient.PurchaseMerch(t, "book", "completelyRandomString")
        if buyResp.code != http.StatusUnauthorized {
            t.Fatalf("expected 401, got %d", buyResp.code)
        }
        if buyResp.Error == nil {
            t.Fatal("expected error response, got nil")
        }
        t.Logf("Error: %s", buyResp.Error.Error)
    })

    })
}

type FTPurchaseMerchResp struct {
    code  int64
    Error *handlers.ErrorResponse
}

func (tc *TestClient) PurchaseMerch(t *testing.T, item, token string) *FTPurchaseMerchResp {
    url := fmt.Sprintf("%s/buy/%s", tc.baseURL, item)
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
    t.Cleanup(func() {
        resp.Body.Close()
    })

    switch resp.StatusCode {
    case http.StatusOK:

        return &FTPurchaseMerchResp{
            code:  200,
            Error: nil,
        }
    case 400, 401, 500:
        var e handlers.ErrorResponse
        if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
            t.Fatal("failed to decode error response:", err)
        }
        return &FTPurchaseMerchResp{
            code:  int64(resp.StatusCode),
            Error: &e,
        }
    default:
        t.Logf("Unhandled status: %d", resp.StatusCode)
        return &FTPurchaseMerchResp{code: int64(resp.StatusCode)}
    }
}