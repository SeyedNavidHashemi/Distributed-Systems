package main

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Success bool
	Message string
}

type AuthService struct {
	Users []User
}

func (a *AuthService) Login(req LoginRequest, res *LoginResponse) error {
	for _, user := range a.Users {
		if user.Username == req.Username {
			err := bcrypt.CompareHashAndPassword(
			    []byte(user.Password),
			    []byte(req.Password),
			)

			if err == nil {
				res.Success = true
				res.Message = "Login successful"
				return nil
			}
		}
	}

	res.Success = false
	res.Message = "Invalid credentials"
	return nil
}

func loadUsers() []User {
	file, err := os.ReadFile("users.json")
	if err != nil {
		log.Fatal(err)
	}

	var users []User

	err = json.Unmarshal(file, &users)
	if err != nil {
		log.Fatal(err)
	}

	return users
}

func main() {
	authService := &AuthService{
		Users: loadUsers(),
	}

	rpc.Register(authService)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Auth RPC Server running on port 9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}