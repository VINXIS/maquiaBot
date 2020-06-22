all: maquiaBot

maquiaBot:
	go build -v

clean:
	rm -f maquiaBot

.PHONY: all clean

