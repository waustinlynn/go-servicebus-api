package main

import (
	sb "github.com/waustinlynn/go-servicebus"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/mendsley/gojwk"
	"os"
)

var SB_URL = os.Getenv("SB_URL")
var SB_KEY = os.Getenv("SB_KEY")
var SB_KEYTYPE = os.Getenv("SB_KEYTYPE")
var AUTH_URL = os.Getenv("AUTH_URL")
var SB_PORT = os.Getenv("SB_PORT")

func main(){
	// setConfig()
	fmt.Println(SB_URL)
	fmt.Println(SB_KEY)
	fmt.Println(SB_KEYTYPE)
	fmt.Println(SB_PORT)
	fmt.Println(AUTH_URL)
	port := ":8001"
	if(len(SB_PORT) > 0){
		port = ":" + SB_PORT
	}
	router := mux.NewRouter()
	router.HandleFunc("/message", ValidateMiddleware(SendMessage)).Methods("POST")
	router.HandleFunc("/status", ApiStatus).Methods("GET")
	http.ListenAndServe(port, router)
}

type Exception struct {
    Message string `json:"message"`
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        authorizationHeader := req.Header.Get("authorization")
        if authorizationHeader != "" {
            bearerToken := strings.Split(authorizationHeader, " ")
            if len(bearerToken) == 2 {
                token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("There was an error")
					}
                    return myLookupKey(token.Header["kid"].(string))
				})
                if err != nil {
                    json.NewEncoder(w).Encode(Exception{Message: err.Error()})
                    return
                }
                if token.Valid {
                    context.Set(req, "decoded", token.Claims)
                    next(w, req)
                } else {
                    json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
                }
            }
        } else {
            json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
        }
    })
}

func myLookupKey(kid string) (interface{}, error) {
	//fmt.Printf("Kid : %v\n", kid)
	var keys struct{ Keys []gojwk.Key }
	getKey(&keys)
	for _, key := range keys.Keys {
		if key.Kid == kid {
			return key.DecodePublicKey()
		}
	}
	return nil, fmt.Errorf("Key not found")
}

func getKey(keys interface{}) {
	resp, err := http.Get(AUTH_URL + "/.well-known/jwks")
	if(err != nil){
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal(body, keys)
}

func SendMessage(w http.ResponseWriter, r *http.Request){
	var msg sb.SbMessage
	body, err := ioutil.ReadAll(r.Body)
	if(err != nil){
		panic(err)
	}

	err = json.Unmarshal(body, &msg)
	if(err != nil){
		panic(err)
	}
	config := &sb.SbConfig{SB_KEY, SB_KEYTYPE, SB_URL}
	_, err = config.Send(&msg)
	if(err != nil){
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Message sent"))
}

func ApiStatus(w http.ResponseWriter, r *http.Request){
	cfg := make(map[string]string)
	cfg["SB_URL"] = SB_URL
	cfg["AUTH_URL"] = AUTH_URL
	cfgBytes, _ := json.Marshal(cfg)
	w.Write(cfgBytes)
}