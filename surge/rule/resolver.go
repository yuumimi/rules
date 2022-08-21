package rule

import (
	"fmt"
	"sort"

	"github.com/kr328/V2rayDomains2Clash/trie"
)

func Resolve(all map[string]*Ruleset, name string) (map[string][]string, error) {
	tags := map[string]*trie.Trie{}

	if err := resolveRecursive(all, name, tags); err != nil {
		return nil, err
	}

	out := map[string][]string{}

	for tag, domains := range tags {
		d := domains.Dump()

		sort.Strings(d)

		out[tag] = d
	}

	return out, nil
}

func resolveRecursive(all map[string]*Ruleset, name string, tags map[string]*trie.Trie) error {
	node := all[name]
	if node == nil {
		return fmt.Errorf("rule %s not found", name)
	}

	for _, rule := range node.Rules {
		switch rule.Type {
		case Include:
			if err := resolveRecursive(all, rule.Payload, tags); err != nil {
				return err
			}
		case Full:
			for _, tag := range rule.Tags {
				_ = getOrPutTag(tags, tag).Insert(rule.Payload, true)
			}

			_ = getOrPutTag(tags, "").Insert(rule.Payload, true)
		case Suffix:
			for _, tag := range rule.Tags {
				_ = getOrPutTag(tags, tag).Insert(rule.Payload, false)
			}

			_ = getOrPutTag(tags, "").Insert(rule.Payload, false)
		}
	}

	return nil
}

func getOrPutTag(tags map[string]*trie.Trie, name string) *trie.Trie {
	tag, ok := tags[name]
	if !ok {
		tag = trie.New()
		tags[name] = tag
	}

	return tag
}
