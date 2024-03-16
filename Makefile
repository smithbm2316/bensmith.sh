webPkg := ./cmd/web
ssgPkg := ./cmd/ssg
appPort := 2323
# set GOBIN to be local to the project
export GOBIN := ${CURDIR}/bin

# ==================================================================================== #
#
# HELP MENU
#
# ==================================================================================== #

#- help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^#-//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
#-:
#- DEVELOPMENT:
#
# ==================================================================================== #
#- dev: runs `wgo` and `templ` w/2 jobs to watch/reload .templ & .go files
.PHONY: dev
dev: clean/ssg
	@./bin/wgo -file=.go -file=.tmpl -file=.templ -xfile=_templ.go \
		./bin/templ generate :: go run ${ssgPkg} --dev

#- dev/ssg: runs `wgo` and `templ` w/2 jobs to watch/reload .templ & .go files
.PHONY: dev/ssg
dev/ssg: clean/ssg
	@$(MAKE) --no-print-directory -j2 dev/templ dev/go

#- dev/hmr: uses `parcel` for our HMR dev server and bundling assets
.PHONY: dev/hmr
dev/hmr: clean/hmr build/public
	@npx parcel 'src/**/*.*' --dist-dir .site --port ${appPort}

# run `wgo` to rebuild our go app when .go & {html,text}/template .tmpl files change
.PHONY: dev/go
dev/go:
	@./bin/wgo run -file=.go -file=.tmpl ${ssgPkg} --dev

# run `templ` generate in watch/dev mode
.PHONY: dev/templ
dev/templ:
	@./bin/templ generate --watch

#- clean: runs `clean/ssg` and `clean/hmr` sequentially
.PHONY: clean
clean:
	@$(MAKE) --no-print-directory clean/ssg clean/hmr

# clean our SSG output directory
.PHONY: clean/ssg
clean/ssg:
	@rm -rf src/*

# clean our bundler output directories
.PHONY: clean/hmr
clean/hmr:
	@rm -rf .parcel-cache/ .site/

#- build: build the application
.PHONY: build
build: clean build/public build/templ build/ssg build/run-ssg build/frontend
	@go build -o=./bin/web ${webPkg}

# build the ssg binary for production
.PHONY: build/ssg
build/ssg:
	@go build -o=./bin/ssg ${ssgPkg}

# execute the production ssg binary to generate our files before processing with Vite
.PHONY: build/run-ssg
build/run-ssg:
	@./bin/ssg

# copy all files recursively from the public directory into our src directory, so that they are available for parcel to process at the root level of our server 
.PHONY: build/public
build/public:
	@cp -r public/* src

# build templ files into go files for production
.PHONY: build/templ
build/templ:
	@./bin/templ generate

# run parcel on our generated static site to optimize for production 
.PHONY: build/frontend
build/frontend:
	@npx parcel build 'src/**/*.*' --dist-dir .site

#- preview: run the production-ready application
.PHONY: preview
preview:
	@./bin/web


# ==================================================================================== #
#-:
#- TESTING:
#
# ==================================================================================== #

#- test: run all tests
.PHONY: test
test:
	@go test -v -race -buildvcs ./...

#- test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	@go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
#-:
#- HELPERS:
#
# ==================================================================================== #

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
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@go test -race -buildvcs -vet=off ./...

# internal target for using a confirmation step in another target
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
