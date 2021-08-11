package v1alpha2

import (
	"context"

	"knative.dev/pkg/apis"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
)

func (cs *ClusterStore) ConvertTo(ctx context.Context, to apis.Convertible) error {
	toClusterStore := to.(*v1alpha1.ClusterStore)
	toClusterStore.ObjectMeta = cs.ObjectMeta
	if err := cs.Spec.ConvertTo(ctx, &toClusterStore.Spec); err != nil {
		return err
	}
	if err := cs.Status.ConvertTo(ctx, &toClusterStore.Status); err != nil {
		return err
	}

	return nil
}

func (cs *ClusterStore) ConvertFrom(ctx context.Context, from apis.Convertible) error {
	fromClusterStore := from.(*v1alpha1.ClusterStore)
	cs.ObjectMeta = fromClusterStore.ObjectMeta
	if err := cs.Spec.ConvertFrom(ctx, &fromClusterStore.Spec); err != nil {
		return err
	}
	if err := cs.Status.ConvertFrom(ctx, &fromClusterStore.Status); err != nil {
		return err
	}

	return nil
}

func (css *ClusterStoreSpec) ConvertTo(ctx context.Context, to *v1alpha1.ClusterStoreSpec) error {
	for _, source := range css.Sources {
		to.Sources = append(to.Sources, corev1alpha1.StoreImage{Image: source.Image})
	}

	return nil
}

func (css *ClusterStoreSpec) ConvertFrom(ctx context.Context, from *v1alpha1.ClusterStoreSpec) error {
	for _, source := range from.Sources {
		css.Sources = append(css.Sources, corev1alpha1.StoreImage{Image: source.Image})
	}

	return nil
}

func (css *ClusterStoreStatus) ConvertFrom(ctx context.Context, from *v1alpha1.ClusterStoreStatus) error {
	css.Status = from.Status
	css.Buildpacks = from.Buildpacks
	return nil
}

func (css *ClusterStoreStatus) ConvertTo(ctx context.Context, to *v1alpha1.ClusterStoreStatus) error {
	to.Status = css.Status
	to.Buildpacks = css.Buildpacks
	return nil
}
