RED=\033[31m
GREEN=\033[32m
RESET=\033[0m

COLORIZE_PASS = sed "s/^\([- ]*\)\(PASS\)/\1$$(printf "$(GREEN)")\2$$(printf "$(RESET)")/g"
COLORIZE_FAIL = sed "s/^\([- ]*\)\(FAIL\)/\1$$(printf "$(RED)")\2$$(printf "$(RESET)")/g"

VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -s -w \
			-X 'github.com/go-to-k/lamver/internal/version.Version=$(VERSION)' \
			-X 'github.com/go-to-k/lamver/internal/version.Revision=$(REVISION)'
GO_FILES := $(shell find . -type f -name '*.go' -print)

DIFF_FILE := "$$(git diff --name-only --diff-filter=ACMRT | grep .go$ | xargs -I{} dirname {} | sort | uniq | xargs -I{} echo ./{})"

TEST_DIFF_RESULT := "$$(go test -race -cover -v $$(echo $(DIFF_FILE)) -coverpkg=./...)"
TEST_FULL_RESULT := "$$(go test -race -cover -v ./... -coverpkg=./...)"
TEST_COV_RESULT := "$$(go test -race -cover -v ./... -coverpkg=./... -coverprofile=cover.out.tmp)"

FAIL_CHECK := "^[^\s\t]*FAIL[^\s\t]*$$"

test_diff:
	@! echo $(TEST_DIFF_RESULT) | $(COLORIZE_PASS) | $(COLORIZE_FAIL) | tee /dev/stderr | grep $(FAIL_CHECK) > /dev/null
test:
	@! echo $(TEST_FULL_RESULT) | $(COLORIZE_PASS) | $(COLORIZE_FAIL) | tee /dev/stderr | grep $(FAIL_CHECK) > /dev/null
test_view:
	@! echo $(TEST_COV_RESULT) | $(COLORIZE_PASS) | $(COLORIZE_FAIL) | tee /dev/stderr | grep $(FAIL_CHECK) > /dev/null
	cat cover.out.tmp | grep -v "**_mock.go" > cover.out
	rm cover.out.tmp
	go tool cover -func=cover.out
	go tool cover -html=cover.out -o cover.html
lint:
	golangci-lint run
lint_diff:
	golangci-lint run $$(echo $(DIFF_FILE))
mockgen:
	go generate ./...
deadcode:
	deadcode ./...
shadow:
	find . -type f -name '*.go' | sed -e "s/\/[^\.\/]*\.go//g" | uniq | xargs shadow
cognit:
	gocognit -top 10 ./ | grep -v "_test.go"
run:
	go mod tidy
	go run -ldflags "$(LDFLAGS)" cmd/lamver/main.go $${OPT}
build: $(GO_FILES)
	go mod tidy
	go build -ldflags "$(LDFLAGS)" -o lamver cmd/lamver/main.go
install:
	go install -ldflags "$(LDFLAGS)" github.com/go-to-k/lamver/cmd/lamver
clean:
	go clean
	rm -f lamver

# Test data generation commands
# ==================================

# Run test functions generator
testgen:
	@echo "Running test functions generator..."
	@cd testdata && go mod tidy && go run cmd/create/main.go $(OPT)

# Clean test functions
testgen_clean:
	@echo "Cleaning test functions..."
	@cd testdata && go mod tidy && go run cmd/delete/main.go $(OPT)

# Help for test functions generation
testgen_help:
	@echo "Test functions generation targets:"
	@echo "  testgen         - Run the test functions generator"
	@echo "  testgen_clean   - Clean test functions"
	@echo ""
	@echo "Example usage:"
	@echo "  make testgen OPT=\"-p myprofile\""
	@echo "  make testgen_clean OPT=\"-p myprofile\""