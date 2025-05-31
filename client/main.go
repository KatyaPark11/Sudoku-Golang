package main

import (
	"log"
	"net/http"

	pbAuth "github.com/KatyaPark11/Sudoku-Golang/generated/auth"
	pbSudoku "github.com/KatyaPark11/Sudoku-Golang/generated/sudoku"

	"google.golang.org/grpc"
)

var (
	authClient   pbAuth.AuthServiceClient
	sudokuClient pbSudoku.SudokuServiceClient
)

func main() {
	var err error

	// Инициализация соединений с gRPC-сервисами
	authConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	authClient = pbAuth.NewAuthServiceClient(authConn)

	sudokuConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to sudoku service: %v", err)
	}
	sudokuClient = pbSudoku.NewSudokuServiceClient(sudokuConn)

	// Запуск HTTP-сервера или другого основного кода
	http.HandleFunc("/register.html", serveRegisterPage)
	http.HandleFunc("/login.html", serveLoginPage)
	http.HandleFunc("/sudoku.html", serveSudokuPage)
	http.HandleFunc("/api/register", handleRegister)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/solve", handleSolve)

	// Статические файлы (HTML)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func serveRegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/register.html")
}

func serveLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login.html")
}

func serveSudokuPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/sudoku.html")
}
