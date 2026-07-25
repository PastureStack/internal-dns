.PHONY: build test integration-test validate package

build:
	./scripts/build

test:
	./scripts/test

integration-test: build
	bash ./scripts/integration-test

validate:
	./scripts/validate

package: build
	./scripts/package
