package main

import (
	"fmt"
	"os"
	"path"

	"github.com/kr328/V2rayDomains2Clash/raw"
	"github.com/kr328/V2rayDomains2Clash/rule"
)

func main() {
	if len(os.Args) < 3 {
		println("Usage: <v2ray-domains-path> <output-path>")

		os.Exit(1)
	}

	data := path.Join(os.Args[1], "data")
	generated := os.Args[2]

	_ = os.MkdirAll(generated, 0755)

	ruleSets, err := rule.ParseDirectory(data)
	if err != nil {
		println("Load domains: " + err.Error())

		os.Exit(1)
	}

	for name := range ruleSets {
		tags, err := rule.Resolve(ruleSets, name)
		if err != nil {
			println("Resolve " + name + ": " + err.Error())

			continue
		}

		for tag, rules := range tags {
			var outputPath string

			if tag == "" {
				outputPath = path.Join(generated, fmt.Sprintf("%s.txt", name))
			} else {
				outputPath = path.Join(generated, fmt.Sprintf("%s@%s.txt", name, tag))
			}

			file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				println("Write file " + outputPath + ": " + err.Error())

				continue
			}

			_, _ = file.WriteString(fmt.Sprintf("# Generated from https://github.com/v2fly/domain-list-community/tree/master/data/%s\n\n", name))
// 			_, _ = file.WriteString(fmt.Sprintf("# Behavior: domain\n\n"))
// 			_, _ = file.WriteString(fmt.Sprintf("payload:\n"))

			for _, domain := range rules {
				_, _ = file.WriteString(fmt.Sprintf("%s\n", domain))
			}

			_ = file.Close()
		}
	}

	raws, err := raw.LoadRawSources()
	if err != nil {
		println("Load raw resources: " + err.Error())

		os.Exit(1)
	}

	for _, r := range raws {
		outputPath := path.Join(generated, r.Name+".txt")

		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			println("Write file " + outputPath + ": " + err.Error())

			continue
		}

		_, _ = file.WriteString(fmt.Sprintf("# Generated from %s\n\n", r.SourceUrl))
// 		_, _ = file.WriteString(fmt.Sprintf("# Behavior: %s\n\n", r.Behavior))
// 		_, _ = file.WriteString(fmt.Sprintf("payload:\n"))

		for _, domain := range r.Rules {
		    _, _ = file.WriteString(fmt.Sprintf("%s\n", domain))
		}

		_ = file.Close()
	}
}
