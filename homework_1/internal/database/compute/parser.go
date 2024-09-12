package compute

import (
	"errors"
	"strings"

	"sadovnikoff/go_concurrency_cource/homework_1/internal/common"
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
		p.logger.DLog.Printf("%s [%s]", errInvalidRequest.Error(), request)
		return Query{}, errInvalidRequest
	}

	query := Query{}
	switch tokens[0] {
	case GetCommand, DelCommand:
		if len(tokens) != 2 {
			p.logger.DLog.Printf("%s [%s]", errInvalidArguments.Error(), request)
			return Query{}, errInvalidArguments
		}
		query = Query{tokens[0], tokens[1], ""}
	case SetCommand:
		if len(tokens) != 3 {
			p.logger.DLog.Printf("%s [%s]", errInvalidArguments.Error(), request)
			return Query{}, errInvalidArguments
		}
		query = Query{tokens[0], tokens[1], tokens[2]}
	default:
		p.logger.DLog.Printf("%s [%s]", errInvalidCommand.Error(), request)
		return Query{}, errInvalidCommand
	}

	return query, nil
}
