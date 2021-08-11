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

func TestBuildConversion(t *testing.T) {
	spec.Run(t, "testBuildConversion", testBuildConversion)
}

func testBuildConversion(t *testing.T, when spec.G, it spec.S) {
	when("converting to and from v1alpha1", func() {
		build := &Build{
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-build",
			},
			Spec:       BuildSpec{
				Tags: []string{"tag"},
				Builder: BuildBuilderSpec{
					Image:            "my-builder=image",
					ImagePullSecrets: []corev1.LocalObjectReference{{
						Name: "secret",
					}},
				},
				ServiceAccount: "default",
				CacheName: "cache",
				Bindings:  nil,
				Env:       []corev1.EnvVar{{
					Name:      "some-var",
					Value:     "some-val",
				}},
				Resources: corev1.ResourceRequirements{
					Limits: map[corev1.ResourceName]resource.Quantity{
						"some-name": resource.MustParse("5M"),
					},
					Requests: map[corev1.ResourceName]resource.Quantity{
						"some-name": resource.MustParse("5M"),
					},
				},
				LastBuild: &LastBuild{
					Image:   "last-image",
					StackId: "my-stack",
				},
				Notary: &NotaryConfig{
					V1: &NotaryV1Config{
						URL: "notary.com",
						SecretRef: NotarySecretRef{
							Name: "notary",
						},
					},
				},
			},
			//todo: fix this
			Status:     BuildStatus{
				Status: corev1alpha1.Status{
					ObservedGeneration: 0,
					Conditions:         nil,
				},
				BuildMetadata: nil,
				Stack: BuildStack{
					RunImage: "",
					ID:       "",
				},
				LatestImage:    "",
				PodName:        "",
				StepStates:     nil,
				StepsCompleted: nil,
			},
		}
		
		it("does not have any data loss with git source", func() {
			build.Spec.Source = corev1alpha1.SourceConfig{
				Git: &corev1alpha1.Git{
					URL:      "my-repo.com",
					Revision: "main",
				},
				SubPath: "my-sub-path",
			}

			v1alpha1Build := &v1alpha1.Build{}
			err := build.DeepCopy().ConvertTo(context.TODO(), v1alpha1Build)
			require.NoError(t, err)

			convertedBackBuild := &Build{}
			err = convertedBackBuild.ConvertFrom(context.TODO(), v1alpha1Build)
			require.NoError(t, err)

			require.Equal(t, build, convertedBackBuild)
		})
		it("does not have any data loss with registry source", func() {
			build.Spec.Source = corev1alpha1.SourceConfig{
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

			v1alpha1Build := &v1alpha1.Build{}
			err := build.DeepCopy().ConvertTo(context.TODO(), v1alpha1Build)
			require.NoError(t, err)

			convertedBackBuild := &Build{}
			err = convertedBackBuild.ConvertFrom(context.TODO(), v1alpha1Build)
			require.NoError(t, err)

			require.Equal(t, build, convertedBackBuild)
		})
		it("does not have any data loss with blob source", func() {
			build.Spec.Source = corev1alpha1.SourceConfig{
				Blob: &corev1alpha1.Blob{
					URL: "my-blob.com",
				},
				SubPath: "my-sub-path",
			}

			v1alpha1Build := &v1alpha1.Build{}
			err := build.DeepCopy().ConvertTo(context.TODO(), v1alpha1Build)
			require.NoError(t, err)

			convertedBackBuild := &Build{}
			err = convertedBackBuild.ConvertFrom(context.TODO(), v1alpha1Build)
			require.NoError(t, err)

			require.Equal(t, build, convertedBackBuild)
		})
	})
}
