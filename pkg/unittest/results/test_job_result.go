package results

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/helm-unittest/helm-unittest/pkg/unittest/coverage"
	"github.com/helm-unittest/helm-unittest/pkg/unittest/printer"
)

// TestJobResult result return by TestJob.Run
type TestJobResult struct {
	DisplayName   string
	Index         int
	Passed        bool
	ExecError     error
	AssertsResult []*AssertionResult
	Duration      time.Duration
}

func (tjr TestJobResult) print(printer *printer.Printer, verbosity int) {
	fmt.Println("COVERAGE")
	p, _ := os.Getwd()
	path := fmt.Sprintf("%s/../../test/data/v3/basic/templates/configmap.yaml", p)

	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	in := coverage.NewInstrumenter(&coverage.TreeStrategy{}, data)
	hm, _ := in.Transform()
	instrumentedYaml := string(hm)
	fmt.Println(instrumentedYaml)
	r := coverage.ResultEntry{}
	r.Extract(hm)

	fmt.Println(r)

	fmt.Println("COVERAGE END. exit....")
	os.Exit(1)
	// Just a testing. Very low confidence
	// r := coverage.NewTemplateResult("/Users/ik/source/self/go-workshop/helm-unittest-tmp/test/data/v3/basic/templates/configmap.yaml")
	// r.Extract("/Users/ik/source/self/go-workshop/helm-unittest-tmp/test/data/v3/basic/templates/configmap.yaml")
	// fmt.Println(r)
	// res := coverage.NewTemplateResult()
	// c := coverage.NewCoverageReporter(coverageResults)
	// c.ComputeCoverage(coverageResults)
	if tjr.Passed {
		return
	}

	if tjr.ExecError != nil {
		printer.Println(printer.Highlight("- %s", tjr.DisplayName), 1)
		printer.Println(printer.Highlight("Error: %s\n", tjr.ExecError.Error()), 2)
		return
	}

	printer.Println(printer.Danger("- %s\n", tjr.DisplayName), 1)
	for _, assertResult := range tjr.AssertsResult {
		assertResult.print(printer, verbosity)
	}

}

// Stringify writing the object to a customized formatted string.
func (tjr TestJobResult) Stringify() string {
	content := ""
	if tjr.ExecError != nil {
		content += tjr.ExecError.Error() + "\n"
	}

	for _, assertResult := range tjr.AssertsResult {
		content += assertResult.stringify()
	}

	return content
}
