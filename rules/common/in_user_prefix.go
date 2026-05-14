package common

import (
	"fmt"
	"strings"

	C "github.com/metacubex/mihomo/constant"
)

type InUserPrefix struct {
	*Base
	prefixes []string
	adapter  string
	payload  string
}

func (u *InUserPrefix) Match(metadata *C.Metadata, helper C.RuleMatchHelper) (bool, string) {
	if metadata.InUser == "" {
		return false, ""
	}
	for _, p := range u.prefixes {
		if strings.HasPrefix(metadata.InUser, p) {
			return true, u.adapter
		}
	}
	return false, ""
}

func (u *InUserPrefix) RuleType() C.RuleType {
	return C.InUserPrefix
}

func (u *InUserPrefix) Adapter() string {
	return u.adapter
}

func (u *InUserPrefix) Payload() string {
	return u.payload
}

func NewInUserPrefix(iPrefixes, adapter string) (*InUserPrefix, error) {
	prefixes := strings.Split(iPrefixes, "/")
	for i, p := range prefixes {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			return nil, fmt.Errorf("in user prefix couldn't be empty")
		}
		prefixes[i] = p
	}

	return &InUserPrefix{
		Base:     &Base{},
		prefixes: prefixes,
		adapter:  adapter,
		payload:  iPrefixes,
	}, nil
}
