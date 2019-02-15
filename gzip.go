package GzipMiddleWare

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type MiddleWare struct {
	Next  http.Handler
	Hosts []string
}

func (gm *MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Header.Get("Origin")

	if gm.Hosts != nil {
		for _, url := range gm.Hosts {
			if url == host {
				w.Header().Set("Access-Control-Allow-Origin", host)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
	}
	if gm.Next == nil {
		gm.Next = http.DefaultServeMux
	}

	encodings := r.Header.Get("Accept-Encoding")
	if !strings.Contains(encodings, "gzip") {
		gm.Next.ServeHTTP(w, r)
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
