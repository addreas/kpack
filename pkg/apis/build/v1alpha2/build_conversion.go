package v1alpha2

import (
	"context"

	"knative.dev/pkg/apis"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
)

func (b *Build) ConvertTo(ctx context.Context, to apis.Convertible) error {
	toBuild := to.(*v1alpha1.Build)
	toBuild.ObjectMeta = b.ObjectMeta
	if err := b.Spec.ConvertTo(ctx, &toBuild.Spec); err != nil {
		return err
	}
	if err := b.Status.ConvertTo(ctx, &toBuild.Status); err != nil {
		return err
	}

	return nil
}

func (b *Build) ConvertFrom(ctx context.Context, from apis.Convertible) error {
	fromBuild := from.(*v1alpha1.Build)
	b.ObjectMeta = fromBuild.ObjectMeta
	if err := b.Spec.ConvertFrom(ctx, &fromBuild.Spec); err != nil {
		return err
	}
	if err := b.Status.ConvertFrom(ctx, &fromBuild.Status); err != nil {
		return err
	}

	return nil
}

func (bs *BuildSpec) ConvertTo(ctx context.Context, to *v1alpha1.BuildSpec) error {
	to.Env = bs.Env
	to.Source = bs.Source
	to.CacheName = bs.CacheName
	to.Resources = bs.Resources
	to.Tags = bs.Tags
	to.LastBuild = &v1alpha1.LastBuild{
		Image:   bs.LastBuild.Image,
		StackId: bs.LastBuild.StackId,
	}
	for _, binding := range bs.Bindings {
		to.Bindings = append(to.Bindings, v1alpha1.Binding{
			Name:        binding.Name,
			MetadataRef: binding.MetadataRef,
			SecretRef:   binding.SecretRef,
		})
	}
	to.ServiceAccount = bs.ServiceAccount
	to.Builder = v1alpha1.BuildBuilderSpec{
		Image:            bs.Builder.Image,
		ImagePullSecrets: bs.Builder.ImagePullSecrets,
	}
	to.Notary = &v1alpha1.NotaryConfig{
		V1: &v1alpha1.NotaryV1Config{
			URL: bs.Notary.V1.URL,
			SecretRef: v1alpha1.NotarySecretRef{
				Name: bs.Notary.V1.SecretRef.Name,
			},
		}}

	return nil
}

func (bs *BuildSpec) ConvertFrom(ctx context.Context, from *v1alpha1.BuildSpec) error {
	bs.Env = from.Env
	bs.Source = from.Source
	bs.CacheName = from.CacheName
	bs.Resources = from.Resources
	bs.Tags = from.Tags
	bs.LastBuild = &LastBuild{
		Image:   from.LastBuild.Image,
		StackId: from.LastBuild.StackId,
	}
	for _, binding := range from.Bindings {
		bs.Bindings = append(bs.Bindings, Binding{
			Name:        binding.Name,
			MetadataRef: binding.MetadataRef,
			SecretRef:   binding.SecretRef,
		})
	}
	bs.ServiceAccount = from.ServiceAccount
	bs.Builder = BuildBuilderSpec{
		Image:            from.Builder.Image,
		ImagePullSecrets: from.Builder.ImagePullSecrets,
	}
	bs.Notary = &NotaryConfig{
		V1: &NotaryV1Config{
			URL: from.Notary.V1.URL,
			SecretRef: NotarySecretRef{
				Name: from.Notary.V1.SecretRef.Name,
			},
		}}
	return nil
}

func (bs *BuildStatus) ConvertFrom(ctx context.Context, from *v1alpha1.BuildStatus) error {

	return nil
}

func (bs *BuildStatus) ConvertTo(ctx context.Context, to *v1alpha1.BuildStatus) error {

	return nil
}
