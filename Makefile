# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/web
TOOL_BIN_DIRNAME := bin
TOOL_BIN := ./${TOOL_BIN_DIRNAME}
BINARY_NAME := bensmith.sh
BINARY_OUTPUT := ${TOOL_BIN}/${BINARY_NAME}
STATIC_SITE_DIR := build
PUBLIC_DIR := public
APP_PORT := 2324
PROXY_PORT := 2323
BROWSER_SYNC_UI_PORT := 2325

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

#- dev/air: run the application with reloading on file changes
.PHONY: dev/air
dev/air: dev/templ dev/css
	${TOOL_BIN}/air

#- dev/serve: use browser-sync for hot-reloading and a dev server
.PHONY: serve
dev/serve:
	npx browser-sync start \
		--files ${BINARY_OUTPUT} \
		--no-open \
		--port ${PROXY_PORT} \
		--proxy 'http://localhost:${APP_PORT}' \
		--ui-port ${BROWSER_SYNC_UI_PORT}

# script for `air`'s .air.toml config to use as the `cmd` to run dev mode with
.PHONY: dev/go
dev/go:
	go build -o=${BINARY_OUTPUT} ${MAIN_PACKAGE_PATH}

# bundle the CSS files in dev mode
.PHONY: dev/css
dev/css:
	npx lightningcss \
		--sourcemap \
		--bundle \
		--custom-media \
		--targets '> 0.5% or last 2 versions' \
		styles/main.css -o ${STATIC_SITE_DIR}/main.css

# run templ generate with local version
.PHONY: dev/templ
dev/templ:
	${TOOL_BIN}/templ generate

# build/css: build CSS for production
.PHONY: build/css
build/css:
	npx lightningcss \
		--minify \
		--bundle \
		--custom-media \
		--targets '> 0.5% or last 2 versions' \
		styles/main.css -o ${STATIC_SITE_DIR}/main.css

# run templ generate with local version
.PHONY: build/templ
build/templ:
	${TOOL_BIN}/templ generate

# build/clean: clean the STATIC_SITE_DIR
.PHONY: build/clean
build/clean:
	rm -rf ./${STATIC_SITE_DIR}

# build/public: copy files from the PUBLIC_DIR to the STATIC_SITE_DIR
.PHONY: build/public
build/public:
	cp -ivr ${PUBLIC_DIR} ${STATIC_SITE_DIR}

#- build: build the application
.PHONY: build
build: build/clean build/public build/templ build/css
	go build -o=${BINARY_OUTPUT} ${MAIN_PACKAGE_PATH}

#- preview: build and run the application
.PHONY: preview
preview: build
	${BINARY_OUTPUT}

#- deploy: deploy the application to production
.PHONY: deploy
deploy: confirm tidy audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	upx -5 ./bin/linux_amd64/${BINARY_NAME}
	# Include additional deployment steps here...


# ==================================================================================== #
#-:
#- TESTING:
#
# ==================================================================================== #

#- test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

#- test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
#-:
#- HELPERS:
#
# ==================================================================================== #

#- tool/install: install tool dependency into TOOL_BIN, pass URL as value to arg "dep=[URL]"
.PHONY: tool/install
tool/install:
	mkdir -pv ${TOOL_BIN}
	GOBIN="$$(pwd)/${TOOL_BIN_DIRNAME}" go install ${dep}

#- audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

#- push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push

#- tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

# internal target for using a confirmation step in another target
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# internal target for making sure that all our changes are committed to git in another target
.PHONY: no-dirty
no-dirty:
	git diff --exit-code
