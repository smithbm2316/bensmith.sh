# bensmith.sh

## Project structure

|Folder     | Description |
|---|---|
|bin/       | our built go app and tool dependency binaries (a project-local `$GOBIN`) |
|cmd/       | our go applications (`package main`) |
|content/   | All markdown files that we will process into routes on our site |
|models/    | Application types |
|views/     | HTML templates, Templ components, and text/templates |
|static/    | The static assets (JS, fonts, images, `robots.txt`, etc) copied into the final build |
|styles/    | CSS files that we bundle into one file for our site, `styles/index.css` is the entrypoint |
|www/       | the final build of the static site |
|Makefile   | Task runner, run `make` to see the help menu for the available task scripts |
|tools.go   | List of tool dependencies for this project |
