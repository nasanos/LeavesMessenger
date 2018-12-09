package main

const hexCharSet = "1234567890ABCDEF"

const (
    ErrorAuthentication = "Unable to authenticate user to branch."
    ErrorDBConnection = "Unable to connect to the database."
    ErrorDBQuery = "An error was encountered in the query. The credentials may be incorrect or invalid."
    ErrorDBUpdate = "The attempt to post the record update failed. Ensure the values were correct and valid."
    ErrorDBOther = "An error occurred while preparing a query."
    ErrorHeaderToken = "An error has been found in the request."
    ErrorJSON = "An error occurred in composing the JSON response."
    ErrorWebsocket = "An error has occurred in attempting to establish a websocket connection."
    ErrorWSParam = "An error was found in the connection parameters."
    SuccessWebsocket = "Successfully established weobsocket connection."
)

type ErrorNotice struct {
	Description    string
}

type LogonToken struct {
	Token      string
    Error      ErrorNotice
}

type User struct {
	LastBranchKey			string
	SubscribedBranchKeys	[]string
	SubscribedBranchNames	[]string
    Error                   ErrorNotice
}

type Leaf struct {
    Id              int
    BranchKey       string
    Body            string
    Username        string
    Datetime        string
    Error           ErrorNotice
}

type Branch struct {
    Key     	string
	Name    	string
	Leaves      []Leaf
    Error       ErrorNotice
}

type Response struct {
	Body	  string
    Error     ErrorNotice
}
