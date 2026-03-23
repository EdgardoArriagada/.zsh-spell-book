package args

type ParsedArgs struct {
	Args []string
	Len  int
}

func (p *ParsedArgs) Get(index int) string {
	if index < p.Len {
		return p.Args[index]
	}
	return ""
}

type ParsedWhole struct {
	Content string
}
