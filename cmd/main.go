package main

import (
	"log"
	"os"

	"github.com/arikui1911/golan"
)

func main() {
	src := `
# comment
hoge = piyo = 111 == -123 + +456 * !789

foo = 100

while (foo == 100) {
    bar = 200
}

if foo >= 10 {
    bar = 20
}

if (
    foo <= 11
)
{
    bar = 22
}
else
{
    bar = 33
}

if hoge > 100 {
	piyo = 200
} elsif huga > 300 {
	piyo = 400
} elsif hoge > 500 {
	piyo = 600
} else {
	piyo = 700
}


print()
    print(666)
print(100 + 200, 300)

print(
    1,
    2,
)

print(1)(2)()

`

	p := &golan.Parser{Buffer: src}
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
