package main

import (
	"fmt"
	"regexp"

	kubeimage "github.com/chenjiandongx/kubectl-images"
	"github.com/spf13/cobra"
)

const version = "0.3.7"

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "kubectl-images [podname-regex]",
		Short: "Show container images used in the cluster.",
		Example: `  # display a table of all images in current namespace using podName/containerName/containerImage as columns.
  kubectl images

  # display a table of images that match 'nginx' podname regex in 'dev' namespace using podName/containerImage as columns.
  kubectl images -n dev nginx -c 1,2`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			var regx *regexp.Regexp
			var err error
			if len(args) > 0 {
				if regx, err = regexp.Compile(args[0]); err != nil {
					fmt.Println("[Oh...] Invalid regex pattern.")
					return
				}
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			columns, _ := cmd.Flags().GetString("columns")
			format, _ := cmd.Flags().GetString("output-format")
			allNamespace, _ := cmd.Flags().GetBool("all-namespaces")
			kubeConfig, _ := cmd.Flags().GetString("kubeConfig")
			context, _ := cmd.Flags().GetString("context")
			kubeImage := kubeimage.NewKubeImage(regx, allNamespace, namespace, columns, kubeConfig, context)
			kubeImage.Render(format)
		},
	}
	rootCmd.Flags().BoolP("all-namespaces", "A", false, "if present, list images in all namespaces.")
	rootCmd.Flags().StringP("namespace", "n", "", "if present, list images in the specified namespace only. Use current namespace as fallback.")
	rootCmd.Flags().StringP("columns", "c", "1,2,3", "specify the columns to display, separated by comma. [0:Namespace, 1:PodName, 2:ContainerName, 3:ContainerImage, 4:ImagePullPolicy]")
	rootCmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file to use for CLI requests.")
	rootCmd.Flags().StringP("output-format", "o", "table", "output format. [json|table]")
	rootCmd.Flags().String("context", "", "The name of the kubeconfig context to use.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
