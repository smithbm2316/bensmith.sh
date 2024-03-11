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
	@$(MAKE) --no-print-directory -j2 dev/app dev/frontend

# uses `templ generate --watch` to live reload our go + templ code
.PHONY: dev/app
dev/app:
	@./bin/templ generate --watch --cmd "go run ${ssgPkg} --dev"

# uses parcel for our dev server and asset processing
.PHONY: dev/frontend
dev/frontend:
	@npx parcel 'src/**/*.html' --dist-dir .site --port ${appPort}

#- clean: clean our output paths
.PHONY: clean
clean:
	@rm -rf .site src

#- build: build the application
.PHONY: build
build: clean build/templ build/ssg build/run-ssg build/frontend
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

# run parcel on our generated static site to optimize for production 
.PHONY: build/frontend
build/frontend:
	@npx parcel build 'src/**/*.html' --dist-dir .site

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
