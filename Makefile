all: maquiaBot

maquiaBot:
	go build -v

test:
	go test -v ./...

checkFmt:
	[ -z "$$(git ls-files | grep '\.go$$' | xargs gofmt -l)" ] || (exit 1)

clean:
	rm -f maquiaBot

.PHONY: all clean test checkFmt

