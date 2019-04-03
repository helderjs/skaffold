/*
Copyright 2019 The Skaffold Authors

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

package context

import (
	"os"

	configutil "github.com/GoogleContainerTools/skaffold/cmd/skaffold/app/cmd/config"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/config"
	kubectx "github.com/GoogleContainerTools/skaffold/pkg/skaffold/kubernetes/context"
	runnerutil "github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/util"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RunContext struct {
	Opts *config.SkaffoldOptions
	Cfg  *latest.Pipeline

	// Plugin is true if at least one artifact is built by a plugin
	// this is temporary - will go away as soon as all the builders are plugins
	Plugin bool

	DefaultRepo string
	KubeContext string
	WorkingDir  string
	Namespaces  []string
}

func GetRunContext(opts *config.SkaffoldOptions, cfg *latest.Pipeline) (*RunContext, error) {
	kubeContext, err := kubectx.CurrentContext()
	if err != nil {
		return nil, errors.Wrap(err, "getting current cluster context")
	}
	logrus.Infof("Using kubectl context: %s", kubeContext)

	// TODO(dgageot): this should be the folder containing skaffold.yaml. Should also be moved elsewhere.
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "finding current directory")
	}

	namespaces, err := runnerutil.GetAllPodNamespaces(opts.Namespace)
	if err != nil {
		return nil, errors.Wrap(err, "getting namespace list")
	}

	defaultRepo, err := configutil.GetDefaultRepo(opts.DefaultRepo)
	if err != nil {
		return nil, errors.Wrap(err, "getting default repo")
	}

	plugin := false
	for _, a := range cfg.Build.Artifacts {
		if a.BuilderPlugin != nil {
			plugin = true
		}
	}
	return &RunContext{
		Opts:        opts,
		Cfg:         cfg,
		Plugin:      plugin,
		WorkingDir:  cwd,
		DefaultRepo: defaultRepo,
		KubeContext: kubeContext,
		Namespaces:  namespaces,
	}, nil
}
