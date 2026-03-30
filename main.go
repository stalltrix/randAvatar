package main

import (
    "log"
	"net/http"
	"github.com/stalltrix/randAvatar/strx_avatar"
	"os"
    "path"
    "strings"
)


func handler(w http.ResponseWriter, r *http.Request) {
    seed := r.URL.Query().Get("seed")
    if seed == "" {
        seed = "default"
    }
	
    svg,err := strx_avatar.NewAvatar(seed)
	if err!=nil{
		http.Error(w,err.Error(),500)
	}

    w.Header().Set("Content-Type", "image/svg+xml")
    w.Write([]byte(svg))
}

func handler_static(w http.ResponseWriter, r *http.Request) {
    name := strings.TrimPrefix(r.URL.Path, "/avatar/")
    if name == "" {
        http.NotFound(w, r)
        return
    }
    ext := path.Ext(name)
    seed := strings.TrimSuffix(name, ext)

    if seed == "" {
        seed = "default"
    }

    svg,err := strx_avatar.NewAvatar(seed)
	if err!=nil{
		http.Error(w,err.Error(),500)
	}

    w.Header().Set("Content-Type", "image/svg+xml")
    w.Write([]byte(svg))
}

func main() {
	argc:=len(os.Args)
	if argc <=1 {
		log.Println("randAvatar [Listen]")
		return
	}
    http.HandleFunc("/avatar", handler)
	http.HandleFunc("/avatar/", handler_static)
	log.Println("Listen on",os.Args[1])
    http.ListenAndServe(os.Args[1], nil)
}