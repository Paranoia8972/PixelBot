package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	bot := flag.Bool("bot", false, "Start the bot")
	server := flag.Bool("server", false, "Start the server")
	all := flag.Bool("all", false, "Start both bot and server")
	flag.Parse()

	var wg sync.WaitGroup

	if !*bot && !*server && !*all {
		wg.Add(2)
		go startBot(&wg)
		go startServer(&wg)
	} else {
		if *bot {
			wg.Add(1)
			go startBot(&wg)
		}
		if *server {
			wg.Add(1)
			go startServer(&wg)
		}
		if *all {
			wg.Add(2)
			go startBot(&wg)
			go startServer(&wg)
		}
	}

	wg.Wait()
}

func startBot(wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("go", "run", "./cmd/bot/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting bot:", err)
	}
}

func startServer(wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("go", "run", "./cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
