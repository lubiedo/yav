package models

import "html/template"

type Sites []SiteFile

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
