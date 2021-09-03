package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"

	"github.com/pivotal/kpack/pkg/apis/validate"
)

func (bbs *BuildBuilderSpec) Validate(ctx context.Context) *apis.FieldError {
	return validate.Image(bbs.Image)
}
