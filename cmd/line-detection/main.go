package main

import (
	"github.com/alistair-english/DRC2019/internal/pkg/driod"
)

func main() {
	droid := droid.New(0, 0, 100, false)
	droid.Move(100, 0)
}
