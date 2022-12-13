package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	__port = 8025
)

var (
	port = fmt.Sprintf(":%d", __port)
)

type ImageFunc func(w http.ResponseWriter, r *http.Request)
type ImageHandler struct {
	*http.ServeMux
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		ServeMux: http.NewServeMux(),
	}
}

func ServeImage(imagefile string) ImageFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join("assets/", imagefile)
		res, err := os.ReadFile(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

func CreateImageFileMap(directoryname string) map[string]ImageFunc {
	filemap := make(map[string]ImageFunc)
	items, err := os.ReadDir(directoryname)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}
	for _, item := range items {
		log.Printf("%s \n", item)
		urlPath := "/" + item.Name()
		filemap[urlPath] = ServeImage(item.Name())
	}
	return filemap
}

func main() {
	imageMap := CreateImageFileMap("assets/")
	imageHandle := NewImageHandler()

	for k, v := range imageMap {
		log.Printf("%+v \n %+v \n", k, v)
		imageHandle.HandleFunc(k, v)
	}
	server := &http.Server{
		Addr:    port,
		Handler: imageHandle,
	}
	go func() {
		log.Printf("We are live at port:%d \n", __port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("There was an error listening to the server:%s \n",
				err.Error())
		}
	}()
	log.Println("image server closed.")
}
