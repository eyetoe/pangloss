// Package to handle and route other web request as a server
package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

func main() {
	http.HandleFunc("/", index)
	//http.HandleFunc("/submit", submit)
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	// may uncomment this line to launch the branches page of the jps_liveops repo
	//openURL("https://github.com/openwager/jps_liveops/branches")
	openURL("http://localhost:8888/")
	log.Println("Listening on port :8888")
	http.ListenAndServe(":8888", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	// template stuff
	renderPage := func() {
		tmpl := make(map[string]*template.Template)
		tmpl["index.html"] = template.Must(template.ParseFiles("templates/index.html"))
		tmpl["index.html"].Execute(w, nil)
	}

	// render page if request is 'GET'
	if r.Method == "GET" {
		log.Println("######################################## GET ##")
		renderPage()
	}

	if r.Method == "POST" {
		log.Println("######################################## POST ##")
		renderPage()
		fmt.Fprint(w, "<pre>")
		fw := flushWriter{w: w}
		if f, ok := w.(http.Flusher); ok {
			fw.f = f
		}
		cmd := exec.Command("/Users/john/bin/streamer.sh")
		//cmd := exec.Command("ansible-playbook", "-e", "'target_env=ow-shared-stage'", "/Users/john/git/ansible/template-cassandra.yml")
		log.Println("Command to backend: ", cmd)

		cmd.Stdout = &fw
		cmd.Stderr = &fw
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, "</pre>")
		log.Println("############################# BACKEND FINISH ##")
	}
}

func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}
