package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// Specifies a value that can be provided via either an environment variable or command line flag, and which will
// default to some particular value if neither are provided.
// Returns a string pointer that will be valid after calling flag.Parse()
func EnvFlagOrDefaultString(flagName string, envName string, desc string, def string) *string {
	var value string
	value, foundEnv := os.LookupEnv(envName)
	if !foundEnv {
		value = def
	}

	return flag.String(flagName, value, desc)
}

func main() {
	// Obtain configuration strings
	listen := EnvFlagOrDefaultString("listen", "LISTEN", "The IP and port to listen on", "0.0.0.0:8080")
	viewDir := EnvFlagOrDefaultString("viewDir", "VIEW_DIR", "The directory that view (.gohtml) files are found in", "./views")
	articleDir := EnvFlagOrDefaultString("articleDir", "ARTICLE_DIR", "The directory that article (.md) files are found in", "./articles")
	staticDir := EnvFlagOrDefaultString("staticDir", "STATIC_DIR", "The directory from which static files are served", "./static")
	updateInterval := EnvFlagOrDefaultString("updateInterval", "UPDATE_INTERVAL", "The interval at which changes to articles and views are checked for.", "5m")
	flag.Parse()

	// Parse the interval at which updates occur
	parsedUpdateInterval, err := time.ParseDuration(*updateInterval)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the blog update and start the repeated update task.
	b := NewBlog(*viewDir, *articleDir)
	go b.UpdateTask(parsedUpdateInterval)

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))
	r.HandleFunc("/", b.GetPostList)
	r.HandleFunc("/{year}/{month}/{day}/{permaTitle}", b.GetPost)

	log.Fatal(http.ListenAndServe(*listen, r))
}
