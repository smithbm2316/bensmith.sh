# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/web
TOOL_BIN := bin
BINARY_NAME := bensmith.sh
STATIC_SITE_DIR := build
PUBLIC_DIR := public

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## get-tool: add new dependency, pass URL to dependency as "dep=[URL]"
.PHONY: get-tool
get-tool:
	mkdir -pv ${TOOL_BIN}
	GOBIN="$$(pwd)/${TOOL_BIN}" go install ${dep}

## gen-views: run templ generate with local version
.PHONY: gen-views
gen-views:
	./${TOOL_BIN}/templ generate

## css/live: live reload building CSS
.PHONY: css/live
css/live:
	npx lightningcss \
		--sourcemap \
		--bundle \
		--custom-media \
		--targets '> 0.5% or last 2 versions' \
		styles/main.css -o ${STATIC_SITE_DIR}/main.css

## css: build CSS for production
.PHONY: css
css:
	npx lightningcss \
		--minify \
		--bundle \
		--custom-media \
		--targets '> 0.5% or last 2 versions' \
		styles/main.css -o ${STATIC_SITE_DIR}/main.css

## build-clean: clean the STATIC_SITE_DIR
.PHONY: build-clean
build-clean:
	rm -rf ./${STATIC_SITE_DIR}

## build: build the application
.PHONY: build
build: build-clean css
	# copy all assets in the PUBLIC_DIR into our build output
	cp -ivr ${PUBLIC_DIR} ${STATIC_SITE_DIR}
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the application
.PHONY: run
run: build
	./bin/${BINARY_NAME}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "./bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: tidy audit no-dirty
	git push

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm tidy audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	upx -5 ./bin/linux_amd64/${BINARY_NAME}
	# Include additional deployment steps here...
