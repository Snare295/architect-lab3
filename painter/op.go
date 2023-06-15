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
	X int
	Y int
}

func (f *Figure) Do(t screen.Texture) bool {
	const HeightHoriz int = 165
	const WidthVert int = 165
	cYellow := color.RGBA{R: 255, G: 255, B: 0, A: 1}

	t.Fill(image.Rect(f.X - WidthVert/2, f.Y - 200, f.X + WidthVert/2, f.Y + 200), cYellow, draw.Src)
	t.Fill(image.Rect(f.X - 200, f.Y - HeightHoriz/2, f.X +200, f.Y + HeightHoriz/2), cYellow, draw.Src)

	return false
}

type Move struct {
	X      int
	Y      int
	Figure []*Figure
}

func (m *Move) Do(t screen.Texture) bool {
	for _, Figure := range m.Figure {
		Figure.X += m.X
		Figure.Y += m.Y
	}
	return false
}

type BackRect struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

func (bcr *BackRect) Do(t screen.Texture) bool {
	c := color.Black
	t.Fill(image.Rect(bcr.X1, bcr.Y1, bcr.X2, bcr.Y2), c, screen.Src)
	return false
}

func Reset(t screen.Texture) {
	c := color.Black
	t.Fill(t.Bounds(), c, draw.Src)
}
