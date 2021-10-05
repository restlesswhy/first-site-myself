package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"html/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showedPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles_2`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("fullText")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "not all data is filled")
	} else {
		db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles_2` (`title`, `anons`, `fullText`) VALUES('%s', '%s', '%s')", title, anons, fullText))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func showPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show_post.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles_2` WHERE `Id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showedPost = Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Anons, &post.Title, &post.FullText)
		if err != nil {
			panic(err)
		}
		
		showedPost = post
	}
	
	t.ExecuteTemplate(w, "show_post", showedPost)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	res, err := db.Query(fmt.Sprintf("DELETE FROM `articles_2` WHERE `Id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {
	r := mux.NewRouter()

	r.HandleFunc("/delete_post/{id:[0-9]+}", deletePost).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")
	r.HandleFunc("/saveArticle", saveArticle).Methods("POST")
	r.HandleFunc("/create", create).Methods("GET")
	r.HandleFunc("/", index).Methods("GET")

	http.Handle("/", r)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}


func main() {
	handleFunc()
}