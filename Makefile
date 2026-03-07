.PHONY: build run clean

build:
	go build -o kase .

run: build
	./kase

clean:
	rm -f kase
