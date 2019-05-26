package GzipMiddleWare

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type MiddleWare struct {
	Next http.Handler
}

func (gm *MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	domain := strings.Join(r.Header["Origin"], " ")
	w.Header().Set("Access-Control-Allow-Origin", domain)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	encodings := r.Header.Get("Accept-Encoding")

	if gm.Next == nil {
		gm.Next = http.DefaultServeMux
	}

	if !strings.Contains(encodings, "gzip") {
		gm.Next.ServeHTTP(w, r)
		return
	}

	// Serve pre compressed files
	// in the /public folder

	if strings.Contains(r.URL.Path, "/public/javascript") || strings.Contains(r.URL.Path, "/public/css") {
		if strings.Contains(r.URL.Path, ".css") {
			w.Header().Add("Content-Type", "text/css")
		} else if strings.Contains(r.URL.Path, ".js") {
			w.Header().Add("Content-Type", "application/javascript")
		}

		if strings.ContainsAny(r.URL.Path, ".js | .css") {
			if strings.Contains(encodings, "br") {
				r.URL.Path = r.URL.Path + ".br"

				w.Header().Add("Content-Encoding", "br")

				gm.Next.ServeHTTP(w, r)

			} else if strings.Contains(encodings, "gzip") {
				r.URL.Path = r.URL.Path + ".gz"

				w.Header().Add("Content-Encoding", "gzip")
				w.Header().Add("Content-Type", "application/javascript")
				gm.Next.ServeHTTP(w, r)

			}
		}

		return
	}

	w.Header().Add("Content-Encoding", "gzip")
	gzipwriter := gzip.NewWriter(w)
	defer gzipwriter.Close()

	grw := GzipResponseWriter{
		ResponseWriter: w,
		Writer:         gzipwriter,
	}
	gm.Next.ServeHTTP(grw, r)
}

type GzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (grw GzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}
