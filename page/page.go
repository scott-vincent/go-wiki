package page

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const pageFolder = "data"

// Page stores a web page
type Page struct {
	Title string
	Body  []byte
	Error string
}

func getFilename(title string) string {
	return pageFolder + "/" + title + ".txt"
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

// GetTitles reads all filenames in the page folder
func GetTitles() []string {
	files, err := ioutil.ReadDir(pageFolder)
	if err != nil {
		panic(err)
	}

	var titles []string
	for _, file := range files {
		title := file.Name()

		// Remove .txt extension
		titles = append(titles, title[0:len(title)-4])
	}

	sort.Strings(titles)
	return titles
}

// Delete the page with the specified title
func Delete(title string) {
	os.Remove(getFilename(title))
}

// ValidateNewPage returns error if the filename is not valid or already exists
func ValidateNewPage(title string) error {
	filename := getFilename(title)

	// Don't allow special chars
	if strings.ContainsAny(title, "#?/\\*\"") {
		return fmt.Errorf("Page name '%s' not allowed: Must not contain special characters", title)
	}

	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return nil
	} else if err == nil {
		return fmt.Errorf("Page '%s' already exists", title)
	} else {
		return fmt.Errorf("Page name '%s' not allowed: %s", title, err.Error())
	}
}
