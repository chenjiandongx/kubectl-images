package kubeimage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
)

const (
	gotemplate = `go-template={{range .items}} {{.metadata.namespace}} {{","}} {{.metadata.name}} {{","}} {{range .spec.containers}} {{.name}} {{","}} {{.image}} {{","}} {{.imagePullPolicy}} {{"\n"}} {{end}} {{range .spec.initContainers}} {{"(init)"}} {{.name}} {{","}} {{.image}} {{","}} {{.imagePullPolicy}} {{"\n"}} {{end}} {{end}}`

	labelNamespace       = "Namespace"
	labelPod             = "Pod"
	labelContainer       = "Container"
	labelImage           = "Image"
	labelImagePullPolicy = "ImagePullPolicy"
)

type Parameters struct {
	AllNamespace bool
	Namespace    string
	Columns      string
	KubeConfig   string
	Context      string
	Unique       bool
}

// KubeImage is the representation of a container image used in the cluster.
type KubeImage struct {
	entities []*ImageEntity
	columns  []string
	regx     *regexp.Regexp
	params   Parameters
}

// NewKubeImage creates a new KubeImage instance.
func NewKubeImage(regx *regexp.Regexp, params Parameters) *KubeImage {
	names := make([]string, 0)
	for _, c := range stringSplit(params.Columns, ",") {
		switch c {
		case "0":
			names = append(names, labelNamespace)
		case "1":
			names = append(names, labelPod)
		case "2":
			names = append(names, labelContainer)
		case "3":
			names = append(names, labelImage)
		case "4":
			names = append(names, labelImagePullPolicy)
		}
	}

	return &KubeImage{
		columns: names,
		params:  params,
		regx:    regx,
	}
}

// ImageEntity is the representation of an entity to be displayed.
type ImageEntity struct {
	Namespace       string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Pod             string `json:"pod,omitempty" yaml:"pod,omitempty"`
	Container       string `json:"container,omitempty" yaml:"container,omitempty"`
	Image           string `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
}

func (ie *ImageEntity) selectBy(columns []string) []string {
	result := make([]string, 0)
	for _, c := range columns {
		switch c {
		case labelNamespace:
			result = append(result, ie.Namespace)
		case labelPod:
			result = append(result, ie.Pod)
		case labelContainer:
			result = append(result, ie.Container)
		case labelImage:
			result = append(result, ie.Image)
		case labelImagePullPolicy:
			result = append(result, ie.ImagePullPolicy)
		}
	}
	return result
}

func (ie *ImageEntity) filterBy(columns []string) ImageEntity {
	var entity ImageEntity
	for _, c := range columns {
		switch c {
		case labelNamespace:
			entity.Namespace = ie.Namespace
		case labelPod:
			entity.Pod = ie.Pod
		case labelContainer:
			entity.Container = ie.Container
		case labelImage:
			entity.Image = ie.Image
		case labelImagePullPolicy:
			entity.ImagePullPolicy = ie.ImagePullPolicy
		}
	}
	return entity
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

func stringSplit(in, sep string) []string {
	out := make([]string, 0)
	for _, s := range strings.Split(in, sep) {
		out = append(out, strings.TrimSpace(s))
	}
	return out
}

// Commands builds the command to be executed based on user input.
func (ki *KubeImage) Commands() []string {
	kubecfg := make([]string, 0)
	if ki.params.KubeConfig != "" {
		kubecfg = append(kubecfg, "--kubeconfig", ki.params.KubeConfig)
	}

	if ki.params.Context != "" {
		kubecfg = append(kubecfg, "--context", ki.params.Context)
	}

	if ki.params.AllNamespace {
		return append([]string{"get", "pods", "--all-namespaces", "-o", gotemplate}, kubecfg...)
	} else if ki.params.Namespace != "" {
		return append([]string{"get", "pods", "-n", ki.params.Namespace, "-o", gotemplate}, kubecfg...)
	}
	return append([]string{"get", "pods", "-o", gotemplate}, kubecfg...)
}

func (ki *KubeImage) exec() {
	process := exec.Command("kubectl", ki.Commands()...)
	bs, err := process.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Oh...] Execute command error: %v, %s", err, string(bs))
		os.Exit(1)
	}

	entities := make([]*ImageEntity, 0)
	for _, line := range stringSplit(string(bs), "\n") {
		items := stringSplit(line, ",")
		entity := &ImageEntity{}

		switch len(items) {
		case 1:
			continue
		case 3:
			entity.Container = items[0]
			entity.Image = items[1]
			entity.ImagePullPolicy = items[2]
		case 5:
			entity.Namespace = items[0]
			entity.Pod = items[1]
			entity.Container = items[2]
			entity.Image = items[3]
			entity.ImagePullPolicy = items[4]
		}
		entities = append(entities, entity)
	}

	for i := 0; i < len(entities); i++ {
		if entities[i].Pod == "" && i > 0 {
			entities[i].Namespace = entities[i-1].Namespace
			entities[i].Pod = entities[i-1].Pod
		}
	}

	for i := 0; i < len(entities); i++ {
		if ki.regx == nil {
			ki.entities = append(ki.entities, entities[i])
			continue
		}
		if ki.regx.Match([]byte(entities[i].Pod)) {
			ki.entities = append(ki.entities, entities[i])
		}
	}
}

func (ki *KubeImage) groupBy() []*ImageEntity {
	if !ki.params.Unique {
		return ki.entities
	}

	set := make(map[string]struct{})
	entities := make([]*ImageEntity, 0)

	for i, entity := range ki.entities {
		k := fmt.Sprintf("%s/%s/%s/%s", entity.Namespace, entity.Container, entity.Image, entity.ImagePullPolicy)
		if _, ok := set[k]; ok {
			continue
		}

		set[k] = struct{}{}
		entities = append(entities, ki.entities[i])
	}
	return entities
}

func (ki *KubeImage) summary() {
	namespaceCnt := NewCounter()
	podCnt := NewCounter()
	imageCnt := NewCounter()
	containerCnt := 0

	for i := 0; i < len(ki.entities); i++ {
		namespaceCnt.add(ki.entities[i].Namespace)
		podCnt.add(ki.entities[i].Pod)
		imageCnt.add(ki.entities[i].Image)
		containerCnt += 1
	}

	fmt.Fprintf(os.Stdout, "[Summary]: %d namespaces, %d pods, %d containers and %d different images\n",
		namespaceCnt.Count(), podCnt.Count(), containerCnt, imageCnt.Count(),
	)
}

func (ki *KubeImage) tableRender() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(ki.columns)
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	entities := ki.groupBy()
	for _, entity := range entities {
		table.Append(entity.selectBy(ki.columns))
	}
	table.Render()
}

func (ki *KubeImage) jsonRender() {
	entities := ki.groupBy()
	records := make([]ImageEntity, 0, len(entities))
	for _, entity := range entities {
		records = append(records, entity.filterBy(ki.columns))
	}

	output, err := json.MarshalIndent(records, "", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Oh...] Failed to marshal JSON data, error: %v", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, string(output))
}

func (ki *KubeImage) yamlRender() {
	entities := ki.groupBy()
	records := make([]ImageEntity, 0, len(entities))
	for _, entity := range entities {
		records = append(records, entity.filterBy(ki.columns))
	}

	output, err := yaml.Marshal(records)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Oh...] Failed to marshal YAML data, error: %v", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, string(output))
}

// Render renders and displays the table output.
func (ki *KubeImage) Render(format string) {
	ki.exec()

	if len(ki.entities) == 0 {
		fmt.Fprintln(os.Stdout, "[Oh...] No images matched!")
		return
	}

	switch format {
	case "json", "j":
		ki.jsonRender()
	case "yaml", "y":
		ki.yamlRender()
	default: // table
		ki.summary()
		ki.tableRender()
	}
}
