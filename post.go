package main

import (
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"os"
	"strconv"
	"time"
)

// Represents the information needed to find a post by a date-based path.
type PermaPath struct {
	Year       string // The year the post was created.
	Month      string // The month the post was created.
	Day        string // The day the post was created.
	PermaTitle string // The permalink-title component of the URL.
}

// Represents a blog post and its metadata.
type Post struct {

	// This section contains metadata you might want to provide in your article front matter:

	Title      string    // The article's real title.
	PermaTitle string    `yaml:"perma_title"` // The permalink title of the article. This should be similar to the real title, but all lowercase with dashes.
	Created    time.Time // The date upon which this article was created.
	Visible    bool      // Whether or not it is possible to view this article. Defaults to true.

	// This section contains metadata that is set by the loader:
	Path         PermaPath
	SourcePath   string
	Body         template.HTML
	LastModified time.Time
}

func (p *Post) LoadPost() error {
	stat, err := os.Stat(p.SourcePath)
	if err != nil {
		return err
	}

	// Don't re-load if it's already loaded!
	if stat.ModTime().Before(p.LastModified) {
		return nil
	}

	f, err := os.Open(p.SourcePath)
	if err != nil {
		return err
	}

	// Load front matter into the post
	remainingBytes, err := SeparateFrontMatter(f, p)
	if err != nil {
		return err
	}

	p.Body = template.HTML(blackfriday.Run(remainingBytes))
	p.LastModified = stat.ModTime()
	p.Path = PermaPath{strconv.Itoa(p.Created.Year()), strconv.Itoa(int(p.Created.Month())), strconv.Itoa(p.Created.Day()), p.PermaTitle}

	return nil
}

func (p Post) GetPermalink() string {
	return "/" + p.Path.Year + "/" + p.Path.Month + "/" + p.Path.Day + "/" + p.Path.PermaTitle
}

func (p Post) GetPrettyDate() string {
	return p.Created.Format("January 2, 2006")
}
