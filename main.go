package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"crud-lengkap/handler"
	"crud-lengkap/repository"
)

func main() {
	connStr := "user=postgres password=db_password dbname=go_crud sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Gagal konfigurasi DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Gagal koneksi ke PostgreSQL: %v", err)
	}

	fmt.Println("Koneksi ke PostgreSQL Berhasil!")

	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)
	http.HandleFunc("/api/products", productHandler.HandleProduct)

	fmt.Println("Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
