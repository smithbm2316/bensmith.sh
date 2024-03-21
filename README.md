# bensmith.sh

## Project structure

|Folder     | Description |
|---|---|
|bin/       | the built go app and tool dependency binaries (a project-local `$GOBIN`) |
|cmd/       | the go applications (`package main`) |
|content/   | All markdown files that we will process into routes on the site |
|models/    | Application types |
|scripts/   | Miscellaneous build scripts, the Makefile uses some of them |
|static/    | The static assets (JS, fonts, images, `robots.txt`, etc) copied into the final build |
|styles/    | CSS files that we bundle into one file for the site, `styles/index.css` is the entrypoint |
|views/     | HTML templates, Templ components, and text/templates |
|www/       | the final build of the static site |
|Makefile   | Task runner, run `make` to see the help menu for the available task scripts |
|tools.go   | List of tool dependencies for this project |
