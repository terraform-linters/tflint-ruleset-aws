default: build

# The generators tag is set so the generator packages and their tests, which
# an untagged build excludes, are covered.
test:
	go test -tags generators $$(go list -tags generators ./... | grep -v integration)

e2e: 
	cd integration && go test && cd ../

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-aws ~/.tflint.d/plugins

release:
	cd tools/release; go run main.go
