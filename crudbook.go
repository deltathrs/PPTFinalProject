package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/addbook", addBookHandler)
	http.HandleFunc("/deletebook", deleteBookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bookTitle := r.Form.Get("book_title")
	authorName := r.Form.Get("author_name")
	publisherName := r.Form.Get("publisher")
	publicationDate := r.Form.Get("publication_date")
	ISBN := r.Form.Get("ISBN")
	price := r.Form.Get("price")
	stockQty := r.Form.Get("stock_qty")

	bookID := generateUniqueID(bookTitle)
	authorID := generateUniqueID(authorName)
	publisherID := generateUniqueID(publisherName)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books_store")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert into authors table
	authorStmt, err := db.Prepare("INSERT INTO authors (author_id, name) VALUES (?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into authors table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer authorStmt.Close()

	_, err = authorStmt.Exec(authorID, authorName)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into authors table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert into publishers table
	publisherStmt, err := db.Prepare("INSERT INTO publishers (publisher_id, name) VALUES (?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into publishers table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer publisherStmt.Close()

	_, err = publisherStmt.Exec(publisherID, publisherName)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into publishers table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert into books table
	bookStmt, err := db.Prepare("INSERT INTO books (book_id, title, author_id, publisher_id, publication_date, ISBN, price, stock_qty) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into books table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer bookStmt.Close()

	_, err = bookStmt.Exec(bookID, bookTitle, authorID, publisherID, publicationDate, ISBN, price, stockQty)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into books table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Book added successfully!")
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bookTitle := r.Form.Get("book_title")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books_store")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM books WHERE title = ?")
	if err != nil {
		log.Println("Failed to prepare SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(bookTitle)
	if err != nil {
		log.Println("Failed to execute SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Book deleted successfully!")
}

func generateUniqueID(input string) string {
	id := ""
	for _, c := range input {
		id += strconv.Itoa(int(c))
	}
	return id
}
