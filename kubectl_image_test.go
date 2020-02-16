package kubeimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_Count(t *testing.T) {
	counter := NewCounter()
	for _, s := range [...]string{"objA", "objB", "objC", "objB"} {
		counter.add(s)
	}

	assert.Equal(t, counter.Count(), 3)
}

func TestImageEntity_format(t *testing.T) {
	entity := &ImageEntity{
		Namespace:      namespace,
		PodName:        podName,
		ContainerName:  containerName,
		ContainerImage: containerImage,
	}

	assert.Equal(t, entity.format([]string{namespace}), []string{namespace})
	assert.Equal(t, entity.format([]string{namespace, "other_column"}), []string{namespace})
}
