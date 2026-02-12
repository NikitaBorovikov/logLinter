package main

import (
	"loglinter/loglinter"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loglinter.Analyzer)
}
