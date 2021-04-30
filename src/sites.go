package main

import "github.com/lubiedo/yav/src/models"

type Sites map[string]*models.SiteFile

func (sites Sites) AddSite(s *models.SiteFile) {
	sites[s.URLPath] = s
}

func (sites Sites) FindSiteByUrl(path string) (*models.SiteFile, bool) {
	s, ok := sites[path]
	return s, ok
}

func (sites Sites) UpdateSite(s *models.SiteFile) {
	sites[s.URLPath] = s
}

func (sites Sites) RemoveSite(s *models.SiteFile) {
	delete(sites, s.URLPath)
}
