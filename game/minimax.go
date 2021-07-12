package game

import "github.com/jinzhu/copier"

type headAndNeighbour struct {
	head      *Tile
	neighbour *Tile
}
type neighbours []headAndNeighbour

type neighbourListWithIterator struct {
	neighbours
	iterator int
}

type listOfNeighbourLists []neighbourListWithIterator
type rounds []neighbours

func Minimax(board board, depth int) {
	for _, round := range allCombinations(board) {

		newBoard := copyBoard(board)
		for _, pair := range round {
			foundFood := newBoard.tiles[pair.neighbour.x][pair.neighbour.y].costIndex == food
			newBoard.tiles[pair.neighbour.x][pair.neighbour.y] = &Tile{x: pair.neighbour.x, y: pair.neighbour.y, board: pair.neighbour.board}

			for i, snake := range board.gameData.Board.Snakes {
				if snake.Head.X == pair.head.x && snake.Head.Y == pair.head.y {
					head := Coord{X: pair.neighbour.x, Y: pair.neighbour.y}
					body := append([]Coord{}, head)
					body = append(body, snake.Body...) //todo: make a copy first?
					newBoard.gameData.Board.Snakes[i].Body = body
					newBoard.gameData.Board.Snakes[i].Head = head
					newBoard.gameData.Board.Snakes[i].Health--
					if foundFood {
						newBoard.gameData.Board.Snakes[i].Length++
					} else {
						newBoard.gameData.Board.Snakes[i].Health--
					}
					if board.gameData.You.Head.X == pair.head.x && board.gameData.You.Head.Y == pair.head.y {
						newBoard.gameData.You = newBoard.gameData.Board.Snakes[i]
					}
				}
			}
		}
		if depth == 0 {
			return // evaluate
		}
		depth--
		Minimax(newBoard, depth)
	}

}

func copyBoard(old board) board {
	tiles := make([][]*Tile, old.gameData.Board.Height)
	for i := range tiles {
		tiles[i] = make([]*Tile, old.gameData.Board.Width)
	}

	for y, yTiles := range old.tiles {
		for x, t := range yTiles {
			tiles[x][y] = &Tile{x: t.x, y: t.y, board: t.board} // todo: snakeTileVanish
		}
	}
	gameRequest := GameRequest{}
	copier.Copy(&gameRequest, old.gameData)
	return board{tiles: tiles, gameData: &gameRequest}
}

func allCombinations(board board) rounds {
	list := makeListOfNeighbourLists(board)

	rounds := rounds{}
	for {
		round := neighbours{}
		for _, comb := range list {
			round = append(round, comb.neighbours[comb.iterator])
		}
		rounds = append(rounds, round)
		for i, _ := range list {
			list[i].iterator++
			if list[i].iterator < len(list[i].neighbours) {
				break
			}
			list[i].iterator = 0
		}
		sum := 0
		for _, comb := range list {
			sum += comb.iterator
		}
		if sum == 0 {
			return rounds
		}
	}

}

func makeListOfNeighbourLists(board board) listOfNeighbourLists {
	listOfListsOfNeighbours := listOfNeighbourLists{}
	for _, snake := range board.gameData.Board.Snakes {
		listOfNeighbours := neighbourListWithIterator{}
		head, ok := board.getTile(snake.Head.X, snake.Head.Y)
		if !ok {
			panic("no head in minimax")
		}

		for _, n := range head.Neighbors() {
			pair := headAndNeighbour{head: head, neighbour: n}
			listOfNeighbours.neighbours = append(listOfNeighbours.neighbours, pair)
		}
		listOfListsOfNeighbours = append(listOfListsOfNeighbours, listOfNeighbours)
	}
	return listOfListsOfNeighbours
}