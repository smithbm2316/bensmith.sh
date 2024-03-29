# set GOBIN to be local to the project (using an absolute path to it)
export GOBIN := ${CURDIR}/bin
# update the Makefile's path with absolute paths to the local $GOBIN and an
# absolute path to the project's lightningcss CLI binary
export PATH := ${CURDIR}/node_modules/lightningcss-cli:$(GOBIN):$(PATH)

# the main package commands
webPkg := ./cmd/web
ssgPkg := ./cmd/ssg
# the output binaries for the main package commands
webBinary := ./bin/web
ssgBinary := ./bin/ssg
# the port to run the dev or prod server on
serverPort := 2323
# the directory of static assets for the site
staticDir := static
# the output directory for the static site
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
# its faster dev mode compilation, so that whenever any of the input files
# (whether being watched by `wgo` or `templ`) we trigger a rebuild of the static
# site with the ${ssgPkg}
# line 5...-1: and run a parallel `wgo` process to update the input css files
# independently of the files above
#- dev: run a file watcher to rebuild the static site automatically
.PHONY: dev
dev: clean
	@wgo -file .go -file .templ -file '.tmpl.*' -file .md \
		-xfile _templ.go -xfile _templ.txt \
		-xdir ${outputDir} -xdir ${staticDir} \
		templ generate \
		:: go run ${ssgPkg} --dev \
		:: wgo -dir styles -file .css \
		-xdir ${outputDir} -xdir ${staticDir} \
		./scripts/build-css.sh -o ${outputDir} -m dev

#- serve: uses `browser-sync` for a auto-reloading dev server
.PHONY: serve
serve:
	@npx browser-sync start \
		--files "${outputDir},${staticDir}" \
		--no-open \
		--port ${serverPort} \
		--serveStatic ${staticDir} \
		--server ${outputDir} \
		--ui-port 2324 $(SERVE_FLAG) $(SERVE_VALUE)

#- serve-local: calls `make serve`, but hides dev server from the local network
.PHONY: serve-local
serve-local:
	@$(MAKE) --no-print-directory serve \
		SERVE_FLAG="--listen" SERVE_VALUE="localhost"

#- clean: remove all build artifacts from the output directory
.PHONY: clean
clean:
	@rm -rf ${outputDir}
	@mkdir ${outputDir}

#- build: run all `build/*` tasks to create the production-ready application
.PHONY: build
build: clean build/assets build/css build/templ build/ssg build/exec-ssg
	@go build -o=${webBinary} ${webPkg}

# transpiles + bundles the css for prod with lightningcss
.PHONY: build/css
build/css:
	@./scripts/build-css.sh -o ${outputDir} -m prod

# build the ssg binary for production
.PHONY: build/ssg
build/ssg:
	@go build -o=${ssgBinary} ${ssgPkg}

# execute the production ssg binary to generate the static HTML
.PHONY: build/exec-ssg
build/exec-ssg:
	@${ssgBinary}

# copy all files recursively from the ${staticDir} directory into the build
# output directory, so that all of the assets that we aren't processing are
# ready for the ${webPkg} server to use in the same ${outputDir} in production
.PHONY: build/assets
build/assets:
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
