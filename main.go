package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type file struct {
	Name     string
	FileType bool
	Size     int64
	Format   string
	Date     string
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	url := r.Form.Get("url")
	if url == "" {
		http.Error(w, "You are missing url parameter.", http.StatusBadRequest)
		return
	}

	resp, err := grab.Get("./downloaded", url)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Download saved to", resp.Filename)

	w.WriteHeader(http.StatusOK)
	jsonresp, _ := json.Marshal("your file downloaded succesfully!")
	w.Write(jsonresp)

}

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	var files []file
	err := filepath.Walk("./downloaded", func(path string, info os.FileInfo, err error) error {
		files = append(files, file{
			Name:     info.Name(),
			FileType: info.IsDir(),
			Size:     info.Size(),
			Format:   filepath.Ext(path),
			Date:     info.ModTime().String(),
		})
		return nil
	})
	if err != nil {
		fmt.Printf("%s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonresp, _ := json.Marshal(files)
	w.Write(jsonresp)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {

	param := mux.Vars(r)
	filename := param["filename"]

	if filename == "" {
		http.Error(w, "You are missing filename parameter.", http.StatusBadRequest)
		return
	}

	err := os.Remove("./downloaded/" + filename)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("File " + filename + " successfully deleted.")

	w.WriteHeader(http.StatusOK)
	resp := make(map[string]string)
	resp["message"] = "The contact has been deleted successfully!"
	jsonresp, _ := json.Marshal(resp)
	w.Write(jsonresp)
}

func Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/ls/", GetFilesHandler).Methods("GET")
	router.HandleFunc("/download/", DownloadFileHandler).Methods("GET")
	router.HandleFunc("/delete/{filename}", DeleteFileHandler).Methods("DELETE")
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	fmt.Println("server is runing on port 8089...")
	http.ListenAndServe(":8000", corsWrapper.Handler(router))
	return nil
}

func main() {
	err := Run()
	if err != nil {
		panic(err.Error())
	}
}
