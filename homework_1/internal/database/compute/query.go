package compute

type Query struct {
	cmd  string
	arg1 string
	arg2 string
}

func NewQuery(cmd, arg1, arg2 string) Query {
	return Query{cmd, arg1, arg2}
}

func (q *Query) Command() string {
	return q.cmd
}

func (q *Query) KeyArgument() string {
	return q.arg1
}

func (q *Query) ValueArgument() string {
	return q.arg2
}
