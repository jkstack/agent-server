package utils

import "github.com/jkstack/jkframe/logging"

// Recover recover func
func Recover(key string) {
	if err := recover(); err != nil {
		logging.Error("%s: %v", key, err)
	}
}
