package multichecker

import (
	"github.com/k0st1a/metrics/internal/pkg/analyzer/osexit"

	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
)

func Run() {
	multichecker.Main(
		printf.Analyzer,
		shift.Analyzer,
		structtag.Analyzer,
		osexit.Analyzer,
	)
}
