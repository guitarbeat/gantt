package main

import (
	"fmt"
	"phd-dissertation-planner/src/core"
)

func main() {
	fmt.Println("PHD PROPOSAL:", core.GenerateCategoryColor("PHD PROPOSAL"))
	fmt.Println("PROPOSAL:", core.GenerateCategoryColor("PROPOSAL"))
	fmt.Println("PhD Proposal:", core.GenerateCategoryColor("PhD Proposal"))
}
