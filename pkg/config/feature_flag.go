package config

import "os"

type FeatureFlag string

const (
	FeatureFlagInsights FeatureFlag = "PIPECD_FEATURE_FLAG_INSIGHTS"
)

func FeatureFlagEnabled(flag FeatureFlag) bool {
	v := os.Getenv(string(flag))
	switch v {
	case "true":
		return true
	case "enabled":
		return true
	case "on":
		return true
	default:
		return false
	}
}
