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

{
    print()
    print(666)
    print(100 + 200, 300)
}

print(
    1,
    2,
)

print(1)(2)()

`

	src = `
print(666 == 111 + 555, !0, 0)
print(1 < 2, 1.001E-3, true, false)
print(1.23 + 4.56, "Hello, " + "world!")
	`

	tree, err := golan.Parse(src)
	if err != nil {
		log.Fatal(err)
	}
	golan.DumpTree(tree, os.Stdout)
	engine := golan.NewEngine()
	val, err := engine.Execute(tree)
	if err != nil {
		log.Fatal(err)
	}
	if !golan.IsUndefined(val) {
		log.Println(val)
	}
}
