// Copied with love from https://github.com/99designs/gqlgen/issues/1171#issuecomment-675045037 as a workaround for the generation error issue explained there
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
)

type ImportsPlugin struct {
}

func (ImportsPlugin) Name() string {
	return "Imports"
}
func (ImportsPlugin) GenerateCode(cfg *codegen.Data) error {
	dir := cfg.Config.Resolver.Dir()

	badLineRegexp := regexp.MustCompile(`\n\s*"errors"\n`)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(info.Name(), ".resolvers.go") {
			return nil
		}
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		if !badLineRegexp.Match(f) {
			return nil
		}
		n := badLineRegexp.ReplaceAll(f, []byte("\n"))
		err = ioutil.WriteFile(path, n, os.ModePerm)
		if err != nil {
			return err
		}
		c := exec.Command("go", "run", "golang.org/x/tools/cmd/goimports", "-w", path)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	})
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	err = api.Generate(cfg,
		api.AddPlugin(ImportsPlugin{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
