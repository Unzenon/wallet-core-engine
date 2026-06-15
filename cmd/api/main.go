package main

import (
	"fmt"
	"log"
	"net/http"

	"belajar-go-docker/internal/config"
	"belajar-go-docker/internal/handler"
	"belajar-go-docker/internal/middleware" // <-- 1. PASTIKAN IMPORT INI ADA
	"belajar-go-docker/internal/repository"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}
	defer db.Close()
	fmt.Println("Sukses terhubung ke PostgreSQL!")

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Route Aplikasi Kita
	http.HandleFunc("/register", userHandler.Register)
	http.HandleFunc("/login", userHandler.Login)
	
	// <-- 2. TAMBAHKAN BARIS TERPROTEKSI INI
	http.HandleFunc("/topup", middleware.AuthMiddleware(userHandler.TopUp)) 

	fmt.Println("Server jalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}