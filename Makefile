all: maquiaBot

maquiaBot:
	go build -v

test:
	go test -v ./...

clean:
	rm -f maquiaBot

.PHONY: all clean test

