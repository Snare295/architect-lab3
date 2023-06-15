package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправленя останнього разу у Receiver

	MQ MessageQueue
	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})

	go func() {
		for !l.stopReq || !l.MQ.empty() {
			op := l.MQ.pull()
			if update := op.Do(l.next); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.MQ.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(screen.Texture) {
		l.stopReq = true
	}))
	<-l.stop
}

type MessageQueue struct {
	Ops     []Operation
	mu      sync.Mutex
	blocked chan struct{}
}

func (MQ *MessageQueue) push(op Operation) {
	MQ.mu.Lock()
	defer MQ.mu.Unlock()

	MQ.Ops = append(MQ.Ops, op)

	if MQ.blocked != nil {
		close(MQ.blocked)
		MQ.blocked = nil
	}
}

func (MQ *MessageQueue) pull() Operation {
	MQ.mu.Lock()
	defer MQ.mu.Unlock()

	for len(MQ.Ops) == 0 {
		MQ.blocked = make(chan struct{})
		MQ.mu.Unlock()
		<-MQ.blocked
		MQ.mu.Lock()
	}

	op := MQ.Ops[0]
	MQ.Ops[0] = nil
	MQ.Ops = MQ.Ops[1:]
	return op
}

func (MQ *MessageQueue) empty() bool {
	MQ.mu.Lock()
	defer MQ.mu.Unlock()

	return len(MQ.Ops) == 0
}