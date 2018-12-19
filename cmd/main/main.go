package main

import (
	"fmt"
	"github.com/alistair-english/DRC2019/internal/pkg/config"
)

func main() {
	var ans config.CVConfig = config.GetCVConfig()
	fmt.Println(ans.LeftLower.H)
}
