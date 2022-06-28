package env

import (
	"fmt"
	"os"
)

func IsRunningInDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println(err)
		return true
	}

	return false
}

func GetHost(localHost string, dockerHost string) string {
	if IsRunningInDockerContainer() {
		return dockerHost
	}
	return localHost
}
