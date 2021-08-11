package v1alpha2

import (
	"context"

	"knative.dev/pkg/apis"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
)

func (i *Image) ConvertTo(ctx context.Context, to apis.Convertible) error {
	toImage := to.(*v1alpha1.Image)
	toImage.ObjectMeta = i.ObjectMeta
	if err := i.Spec.ConvertTo(ctx, &toImage.Spec); err != nil {
		return err
	}
	if err := i.Status.ConvertTo(ctx, &toImage.Status); err != nil {
		return err
	}

	return nil
}

func (i *Image) ConvertFrom(ctx context.Context, from apis.Convertible) error {
	fromImage := from.(*v1alpha1.Image)
	i.ObjectMeta = fromImage.ObjectMeta
	if err := i.Spec.ConvertFrom(ctx, &fromImage.Spec); err != nil {
		return err
	}
	if err := i.Status.ConvertFrom(ctx, &fromImage.Status); err != nil {
		return err
	}

	return nil
}

func (is *ImageSpec) ConvertTo(ctx context.Context, to *v1alpha1.ImageSpec) error {
	to.Tag = is.Tag
	to.Builder = is.Builder
	to.ServiceAccount = is.ServiceAccount
	to.CacheSize = is.CacheSize
	to.FailedBuildHistoryLimit = is.FailedBuildHistoryLimit
	to.SuccessBuildHistoryLimit = is.SuccessBuildHistoryLimit
	to.ImageTaggingStrategy = v1alpha1.ImageTaggingStrategy(is.ImageTaggingStrategy)
	to.Source = is.Source
	to.Build = &v1alpha1.ImageBuild{
		Env: is.Build.Env,
	}
	to.Notary = &v1alpha1.NotaryConfig{
		V1: &v1alpha1.NotaryV1Config{
			URL: is.Notary.V1.URL,
			SecretRef: v1alpha1.NotarySecretRef{
				Name: is.Notary.V1.SecretRef.Name,
			},
		},
	}

	return nil
}

func (is *ImageSpec) ConvertFrom(ctx context.Context, from *v1alpha1.ImageSpec) error {
	is.Tag = from.Tag
	is.Builder = from.Builder
	is.ServiceAccount = from.ServiceAccount
	is.Source = from.Source
	is.CacheSize = from.CacheSize
	is.FailedBuildHistoryLimit = from.FailedBuildHistoryLimit
	is.SuccessBuildHistoryLimit = from.SuccessBuildHistoryLimit
	is.ImageTaggingStrategy = ImageTaggingStrategy(from.ImageTaggingStrategy)
	is.Build = &ImageBuild{
		Env: from.Build.Env,
	}
	is.Notary = &NotaryConfig{
		V1: &NotaryV1Config{
			URL: from.Notary.V1.URL,
			SecretRef: NotarySecretRef{
				Name: from.Notary.V1.SecretRef.Name,
			},
		},
	}
	return nil
}

func (is *ImageStatus) ConvertFrom(ctx context.Context, from *v1alpha1.ImageStatus) error {
	is.LatestBuildImageGeneration = from.LatestBuildImageGeneration
	is.BuildCounter = from.BuildCounter
	is.BuildCacheName = from.BuildCacheName
	is.LatestBuildReason = from.LatestBuildReason
	is.LatestBuildRef = from.LatestBuildRef
	is.LatestImage = from.LatestImage
	is.LatestStack = from.LatestStack
	is.Status = from.Status

	return nil
}

func (is *ImageStatus) ConvertTo(ctx context.Context, to *v1alpha1.ImageStatus) error {
	to.LatestBuildImageGeneration = is.LatestBuildImageGeneration
	to.BuildCounter = is.BuildCounter
	to.BuildCacheName = is.BuildCacheName
	to.LatestBuildReason = is.LatestBuildReason
	to.LatestBuildRef = is.LatestBuildRef
	to.LatestImage = is.LatestImage
	to.LatestStack = is.LatestStack
	to.Status = is.Status

	return nil
}
