package week5

import (
	"fmt"
	"net/http"
	"strconv"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Year   int
	ISBN   string
}

var books = []Book{
	{1, "The Go Programming Language", "Alan Donovan", 2015, "9780134190440"},
	{2, "Clean Code", "Robert C. Martin", 2008, "9780132350884"},
	{3, "Design Patterns", "GoF", 1994, "9780201633610"},
}

var favorites = map[int]Book{}

func getBookByID(id int) (Book, bool) {
	for _, b := range books {
		if b.ID == id {
			return b, true
		}
	}
	return Book{}, false
}

func isFavorite(id int) bool {
	_, ok := favorites[id]
	return ok
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>All Books</h1>")
	fmt.Fprintln(w, `<a href="/favorites">Favorites</a><hr>`)

	for _, b := range books {
		fmt.Fprintf(w, "<b>%s</b> by %s<br>", b.Title, b.Author)

		fmt.Fprintf(w, `
		<form style="display:inline" method="POST" action="/add">
			<input type="hidden" name="id" value="%d">
			<button>Add to Favorites</button>
		</form>
		<a href="/book?id=%d"><button>Details</button></a>
		<br><br>`, b.ID, b.ID)
	}
}
func addFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))

	if book, ok := getBookByID(id); ok {
		favorites[id] = book
	}
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
func favoritesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Favorites</h1>")
	fmt.Fprintln(w, `<a href="/books"><button>All Books</button></a><hr>`)

	if len(favorites) == 0 {
		fmt.Fprintln(w, "No favorite books.")
		return
	}

	for _, b := range favorites {
		fmt.Fprintf(w, "<b>%s</b><br>", b.Title)

		fmt.Fprintf(w, `
		<form style="display:inline" method="POST" action="/remove">
			<input type="hidden" name="id" value="%d">
			<button>Remove</button>
		</form>
		<a href="/book?id=%d"><button>Details</button></a>
		<br><br>`, b.ID, b.ID)
	}
}
func removeFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))
	delete(favorites, id)
	http.Redirect(w, r, "/favorites", http.StatusSeeOther)
}

func bookDetailsHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	book, ok := getBookByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, `<h1>%s</h1>`, book.Title)
	fmt.Fprintf(w, `
	Author: %s<br>
	Year: %d<br>
	ISBN: %s<br><br>
	`, book.Author, book.Year, book.ISBN)

	if isFavorite(id) {
		fmt.Fprintf(w, `
		<form method="POST" action="/remove">
			<input type="hidden" name="id" value="%d">
			<button>Remove from Favorites</button>
		</form>`, id)
	} else {
		fmt.Fprintf(w, `
		<form method="POST" action="/add">
			<input type="hidden" name="id" value="%d">
			<button>Add to Favorites</button>
		</form>`, id)
	}

	fmt.Fprintln(w, `<br><a href="/books"><button>All Books</button></a>`)
}

func main() {
	http.HandleFunc("/books", booksHandler)
	http.HandleFunc("/favorites", favoritesHandler)
	http.HandleFunc("/book", bookDetailsHandler)
	http.HandleFunc("/add", addFavoriteHandler)
	http.HandleFunc("/remove", removeFavoriteHandler)

	fmt.Println("Server running at http://localhost:8080/books")
	http.ListenAndServe(":8080", nil)
}
