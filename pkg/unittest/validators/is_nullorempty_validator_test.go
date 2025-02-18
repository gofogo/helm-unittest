package validators_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/helm-unittest/helm-unittest/internal/common"
	. "github.com/helm-unittest/helm-unittest/pkg/unittest/validators"
)

var docWithEmptyElements = `
a:
b: ""
c: 0
d: null
e: []
f: {}
`

var docWithNonEmptyElement = `
a: {a: A}
b: "b"
c: 1
d: [d]
`

func TestIsNullOrEmptyValidatorWhenOk(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	for key := range manifest {
		validator := IsNullOrEmptyValidator{key}
		pass, diff := validator.Validate(&ValidateContext{
			Docs: []common.K8sManifest{manifest},
		})

		assert.True(t, pass)
		assert.Equal(t, []string{}, diff)
	}
}

func TestIsNullOrEmptyValidatorWhenNegativeAndOk(t *testing.T) {
	manifest := makeManifest(docWithNonEmptyElement)

	for key := range manifest {
		validator := IsNullOrEmptyValidator{key}
		pass, diff := validator.Validate(&ValidateContext{
			Docs:     []common.K8sManifest{manifest},
			Negative: true,
		})

		assert.True(t, pass)
		assert.Equal(t, []string{}, diff)
	}
}

func TestIsNullOrEmptyValidatorWhenFail(t *testing.T) {
	manifest := makeManifest(docWithNonEmptyElement)

	log.SetLevel(log.DebugLevel)

	for key, value := range manifest {
		validator := IsNullOrEmptyValidator{key}
		valueYAML := common.TrustedMarshalYAML(value)
		pass, diff := validator.Validate(&ValidateContext{
			Docs: []common.K8sManifest{manifest},
		})
		assert.False(t, pass)
		assert.Equal(t, []string{
			"DocumentIndex:	0",
			"ValuesIndex:	0",
			"Path:	" + key,
			"Expected to be null or empty, got:",
			"\t" + string(valueYAML)[:len(valueYAML)-1],
		}, diff)
	}
}

func TestIsNullOrEmptyValidatorWhenNegativeAndFail(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	for key, value := range manifest {
		validator := IsNullOrEmptyValidator{key}
		pass, diff := validator.Validate(&ValidateContext{
			Docs:     []common.K8sManifest{manifest},
			Negative: true,
		})

		valueYAML := common.TrustedMarshalYAML(value)

		assert.False(t, pass)
		assert.Equal(t, []string{
			"DocumentIndex:	0",
			"ValuesIndex:	0",
			"Path:	" + key,
			"Expected NOT to be null or empty, got:",
			"\t" + string(valueYAML)[:len(valueYAML)-1],
		}, diff)
	}
}

func TestIsNullOrEmptyValidatorWhenInvalidPath(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	validator := IsNullOrEmptyValidator{"x.a"}
	pass, diff := validator.Validate(&ValidateContext{
		Docs: []common.K8sManifest{manifest},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:	0",
		"Error:",
		"	unknown path x.a",
	}, diff)
}

func TestIsNullOrEmptyValidatorWhenInvalidPathNegative(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	validator := IsNullOrEmptyValidator{"x.a"}
	pass, diff := validator.Validate(&ValidateContext{
		Docs:     []common.K8sManifest{manifest},
		Negative: true,
	})

	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}

func TestIsNullOrEmptyValidatorWhenInvalidPathFailFast(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	validator := IsNullOrEmptyValidator{"x.a"}
	pass, diff := validator.Validate(&ValidateContext{
		FailFast: true,
		Docs:     []common.K8sManifest{manifest, manifest},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:	0",
		"Error:",
		"	unknown path x.a",
	}, diff)
}

func TestIsNullOrEmptyValidatorWhenFailFast(t *testing.T) {
	manifest := makeManifest(docWithNonEmptyElement)

	log.SetLevel(log.DebugLevel)

	for key, value := range manifest {
		validator := IsNullOrEmptyValidator{key}
		valueYAML := common.TrustedMarshalYAML(value)
		pass, diff := validator.Validate(&ValidateContext{
			FailFast: true,
			Docs:     []common.K8sManifest{manifest, manifest},
		})
		assert.False(t, pass)
		assert.Equal(t, []string{
			"DocumentIndex:	0",
			"ValuesIndex:	0",
			"Path:	" + key,
			"Expected to be null or empty, got:",
			"\t" + string(valueYAML)[:len(valueYAML)-1],
		}, diff)
	}
}

func TestFailWhenInvalidJsonPath(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	validator := IsNullOrEmptyValidator{"x[b]"}
	pass, diff := validator.Validate(&ValidateContext{
		Docs: []common.K8sManifest{manifest, manifest},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t0",
		"Error:",
		"\tinvalid array index [b] before position 4: non-integer array index",
		"DocumentIndex:\t1",
		"Error:",
		"\tinvalid array index [b] before position 4: non-integer array index",
	}, diff)
}

func TestFailWhenInvalidJsonPathFailFast(t *testing.T) {
	manifest := makeManifest(docWithEmptyElements)

	validator := IsNullOrEmptyValidator{"x[b]"}
	pass, diff := validator.Validate(&ValidateContext{
		FailFast: true,
		Docs:     []common.K8sManifest{manifest, manifest},
	})

	assert.False(t, pass)
	assert.Equal(t, []string{
		"DocumentIndex:\t0",
		"Error:",
		"\tinvalid array index [b] before position 4: non-integer array index",
	}, diff)
}

func TestIsNullOrEmptyValidatorWhenNoManifestFail(t *testing.T) {
	validator := IsNullOrEmptyValidator{"key"}
	pass, diff := validator.Validate(&ValidateContext{
		Docs: []common.K8sManifest{},
	})
	assert.False(t, pass)
	assert.Equal(t, []string{
		"Path:\tkey",
		"Expected to be null or empty, got:",
		"\tno manifest found",
	}, diff)
}

func TestIsNullOrEmptyValidatorWhenNoManifestNegativeOk(t *testing.T) {
	validator := IsNullOrEmptyValidator{"key"}
	pass, diff := validator.Validate(&ValidateContext{
		Docs:     []common.K8sManifest{},
		Negative: true,
	})
	assert.True(t, pass)
	assert.Equal(t, []string{}, diff)
}
