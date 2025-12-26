package choices

import (
	"strings"
)

type Choice struct {
	Name    string
	Variant string
	Group   string
}

func (c Choice) Format() string {
	s := c.Name
	if c.Variant != "" {
		s += "/" + c.Variant
	}
	if c.Group != "" {
		s += "@" + c.Group
	}
	return s
}

func (c *Choice) Match(other Choice) bool {
	if c == nil {
		return false
	}

	switch {
	case c.Name != "" && c.Name != other.Name,
		c.Variant != "" && c.Variant != other.Variant,
		c.Group != "" && c.Group != other.Group:
		return false
	}
	return true
}

func (c *Choice) MatchS(s string) bool {
	if c == nil {
		return false
	}
	return c.Match(Parse(s))
}

func Parse(s string) Choice {
	nameVariant, group, _ := strings.Cut(s, "@")
	cName, variant, _ := strings.Cut(nameVariant, "/")
	return Choice{
		Name:    cName,
		Variant: variant,
		Group:   group,
	}
}
