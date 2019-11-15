package main

import (
	"crypto/rand"
	tl "github.com/JoelOtter/termloop"
	"math/big"
)

type Room struct {
	*tl.Rectangle
}

func (r Room) Overlaps(s Room) bool {
	px, py := r.Position()
	cx, cy := s.Position()
	pw, ph := r.Size()
	cw, ch := s.Size()
	if px < cx+cw && px+pw > cx &&
		py < cy+ch && py+ph > cy {
		return true
	}

	return false
}

func generateRoom(width, height int) *Room {
	size, _ := rand.Int(rand.Reader, big.NewInt(8))
	s := int(size.Int64())
	drift, _ := rand.Int(rand.Reader, big.NewInt(int64(1+(s/2))))
	d := int(drift.Int64())
	w := s
	h := s

	if check, ok := rand.Int(rand.Reader, big.NewInt(1)); ok != nil {
		if check.Int64() == 0 {
			w += d
		} else {
			h += d
		}
	}

	px := int64(((width-w)/2)*2 + 1)
	py := int64(((height-h)/2)*2 + 1)

	x, _ := rand.Int(rand.Reader, big.NewInt(px))
	y, _ := rand.Int(rand.Reader, big.NewInt(py))

	room := Room{tl.NewRectangle(int(x.Int64()), int(y.Int64()), w, h, tl.ColorBlack)}

	return &room
}

func generateFloor(w, h int) [][]rune {

	// fill with wall tiles
	floor := make([][]rune, w)
	for row := range floor {
		floor[row] = make([]rune, h)
		for ch := range floor[row] {
			floor[row][ch] = '#'
		}
	}

	rooms := make([]*Room, 0, 20)

	// cut out one room
	for a := 0; a < 100; a++ {
		room := generateRoom(w, h)

		// test for overlap
		overlap := false
		for _, r := range rooms {
			if room.Overlaps(*r) {
				overlap = true
				break
			}
		}

		if overlap {
			continue
		}
		rooms = append(rooms, room)

		x, y := room.Position()
		rw, rh := room.Size()

		for i := x; i < x+rw; i++ {
			for j := y; j < y+rh; j++ {
				floor[i][j] = ' '
			}
		}
	}

	return floor
}

type Player struct {
	*tl.Entity
	prevX int
	prevY int
	level *tl.BaseLevel
}

func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		player.prevX, player.prevY = player.Position()
		switch event.Key {
		case tl.KeyArrowRight:
			player.SetPosition(player.prevX+1, player.prevY)
		case tl.KeyArrowLeft:
			player.SetPosition(player.prevX-1, player.prevY)
		case tl.KeyArrowUp:
			player.SetPosition(player.prevX, player.prevY-1)
		case tl.KeyArrowDown:
			player.SetPosition(player.prevX, player.prevY+1)
		}
	}
}

// func (player *Player) Size(int, int) {
// 	return player.Size()
// }

// func (player *Player) Position(int, int) {
// 	return player.Position()
// }

func (player *Player) Collide(collision tl.Physical) {
	// are we colliding with a rectangle?
	if _, ok := collision.(*tl.Rectangle); ok {
		player.SetPosition(player.prevX, player.prevY)
	}
}

func (player Player) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	x, y := player.Position()
	player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
	player.Entity.Draw(screen)
}

type Player struct {
	*tl.Entity
	prevX int
	prevY int
	level *tl.BaseLevel
}

func (player *Player) Draw(screen *tl.Screen) {
	screenWidth, screenHeight := screen.Size()
	x, y := player.Position()
	player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
	player.Entity.Draw(screen)
}

func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		player.prevX, player.prevY = player.Position()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.SetPosition(player.prevX+1, player.prevY)
		case tl.KeyArrowLeft:
			player.SetPosition(player.prevX-1, player.prevY)
		case tl.KeyArrowUp:
			player.SetPosition(player.prevX, player.prevY-1)
		case tl.KeyArrowDown:
			player.SetPosition(player.prevX, player.prevY+1)
		}
	}
}

func (player *Player) Collide(collision tl.Physical) {
	// Check if it's a Rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		player.SetPosition(player.prevX, player.prevY)
	}
}

func main() {
	game := tl.NewGame()

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
		Ch: '-',
	})

	w := 100
	h := 51
	floor := generateFloor(w, h)
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if floor[i][j] == '#' {
				level.AddEntity(tl.NewRectangle(i, j, 1, 1, tl.ColorWhite))
			}
		}
	}

	// level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))

	player := Player{
		Entity: tl.NewEntity(1, 1, 1, 1),
		level:  level,
	}

	// Set the character at position (0, 0) on the entity.
	player.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: '@'})
	level.AddEntity(&player)

	game.Screen().SetLevel(level)
	game.Start()
}
