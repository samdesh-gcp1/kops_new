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

package model

import (
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/nodeup/nodetasks"
)

// KubeRouterBuilder installs kube-router
type KubeRouterBuilder struct {
	*NodeupModelContext
}

var _ fi.ModelBuilder = &KubeRouterBuilder{}

// Build is responsible for configuring the kube-router
func (b *KubeRouterBuilder) Build(c *fi.ModelBuilderContext) error {
	{
		kubeconfig, err := b.BuildPKIKubeconfig("kube-router")
		if err != nil {
			return err
		}

		c.AddTask(&nodetasks.File{
			Path:     "/var/lib/kube-router/kubeconfig",
			Contents: fi.NewStringResource(kubeconfig),
			Type:     nodetasks.FileType_File,
			Mode:     s("0400"),
		})
	}

	return nil
}
