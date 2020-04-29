package main

import (
	"log"
	"os"
	"strconv"
)

// Manages a list of articles stored in memory and provides the ability to update this list from a directory.
type PostManager struct {
	ArticleDir   string
	Posts        map[PermaPath]Post
	IndexedPosts []Post
}

// Create a new post manager. Note that posts are not initially loaded.
func NewPostManager(articleDir string) PostManager {
	return PostManager{ArticleDir: articleDir, Posts: map[PermaPath]Post{}, IndexedPosts: []Post{}}
}

// Load posts into the post manager - unchanged posts will not be reloaded.
func (p *PostManager) UpdatePosts() {
	var newPosts []Post
	postMap := map[PermaPath]Post{}

	// Count up from 0 and search for markdown files `<i>.md`.
	for i := 0; true; i++ {

		// Initialize a new post with default values.
		post := Post{SourcePath: p.ArticleDir + "/" + strconv.Itoa(i) + ".md", Visible: true}

		// If there's an item in the post array at this index, use that instead.
		if i < len(p.IndexedPosts) {
			post = p.IndexedPosts[i]
		}

		// Try to load the post and handle any errors that occur.
		if err := post.LoadPost(); err != nil {
			if os.IsNotExist(err) {
				// If the file simply didn't exist, that means we can stop searching for posts.
				break
			} else {
				log.Fatal(err)
			}
		}

		// Add this post to the list and map.
		newPosts = append(newPosts, post)
		postMap[post.Path] = post
	}

	p.IndexedPosts = newPosts
	p.Posts = postMap
}
