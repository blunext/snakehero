package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "gerard",
		Color:      "#880074",
		Head:       "bendr",
		Tail:       "freckled",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("START\n")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
// TODO: Use the information in the GameRequest object to determine your next move.
func HandleMove(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("Next move\n")

	game := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		log.Fatal(err)
	}

	board := makeBoard(game)
	//availableMoves := avoidTakenSpace(game.You.Head, newMoves(), board)
	//
	//move := newMoves()[0] // in
	//if len(availableMoves) > 0 {
	//	move = availableMoves[rand.Intn(len(availableMoves))]
	//}
	moves := rankSpace(game.You.Head, board)
	best := findBest(moves)

	//sort.Slice(best, func(i, j int) bool {
	//	return best.rank < best.rank
	//})

	response := MoveResponse{
		Move: best.heading,
	}

	fmt.Printf("move: %s, latency: %s\n", response.Move, game.You.Latency)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("END\n")
}