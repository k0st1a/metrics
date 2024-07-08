package pkg1 // want "at least one file in a package should have a package comment"

import "os"

func main() {
	os.Exit(0)
}
