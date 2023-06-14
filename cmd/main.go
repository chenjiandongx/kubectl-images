package main

import (
	"fmt"
	"os"
	"regexp"

	kubeimages "github.com/chenjiandongx/kubectl-images"
	"github.com/spf13/cobra"
)

const version = "0.6.3"

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "kubectl-images [podname-regex]",
		Short: "Show container images used in the cluster.",
		Example: `  # display a table of all images in current namespace using podName/containerName/containerImage as columns.
  kubectl images

  # display images info in yaml format
  kubectl images -oy

  # display a table of images that match 'nginx' podname regex in 'dev' namespace using podName/containerImage as columns.
  kubectl images -n dev nginx -c 1,2`,
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			var regx *regexp.Regexp
			var err error
			if len(args) > 0 {
				if regx, err = regexp.Compile(args[0]); err != nil {
					fmt.Fprintf(os.Stderr, "[Oh...] Invalid regex pattern (%q)", args[0])
					return
				}
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			columns, _ := cmd.Flags().GetString("columns")
			format, _ := cmd.Flags().GetString("output-format")
			allNamespace, _ := cmd.Flags().GetBool("all-namespaces")
			kubeConfig, _ := cmd.Flags().GetString("kubeConfig")
			context, _ := cmd.Flags().GetString("context")
			unique, _ := cmd.Flags().GetBool("unique")
			kubeImage := kubeimages.NewKubeImage(regx, kubeimages.Parameters{
				AllNamespace: allNamespace,
				Namespace:    namespace,
				Columns:      columns,
				KubeConfig:   kubeConfig,
				Context:      context,
				Unique:       unique,
			})
			kubeImage.Render(format)
		},
	}
	rootCmd.Flags().BoolP("all-namespaces", "A", false, "if present, list images in all namespaces.")
	rootCmd.Flags().StringP("namespace", "n", "", "if present, list images in the specified namespace only. Use current namespace as fallback.")
	rootCmd.Flags().StringP("columns", "c", "1,2,3", "specify the columns to display, separated by comma. [0:Namespace, 1:PodName, 2:ContainerName, 3:ContainerImage, 4:ImagePullPolicy, 5:ImageSize]")
	rootCmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file to use for CLI requests.")
	rootCmd.Flags().StringP("output-format", "o", "table", "output format. [json(j)|table(t)|yaml(y)]")
	rootCmd.Flags().StringP("context", "C", "", "The name of the kubeconfig context to use.")
	rootCmd.Flags().BoolP("unique", "u", false, "Unique images group by namespace/container/images/pullPolicy.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[Oh...] Failed to exec command: %v", err)
	}
}
