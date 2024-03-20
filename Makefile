# set GOBIN to be local to the project
export GOBIN := ${CURDIR}/bin

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
#- dev/ssg: runs `wgo` & `templ generate --watch` together for watching .go files
.PHONY: dev/ssg
dev/ssg:
	@./bin/wgo \
		-file=.go -file=.tmpl.xml -file=.tmpl.json \
		-xfile=_templ.go -xfile=_templ.txt -xfile=.templ \
		-xdir=bin -xdir=styles -xdir=${staticDir} \
		./bin/templ generate --watch --cmd "go run ${ssgPkg} --dev"

#- dev/css: transpiles + bundles our css in dev mode with lightningcss
.PHONY: dev/css
dev/css:
	@cd styles && ../bin/wgo -file .css \
		npx lightningcss \
		--sourcemap \
		--bundle \
		--custom-media \
		--targets 'defaults' \
		index.css -o ../${outputDir}/styles.css \
		:: echo "bundled and transpiled CSS"

#- dev/hmr: uses `browser-sync` for a auto-reloading dev server
.PHONY: dev/hmr
dev/hmr:
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
	@npx lightningcss \
		--minify \
		--bundle \
		--custom-media \
		--targets 'defaults' \
		styles/index.css -o ${outputDir}/styles.css

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
	@./bin/templ generate

#- preview: run the production-ready application
.PHONY: preview
preview:
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
