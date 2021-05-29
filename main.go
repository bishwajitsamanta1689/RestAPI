package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const Port = ":8080"

func serveIndex(w http.ResponseWriter, r *http.Request) {
	//http.FileServer(http.Dir("./static"))
	http.ServeFile(w, r, "index.html")
	return
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File name:: %+v\n", handler.Filename)
	log.Printf("Uploaded File Size:: %+v\n", handler.Size)
	log.Printf("MIME Header:: %+v\n", handler.Header)

	// Creating Directory if doesn't exist
	_, err = os.Stat("Temp-Storage")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("Temp-Storage", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	// Write Temporary File on our Server
	tempFile, err := ioutil.TempFile("Temp-Storage", "upload-*.csv")
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/v1/upload", uploadFile)
	log.Printf("Application Started on %v", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}

func main() {
	setupRoutes()
}
