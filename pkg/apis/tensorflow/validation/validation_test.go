package validation

import (
	"testing"

	tfv1 "github.com/tensorflow/k8s/pkg/apis/tensorflow/v1alpha1"

	"github.com/gogo/protobuf/proto"
	"k8s.io/api/core/v1"
)

func TestValidate(t *testing.T) {
	type testCase struct {
		in             *tfv1.TFJobSpec
		expectingError bool
	}

	testCases := []testCase{
		{
			in: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
									},
								},
							},
						},
						TFReplicaType: tfv1.MASTER,
						Replicas:      proto.Int32(1),
					},
				},
				TFImage: "tensorflow/tensorflow:1.3.0",
			},
			expectingError: false,
		},
		{
			in: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
									},
								},
							},
						},
						TFReplicaType: tfv1.WORKER,
						Replicas:      proto.Int32(1),
					},
				},
				TFImage: "tensorflow/tensorflow:1.3.0",
			},
			expectingError: true,
		},
		{
			in: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
									},
								},
							},
						},
						TFReplicaType: tfv1.WORKER,
						Replicas:      proto.Int32(1),
					},
				},
				TFImage: "tensorflow/tensorflow:1.3.0",
				TerminationPolicy: &tfv1.TerminationPolicySpec{
					Chief: &tfv1.ChiefSpec{
						ReplicaName:  "WORKER",
						ReplicaIndex: 0,
					},
				},
			},
			expectingError: false,
		},
	}

	for _, c := range testCases {
		job := &tfv1.TFJob{
			Spec: *c.in,
		}
		tfv1.SetObjectDefaults_TFJob(job)
		if err := ValidateTFJobSpec(&job.Spec); (err != nil) != c.expectingError {
			t.Errorf("unexpected validation result: %v", err)
		}
	}
}
