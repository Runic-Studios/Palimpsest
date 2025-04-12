package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Runic-Studios/Palimpsest/internal/walker"
)

func main() {
	var overlays multiFlag
	var output string

	flag.Var(&overlays, "overlay", "Overlay directory to apply (can specify multiple)")
	flag.Var(&overlays, "o", "Alias for --overlay")

	flag.StringVar(&output, "target", "", "Output directory")
	flag.StringVar(&output, "t", "", "Alias for --target")

	flag.Usage = func() {
		_ = fmt.Errorf("usage: palimpsest -o overlay1 -o overlay2 -t output_dir\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(overlays) == 0 || output == "" {
		_ = fmt.Errorf("error: must specify at least one overlay (-o) and an output directory (-t)\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := walker.Walk(overlays, output); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("overlays %s applied to %s\n", strings.Join(overlays, ", "), output)
	}
}

// multiFlag allows -o flag to be specified multiple times
type multiFlag []string

func (m *multiFlag) String() string {
	return fmt.Sprint(*m)
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}
