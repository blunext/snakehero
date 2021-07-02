package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"starter-snake-go/game"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				os.Exit(0)
			}
		}
	}()
	go func() {
		time.Sleep(300 * time.Second)
		done <- true
	}()

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	//moves := newMoves()
	g := game.GameRequest{}
	g.Board.Width = 11
	g.Board.Height = 11
	//g.You.Head.X = 5
	//g.You.Head.Y = 5

	snake := []pairs{
		{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {5, 2}, {5, 1},
	}
	addSnake(&g, snake)
	snake = []pairs{
		{1, 0}, {0, 0},
	}
	addYou(&g, snake)
	food := []game.Coord{
		{8, 8}, {7, 6},
	}
	g.Board.Food = food
	//fmt.Println(s)
	board := game.MakeBoard(g)

	game.FindFood(g.You.Head, board, g.Board.Food)
	moves := game.RankSpace(g.You.Head, board)

	fmt.Println(game.FindBest(moves))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
