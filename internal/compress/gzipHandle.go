package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandleCompress(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") != "application/json" && r.Header.Get("Content-type") != "text/html" {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(GzipWriter{ResponseWriter: w, Writer: gz}, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type GzipReader struct {
	io.ReadCloser
}

func (gr *GzipReader) Read(p []byte) (int, error) {
	return gr.ReadCloser.Read(p)
}

func GzipHandleDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer gr.Close()

			r.Body = &GzipReader{ReadCloser: gr}
		}

		next.ServeHTTP(w, r)
	})
}

func GzipDecompress(req *http.Request) (*http.Request, error) {
	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
		gr, err := gzip.NewReader(req.Body)
		if err != nil {
			return req, err
		}
		defer gr.Close()

		req.Body = &GzipReader{ReadCloser: gr}
	}

	return req, nil
}
