package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/cgxeiji/inventory"
	"github.com/skip2/go-qrcode"
)

func rootH(w http.ResponseWriter, r *http.Request) {
	items, err := inventory.SortedItems(inventory.ByInUseDate, true)
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	if err := templates.ExecuteTemplate(w, "inventory",
		&struct {
			Items []*inventory.Item
		}{
			Items: items,
		},
	); err != nil {
		log.Println("[ERR]", err)
		return
	}
}

func updateH(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	items, err := inventory.Items()
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	for _, item := range items {
		if item.ID == id {
			switch r.Method {
			case "POST":
				who := r.FormValue("who")
				item.Use(who)
				if item.InUse {
					log.Println("[USE]", item)
				} else {
					img, _, err := r.FormFile("image")
					if err != nil {
						log.Println("[ERR]", err)
						return
					}
					defer img.Close()

					if err := item.SetLocationPicture(img); err != nil {
						log.Println("[ERR]", err)
						return
					}
					log.Println("[RET]", item)
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)

			case "GET":
				if item.InUse {
					if err := templates.ExecuteTemplate(w, "return",
						&struct {
							Item *inventory.Item
						}{
							Item: item,
						},
					); err != nil {
						log.Println("[ERR]", err)
						return
					}
				} else {
					if err := templates.ExecuteTemplate(w, "use",
						&struct {
							Item *inventory.Item
						}{
							Item: item,
						},
					); err != nil {
						log.Println("[ERR]", err)
						return
					}
				}
			}
		}
	}
}

func addH(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		name := r.FormValue("name")
		item, err := inventory.Add(name)
		if err != nil {
			log.Println("[ERR]", err)
			return
		}

		img, _, err := r.FormFile("image")
		if err != nil {
			log.Println("[ERR]", err)
			return
		}
		defer img.Close()
		if err := item.SetPicture(img); err != nil {
			log.Println("[ERR]", err)
			return
		}

		log.Println("[ADD]", item)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	case "GET":
		if err := templates.ExecuteTemplate(w, "add", nil); err != nil {
			log.Println("[ERR]", err)
			return
		}
	}
}

func qrH(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	qr, err := qrcode.Encode(fmt.Sprintf("http://%s/update?id=%s", r.Host, id), qrcode.Medium, 256)
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	io.Copy(w, bytes.NewReader(qr))
}

func locationH(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	items, err := inventory.Items()
	if err != nil {
		log.Println("[ERR]", err)
		return
	}

	for _, item := range items {
		if item.ID == id {
			img, err := item.LocationPicture()
			if err != nil {
				log.Println("[ERR]", err)
				return
			}
			jpeg.Encode(w, img, nil)
			break
		}
	}
}
