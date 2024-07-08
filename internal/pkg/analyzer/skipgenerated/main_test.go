package skipgenerated

import (
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"honnef.co/go/tools/analysis/facts/generated"
	"honnef.co/go/tools/config"
	"honnef.co/go/tools/stylecheck"
)

func TestAnalyzerSkipGenerated(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор ErrCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	tests := []struct {
		name     string
		analyzer analysis.Analyzer
		paths    []string
	}{
		{
			name: "check skipgenerated via ST1000",
			analyzer: analysis.Analyzer{
				Name: "ST1000",
				Doc:  "test skipgenerated via ST1000",
				Run:  stylecheck.CheckPackageComment,
			},
			paths: []string{
				"./pkg1",
				"./pkg2",
			},
		},
		{
			name: "check skipgenerated via ST1001",
			analyzer: analysis.Analyzer{
				Name:     "ST1001",
				Doc:      "test skipgenerated via ST1001",
				Run:      stylecheck.CheckDotImports,
				Requires: []*analysis.Analyzer{generated.Analyzer, config.Analyzer},
			},
			paths: []string{
				"./pkg3",
			},
		},
	}

	for _, test := range tests {
		analyzers := []*analysis.Analyzer{&test.analyzer}
		WrapList(analyzers)

		for _, path := range test.paths {
			analysistest.Run(t, analysistest.TestData(), analyzers[0], path)
		}
	}
}
