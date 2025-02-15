package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/titoffon/merch-store/internal/config"
	"github.com/titoffon/merch-store/internal/delivery/handlers"
)

type TestClient struct{
	baseURL string
}

func TestE2E(t *testing.T) {

	cfg := config.LoadConfig()
	
	go func(){
		err := Run(cfg)
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
			
		case 400, 500:
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