package validators

import (
	"fmt"

	"github.com/helm-unittest/helm-unittest/internal/common"
	log "github.com/sirupsen/logrus"
	"github.com/yannh/kubeconform/pkg/resource"

	// "github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
)

// IsValidSchemaValidator validate manifest against a valid schema
type IsValidSchemaValidator struct {
	// Path string
	Schemas []string
}

func (v IsValidSchemaValidator) failInfo(manifestIndex, actualIndex int, not bool) []string {
	format := "Path:%s expected to "

	if not {
		format = format + "NOT "
	}

	format = format + "exists"

	return splitInfof(
		format,
		manifestIndex,
		actualIndex,
	)
}

// Validate implement Validatable
func (v IsValidSchemaValidator) Validate(context *ValidateContext) (bool, []string) {
	manifests := context.getManifests()

	validateSuccess := false
	validateErrors := make([]string, 0)

	// TODO: schema in yaml -> translate to json first
	// TODO: add caching
	// TODO: common.K8sManifest should contain byte array and path to file

	for idx, manifest := range manifests {
		// we could provide a schema to validate against
		// Convert K8sManifest to byte array as temporary solution
		m, err := common.YmlMarshall(manifest)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		// support schema local or remote
		vr, err := validator.New(v.Schemas, validator.Opts{Strict: true})
		if err != nil {
			log.Fatalf("failed initializing validator: %s", err)
		}

		// make sure we are using validate with context
		r := resource.Resource{
			Bytes: []byte(m),
		}

		res := vr.ValidateResource(r)
		fmt.Println("comport result", res.Status, res.ValidationErrors)
		if res.Status == validator.Invalid || res.Status == validator.Error {
			// log.Fatalf("resource %d in file is not valid: %s", idx, res.Err)
			fmt.Println(fmt.Sprintf("validator.Invalid resource %d in file is not valid: %s", idx, res.Err))
		}
		if res.Status == validator.Valid {
			fmt.Println(fmt.Sprintf("validator.Valid resource '%d' and manifest '%v' is valid", idx, manifest))
		}

		// actual, err := valueutils.GetValueOfSetPath(manifest, v.Path)
		// if err != nil {
		// 	validateSuccess = false
		// 	errorMessage := splitInfof(errorFormat, idx, -1, err.Error())
		// 	validateErrors = append(validateErrors, errorMessage...)
		// 	if context.FailFast {
		// 		break
		// 	}
		// 	continue
		// }

		// if len(actual) > 0 == context.Negative {
		// 	validateSuccess = false
		// 	errorMessage := v.failInfo(idx, -1, context.Negative)
		// 	validateErrors = append(validateErrors, errorMessage...)
		// 	continue
		// }

		validateSuccess = determineSuccess(idx, validateSuccess, true)
	}

	// if len(manifests) == 0 && !context.Negative {
	// 	errorMessage := v.failInfo(-1, -1, context.Negative)
	// 	validateErrors = append(validateErrors, errorMessage...)
	// } else if len(manifests) == 0 && context.Negative {
	// 	validateSuccess = true
	// }

	return validateSuccess, validateErrors
}
