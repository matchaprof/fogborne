package render

import "fmt"

type Tile rune

const (
	CeilingTile Tile = '▀'
	WallTile    Tile = '█'
	LowerTile   Tile = '▄'
	FloorTile   Tile = '·'
	EmptyTile   Tile = ' '
	PlayerTile  Tile = '@'
)

type GameMap struct {
	Width, Height int
	Tiles         [][]Tile
}

func NewGameMap(width, height int) *GameMap {
	tiles := make([][]Tile, height)

	for i := range tiles {
		tiles[i] = make([]Tile, width)
		for j := range tiles[i] {
			if i%2 == 0 && j%2 == 0 {
				tiles[i][j] = FloorTile
			} else {
				tiles[i][j] = EmptyTile
			}
		}
	}

	return &GameMap{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

func (m *GameMap) Draw(offsetX, offsetY int) {
	for y := 0; y < m.Height; y++ {
		for i := 0; i < offsetX; i++ {
			fmt.Print(" ")
		}
		for x := 0; x < m.Width; x++ {
			fmt.Print(string(m.Tiles[y][x]))
		}

		fmt.Println()
	}
}
