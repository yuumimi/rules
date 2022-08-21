package rule

import (
	"os"
	"path"
	"strings"
)

func ParseFile(file string) (*Ruleset, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	set := &Ruleset{}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		rule := &Rule{}

		descriptor := fields[0]

		switch {
		case !strings.Contains(descriptor, ":"):
			rule.Type = Suffix
			rule.Payload = descriptor
		case strings.HasPrefix(descriptor, "include:"):
			rule.Type = Include
			rule.Payload = descriptor[len("include:"):]
		case strings.HasPrefix(descriptor, "full:"):
			rule.Type = Full
			rule.Payload = descriptor[len("full:"):]
		case strings.HasPrefix(descriptor, "domain:"):
			rule.Type = Full
			rule.Payload = descriptor[len("domain:"):]
		default:
			println("Unsupported rule: " + line)
			continue
		}

		var tags []string

		for i := 1; i < len(fields); i++ {
			if strings.HasPrefix(fields[i], "@") {
				tags = append(tags, fields[i][len("@"):])
			}
		}

		rule.Tags = tags

		set.Rules = append(set.Rules, rule)
	}

	return set, nil
}

func ParseDirectory(directory string) (map[string]*Ruleset, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	r := map[string]*Ruleset{}

	for _, file := range files {
		entry, err := ParseFile(path.Join(directory, file.Name()))
		if err != nil {
			return nil, err
		}

		r[file.Name()] = entry
	}

	return r, nil
}
