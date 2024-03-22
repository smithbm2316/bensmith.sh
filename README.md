# bensmith.sh

## Project structure

| Folder      | Description |
|-------------|-------------|
| bin/        | the built go app and tool dependency binaries (a project-local `$GOBIN`) |
| cmd/        | the go applications (`package main`) |
| components/ | Templ components used in the app's routes |
| content/    | All markdown files that we will process into routes on the site |
| docs/       | Any extra documentation or notes for the project |
| routes/     | Templ components that correspond to a specific route on the site |
| scripts/    | Miscellaneous build scripts, the Makefile uses some of them |
| static/     | The static assets (JS, fonts, images, `robots.txt`, etc) copied into the final build |
| styles/     | CSS files that we bundle into one file for the site, `styles/index.css` is the entrypoint |
| templates/  | text/template and html/template files used to generate feeds |
| www/        | the final build of the static site |
| Makefile    | Task runner, run `make` to see the help menu for the available task scripts |
| tools.go    | List of tool dependencies for this project |


## Go packages

| Package     | Description |
|-------------|-------------|
| bs          | the top-level package where most files should live |
| components  | Templ components used in the app's routes |
| routes      | Templ components that correspond to a specific route on the site | 
