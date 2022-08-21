package rule

const (
	Include LineType = iota
	Full
	Suffix
)

type LineType int

type Rule struct {
	Type    LineType
	Payload string
	Tags    []string
}

type Ruleset struct {
	Rules []*Rule
}
