package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pbAuth "github.com/KatyaPark11/Sudoku-Golang/generated/auth"
	pbSudoku "github.com/KatyaPark11/Sudoku-Golang/generated/sudoku"
)

// Структуры для JSON-запросов/ответов
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}

type SolveRequest struct {
	Puzzle  string `json:"puzzle"`
	IsSteps bool   `json:"isSteps"`
}

type SolveResponse struct {
	Solution string `json:"solution,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Обработчик регистрации
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Register(ctx, &pbAuth.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		log.Println("Auth Register error:", err)
		http.Error(w, "Auth service error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": resp.Success})
}

// Обработчик входа
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Login(ctx, &pbAuth.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		log.Println("Auth Login error:", err)
		json.NewEncoder(w).Encode(LoginResponse{Success: false, Message: "Login failed"})
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Success: resp.Success,
		Token:   resp.Token,
		Message: "",
	})
}

// Обработчик решения судоку (требует токена в заголовке Authorization)
func handleSolve(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := sudokuClient.Solve(ctx, &pbSudoku.SudokuRequest{Puzzle: req.Puzzle, IsSteps: req.IsSteps})
	if err != nil {
		log.Println("Sudoku solve error:", err)
		json.NewEncoder(w).Encode(SolveResponse{Error: "Данное судоку не имеет решения"})
		return
	}

	json.NewEncoder(w).Encode(SolveResponse{Solution: resp.Solution})
}
