package v1alpha2

import (
	"context"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	corev1alpha1 "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
)

func TestImageConversion(t *testing.T) {
	spec.Run(t, "Test Image Conversion", testImageConversion)
}

func testImageConversion(t *testing.T, when spec.G, it spec.S) {
	when("converting to v1alpha1 and back without data loss", func() {
		cacheSize := resource.MustParse("5G")
		var buildHistoryLimit int64 = 5

		image := &Image{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-super-convertable-image",
			},
			Spec: ImageSpec{
				Tag:            "my-tag",
				Builder:        corev1.ObjectReference{},
				ServiceAccount: "service-account",
				Source: corev1alpha1.SourceConfig{
					Git: &corev1alpha1.Git{
						URL:      "https://my-github.com/git",
						Revision: "main",
					},
					SubPath: "sub-path",
				},
				CacheSize:                &cacheSize,
				FailedBuildHistoryLimit:  &buildHistoryLimit,
				SuccessBuildHistoryLimit: &buildHistoryLimit,
				ImageTaggingStrategy:     BuildNumber,
				Build: &ImageBuild{
					Env: []corev1.EnvVar{
						{
							Name:  "blah",
							Value: "env",
						},
					},
				},
				Notary: &NotaryConfig{
					V1: &NotaryV1Config{
						URL: "notary.com",
						SecretRef: NotarySecretRef{
							Name: "shhh",
						},
					},
				},
			},
			Status: ImageStatus{
				Status: corev1alpha1.Status{
					ObservedGeneration: 1,
					Conditions: []corev1alpha1.Condition{{
						Type:               corev1alpha1.ConditionReady,
						Status:             "True",
						Severity:           "tornado-warning",
						LastTransitionTime: corev1alpha1.VolatileTime{},
						Reason:             "executive-order",
						Message:            "it-is-too-late",
					}},
				},
				LatestBuildRef:             "some-build",
				LatestBuildImageGeneration: 1,
				LatestImage:                "my-repo/my-image",
				LatestStack:                "io.buildpacks.stacks.full",
				BuildCounter:               1,
				BuildCacheName:             "build-pvc",
				LatestBuildReason:          "COMMIT",
			},
		}

		it("can convert a git source image", func() {
			v1alpha1Image := &v1alpha1.Image{}
			err := image.DeepCopy().ConvertTo(context.TODO(), v1alpha1Image)
			require.NoError(t, err)

			convertedBackImage := &Image{}
			err = convertedBackImage.ConvertFrom(context.TODO(), v1alpha1Image)
			require.NoError(t, err)

			require.Equal(t, image, convertedBackImage)
		})

		it("can convert a registry source image", func() {
			image.Spec.Source = corev1alpha1.SourceConfig{
				Registry: &corev1alpha1.Registry{
					Image: "my-registry.com/image",
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "my-image-secret",
						},
					},
				},
				SubPath: "my-sub-path",
			}

			v1alpha1Image := &v1alpha1.Image{}
			err := image.DeepCopy().ConvertTo(context.TODO(), v1alpha1Image)
			require.NoError(t, err)

			convertedBackImage := &Image{}
			err = convertedBackImage.ConvertFrom(context.TODO(), v1alpha1Image)
			require.NoError(t, err)

			require.Equal(t, image, convertedBackImage)
		})

		it("can convert a blob source image", func() {
			image.Spec.Source = corev1alpha1.SourceConfig{
				Blob: &corev1alpha1.Blob{
					URL: "blob.com",
				},
				SubPath: "my-sub-path",
			}

			v1alpha1Image := &v1alpha1.Image{}
			err := image.DeepCopy().ConvertTo(context.TODO(), v1alpha1Image)
			require.NoError(t, err)

			convertedBackImage := &Image{}
			err = convertedBackImage.ConvertFrom(context.TODO(), v1alpha1Image)
			require.NoError(t, err)

			require.Equal(t, image, convertedBackImage)
		})
	})
}
