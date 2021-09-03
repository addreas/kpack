package v1alpha2

import (
	"context"
	"encoding/json"
	"fmt"

	"knative.dev/pkg/apis"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
)

func (i *Image) ConvertTo(_ context.Context, to apis.Convertible) error {
	switch toImage := to.(type) {
	case *v1alpha1.Image:
		toImage.ObjectMeta = i.ObjectMeta
		i.Spec.convertTo(&toImage.Spec)
		i.Status.convertTo(&toImage.Status)

		if annotation, ok := i.Annotations[V1Alpha1BindingsAnnotation]; ok {
			var v1alpha1Bindings v1alpha1.Bindings

			err := json.Unmarshal([]byte(annotation), &v1alpha1Bindings)
			if err != nil {
				return err
			}

			toImage.Spec.Build.Bindings = v1alpha1Bindings
		}

	default:
		return fmt.Errorf("unknown version, got: %T", toImage)
	}

	return nil
}

func (i *Image) ConvertFrom(_ context.Context, from apis.Convertible) error {
	switch fromImage := from.(type) {
	case *v1alpha1.Image:
		i.ObjectMeta = fromImage.ObjectMeta
		i.Spec.convertFrom(&fromImage.Spec)
		i.Status.convertFrom(&fromImage.Status)
		if bindings := fromImage.Spec.Build.Bindings; bindings != nil {
			if i.Annotations == nil {
				i.Annotations = map[string]string{}
			}

			bytes, err := json.Marshal(bindings)
			if err != nil {
				return err
			}

			i.Annotations[V1Alpha1BindingsAnnotation] = string(bytes)
		}
	default:
		return fmt.Errorf("unknown version, got: %T", fromImage)
	}

	return nil
}

func (is *ImageSpec) convertTo(to *v1alpha1.ImageSpec) {
	to.Tag = is.Tag
	to.Builder = is.Builder
	to.ServiceAccount = is.ServiceAccount
	if is.Cache.Volume != nil {
		to.CacheSize = is.Cache.Volume.Size
	}
	to.FailedBuildHistoryLimit = is.FailedBuildHistoryLimit
	to.SuccessBuildHistoryLimit = is.SuccessBuildHistoryLimit
	to.ImageTaggingStrategy = is.ImageTaggingStrategy
	to.Source = is.Source
	to.Build = &v1alpha1.ImageBuild{
		Env:       is.Build.Env,
		Resources: is.Build.Resources,
	}
	to.Notary = is.Notary
}

func (is *ImageSpec) convertFrom(from *v1alpha1.ImageSpec) {
	is.Tag = from.Tag
	is.Builder = from.Builder
	is.ServiceAccount = from.ServiceAccount
	is.Source = from.Source
	is.Cache = &ImageCacheConfig{
		Volume: &ImagePersistentVolumeCache{
			Size: from.CacheSize,
		},
	}
	is.FailedBuildHistoryLimit = from.FailedBuildHistoryLimit
	is.SuccessBuildHistoryLimit = from.SuccessBuildHistoryLimit
	is.ImageTaggingStrategy = from.ImageTaggingStrategy
	is.Build = &ImageBuild{
		Env:       from.Build.Env,
		Resources: from.Build.Resources,
	}
	is.Notary = from.Notary
}

func (is *ImageStatus) convertFrom(from *v1alpha1.ImageStatus) {
	is.LatestBuildImageGeneration = from.LatestBuildImageGeneration
	is.BuildCounter = from.BuildCounter
	is.BuildCacheName = from.BuildCacheName
	is.LatestBuildReason = from.LatestBuildReason
	is.LatestBuildRef = from.LatestBuildRef
	is.LatestImage = from.LatestImage
	is.LatestStack = from.LatestStack
	is.Status = from.Status
}

func (is *ImageStatus) convertTo(to *v1alpha1.ImageStatus) {
	to.LatestBuildImageGeneration = is.LatestBuildImageGeneration
	to.BuildCounter = is.BuildCounter
	to.BuildCacheName = is.BuildCacheName
	to.LatestBuildReason = is.LatestBuildReason
	to.LatestBuildRef = is.LatestBuildRef
	to.LatestImage = is.LatestImage
	to.LatestStack = is.LatestStack
	to.Status = is.Status
}
