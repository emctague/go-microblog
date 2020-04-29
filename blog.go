package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Represents a microblog - deals with articles and updating them.
type Blog struct {
	ViewDir string
	Views   *template.Template
	Pm      PostManager
}

// Create a new blog that reads views and articles from the given directories.
func NewBlog(viewDir string, articleDir string) Blog {
	return Blog{
		ViewDir: viewDir,
		Views:   nil,
		Pm:      NewPostManager(articleDir),
	}
}

// Update the state of the blog at the given interval of time. This should be run as a goroutine.
func (b *Blog) UpdateTask(updateInterval time.Duration) {
	for {
		var err error

		b.Views, err = template.New("").Funcs(TemplateFunctions).ParseGlob(b.ViewDir + "/*")
		if err != nil {
			log.Fatal(err)
		}

		b.Pm.UpdatePosts()
		time.Sleep(updateInterval)
	}
}

// Handle the route for the blog's index.
func (b *Blog) GetPostList(w http.ResponseWriter, _ *http.Request) {
	if err := b.Views.ExecuteTemplate(w, "index.gohtml", b.Pm.IndexedPosts); err != nil {
		log.Println(err)
		http.Error(w, "An internal server error occurred.", 500)
	}
}

// Handle routes for posts on the blog.
func (b *Blog) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Find the post based on the permalink path
	post, ok := b.Pm.Posts[PermaPath{vars["year"], vars["month"], vars["day"], vars["permaTitle"]}]
	if !ok || !post.Visible {
		http.Error(w, "Post not found", 404)
	}

	if err := b.Views.ExecuteTemplate(w, "post.gohtml", post); err != nil {
		log.Println(err)
		http.Error(w, "An internal server error occurred.", 500)
	}
}
