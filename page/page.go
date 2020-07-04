package page

import "io/ioutil"

// Page stores a web page
type Page struct {
	Title string
	Body  []byte
}

func getFilename(title string) string {
	return "data/" + title + ".txt"
}

// Save page
func (p *Page) Save() error {
	filename := getFilename(p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Load page
func Load(title string) (*Page, error) {
	body, err := ioutil.ReadFile(getFilename(title))
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
