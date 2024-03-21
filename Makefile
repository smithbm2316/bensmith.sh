# set GOBIN to be local to the project (using an absolute path to it)
export GOBIN := ${CURDIR}/bin
# update our Makefile's path with absolute paths to our local $GOBIN and an
# absolute path to our project's lightningcss CLI binary
export PATH := ${CURDIR}/node_modules/lightningcss-cli:$(GOBIN):$(PATH)

# our main package commands
webPkg := ./cmd/web
ssgPkg := ./cmd/ssg
# the output binaries for our main package commands
webBinary := ./bin/web
ssgBinary := ./bin/ssg
# the port to run our dev or prod server on
serverPort := 2323
# the directory of static assets for our site
staticDir := static
# the output directory for our static site
outputDir := www

# ============================================================================ #
#
# HELP MENU
#
# ============================================================================ #

#- help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^#-//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ============================================================================ #
#-:
#- DEVELOPMENT:
#
# ============================================================================ #
# line 1: watch .go, .tmpl.{xml,json}, and .md files
# line 2: ignore watch-mode artifacts from templ
# line 3: ignore ${outputDir} and ${staticDir}
# line 4: run templ generate in watch mode to handle watching .templ files with
# its faster dev mode compilation, so that whenever any of our input files
# (whether being watched by `wgo` or `templ`) we trigger a rebuild of our static
# site with our ${ssgPkg}
# line 5...-1: and run a parallel `wgo` process to update our input css files
# independently of the files above
#- dev: run a file watcher to rebuild our static site automatically
.PHONY: dev
dev:
	@wgo -file .go -file '.tmpl.*' -file .md \
		-xfile _templ.go -xfile _templ.txt -xfile .templ \
		-xdir ${outputDir} -xdir ${staticDir} \
		templ generate --watch --cmd="go run ./cmd/ssg --dev" \
		:: wgo -dir styles -file .css lightningcss \
			--sourcemap \
			--error-recovery \
			--bundle \
			--custom-media \
			--targets 'defaults' \
			styles/index.css -o ${outputDir}/styles.css \
			:: echo "bundled and transpiled CSS"

#- dev/serve: uses `browser-sync` for a auto-reloading dev server
.PHONY: serve
serve:
	@npx browser-sync start \
		--server ${outputDir} \
		--port ${serverPort} \
		--ui-port 2324 \
		--no-open \
		--serveStatic ${staticDir}

#- clean: remove all build artifacts from the output directory
.PHONY: clean
clean:
	@rm -rf ${outputDir}/*

#- build: run all `build/*` tasks to create the production-ready application
.PHONY: build
build: clean build/assets build/css build/templ build/ssg build/run-ssg
	@go build -o=${webBinary} ${webPkg}

# transpiles + bundles our css for prod with lightningcss
.PHONY: build/css
build/css:
	@lightningcss \
		--minify \
		--bundle \
		--custom-media \
		--targets 'defaults' \
		styles/index.css -o ${outputDir}/styles.css
	@echo "Bundled and transpiled CSS"

# build the ssg binary for production
.PHONY: build/ssg
build/ssg:
	@go build -o=${ssgBinary} ${ssgPkg}

# execute the production ssg binary to generate our static HTML
.PHONY: build/run-ssg
build/run-ssg:
	@${ssgBinary}

# copy all files recursively from the ${staticDir} directory into our build
# output directory, so that all of our assets that we aren't processing are
# ready for our ${webPkg} server to use in the same ${outputDir} in production
.PHONY: build/assets
build/assets:
	@mkdir -pv ${outputDir}
	@cp -r ${staticDir}/* ${outputDir}

# build templ files into go files for production
.PHONY: build/templ
build/templ:
	@templ generate

#- preview: run the production-ready application
.PHONY: preview
preview: build
	@${webBinary}


# ============================================================================ #
#-:
#- TESTING:
#
# ============================================================================ #

#- test: run all tests
.PHONY: test
test:
	@go test -v -race -buildvcs ./...

#- test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	@go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out


# ============================================================================ #
#-:
#- HELPERS:
#
# ============================================================================ #

#- install-tools: install all tool dependencies in the `tools.go` file
.PHONY: install-tools
install-tools:
	@echo "installing tool dependencies from 'tools.go' into local GOBIN './bin'"
	@mkdir -pv bin
	@sed -n 's/^.*"\(.*\)".*$$/\1/p' tools.go | xargs -t -I _ go install _

#- audit: run quality control checks
.PHONY: audit
audit:
	@go mod verify
	@go vet ./...
	@go run honnef.co/go/tools/cmd/staticcheck@latest \
		-checks=all,-ST1000,-U1000 ./...
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@go test -race -buildvcs -vet=off ./...

# internal target for using a confirmation step in another target
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
