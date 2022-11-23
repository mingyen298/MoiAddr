package moi

import "regexp"

type ContentExtractor struct {
	f1 *regexp.Regexp
	f2 *regexp.Regexp
}

type Extractor struct {
	filter *regexp.Regexp
}

func NewExtractor() *Extractor {
	extractor := Extractor{}
	extractor.filter = regexp.MustCompile(`({.*})`)
	return &extractor
}

func (m *Extractor) GetJson(data []byte) []byte {
	sub := m.filter.FindSubmatch(data)

	if len(sub) == 0 {
		return nil
	}
	return sub[1]
}
