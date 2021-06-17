package models

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path"
	"strings"
)

type SiteFile struct {
	FileName   string
	FileDir    string
	URLPath    string
	MimeType   string
	Data       []byte
	Rendered   []byte
	IsMarkdown bool
	Checksum   [32]byte
	Attrs      SiteFileAttr
}

type SiteFileAttr struct {
	Template    string                 `yaml:"template"`         /* use template */
	Render      bool                   `yaml:"render,omitempty"` /* render or plain */
	ExtraFields map[string]interface{} `yaml:",inline,omitempty"`
}

type SiteTemplate struct {
	FileName string
	FileDir  string
	Tpl      *template.Template
}

type SiteFiles map[string]*SiteFile

func (sf *SiteFiles) Keys() (k []string) {
	for f := range *sf {
		k = append(k, f)
	}
	return
}

func (sf SiteFiles) AddFile(s *SiteFile)    { sf[s.URLPath] = s }
func (sf SiteFiles) UpdateFile(s *SiteFile) { sf[s.URLPath] = s }
func (sf SiteFiles) RemoveFile(s *SiteFile) { delete(sf, s.URLPath) }

func (sf SiteFiles) FindFileByUrl(path string) (*SiteFile, bool) {
	s, ok := sf[path]
	return s, ok
}

type SitePlace struct {
	Level    int          `json:"level"`
	Parent   *SitePlace   `json:"parent,omitempty"`
	Self     *SiteFile    `json:"self"`
	Children []*SitePlace `json:"children,omitempty"`
}
type SiteMap struct {
	Places *SitePlace `json:"places"`
}

func (sf SiteFiles) GenerateSiteMap() (sm SiteMap) {
	/* head of site */
	sm = SiteMap{
		Places: &SitePlace{
			Level: 0,
			Self: &SiteFile{
				URLPath: "/",
			},
		},
	}

	addSiteMapChildren(sm.Places, &sf, 0)
	// fmt.Println(sm.Places.Self.URLPath)
	// sm.Places.Describe()
	return
}

func (sm *SiteMap) ToJSON() ([]byte, error) {
	tmp := *sm.Places // JSON doesn't support cycles so Parents must go...
	tmp.Abandon()
	return json.Marshal(tmp)
}

func (sp *SitePlace) Abandon() {
	sp.Parent = nil
	for _, c := range sp.Children {
		c.Abandon()
	}
	return
}

func (sp *SitePlace) Describe() {
	indent(sp.Level)
	for _, c := range sp.Children {
		indent(sp.Level + 1)
		fmt.Printf("%s\n", c.Self.URLPath)
		c.Describe()
	}
}

func addSiteMapChildren(sp *SitePlace, sf *SiteFiles, l int) {

	for _, p := range *sf {

		dir := path.Dir(p.URLPath)
		dir_level := getLevel(dir)

		place := &SitePlace{
			Level:  l,
			Parent: sp,
			Self:   p,
		}
		if l == dir_level && dir == place.Parent.Self.URLPath { // child files
			if _, exists := place.Parent.childExists(place.Self); !exists {
				place.Parent.Children = append(place.Parent.Children, place)
			}
		} else if (l+1) == dir_level && path.Dir(dir) == place.Parent.Self.URLPath { // child directories
			place.Level += 1
			place.Self = &SiteFile{
				URLPath: dir,
			}

			if _, exists := place.Parent.childExists(place.Self); !exists {
				place.Parent.Children = append(place.Parent.Children, place)

				addSiteMapChildren(place, sf, l+1)
			}
		}
	}
}

func (sp *SitePlace) childExists(s *SiteFile) (*SitePlace, bool) {
	for _, c := range sp.Children {
		if c.Self.URLPath == s.URLPath {
			return c, true
		}
	}
	return nil, false
}

func getLevel(path string) int {
	if path == "/" {
		return 0
	}
	return len(strings.Split(path, "/")) - 1
}

func indent(levels int) {
	for i := 0; i < levels; i++ {
		fmt.Printf("  ") // 2 spaces
	}
}
