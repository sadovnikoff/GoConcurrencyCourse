package compute

import (
	"errors"
	"strings"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
)

const (
	SetCommand = "SET"
	GetCommand = "GET"
	DelCommand = "DEL"
)

type Parser struct {
	logger *common.Logger
}

var (
	errInvalidRequest   = errors.New("invalid request")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)

func NewParser(logger *common.Logger) (*Parser, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Parser{logger: logger}, nil
}

func (p *Parser) Parse(request string) (Query, error) {
	tokens := strings.Fields(request)

	if len(tokens) < 2 {
		p.logger.Debug("%s [%s]", errInvalidRequest.Error(), request)
		return Query{}, errInvalidRequest
	}

	var query Query
	switch tokens[0] {
	case GetCommand, DelCommand:
		if len(tokens) != 2 {
			p.logger.Debug("%s [%s]", errInvalidArguments.Error(), request)
			return Query{}, errInvalidArguments
		}
		query = NewQuery(tokens[0], tokens[1], "")
	case SetCommand:
		if len(tokens) != 3 {
			p.logger.Debug("%s [%s]", errInvalidArguments.Error(), request)
			return Query{}, errInvalidArguments
		}
		query = NewQuery(tokens[0], tokens[1], tokens[2])
	default:
		p.logger.Debug("%s [%s]", errInvalidCommand.Error(), request)
		return Query{}, errInvalidCommand
	}

	return query, nil
}
