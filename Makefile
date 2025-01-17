build:
	go build -o bin/PixelBot

clean:
	rm -f bin/PixelBot

run: build
	./bin/PixelBot