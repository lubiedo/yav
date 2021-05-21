package models

import (
	"html/template"
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

func (sf SiteFiles) AddFile(s *SiteFile)    { sf[s.URLPath] = s }
func (sf SiteFiles) UpdateFile(s *SiteFile) { sf[s.URLPath] = s }
func (sf SiteFiles) RemoveFile(s *SiteFile) { delete(sf, s.URLPath) }

func (sf SiteFiles) FindFileByUrl(path string) (*SiteFile, bool) {
	s, ok := sf[path]
	return s, ok
}
