build:
    go build -o bin/PixelBot .

clean:
    rm -f bin/bot

run: build
    ./bin/PixelBot