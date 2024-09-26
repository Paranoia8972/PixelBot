build:
    go build -o bin/PixelBot ./cmd/bot/main.go

clean:
    rm -f bin/bot

run: build
    ./bin/PixelBot