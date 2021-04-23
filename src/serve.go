package main

import (
	"crypto/tls"
	"html/template"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/lubiedo/yav/src/models"
	"github.com/lubiedo/yav/src/utils"
)

// Main rendering process
// This is the only handler for HTTP(S) requests. Will get all client requests
// and render the content using templates. It will also serve files as well.
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
				file, _ = ReturnSiteFile(newlocation)
			} else {
				http.ServeFile(w, req, urlpath)
				return
			}
		} else if utils.FileExist(urlpath) || utils.FileExist(FromTemplateExt(urlpath)) {
			/* revert to markdown ext to use real path if necessary */
			if urlpath[len(urlpath)-len(templateext):] == templateext {
				urlpath = FromTemplateExt(urlpath)
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
			NotFound(w, req)
			return
		}
	} else {
		/* it exists in filesystem? */
		path := GetSiteFilePath(file)

		if !utils.FileExist(path) {
			RemoveSiteFile(file)
			NotFound(w, req)
			return
		}
	}

	/* default http headers... */
	H := w.Header()
	H.Set("Content-Type", file.MimeType)
	if len(headers) > 0 { /* ... and extra headers */
		headers.AddToHttp(&H)
	}

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

	tplvariables := GenerateTemplateVars(req, file)
	tpls.ExecuteTemplate(w, file.Attrs.Template, tplvariables)
}

// Listen and serve via HTTP(S)
func Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Render)

	fulladdr := config.Addr + ":" + config.Port
	listener, err := net.Listen("tcp", fulladdr)
	if err != nil {
		config.Log.Fatal("%s", err)
	}

	server := &http.Server{
		Addr:         fulladdr,
		Handler:      mux,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	if !config.UseHTTPS {
		err := server.Serve(listener)
		if err != nil {
			config.Log.Fatal("%s", err)
		}
	} else {
		server.ConnState = CheckHTTPSConnState
		server.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
		}
		err := server.ServeTLS(listener, config.CertPath, config.KeyPath)
		if err != nil {
			config.Log.Fatal("%s", err)
		}
	}
}

// Finds the SiteFile struct depending on file URL path
func ReturnSiteFile(path string) (models.SiteFile, bool) {
	for _, f := range files {
		if path == f.URLPath {
			return f, true
		}
	}
	return models.SiteFile{}, false
}

// 404 not found HTTP response
func NotFound(w http.ResponseWriter, r *http.Request) {
	if config.TplErroPage == "" {
		http.NotFound(w, r)
		return
	}
	ServeErrorTemplate(w, r, http.StatusNotFound, "404 not found")
}

// 500 server error HTTP response
func ServerError(w http.ResponseWriter, r *http.Request) {
	if config.TplErroPage == "" {
		http.Error(w, "500 server error", http.StatusInternalServerError)
		return
	}
	ServeErrorTemplate(w, r, http.StatusInternalServerError, "500 server error")
}

// Send `Location` header to redirect user
func Location(url string, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", url)
	http.Error(w, "301 moved permanently", http.StatusMovedPermanently)
}

// Serve the client the error in template
func ServeErrorTemplate(w http.ResponseWriter, r *http.Request, code int, msg string) {
	vars := GenerateTemplateVars(r, models.SiteFile{Rendered: []byte(msg)})

	/* extra error vars */
	vars["code"] = code

	w.WriteHeader(code)
	tpls.ExecuteTemplate(w, config.TplErroPage, vars)
}

// Returns variables to be passed to the templates
func GenerateTemplateVars(r *http.Request, f models.SiteFile) (vars map[string]interface{}) {
	vars = make(map[string]interface{})
	vars["url"] = r.URL
	vars["query"] = r.URL.Query()
	vars["template"] = f.Attrs.Template
	vars["content"] = template.HTML(string(f.Rendered))
	if len(tplvars) > 0 {
		for k, v := range tplvars {
			vars[k] = v
		}
	}
	for k, v := range f.Attrs.ExtraFields {
		/* extra fields will override previously defined variables */
		vars[k] = v
	}
	return
}

/*
this function exists to forcefully redirect users from HTTP to HTTPS
*/
func CheckHTTPSConnState(c net.Conn, s http.ConnState) {
	if s != http.StateNew {
		return
	}
	tlsConn, _ := c.(*tls.Conn)
	hs := tlsConn.Handshake()

	if hs == nil || hs == io.EOF {
		return
	}
	rh, ok := hs.(tls.RecordHeaderError)
	if !ok {
		return
	}

	if tlsConn.ConnectionState().CipherSuite == 0 && tlsRecordHeaderLooksLikeHTTP(rh.RecordHeader) {
		if config.Verbose {
			config.Log.Info("User \"%s\" redirected from HTTP to HTTPS", c.RemoteAddr())
		}
		_, _ = rh.Conn.Write([]byte("HTTP/1.0 301 Moved Permanently\r\n" +
			"Location: https://" + config.Addr + ":" + config.Port + "\r\n\r\n" +
			"Redirecting you to HTTPS\n"))
		rh.Conn.Close()
		return
	}
}

// taken from go/src/net/http/server.go
func tlsRecordHeaderLooksLikeHTTP(hdr [5]byte) bool {
	switch string(hdr[:]) {
	case "GET /", "HEAD ", "POST ", "PUT /", "OPTIO":
		return true
	}
	return false
}
