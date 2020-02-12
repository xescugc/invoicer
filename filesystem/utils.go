package filesystem

import "fmt"

func jsonFilename(name string) string {
	return fmt.Sprintf("%s.json", name)
}
