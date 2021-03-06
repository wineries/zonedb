package build

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
)

type Zone struct {
	Domain      string    `json:"domain,omitempty"`
	InfoURL     string    `json:"infoURL,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Locations   []string  `json:"locations,omitempty"`
	WhoisServer string    `json:"whoisServer,omitempty"`
	WhoisURL    string    `json:"whoisURL,omitempty"`
	NameServers []string  `json:"nameServers,omitempty"`
	CodePoints  CodeTable `json:"codePoints,omitempty"`
	Subdomains  []string  `json:"-"`

	// Exported for use in text/template
	POffset, SOffset, SEnd int    `json:"-"`
	CPOffset, CPEnd        int    `json:"-"`
	TagBits                uint64 `json:"-"`
}

func (z *Zone) Normalize() {
	z.Domain = Normalize(z.Domain)
	var tags []string
	for _, t := range z.Tags {
		tags = append(tags, t)
	}
	z.Tags = NewSet(tags...).Values()
	z.NameServers = NewSet(z.NameServers...).Values()
	z.CodePoints.Compress()
}

func (z *Zone) IsIDN() bool {
	for i := 0; i < len(z.Domain); i++ {
		if z.Domain[i] >= utf8.RuneSelf {
			return true
		}
	}
	return false
}

func (z *Zone) HasMetadata() bool {
	z2 := *z
	z2.Domain = ""
	if j, _ := json.Marshal(&z2); string(j) == "{}" {
		return false
	}
	return true
}

func (z *Zone) ASCII() string {
	s, _ := idna.ToASCII(z.Domain)
	return s
}

func (z *Zone) ParentDomain() string {
	labels := strings.Split(z.Domain, ".")
	if len(labels) == 1 {
		return ""
	}
	return strings.Join(labels[1:], ".")
}
