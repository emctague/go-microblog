# go-microblog

<small>[demo of default source code](https://go-microblog-demo.herokuapp.com/)
</small>

This is a small, extremely simple platform for blogging via go.

Similarly to Jekyll, you author markdown files with YAML front matter in order to create posts.
However, `go-microblog` is *not* a static site generator - it acts as a server.

## Installation

1. Run `go get -u github.com/emctague/go-microblog`
2. Install it with `go install github.com/emctague/go-microblog`
3. Create a directory with `articles`, `views`, and `static` directories. You can either author your own theme, or copy
   the files from `go-microblog`'s own `views` and `static` directories.
4. Write articles under the `articles` directory, as markdown files numbered in sequential order.
5. Run `go-microblog`! Visit `localhost:8080` to check out your blog!

For more information on article authoring, customization and configuration, read on...

## Authoring Articles

Markdown files live under `articles/` (or the specified `articleDir`, see the Configuration section below) and are
numbered sequentially starting at `0.md`.
Be warned: `go-microblog` will stop searching for new articles at the first number it fails to find, so make sure to
keep your article numbers in order!

An article looks like this:

```markdown
title: The article's title!
perma_title: the-title-part-of-the-permalink
created: 2020-04-22
---
This is the markdown content of the article.
The article's main title is automatically created, so don't bother writing the article's
title as a markdown header!
```

`go-microblog` stores information about articles in memory, and searches for new articles and changes to existing
articles every five minutes. If you want to preview your article ASAP, I highly recommend using a local 'staging' server
to preview before pushing article changes to your real blog.

If you want to save an article without it actually being accessible, use `visible: false` in the front matter.

## Static Files and Theming

The `views` directory (or whatever directory you specify with `viewsDir`, see the Configuration section below) contains
go HTML template files defining the markup of pages.
Rather than authoring new templates from scratch, I recommend modifying the existing templates, which are simple and
easy to read. 

`views/base.gohtml` contains two sub-templates, `header` and `footer`, that are included from the other templates to
provide the common content at the beginning and end of the page. If you plan on making a custom theme, this is the file
you should look into first! This happens to be the only view file used for utility features at the moment, but you can
author as many as you please, as long as they're in `.gohtml` files! 

`views/index.gohtml` provides the main page of the blog, which presents a list of articles. These are provided as an
array of articles, starting with the first article created - to iterate over them such that the latest article appears
on top, you can use `{{range reverse .}}`. Within this range, you can refer to all the same properties of a post listed
in the section about `views/post.gohtml` below.

`views/post.gohtml` provides the contents of a single post. A post consists of:

- `{{ .Title }}` - The title of the article
- `{{ .Body }}` - The body of the article
- `{{ .Created }}` - When the article was created, in a not-particularly-friendly date format.
- `{{ prettyDate . }}` - A human-readable creation date, e.g. `January 2, 2002`.
- `{{ permalink . }}` - The path to this article, relative to the blog's root, e.g. `/2002/1/2/article-perma-title`.

The server reloads these view files every five minutes (or the specified `updateInterval`, see the configuration section
below). While designing a theme, I highly recommend running a staging server locally and simply restarting it to preview
changes.

Static files go in the `static` directory (or whatever directory you specify as `staticDir`, see Configuration.)
These files receive URLs under the path `/static/` on the server.

## Configuration

The following values can be provided as either command-line flags in the format `-key=value`, or via environment
variables. The command-line flag takes precedence over the environment variable if both are provided.

|Command-Line Flag|Environment Variable|Default|Description|
|---|---|---|---|
|ip|IP|0.0.0.0|The IP the server listens on.|
|port|PORT|8080|The port the server listens on.|
|viewDir|VIEW_DIR|./views|The directory from which views (theme files) are read.|
|articleDir|ARTICLE_DIR|./articles|The directory from which articles are read.|
|staticDir|STATIC_DIR|./static|The directory from which static files are read.|
|updateInterval|UPDATE_INTERVAL|5m|The interval at which articles and view files should be reloaded from memory (Articles that have not changed are not reloaded.)|
