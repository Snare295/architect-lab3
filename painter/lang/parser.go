package lang

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/Snare295/architect-lab3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	bgRect  *painter.BackRect
	backOp  painter.Operation
	move    []painter.Operation
	figures []*painter.Figure
	res     []painter.Operation
	upNow  painter.Operation
	updated bool
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.upNow = nil
	p.res = nil
	
	if p.backOp == nil {
		p.backOp = painter.OperationFunc(painter.WhiteFill)
	}

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		if len(commandLine) == 0 {
			continue
		}
		err := p.parse(commandLine)
		if err != nil {
			return nil, err
		}
		if p.updated {
			p.updated = false
			p.res = p.generateOperation()
		}
	}

	return p.res, scanner.Err()
}

func (p *Parser) parse(commandLine string) error {
	var args []int
	f := strings.Fields(commandLine)
	op := f[0]

	for i := 1; i < len(f); i++ {
		arg, err := strconv.ParseFloat(f[i], 64)
		
		if err != nil {
			return err
		}
		
		if arg > 0 && arg < 1 {
			arg = arg * 800.0
		}
	
		args = append(args, int(arg))
	}

	switch op {
	case "white": p.backOp = painter.OperationFunc(painter.WhiteFill)
	case "green": p.backOp = painter.OperationFunc(painter.GreenFill)
	case "update": p.updated, p.upNow = !p.updated, painter.UpdateOp
	case "bgrect": p.bgRect = &painter.BackRect{X1: args[0], Y1: args[1], X2: args[2], Y2: args[3]}
	case "figure": p.figures = append(p.figures, &painter.Figure{X: args[0], Y: args[1]})
	case "move": p.move = append(p.move, &painter.Move{X: args[0], Y: args[1], Figure: p.figures})
	case "reset": p.figures, p.bgRect, p.move, p.backOp = nil, nil, nil, painter.OperationFunc(painter.Reset)
	default: return errors.New("Failed")
	}

	return nil
}

func (p *Parser) generateOperation() []painter.Operation {
	var res []painter.Operation

	if p.backOp != nil { res = append(res, p.backOp) }
	if p.bgRect != nil { res = append(res, p.bgRect) }
	if p.move != nil { res = append(res, p.move...); p.move = nil }
	for _, figure := range p.figures { res = append(res, figure) }
	if p.upNow != nil { res = append(res, p.upNow) }

	return res
}