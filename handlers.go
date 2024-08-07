package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.FileServer(http.FS(stripSiteFS)).ServeHTTP(w, r)
		return
	}

	type FileInfo struct {
		Ext  string
		Name string
		Size string
		Date string
	}

	var data struct {
		Files           []FileInfo
		StorageDuration uint
	}

	data.StorageDuration = config.RemoveFilePeriod

	err := filepath.Walk(config.StoragePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileInfo := FileInfo{
			Ext:  nvl(strings.ToLower(strings.TrimPrefix(filepath.Ext(path), ".")), "dat"),
			Name: filepath.Base(path),
			Size: sizeToApproxHuman(info.Size()),
			Date: info.ModTime().Format("02.01.2006")}

		data.Files = append(data.Files, fileInfo)

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

	var errors []string
	for _, header := range r.MultipartForm.File["file"] {
		file, err := header.Open()
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		defer file.Close()

		filePath := filepath.Join(config.StoragePath, header.Filename)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("файл с именем %s уже существует", header.Filename))
			continue
		}

		f, err := os.Create(filePath)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
	}

	if len(errors) > 0 {
		http.Error(w, strings.Join(errors, "\n"), http.StatusBadRequest)
	}
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.FormValue("filename"))

	if filename == "" {
		http.Error(w, `"filename" field can't be empty`, http.StatusBadRequest)
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

	_, _ = io.CopyBuffer(w, f, make([]byte, 4096))
}

func HandleStream(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.FormValue("filename"))

	if filename == "" {
		http.Error(w, `"filename" field can't be empty`, http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, filepath.Join(config.StoragePath, filename))
}

func HandleIcon(w http.ResponseWriter, r *http.Request) {
	ext := r.FormValue("ext")

	if ext == "" {
		http.Error(w, `"ext" field can't be empty`, http.StatusBadRequest)
		return
	}

	f, err := iconsFS.Open(filepath.ToSlash(filepath.Join("icons", ext+".svg")))
	if err != nil {
		f, _ = iconsFS.Open(filepath.ToSlash(filepath.Join("icons", "bin.svg")))
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=31557600")

	_, _ = io.Copy(w, f)
}
