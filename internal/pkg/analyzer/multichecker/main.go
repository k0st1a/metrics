// Package multichecker for static code analys.
package multichecker

import (
	"github.com/k0st1a/metrics/internal/pkg/analyzer/osexit"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"github.com/kyoh86/nolint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"go.uber.org/nilaway"
	nilawayconfig "go.uber.org/nilaway/config"
)

func Run() {
	analyzers := append(
		staticcheckAnalyzers(),
		printf.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		osexit.Analyzer,
		nolint.Analyzer,
		nilawayconfig.Analyzer,
		nilaway.Analyzer,
	)

	multichecker.Main(nolint.WrapAll(analyzers...)...)
}

func staticcheckAnalyzers() []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range stylecheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range quickfix.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	return checks
}
