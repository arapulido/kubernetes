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

package hellokubernetes

import (
	"fmt"

	"github.com/spf13/cobra"

	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
	"k8s.io/kubernetes/pkg/kubectl/util/templates"
)

// HelloKubernetesOptions is the commandline options for 'hello-kubernetes' sub command
type HelloKubernetesOptions struct {
	FilenameOptions resource.FilenameOptions

	PrintObj func(obj kruntime.Object) error

	genericclioptions.IOStreams
}

var (
	helloKubernetesLong = templates.LongDesc(i18n.T(`
		Print the type and name from a resource from a file or from stdin.

		JSON and YAML formats are accepted.`))

	helloKubernetesExample = templates.Examples(i18n.T(`
		# Print the name and type of a resource using the data in pod.json.
		kubectl hello-kubernetes -f ./pod.json

		# Print the name and type of a resource based on the JSON passed into stdin.
		cat pod.json | kubectl hello-kubernetes -f -`))
)

// NewHelloKubernetesOptions returns an initialized HelloKubernetesOptions instance
func NewHelloKubernetesOptions(ioStreams genericclioptions.IOStreams) *HelloKubernetesOptions {
	return &HelloKubernetesOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdHelloKubernetes returns new initialized instance of hello-kubernetes sub command
func NewCmdHelloKubernetes(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewHelloKubernetesOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "hello-kubernetes -f FILENAME",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Print name and type of a resource from a file or from stdin."),
		Long:                  helloKubernetesLong,
		Example:               helloKubernetesExample,
		Run: func(cmd *cobra.Command, args []string) {
			if cmdutil.IsFilenameSliceEmpty(o.FilenameOptions.Filenames, o.FilenameOptions.Kustomize) {
				ioStreams.ErrOut.Write([]byte("Error: must specify one of -f and -k\n\n"))
				defaultRunFunc := cmdutil.DefaultSubCommandRun(ioStreams.ErrOut)
				defaultRunFunc(cmd, args)
				return
			}
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.ValidateArgs(cmd, args))
			cmdutil.CheckErr(o.RunHelloKubernetes(f, cmd))
		},
	}

	usage := "to use to print the name and type of the resource"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	//	cmdutil.AddValidateFlags(cmd)

	return cmd
}

// ValidateArgs makes sure there is no discrepency in command options
func (o *HelloKubernetesOptions) ValidateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	return nil
}

// Complete completes all the required options
func (o *HelloKubernetesOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	var err error

	template := []byte("{{printf \"Hello %s %s\\n\" .metadata.name .kind}}")
	printer, err := printers.NewGoTemplatePrinter([]byte(template))
	if err != nil {
		return err
	}

	o.PrintObj = func(obj kruntime.Object) error {
		return printer.PrintObj(obj, o.Out)
	}

	return nil
}

// RunHelloKubernetes performs the print
func (o *HelloKubernetesOptions) RunHelloKubernetes(f cmdutil.Factory, cmd *cobra.Command) error {

	cmdNamespace, enforceNamespace, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	r := f.NewBuilder().
		Unstructured().
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, &o.FilenameOptions).
		Flatten().
		Do()
	err = r.Err()
	if err != nil {
		return err
	}

	count := 0
	err = r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		count++

		return o.PrintObj(info.Object)
	})
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no objects passed to print")
	}
	return nil
}
