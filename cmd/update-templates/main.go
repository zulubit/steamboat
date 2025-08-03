package main

import (
	"fmt"
	"os"

	"github.com/zulubit/steamboat/internal/update-templates"
)

func main() {
	if err := updatetemplates.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Templates updated successfully!")
}

