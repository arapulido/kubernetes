/*
Copyright 2019 The Kubernetes Authors.

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

package helloworld

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
	"k8s.io/kubernetes/pkg/kubectl/util/templates"
)

var (
	helloWorldExample = templates.Examples(i18n.T(`
		# Print Hello World and exit 
		kubectl hello-world`))
)

// NewCmdHelloWorld returns a cobra command to say hello to the world
func NewCmdHelloWorld(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hello-world",
		Short:   i18n.T("Print Hello World"),
		Long:    "Print Hello World and exit",
		Example: helloWorldExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(Run())
		},
	}
	return cmd
}

// Run executes version command
func Run() error {

	fmt.Println("Hello World!")
	return nil
}
