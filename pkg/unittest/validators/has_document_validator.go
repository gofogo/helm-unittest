package validators

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

// HasDocumentsValidator validate whether the count of manifests rendered form template is Count
type HasDocumentsValidator struct {
	Count       int
	FilterAware bool
}

func (v HasDocumentsValidator) failInfo(actual int, not bool) []string {
	expectedCount := strconv.Itoa(v.Count)
	actualCount := strconv.Itoa(actual)
	customMessage := " documents count to be"

	log.WithField("validator", "has_document").Debugln("expected content:", expectedCount)
	log.WithField("validator", "has_document").Debugln("actual content:", actualCount)

	if not {
		return splitInfof(
			setFailFormat(not, false, false, false, customMessage),
			-1,
			-1,
			expectedCount,
		)
	}
	return splitInfof(
		setFailFormat(not, false, true, false, customMessage),
		-1,
		-1,
		expectedCount,
		actualCount,
	)
}

// Validate implement Validatable
func (v HasDocumentsValidator) Validate(context *ValidateContext) (bool, []string) {
	documentsLength := len(context.Docs)
	if v.FilterAware {
		documentsLength = len(context.getManifests())
	}
	if documentsLength == v.Count != context.Negative {
		return true, []string{}
	}
	return false, v.failInfo(documentsLength, context.Negative)
}
