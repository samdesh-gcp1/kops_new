/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/upup/pkg/fi"
)

func s(v string) *string {
	return fi.String(v)
}
func TestValidateInstanceProfile(t *testing.T) {
	grid := []struct {
		Input          *kops.IAMProfileSpec
		ExpectedErrors []string
		ExpectedDetail string
	}{
		{
			Input: &kops.IAMProfileSpec{
				Profile: s("arn:aws:iam::123456789012:instance-profile/S3Access"),
			},
		},
		{
			Input: &kops.IAMProfileSpec{
				Profile: s("arn:aws:iam::123456789012:instance-profile/has/path/S3Access"),
			},
		},
		{
			Input: &kops.IAMProfileSpec{
				Profile: s("arn:aws-cn:iam::123456789012:instance-profile/has/path/S3Access"),
			},
		},
		{
			Input: &kops.IAMProfileSpec{
				Profile: s("arn:aws-us-gov:iam::123456789012:instance-profile/has/path/S3Access"),
			},
		},
		{
			Input: &kops.IAMProfileSpec{
				Profile: s("42"),
			},
			ExpectedErrors: []string{"Invalid value::IAMProfile.Profile"},
			ExpectedDetail: "Instance Group IAM Instance Profile must be a valid aws arn such as arn:aws:iam::123456789012:instance-profile/KopsExampleRole",
		},
		{
			Input: &kops.IAMProfileSpec{
				Profile: s("arn:aws:iam::123456789012:group/division_abc/subdivision_xyz/product_A/Developers"),
			},
			ExpectedErrors: []string{"Invalid value::IAMProfile.Profile"},
			ExpectedDetail: "Instance Group IAM Instance Profile must be a valid aws arn such as arn:aws:iam::123456789012:instance-profile/KopsExampleRole",
		},
	}

	for _, g := range grid {
		err := validateInstanceProfile(g.Input, field.NewPath("IAMProfile"))
		allErrs := field.ErrorList{}
		if err != nil {
			allErrs = append(allErrs, err)
		}
		testErrors(t, g.Input, allErrs, g.ExpectedErrors)

		if g.ExpectedDetail != "" {
			found := false
			for _, err := range allErrs {
				if err.Detail == g.ExpectedDetail {
					found = true
				}
			}
			if !found {
				for _, err := range allErrs {
					t.Logf("found detail: %q", err.Detail)
				}

				t.Errorf("did not find expected error %q", g.ExpectedDetail)
			}
		}
	}
}
