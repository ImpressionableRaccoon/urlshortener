package osexit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	tests := []string{"a", "b", "c"}
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, tests...)
}
