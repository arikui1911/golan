package main

import (
	"log"
	"os"

	"github.com/arikui1911/golan"
)

func main() {
	p := &golan.Parser{Buffer: `hoge = piyo = 123 + 456 * 789`}
	defer func() {
		p.Recover(recover())
		if p.Err() != nil {
			log.Fatal(p.Err())
		}
	}()
	p.Init()
	p.ASTBuilderInit(p.Buffer)
	if err := p.Parse(); err != nil {
		p.Raise(err)
	}
	p.Execute()
	golan.DumpTree(p.Finish(), os.Stdout)
}
