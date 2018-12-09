package main

import (
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	http.HandleFunc("/logon/", handleLogon)
	http.HandleFunc("/u/", handleUser)
	http.HandleFunc("/b/", handleBranch)
	http.HandleFunc("/l/", handlePostLeaf)
    http.HandleFunc("/ws/", handleWebsocket)

    http.Handle("/", http.FileServer(http.Dir("../public")))

	http.ListenAndServe(":8080", nil)
}
