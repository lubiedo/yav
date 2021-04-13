package main

import (
	"crypto/tls"
	"html/template"
	"net/http"
	"time"

	"github.com/lubiedo/yav/src/models"
	"github.com/lubiedo/yav/src/utils"
)

func Render(w http.ResponseWriter, req *http.Request) {
	var (
		newfile models.SiteFile
		err     error
	)

	config.Log.Access(req)
	file, found := ReturnSiteFile(req.URL.Path)

	if !found {
		/*
			TODO:
			 * Fix file extensions (markdown -> template) shown in the file list
			   when ServeFile
		*/
		urlpath := sitedir + req.URL.Path

		/* are we missing a new file or dir? */
		if utils.FileExist(urlpath) && utils.FileIsDir(urlpath) { /* serve a dir */
			if utils.FileExist(urlpath + "/" + "index" + markdownext) {
				/* redirect to index if exists */
				newlocation := req.URL.Path

				if newlocation[len(newlocation)-1] != '/' {
					newlocation += "/"
				}
				newlocation += "index" + templateext
				Location(newlocation, w, req)
				return
			} else {
				http.ServeFile(w, req, urlpath)
				return
			}
		} else if utils.FileExist(urlpath) || utils.FileExist(RevertTemplateExt(urlpath)) {
			/* revert to markdown ext to use real path if necessary */
			if urlpath[len(urlpath)-len(templateext):] == templateext {
				urlpath = RevertTemplateExt(urlpath)
			}

			newfile, err = ProcessSiteFile(urlpath)
			if err != nil {
				config.Log.Error("Error processing new file \"%s\"",
					urlpath)
				ServerError(w, req)
				return
			}

			file = newfile
		} else {
			http.NotFound(w, req)
			return
		}
	} else {
		/* it exists in filesystem? */
		path := GetSiteFilePath(file)

		if !utils.FileExist(path) {
			RemoveSiteFile(file)
			http.NotFound(w, req)
			return
		}
	}

	w.Header().Add("Content-Type", file.MimeType)

	filepath := GetSiteFilePath(file)

	if !file.IsMarkdown {
		http.ServeFile(w, req, filepath)
		return
	}

	if file.Checksum != utils.FileChecksum(filepath) {
		config.Log.Info("File changed! Updating content of \"%s\" in cache",
			filepath)

		newfile, err = UpdateSiteFile(file)
		if err != nil {
			config.Log.Error("Error updating content of \"%s\" in cache",
				filepath)
			ServerError(w, req)
			return
		}

		file = newfile
	}

	tplvariables := make(map[string]interface{})
	tplvariables["content"] = template.HTML(string(file.Rendered))
	for k, v := range file.Attrs.ExtraFields {
		tplvariables[k] = v.(string)
	}

	tpls.ExecuteTemplate(w, file.Attrs.Template, tplvariables)
}

func Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Render)

	server := &http.Server{
		Addr:         config.Addr + ":" + config.Port,
		Handler:      mux,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	if !config.UseHTTPS {
		err := server.ListenAndServe()
		if err != nil {
			config.Log.Fatal("%s", err)
		}
	} else {
		server.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		}
		err := server.ListenAndServeTLS(config.CertPath, config.KeyPath)
		if err != nil {
			config.Log.Fatal("%s", err)
		}
	}
}

func ServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "500 server error", http.StatusInternalServerError)
}
func Location(url string, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", url)
	http.Error(w, "301 moved permanently", http.StatusMovedPermanently)
}

func ReturnSiteFile(path string) (models.SiteFile, bool) {
	if path == "/" { /* redirect to index */
		path += "index" + templateext
	}

	for _, f := range files {
		if path == f.URLPath {
			return f, true
		}
	}
	return models.SiteFile{}, false
}
