package kubeimage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	gotemplate = `go-template={{range .items}} {{.metadata.namespace}} {{","}} {{.metadata.name}} {{","}} {{range .spec.containers}} {{.name}} {{","}} {{.image}} {{"\n"}} {{end}} {{range .spec.initContainers}} {{"(init)"}} {{.name}} {{","}} {{.image}} {{"\n"}} {{end}} {{end}}`

	namespace      = "Namespace"
	podName        = "PodName"
	containerName  = "ContainerName"
	containerImage = "ContainerImage"
)

// KubeImage is the representation of a container image used in the cluster.
type KubeImage struct {
	entities     []*ImageEntity
	allNamespace bool
	namespace    string
	columns      string
	kubeconfig   string
	regx         *regexp.Regexp
}

// NewKubeImage creates a new KubeImage instance.
func NewKubeImage(regx *regexp.Regexp, allNamespace bool, namespace, columns, kubeconfig string) *KubeImage {
	return &KubeImage{
		allNamespace: allNamespace,
		columns:      columns,
		namespace:    namespace,
		kubeconfig:   kubeconfig,
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
	kubecfg := make([]string, 0)
	if ki.kubeconfig != "" {
		kubecfg = append(kubecfg, "--kubeconfig", ki.kubeconfig)
	}

	if ki.allNamespace {
		return append([]string{"get", "pods", "--all-namespaces", "-o", gotemplate}, kubecfg...)
	} else if ki.namespace != "" {
		return append([]string{"get", "pods", "-n", ki.namespace, "-o", gotemplate}, kubecfg...)
	}
	return append([]string{"get", "pods", "-o", gotemplate}, kubecfg...)
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
func (ki *KubeImage) Render(format string) {
	ki.run()

	if len(ki.entities) == 0 {
		fmt.Println("[Oh...] No images matched!")
		return
	}

	switch format {
	case "table":
		{
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
	case "json":
		{
			type PodRecord struct {
				Namespace string
				Pod       string
				Container string
				Image     string
			}
			var rec []PodRecord
			for _, v := range ki.entities {
				rec = append(rec, PodRecord{v.Namespace, v.PodName, v.ContainerName, v.ContainerImage})
			}
			output, err := json.Marshal(rec)
			if err != nil {
				fmt.Println("[Oh...] failed to marshal JSON data")
				return
			}
			fmt.Println(string(output))
		}
	}
}
