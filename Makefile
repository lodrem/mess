DML=insert

.PHONY: build
build:
	go build -o mess

.PHONY: run-example
run-example: build
	./mess generate \
		--dml=$(DML) \
		--num-rows=10 \
		--schema-path=./examples/schema.json \
		--metadata-path=./metadata.json \
		--output-path=./output.sql