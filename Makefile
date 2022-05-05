default: build

test:
	go test $$(go list ./... | grep -v integration)

e2e: 
	cd integration && go test && cd ../

build:
	go build

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-aws ~/.tflint.d/plugins
