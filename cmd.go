package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	gotemplate = `go-template={{range .items}} {{.metadata.namespace}} {{","}} {{.metadata.name}} {{","}} {{range .spec.containers}} {{.name}} {{","}} {{.image}} {{"\n"}} {{end}} {{end}}`

	version = "0.1.0"

	namespace      = "Namespace"
	podName        = "PodName"
	containerName  = "ContainerName"
	containerImage = "ContainerImage"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "kubectl-image [podname-regex]",
		Short: "Show container images used in the cluster.",
		Example: `  # display a table of all images in current namespace using podName/containerName/containerImage as columns.
  kubectl image

  # display a table of images that match 'nginx' podname regex in 'dev' namespace using podName/containerImage as columns.
  kubectl image -n dev nginx -c 1,2`,
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
			allNamespace, _ := cmd.Flags().GetBool("all-namespaces")
			kubeImage := NewKubeImage(regx, allNamespace, namespace, columns)
			kubeImage.Render()
		},
	}
	rootCmd.Flags().BoolP("all-namespaces", "A", false, "if present, list images in all namespaces.")
	rootCmd.Flags().StringP("namespace", "n", "", "if present, list images in the specified namespace only. Use current namespace as fallback.")
	rootCmd.Flags().StringP("columns", "c", "1,2,3", "specify the columns to display, separated by comma. [0:Namespace, 1:PodName, 2:ContainerName, 3:ContainerImage]")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// KubeImage is the representation of a container image used in the cluster.
type KubeImage struct {
	entities     []*ImageEntity
	allNamespace bool
	namespace    string
	columns      string
	regx         *regexp.Regexp
	command      []string
}

// NewKubeImage creates a new KubeImage instance.
func NewKubeImage(regx *regexp.Regexp, allNamespace bool, namespace, columns string) *KubeImage {
	return &KubeImage{
		allNamespace: allNamespace,
		columns:      columns,
		namespace:    namespace,
		regx:         regx,
	}
}

// ImageEntity is the representation of an entity to be displayed.
type ImageEntity struct {
	Namespace      string
	PodName        string
	ContainerName  string
	ContainerImage string
}

func (ie *ImageEntity) format(columns []string) []string {
	result := make([]string, 0)
	for _, c := range columns {
		switch c {
		case namespace:
			result = append(result, ie.Namespace)
		case podName:
			result = append(result, ie.PodName)
		case containerName:
			result = append(result, ie.ContainerName)
		case containerImage:
			result = append(result, ie.ContainerImage)
		}
	}
	return result
}

// Counter is a simple counter.
type Counter struct {
	cnt   int
	items map[string]bool
}

// NewCounter creates a new Counter instance.
func NewCounter() *Counter {
	return &Counter{items: make(map[string]bool)}
}

func (c *Counter) add(obj string) {
	if !c.items[obj] {
		c.cnt += 1
		c.items[obj] = true
	}
}

// Count returns current counter reading.
func (c *Counter) Count() int {
	return c.cnt
}

func (ki *KubeImage) stringSplit(in, sep string) []string {
	out := make([]string, 0)
	for _, s := range strings.Split(in, sep) {
		out = append(out, strings.TrimSpace(s))
	}
	return out
}

// Columns builds a display row with specified columns.
func (ki *KubeImage) Columns() []string {
	result := make([]string, 0)
	for _, c := range ki.stringSplit(ki.columns, ",") {
		switch c {
		case "0":
			result = append(result, namespace)
		case "1":
			result = append(result, podName)
		case "2":
			result = append(result, containerName)
		case "3":
			result = append(result, containerImage)
		}
	}
	return result
}

// Commands builds the command to be executed based on user input.
func (ki *KubeImage) Commands() []string {
	if ki.allNamespace {
		return []string{"get", "pods", "-A", "-o", gotemplate}
	} else if ki.namespace != "" {
		return []string{"get", "pods", "-n", ki.namespace, "-o", gotemplate}
	}
	return []string{"get", "pods", "-o", gotemplate}
}

func (ki *KubeImage) run() {
	process := exec.Command("kubectl", ki.Commands()...)
	bs, err := process.CombinedOutput()
	if err != nil {
		fmt.Printf("[Oh...] Execute command error: %s\n", err.Error())
		return
	}
	output := string(bs)

	entities := make([]*ImageEntity, 0)
	for _, line := range ki.stringSplit(output, "\n") {
		items := ki.stringSplit(line, ",")
		entity := &ImageEntity{}

		switch len(items) {
		case 1:
			continue
		case 2:
			entity.ContainerName = items[0]
			entity.ContainerImage = items[1]
		case 4:
			entity.Namespace = items[0]
			entity.PodName = items[1]
			entity.ContainerName = items[2]
			entity.ContainerImage = items[3]
		}
		entities = append(entities, entity)
	}

	for i := 0; i < len(entities); i++ {
		if entities[i].PodName == "" && i > 0 {
			entities[i].Namespace = entities[i-1].Namespace
			entities[i].PodName = entities[i-1].PodName
		}
	}

	for i := 0; i < len(entities); i++ {
		if ki.regx == nil {
			ki.entities = append(ki.entities, entities[i])
			continue
		}
		if ki.regx.Match([]byte(entities[i].PodName)) {
			ki.entities = append(ki.entities, entities[i])
		}
	}
}

func (ki *KubeImage) summary() {
	namespaceCnt := NewCounter()
	podCnt := NewCounter()
	imageCnt := NewCounter()
	containerCnt := 0

	for i := 0; i < len(ki.entities); i++ {
		namespaceCnt.add(ki.entities[i].Namespace)
		podCnt.add(ki.entities[i].PodName)
		imageCnt.add(ki.entities[i].ContainerImage)
		containerCnt += 1
	}

	fmt.Println(fmt.Sprintf("[Summary]: %d namespaces, %d pods, %d containers and %d different images",
		namespaceCnt.Count(), podCnt.Count(), containerCnt, imageCnt.Count(),
	))
}

// Render renders and displays the table output.
func (ki *KubeImage) Render() {
	ki.run()

	if len(ki.entities) == 0 {
		fmt.Println("[Oh...] No images matched!")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(ki.Columns())
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, v := range ki.entities {
		table.Append(v.format(ki.Columns()))
	}

	ki.summary()
	table.Render()
}
