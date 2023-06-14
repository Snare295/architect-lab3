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
	BgRect  *painter.BackRect
	BackOp  painter.Operation
	Move    painter.Operation
	Figures []*painter.Figure
	res     []painter.Operation
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	p.BgRect, p.BackOp, p.Figures, p.res = nil, nil, nil, nil

	for scanner.Scan() {
		commandLine := scanner.Text()
		if len(commandLine) == 0 {
			continue
		}
		err := p.parse(commandLine)
		if err != nil {
			return nil, err
		}

	}
	p.res = append(p.res, p.BackOp)
	p.res = append(p.res, p.BgRect)
	p.res = append(p.res, p.Move)
	
	for _, figure := range p.Figures {
		p.res = append(p.res, figure)
	}
	
	return p.res, scanner.Err()
}

func (p *Parser) parse(commandLine string) error {
	fields := strings.Fields(commandLine)
	operation := fields[0]
	var args []int

	for i := 1; i < len(fields); i++ {
		arg, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return err
		}
		arg = arg * 800.0
		args = append(args, int(arg))
	}
	switch operation {
	case "white":
		p.BackOp = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.BackOp = painter.OperationFunc(painter.GreenFill)
	case "update":
		p.res = append(p.res, painter.UpdateOp)
	case "figure":
		figure := &painter.Figure{X: args[0], Y: args[1]}
		p.Figures = append(p.Figures, figure)
	case "Move":
		p.Move = &painter.Move{X: args[0], Y: args[1], Figure: p.Figures}
	case "BgRect":
		p.BgRect = &painter.BackRect{X1: args[0], Y1: args[1], X2: args[2], Y2: args[3]}
	case "reset":
		p.Figures = p.Figures[:0]
		p.BgRect = nil
		p.BackOp = painter.OperationFunc(painter.Reset)
	default:
		return errors.New("Failed")
	}
	return nil
}
