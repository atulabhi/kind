/*
Copyright 2018 The Kubernetes Authors.

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

package kube

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/kind/pkg/log"

)

// TODO(bentheelder): plumb through arch

// dockerBuilder implements Bits for a local docker-ized make / bash build
type dockerBuilder struct {
	kubeRoot string
	arch     string
	logger   log.Logger
}

var _ Builder = &dockerBuilder{}

// NewDockerBuilder returns a new Bits backed by the docker-ized build,
// given kubeRoot, the path to the kubernetes source directory
func NewDockerBuilder(logger log.Logger, kubeRoot, arch string) (Builder, error) {
	return &dockerBuilder{
		kubeRoot: kubeRoot,
		arch:     arch,
		logger:   logger,
	}, nil
}

// Build implements Bits.Build
func (b *dockerBuilder) Build() (Bits, error) {
	// cd to k8s source
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// make sure we cd back when done
	defer func() {
		// TODO(bentheelder): set return error?
		_ = os.Chdir(cwd)
	}()
	if err := os.Chdir(b.kubeRoot); err != nil {
		return nil, err
	}

	// capture version info
	_, err = sourceVersion(b.kubeRoot)
	if err != nil {
		return nil, err
	}


	// NOTE: currently there are no defaults so this is essentially a deep copy
	// binaries we want to build
	binDir := filepath.Join(b.kubeRoot,
		"_output", "dockerized", "bin", "linux", b.arch,
	)
	imageDir := filepath.Join(b.kubeRoot,
		"_output", "release-images", b.arch,
	)
        fmt.Println( "Imagedir ------------------------------------------------------------------", imageDir)
	return &bits{
		binaryPaths: []string{
			filepath.Join(binDir, "kubeadm"),
			filepath.Join(binDir, "kubelet"),
			filepath.Join(binDir, "kubectl"),
		},
		imagePaths: []string{
			filepath.Join(imageDir, "kube-apiserver.tar"),
			filepath.Join(imageDir, "kube-controller-manager.tar"),
			filepath.Join(imageDir, "kube-scheduler.tar"),
			filepath.Join(imageDir, "kube-proxy.tar"),
			filepath.Join(imageDir, "etcd.tar"),
		},
		version: "v1.22.8",
	}, nil
}

func dockerBuildOsAndArch(arch string) string {
	return "linux/" + arch
}
