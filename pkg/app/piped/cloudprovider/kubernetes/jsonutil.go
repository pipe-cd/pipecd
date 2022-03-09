package kubernetes

// This file is copied from argocd/gitops-engine and modified
// https://github.com/argoproj/gitops-engine/blob/0bc2f8c395f67123156d4ce6b667bf730618307f/pkg/utils/json/json.go

// https://github.com/ksonnet/ksonnet/blob/master/pkg/kubecfg/diff.go
func removeFields(config, live interface{}) interface{} {
	switch c := config.(type) {
	case map[string]interface{}:
		l, ok := live.(map[string]interface{})
		if ok {
			return removeMapFields(c, l)
		} else {
			return live
		}
	case []interface{}:
		l, ok := live.([]interface{})
		if ok {
			return removeListFields(c, l)
		} else {
			return live
		}
	default:
		return live
	}

}

// removeMapFields remove all non-existent fields in the live that don't exist in the config
func removeMapFields(config, live map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for k, v1 := range config {
		v2, ok := live[k]
		if !ok {
			continue
		}
		if v2 != nil {
			v2 = removeFields(v1, v2)
		}
		result[k] = v2
	}
	return result
}

func removeListFields(config, live []interface{}) []interface{} {
	// If live is longer than config, then the extra elements at the end of the
	// list will be returned as-is so they appear in the diff.
	result := make([]interface{}, 0, len(live))
	for i, v2 := range live {
		if len(config) > i {
			if v2 != nil {
				v2 = removeFields(config[i], v2)
			}
			result = append(result, v2)
		} else {
			result = append(result, v2)
		}
	}
	return result
}
