package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	var testOps []string

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", WhiteFill))
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", GreenFill))
	}

	l.Post(OperationFunc(func(tx screen.Texture) {
		testOps = append(testOps, "op 1")
	}))

	l.Post(OperationFunc(func(tx screen.Texture) {
			testOps = append(testOps, "op 2")
	
	}))

	l.Post(OperationFunc(func(tx screen.Texture) {
		testOps = append(testOps, "op 3")
	}))

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) != 2 {
		t.Error("Unexpected size of colors:", mt.Colors)
	}

	expectedOps := []string{"op 1", "op 2", "op 3"}
	if !reflect.DeepEqual(testOps, expectedOps) {
		t.Errorf("Unexpected order of operations. Got: %v, Expected: %v", testOps, expectedOps)
	}
}

func TestLoop_Pull(t *testing.T) {
	mq := &MessageQueue{}

	operationA := &queueMock{}
	go func() {
		time.Sleep(50 * time.Millisecond)
		mq.push(operationA)
	}()

	start := time.Now()
	pulledOp := mq.pull()
	elapsed := time.Since(start)

	if !reflect.DeepEqual(pulledOp, operationA) {
		t.Errorf("Expected pulled operation to be equal to the pushed operation")
	}
	if elapsed < 50*time.Millisecond {
		t.Errorf("Expected Pull to block when pulling from an empty queue")
	}

	operationB := &queueMock{}
	operationC := &queueMock{}
	mq.push(operationB)
	mq.push(operationC)
	pulledOp = mq.pull()

	if !reflect.DeepEqual(pulledOp, operationB) {
		t.Error("Expected pulled operation to be the first pushed operation")
	}
	if mq.empty() {
		t.Errorf("Expected queue to be non-empty after pulling an operation")
	}
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(tx screen.Texture) {
		t.Log(msg)
		op(tx)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type queueMock struct{}

func (m *queueMock) Do(t screen.Texture) bool {
	return false
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Pt(800,800)}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}