package painter

import (
	"image"
	"image/color"
	"image/draw"
	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type Figure struct {
	x int
	y int
}

func (f *Figure) Do(t screen.Texture) bool {
	const heightHoriz int = 165
	const widthVert int = 165
	cYellow := color.RGBA{R: 255, G: 255, B: 0, A: 1}

	t.Fill(image.Rect(f.x-widthVert/2, f.y, f.x+widthVert/2, f.y+400), cYellow, draw.Src)
	t.Fill(image.Rect(f.x, f.y-heightHoriz/2, f.x+400, f.y+heightHoriz/2), cYellow, draw.Src)

	return false
}

type Move struct {
	x int 
	y int
	Figure []*Figure
}

func (m *Move) Do(t screen.Texture) bool {
	for _, figure := range m.Figure {
		figure.x += m.x
		figure.y += m.y
	}
	return false
}

type BackRect struct {
	x1 int
	y1 int
	x2 int 
	y2 int
}

func (bcr *BackRect) Do(t screen.Texture) bool {
	c := color.Black
	t.Fill(image.Rect(bcr.x1, bcr.y1, bcr.x2, bcr.y2), c, screen.Src)
	return false
}

func Reset(t screen.Texture) {
	c := color.Black
	t.Fill(t.Bounds(), c, draw.Src)
}
