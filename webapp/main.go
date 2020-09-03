package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	auth "github.com/abbot/go-http-auth"
	"github.com/cgxeiji/inventory"
	"github.com/markbates/pkger"
)

var templates *template.Template

func main() {
	port := flag.Int("p", 8080, "port to sever the inventory")
	path := flag.String("d", inventory.Path(), "path to inventory directory")
	credentials := flag.String("c", "", "filepath to credentials in realm [inventory] (file name should end in either .htdigest or .htpasswd)")
	flag.Parse()

	if *path != inventory.Path() {
		if _, err := os.Stat(*path); os.IsNotExist(err) {
			log.Fatalf("error with inventory path: %v", err)
		}
		inventory.CustomPath = *path
	}

	inventory.Items()

	f, err := os.OpenFile(filepath.Join(*path, "log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	templates, err = initTemplates(pkger.Include("/templates"))
	if err != nil {
		log.Fatalf("error reading templates: %v", err)
	}

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(pkger.Dir("/styles"))))
	http.Handle("/inventory/", http.StripPrefix("/inventory/", http.FileServer(http.Dir(inventory.Path()))))

	http.HandleFunc("/update", updateH)
	http.HandleFunc("/qr", qrH)
	http.HandleFunc("/location", locationH)
	http.HandleFunc("/add", checkAuth(*credentials, addH))
	http.HandleFunc("/", checkAuth(*credentials, rootH))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func initTemplates(dir string) (*template.Template, error) {
	t := template.New("")

	err := pkger.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		f, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		_, err = t.Parse(string(data))
		if err != nil {
			return err
		}

		return nil
	})

	return t, err
}

func checkAuth(path string, hf http.HandlerFunc) http.HandlerFunc {
	switch filepath.Ext(path) {
	case ".htdigest":
		a := auth.NewDigestAuthenticator("inventory", auth.HtdigestFileProvider(path))
		return a.JustCheck(hf)

	case ".htpasswd":
		a := auth.NewBasicAuthenticator("inventory", auth.HtpasswdFileProvider(path))
		return auth.JustCheck(a, hf)
	}

	return rootH
}
