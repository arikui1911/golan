package golan

func Parse(src string) (tree Node, err error) {
	p := &Parser{Buffer: src}
	defer func() {
		p.Recover(recover())
		if p.Err() != nil {
			err = p.Err()
		}
	}()
	p.Init()
	p.ASTBuilderInit(p.Buffer)
	if err := p.Parse(); err != nil {
		p.Raise(err)
	}
	p.Execute()
	tree = p.Finish()
	return
}
