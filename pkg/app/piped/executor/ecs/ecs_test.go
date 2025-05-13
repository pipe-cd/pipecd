package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"

	provider "github.com/pipe-cd/pipecd/pkg/app/piped/platformprovider/ecs"
)

func TestFindRemovedTags(t *testing.T) {
	currentTags := []types.Tag{
		{Key: strPtr(provider.LabelManagedBy), Value: strPtr("piped")},
		{Key: strPtr("region"), Value: strPtr("us-west-1")},
		{Key: strPtr("project"), Value: strPtr("abc")},
	}

	desiredTags := []types.Tag{
		{Key: strPtr("project"), Value: strPtr("abc")},
	}

	got := findRemovedTags(currentTags, desiredTags)
	assert.ElementsMatch(t, []string{"region"}, got)
}

func strPtr(s string) *string {
	return &s
}
