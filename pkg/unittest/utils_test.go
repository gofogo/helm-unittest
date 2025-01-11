package unittest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helm-unittest/helm-unittest/internal/common"
	"github.com/helm-unittest/helm-unittest/internal/printer"
	. "github.com/helm-unittest/helm-unittest/pkg/unittest"
	"github.com/stretchr/testify/assert"
)

// unmarshalJobTestHelper unmarshall a YAML-encoded string into a TestJob struct.
// It extracts the majorVersion, minorVersion, and apiVersions fields from
// CapabilitiesFields and populates the corresponding fields in Capabilities.
// If apiVersions is nil, it sets APIVersions to nil. If it's a slice,
// it appends string values to APIVersions. Returns an error if unmarshaling
// or processing fails.
func unmarshalJobTestHelper(input string, out *TestJob, t *testing.T) {
	t.Helper()
	err := common.YmlUnmarshal(input, &out)
	assert.NoError(t, err)
	out.SetCapabilities()
}

// writeToFile writes the provided string data to a file with the given filename.
// It returns an error if the file cannot be created or if there is an error during writing.
func writeToFile(data string, filename string) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	// Create the file with an absolute path
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

// This file contains unit tests for the TestRunner functionality in the unittest package.
// The purpose of these tests is to verify that the TestRunner behaves correctly when running
// tests on Helm charts, especially when handling multiple complex cases in the test files.
// The tests ensure that the TestRunner can correctly process and report the results of the tests.

// How to add more end-2-end tests
// 1. Create a new function called TestV3RunnerWith_Fixture_Chart_<Context>
// 2. Create fixtures in the `testdata` directory. Example `testdata/chart<number>`
// 3. Create test files in the `tests` directory. Example `testdata/chart<number>/<name>_test.yaml`
// 4. Add metadata information to the test files. Example `testdata/chart<number>/Chart.yaml`

func TestV3RunnerWith_Fixture_Chart_ErrorWhenMetaCharacters(t *testing.T) {
	// buffer := new(bytes.Buffer)
	runner := TestRunner{
		Printer:   printer.NewPrinter(os.Stdout, nil),
		TestFiles: []string{"tests/*_test.yaml"},
	}
	_ = runner.RunV3([]string{"testdata/chart01"})
	// assert.True(t, passed, buffer.String())
}
