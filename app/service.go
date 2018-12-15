// ----------
// Next steps:
//  - Implement websocket disconnects;
//  - Fix Logon button to accept enter key;
//  - Implement branch change on frontend --
//    - Ensure this functions with the websocket disconnect.

package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
    "github.com/gorilla/websocket"
    "path"
	"net/http"
    "time"
)

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Leaf)
var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func handleLogon(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	token := updateLogonToken(username, password)

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		fmt.Fprintf(writer, ErrorJSON)
        return
	}

    fmt.Fprintf(writer, string(tokenJSON))
}

func handleUser(writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Authentication")

    authentication := checkAuthentication(token, "")

    if authentication != "" {
        fmt.Fprintf(writer, authentication)
        return
    }
    if token == "" || len(token) != 32 {
        fmt.Fprintf(writer, ErrorHeaderToken)
        return
    }

	user := selectUser(token)

	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Fprintf(writer, ErrorJSON)
        return
	}

    fmt.Fprintf(writer, string(userJSON))
}

func handleBranch(writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("Authentication")
    requestUrl := request.URL.String()
	branchKey := path.Base(requestUrl)

    if len(token) != 32 || len(branchKey) != 16 {
        fmt.Fprintf(writer, ErrorHeaderToken)
        return
    }

    authentication := checkAuthentication(token, branchKey)

    if authentication != "" {
        fmt.Fprintf(writer, authentication)
        return
    }

	branch := selectBranch(token, branchKey)

	branchJSON, err := json.Marshal(branch)
	if err != nil {
		fmt.Fprintf(writer, ErrorJSON)
        return
	}

    fmt.Fprintf(writer, string(branchJSON))

    updateLastBranch(token, branchKey)
}

func handlePostLeaf(writer http.ResponseWriter, request *http.Request) {
    token := request.Header.Get("Authentication")
    branchKey := request.FormValue("branchKey")
    leafBody := request.FormValue("body")

    authentication := checkAuthentication(token, branchKey)

    if authentication != "" {
        fmt.Fprintf(writer, authentication)
        return
    }
    if token == "" || len(token) != 32 || branchKey == "" || leafBody == "" {
        fmt.Fprintf(writer, ErrorHeaderToken)
        return
    }

    datetime := time.Now().Format("2006-01-02 15:04")

    requestLeaf := Leaf{BranchKey: branchKey, Body: leafBody, Datetime: datetime}
	resultLeaf := insertLeaf(token, requestLeaf)

	resultLeafJSON, err := json.Marshal(resultLeaf)
	if err != nil {
		fmt.Fprintf(writer, ErrorJSON)
        return
	}

    fmt.Fprintf(writer, string(resultLeafJSON))

    handleLeaves(resultLeaf, branchKey)
}

func handleWebsocket(writer http.ResponseWriter, request *http.Request) {
    connection, err := upgrader.Upgrade(writer, request, nil)
    if err != nil {
		fmt.Fprintf(writer, ErrorWebsocket)
        return
	}
    //defer connection.Close()

    branchKey := request.URL.Query()["branch"]
    if branchKey == nil {
        fmt.Fprintf(writer, ErrorWSParam)
        return
    }

    token := request.URL.Query()["token"]
    if token == nil {
        fmt.Fprintf(writer, ErrorWSParam)
        return
    }

    authentication := checkAuthentication(token[0], branchKey[0])
    if authentication != "" {
        fmt.Fprintf(writer, authentication)
        return
    }

    clients[connection] = branchKey[0]

    fmt.Println("Websocket connection made: " + branchKey[0])

    fmt.Fprintf(writer, SuccessWebsocket)
}

func handleLeaves(leaf *Leaf, branchKey string) {
    for client, clientBranch := range clients {
        if clientBranch == branchKey {
            leafJSON, err := json.Marshal(leaf)
        	if err != nil {
            	fmt.Println(ErrorJSON)
                return
        	}

            err = client.WriteMessage(websocket.TextMessage, leafJSON)
            if err != nil {
    		    fmt.Println(ErrorWebsocket + ": %v", err)
                client.Close()
                delete(clients, client)
                return
    	    }
        }
    }
}
