package main

// Provides a channel that can be used to iterate through the given post array in reverse order.
func ReverseIteratePosts(p []Post) chan Post {
	c := make(chan Post)

	go func() {
		for i := len(p) - 1; i >= 0; i-- {
			if p[i].Visible {
				c <- p[i]
			}
		}

		close(c)
	}()

	return c
}

// Provides functions that are available from within views / templates.
var TemplateFunctions = map[string]interface{}{
	"permalink":  Post.GetPermalink,
	"prettyDate": Post.GetPrettyDate,
	"reverse":    ReverseIteratePosts,
}
