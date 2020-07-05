package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/scott-vincent/go-wiki/page"
)

var templates = template.Must(template.ParseFiles(
	"tmpl/home.html",
	"tmpl/missingPage.html",
	"tmpl/edit.html",
	"tmpl/view.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, p *page.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderTemplate(w, "missingPage", nil)
		return
	}

	// Find all existing pages
	err := templates.ExecuteTemplate(w, "home.html", page.GetTitles())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	if title == "" {
		renderTemplate(w, "missingPage", nil)
		return
	}

	p, err := page.Load(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := page.Load(title)
	if err != nil {
		p = &page.Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	oldTitle := r.URL.Path[len("/save/"):]
	newTitle := strings.TrimSpace(r.FormValue("title"))
	body := r.FormValue("body")

	if newTitle == "" {
		p := &page.Page{Body: []byte(body), Error: "Page must have a title"}
		renderTemplate(w, "edit", p)
		return
	}

	// If page title has changed, make sure it is valid
	if newTitle != oldTitle {
		err := page.ValidateNewPage(newTitle)
		if err != nil {
			// Redisplay the edit page and show the error
			p := &page.Page{Title: oldTitle, Body: []byte(body), Error: err.Error()}
			renderTemplate(w, "edit", p)
			return
		}

		// Delete the old page
		page.Delete(oldTitle)
	}

	p := &page.Page{Title: newTitle, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+newTitle, http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/delete/"):]
	page.Delete(title)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/delete/", deleteHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
