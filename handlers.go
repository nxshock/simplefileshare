package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	type FileInfo struct {
		Name string
		Size string
		Date string
	}

	var data []FileInfo

	err := filepath.Walk(config.StoragePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data = append(data, FileInfo{filepath.Base(path), sizeToApproxHuman(info.Size()), info.ModTime().Format("02.01.2006 15:04")})

		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.htm", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join(config.StoragePath, header.Filename)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		http.Error(w, "файл с таким именем уже существует", http.StatusBadRequest)
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.FormValue("filename"))

	if filename == "" {
		http.Error(w, `"filename" field can not be empty`, http.StatusBadRequest)
		return
	}

	f, err := os.Open(filepath.Join(config.StoragePath, filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	fileStat, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	w.Header().Set("Accept-Ranges", "none")
	w.Header().Set("Content-Length", strconv.Itoa(int(fileStat.Size())))

	io.CopyBuffer(w, f, make([]byte, 4096))
}
