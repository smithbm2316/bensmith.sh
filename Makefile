webPkg := ./cmd/web
ssgPkg := ./cmd/ssg
appPort := 2323

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

#- dev: run our app and css file watchers in dev mode
.PHONY: dev
dev: clean
	@$(MAKE) --no-print-directory -j2 dev/app dev/vite

# uses `templ generate --watch` to live reload our go + templ code
.PHONY: dev/app
dev/app:
	@./bin/templ generate --watch --cmd "go run ${ssgPkg} --dev"

# uses vite for our dev server and asset processing
.PHONY: dev/vite
dev/vite:
	@npx vite

#- clean: clean our output paths
# this will remove the .site directory and everything in the src directory that isn't in the src/public or src/styles folders or isn't the src/vite.config.js config file
.PHONY: clean
clean:
	@rm -rf .site
	@find src/ -mindepth 1 -path "src/styles" -prune -o -exec rm -rf {} +

#- build: build the application
.PHONY: build
build: clean build/templ build/ssg build/run-ssg build/vite
	@go build -o=./bin/web ${webPkg}

# build the ssg binary for production
.PHONY: build/ssg
build/ssg:
	@go build -o=./bin/ssg ${ssgPkg}

# execute the production ssg binary to generate our files before processing with Vite
.PHONY: build/run-ssg
build/run-ssg:
	@./bin/ssg

# build templ files into go files for production
.PHONY: build/templ
build/templ:
	@./bin/templ generate

# run vite on our generated static site to optimize for production 
.PHONY: build/vite
build/vite:
	@npx vite build

#- preview: run the production-ready application
.PHONY: preview
preview:
	@./bin/web

# run a production server with vite instead of our webPkg
.PHONY: preview/vite
preview/vite:
	@npx vite preview


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

#- tool/install: install tool dependency into ./bin, pass URL as value to arg "dep=[URL]"
.PHONY: tool/install
tool/install:
	@mkdir -pv bin
	GOBIN="$$(pwd)/bin" @go install ${dep}

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
