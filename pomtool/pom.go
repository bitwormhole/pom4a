package pomtool

import (
	"encoding/xml"
	"strings"
)

type POM struct {
}

func (inst *POM) Parse(data []byte) (*PomProject, error) {
	project := &PomProject{}
	err := xml.Unmarshal(data, &project)
	if err != nil {
		return nil, err
	}
	project.Properties.Table()
	return project, nil
}

////////////////////////////////////////////////////////////////////////////////

type PomProject struct {
	ModelVersion string `xml:"modelVersion"`

	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
	Packaging  string `xml:"packaging"`

	Properties PomProperties `xml:"properties"`
}

////////////////////////////////////////////////////////////////////////////////

type PomProperties struct {
	InnerXML string `xml:",innerxml"`

	cache     map[string]string
	prevName  string
	prevValue string
}

// GetProperty get required property
func (inst *PomProperties) GetProperty(name string) string {
	tab := inst.Table()
	value := tab[name]
	if value == "" {
		panic("no required property named : " + name)
	}
	return value
}

func (inst *PomProperties) Table() map[string]string {
	c := inst.cache
	if c == nil {
		c = inst.parse()
		inst.cache = c
	}
	return c
}

func (inst *PomProperties) parse() map[string]string {
	text := inst.InnerXML
	chs := []byte(text)
	buffer := &strings.Builder{}
	tab := make(map[string]string)
	for _, b := range chs {
		if b == '<' {
			inst.handleFragmentTag1(buffer, tab)
			buffer.Reset()
		} else if b == '>' {
			inst.handleFragmentTag2(buffer, tab)
			buffer.Reset()
		} else {
			buffer.WriteByte(b)
		}
	}
	return tab
}

func (inst *PomProperties) handleFragmentTag1(buffer *strings.Builder, dst map[string]string) {
	text := strings.TrimSpace(buffer.String())
	inst.prevValue = text
	dst[""] = text
}

func (inst *PomProperties) handleFragmentTag2(buffer *strings.Builder, dst map[string]string) {
	tag := strings.TrimSpace(buffer.String())
	if strings.HasPrefix(tag, "/") {
		n1 := inst.prevName
		n2 := "/" + n1
		// end
		if tag == n2 {
			dst[n1] = inst.prevValue
		}
	} else {
		// begin
		inst.prevName = tag
	}
}

////////////////////////////////////////////////////////////////////////////////
