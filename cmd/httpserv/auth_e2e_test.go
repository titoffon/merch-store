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

type TestClient struct{
	baseURL string
}

func TestE2EAuth(t *testing.T) {

	cfg := config.LoadConfig()
	
	go func(){
		err := server.Run(cfg)
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)
	if t.Failed(){
		t.Fatal("fail to start")
	}

	tClient := TestClient{
		baseURL: "http://localhost:8080/api",
	}

	t.Run("Auth", func(t *testing.T){
		t.Run("Green", func(t *testing.T){
				response := tClient.Auth(t, handlers.AuthRequest{
				Username: "1234",
				Password: "23423421",
			})

			if response.code != 200 {
				t.Fatal(response.code)
			}
			t.Log(response.Token.Token)		

		})

	})

	 t.Run("Green (Existing user)", func(t *testing.T) {
            resp := tClient.Auth(t, handlers.AuthRequest{
                Username: "someNewUser",
                Password: "somePassword123",
            })
            if resp == nil {
                t.Fatal("response is nil")
            }
            if resp.code != http.StatusOK {
                t.Fatalf("expected 200, got %d", resp.code)
            }
            t.Logf("Existing user token: %s", resp.Token.Token)
        })

		t.Run("Empty username => 400", func(t *testing.T) {
            resp := tClient.Auth(t, handlers.AuthRequest{
                Username: "",
                Password: "somePassword123",
            })
            if resp == nil {
                t.Fatal("response is nil")
            }
            if resp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", resp.code)
            }
            if resp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", resp.Error.Error)
        })

		t.Run("Empty password => 400", func(t *testing.T) {
            resp := tClient.Auth(t, handlers.AuthRequest{
                Username: "userWithoutPassword",
                Password: "",
            })
            if resp == nil {
                t.Fatal("response is nil")
            }
            if resp.code != http.StatusBadRequest {
                t.Fatalf("expected 400, got %d", resp.code)
            }
            if resp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", resp.Error.Error)
        })

		t.Run("Wrong password => 401", func(t *testing.T) {
            
            resp := tClient.Auth(t, handlers.AuthRequest{
                Username: "someNewUser",
                Password: "wrongPassword", 
            })
            if resp == nil {
                t.Fatal("response is nil")
            }
            if resp.code != http.StatusUnauthorized {
                t.Fatalf("expected 401, got %d", resp.code)
            }
            if resp.Error == nil {
                t.Fatal("expected error message, got nil")
            }
            t.Logf("Error: %s", resp.Error.Error)
        })

		
}

type FTAuthResponse struct{
	code int64
	Token *handlers.AuthResponse
	Error *handlers.ErrorResponse
}

func (tc *TestClient) Auth(t *testing.T, r handlers.AuthRequest) *FTAuthResponse{

	reqBody, err := json.Marshal(r)
		if err != nil{
			t.Fatal(err)
		}
		
		reqBodyReader := bytes.NewReader(reqBody)

		req, err := http.NewRequest("POST", tc.baseURL + "/auth", reqBodyReader)
		if err != nil {
			t.Fatal("failed to req")
		}
		


		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("faile do req")
		}
		
		t.Cleanup(func (){
			resp.Body.Close()
			})
		
		switch resp.StatusCode{
		case 200: 	
			
			var response handlers.AuthResponse

			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil{
				t.Fatal(err)
			}
			return &FTAuthResponse{
				code: 200,
				Token: &response,
				Error: nil,
			}
			
		case 400, 401, 500:
			var response handlers.ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			if err != nil{
				t.Fatal(err)
			}
			return &FTAuthResponse{
				code: int64(resp.StatusCode),
				Token: nil,
				Error: &response,
			}
			
		default:
			t.Log(resp.Status)
		}
		return nil
}
