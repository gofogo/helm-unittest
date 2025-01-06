package validators

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/helm-unittest/helm-unittest/internal/common"
	log "github.com/sirupsen/logrus"
	// "github.com/yannh/kubeconform/pkg/resource"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	kubecmd "sigs.k8s.io/kubectl-validate/pkg/cmd"
	kubevalidator "sigs.k8s.io/kubectl-validate/pkg/validator"
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

// TODO: schema in yaml -> translate to json first
// TODO: add caching
// TODO: common.K8sManifest should contain byte array and path to file

// Validate implement Validatable
func (v IsValidSchemaValidator) Validate(context *ValidateContext) (bool, []string) {
	manifests := context.getManifests()

	validateSuccess := false
	validateErrors := make([]string, 0)

	var schemaPatchesFs, localSchemasFs fs.FS
	schemaPatchesFs = os.DirFS(".")
	localSchemasFs = os.DirFS(".")

	localCRDsDir := []string{"./Users/ik/source/self/go-workshop/helm-unittest-tmp/pkg/unittest/testdata/chart02/crds"}
	var localCRDsFileSystems []fs.FS
	for _, current := range localCRDsDir {
		localCRDsFileSystems = append(localCRDsFileSystems, os.DirFS(current))
	}

	factory, err := kubevalidator.New(
		openapiclient.NewOverlay(
			// apply user defined patches on top of the final schema
			openapiclient.PatchLoaderFromDirectory(schemaPatchesFs),
			openapiclient.NewComposite(
				// consult local OpenAPI
				openapiclient.NewLocalSchemaFiles(localSchemasFs),
				// consult local CRDs
				openapiclient.NewLocalCRDFiles(localCRDsFileSystems...),
			),
		))
	if err != nil {
		log.Fatalf("failed initializing validator: %s", err)
	}

	for idx, manifest := range manifests {
		// we could provide a schema to validate against
		// Convert K8sManifest to byte array as temporary solution
		m, err := common.YmlMarshall(manifest)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Println("validating with kubevalidator", manifest)
		err = kubecmd.ValidateDocument([]byte(m), factory)
		if err != nil {
			log.Fatalf("failed validating document with kubectl-validator: %s", err)
		} else {
			fmt.Println("kubevalidator success for manifest", manifest)
		}

		// kubeconform example

		// support schema local or remote
		// vr, err := validator.New(v.Schemas, validator.Opts{Strict: true})
		// if err != nil {
		// 	log.Fatalf("failed initializing validator: %s", err)
		// }

		// make sure we are using validate with context
		// r := resource.Resource{
		// 	Bytes: []byte(m),
		// }

		// res := vr.ValidateResource(r)
		// fmt.Println("comport result", res.Status, res.ValidationErrors)
		// if res.Status == validator.Invalid || res.Status == validator.Error {
		// 	// log.Fatalf("resource %d in file is not valid: %s", idx, res.Err)
		// 	fmt.Println(fmt.Sprintf("validator.Invalid resource %d in file is not valid: %s", idx, res.Err))
		// }
		// if res.Status == validator.Valid {
		// 	fmt.Println(fmt.Sprintf("validator.Valid resource '%d' and manifest '%v' is valid", idx, manifest))
		// }

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
