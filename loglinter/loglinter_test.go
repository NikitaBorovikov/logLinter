package loglinter

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLogLinter(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "lintertest")
}
