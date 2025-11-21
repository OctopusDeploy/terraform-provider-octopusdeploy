package internal

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	DeprecationReversalsEnvVar = "TF_OCTOPUS_DEPRECATION_REVERSALS"
	DeprecationKeyProcess      = "Process_v1.0.0"
)

var (
	deprecationCache     map[string]bool
	deprecationCacheOnce sync.Once
)

func getDeprecationCache() map[string]bool {
	deprecationCacheOnce.Do(func() {
		deprecationCache = make(map[string]bool)
		envValue := os.Getenv(DeprecationReversalsEnvVar)
		if envValue != "" {
			enabledKeys := strings.Split(envValue, ",")
			for _, key := range enabledKeys {
				trimmedKey := strings.TrimSpace(key)
				if trimmedKey != "" {
					deprecationCache[trimmedKey] = true
				}
			}
		}
	})
	return deprecationCache
}

func IsDeprecatedResourceEnabled(deprecationKey string) bool {
	cache := getDeprecationCache()
	return cache[deprecationKey]
}

func GetDeprecatedResourceError(resourceName, deprecationKey string) error {
	return fmt.Errorf("the '%s' resource is deprecated and disabled. This resource will be permanently removed in a future version. To temporarily enable it, set the environment variable TF_OCTOPUS_DEPRECATION_REVERSALS=%s", resourceName, deprecationKey)
}
