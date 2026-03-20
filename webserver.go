package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func webserver() {
	http.Handle("/build/", http.FileServer(http.FS(staticFiles)))
	ext := make(map[string]string)
	ext["html"] = "text/html"
	ext["json"] = "application/json"
	ext["css"] = "text/css"
	ext["js"] = "application/javascript"
	ext["gif"] = "image/gif"
	ext["svg"] = "image/svg+xml"
	ext["png"] = "image/png"
	ext["jpg"] = "image/jpeg"
	ext["jpeg"] = "image/jpeg"
	ext["ico"] = "image/x-icon"
	ext["woff"] = "font/woff"
	ext["woff2"] = "font/woff2"
	ext["ttf"] = "font/ttf"
	ext["eot"] = "application/vnd.ms-fontobject"

	// Получаем содержимое корневой директории
	contents, err := readDirRecursively("build")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range contents {
		ff := item.Path
		http.HandleFunc(strings.ReplaceAll(ff, "build", ""), func(w http.ResponseWriter, r *http.Request) {
			file, err := staticFiles.ReadFile(ff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", ext[getFileExtension(ff)])
			_, err = w.Write(file)
			if err != nil {
				log.Println(err)
			}
		})

	}

	//
	//
	//
	//
	//
	//

	//  █████╗ ██████╗ ██╗
	// ██╔══██╗██╔══██╗██║
	// ███████║██████╔╝██║
	// ██╔══██║██╔═══╝ ██║
	// ██║  ██║██║     ██║
	// ╚═╝  ╚═╝╚═╝     ╚═╝

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if debug {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			reqHdrs := r.Header.Get("Access-Control-Request-Headers")
			if reqHdrs != "" {
				w.Header().Set("Access-Control-Allow-Headers", reqHdrs)
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			}
			// w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "600")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.Method != http.MethodPost {
			w.Header().Set("Content-type", "application/json")
			http.Error(w, "{\"error\":\"Only POST data!\"}", http.StatusMethodNotAllowed)
			return
		} else {
			w.Header().Set("Content-type", "application/json")
			out := map[string]any{}
			in := map[string]any{}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error read body request", http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(body, &in)
			if err != nil {
				fmt.Println("Error translate income JSON:", err)
				return
			}

			if _, ok := in["op"]; ok { // все остальные случаи только с проверкой подлинности пользователя
				out["data"] = api(in, r, w)
			} else {
				out["error"] = "No operation for API"
			}

			json, err := json.Marshal(out)
			if err != nil {
				slog("Fail transfer to JSON", "err")
				return
			}
			w.Write([]byte(json))
		}
	})

	/*	███╗   ███╗ █████╗ ██╗███╗   ██╗
		████╗ ████║██╔══██╗██║████╗  ██║
		██╔████╔██║███████║██║██╔██╗ ██║
		██║╚██╔╝██║██╔══██║██║██║╚██╗██║
		██║ ╚═╝ ██║██║  ██║██║██║ ╚████║
		╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝ */

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		file := []byte{}
		ip := getClientIP(r)

		if checkAccessList(ip) {
			file, err = staticFiles.ReadFile("build/index.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			msg := fmt.Sprintf("Access denied. IP: %s (%s)", ip, geoIP(ip))
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write(file)
		if err != nil {
			log.Println(err)
		}
	})

	// Запуск сервера на порту
	webport := getSetting("port")
	slog(APPNAME + " ver." + VERSION + " was run on the " + webport)
	log.Fatal(http.ListenAndServe(webport, nil))
}
