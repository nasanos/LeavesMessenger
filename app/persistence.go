// ----------
// To do:
//  - Abstract the logic to select a user ID from a token into its own function.

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

func selectAuthentication(token string) string {
	var userId string = ""

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		return ErrorDBConnection
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT id
		FROM users
		WHERE user_token=?
		;
	`)
	if err != nil {
		return ErrorDBOther
	}

	row := stmt.QueryRow(token)
	err = row.Scan(&userId)
	if err != nil {
		return ErrorDBQuery
	}

	if userId == "" {
		return ErrorAuthentication
	} else {
		return ""
	}
}

func selectAuthenticationForBranch(token string, branchKey string) string {
	var userId string = ""
	var branchId string = ""

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		return ErrorDBConnection
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT id
		FROM users
		WHERE user_token=?
		;
	`)
	if err != nil {
		return ErrorDBOther
	}

	row := stmt.QueryRow(token)
	err = row.Scan(&userId)
	if err != nil {
		return ErrorDBQuery
	}

	stmt, err = db.Prepare(`
        SELECT b.id
        FROM branches b
         INNER JOIN user_branches ub
            ON b.id = ub.branch_id
        WHERE ub.user_id=?
        AND b.branch=?
        ;
    `)
	if err != nil {
		return ErrorDBOther
	}

	row = stmt.QueryRow(userId, branchKey)
	err = row.Scan(&branchId)
	if err != nil {
		fmt.Println(err)
		return ErrorDBQuery
	}

	if userId == "" || branchId == "" {
		return ErrorAuthentication
	} else {
		return ""
	}
}

func updateLogonToken(username string, password string) *LogonToken {
	var logonToken LogonToken = LogonToken{Token: "", Error: ErrorNotice{Description: ""}}

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		logonToken.Error = ErrorNotice{Description: ErrorDBConnection}
		return &logonToken
	}
	defer db.Close()

	var token string = generateHex(32)

	stmt, err := db.Prepare(`
		UPDATE users
		SET user_token=?
		WHERE username=?
			AND password=?
		;
	`)
	if err != nil {
		logonToken.Error = ErrorNotice{Description: ErrorDBQuery}
		return &logonToken
	}

	_, err = stmt.Exec(token, username, password)
	if err != nil {
		logonToken.Error = ErrorNotice{Description: ErrorDBUpdate}
		return &logonToken
	}

	fmt.Println("Provided authentication token to user " + username + ": " + token + ".")

	logonToken.Token = token
	return &logonToken
}

func selectUser(token string) *User {
	var user User = User{LastBranchKey: "", Error: ErrorNotice{Description: ""}}

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		user.Error = ErrorNotice{Description: ErrorDBConnection}
		return &user
	}
	defer db.Close()

	// Get the key for the user's last visited branch.
	stmt, err := db.Prepare(`
		SELECT b.branch
		FROM branches b
			INNER JOIN users u
				ON b.id=u.last_branch_id
		WHERE u.user_token=?
		;
	`)
	if err != nil {
		user.Error = ErrorNotice{Description: ErrorDBOther}
		return &user
	}

	row := stmt.QueryRow(token)
	err = row.Scan(&user.LastBranchKey)
	if err != nil {
		user.Error = ErrorNotice{Description: ErrorDBQuery}
		return &user
	}

	// Get the keys and names for the branches the user is subscribed to.
	stmt, err = db.Prepare(`
		SELECT b.branch, b.name
		FROM branches b
			INNER JOIN user_branches ub
				ON ub.branch_id=b.id
			INNER JOIN users u
				ON u.id=ub.user_id
		WHERE u.user_token=?
		;
	`)
	if err != nil {
		user.Error = ErrorNotice{Description: ErrorDBOther}
		return &user
	}

	rows, err := stmt.Query(token)
	if err != nil {
		user.Error = ErrorNotice{Description: ErrorDBQuery}
		return &user
	}
	defer rows.Close()

	for rows.Next() {
		var branch string
		var name string

		err := rows.Scan(&branch, &name)
		if err != nil {
			user.Error = ErrorNotice{Description: ErrorDBQuery}
			return &user
		}

		user.SubscribedBranchKeys = append(user.SubscribedBranchKeys, branch)
		user.SubscribedBranchNames = append(user.SubscribedBranchNames, name)
	}

	fmt.Println("Queried user for: " + token + ".")

	return &user
}

func selectBranch(token string, branchKey string) *Branch {
	var branch Branch = Branch{Key: branchKey, Name: "", Leaves: []Leaf{}, Error: ErrorNotice{Description: ""}}

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		branch.Error = ErrorNotice{Description: ErrorDBConnection}
		return &branch
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT b.name, l.id, l.body, u.username, l.datetime
		FROM leaves l
			INNER JOIN branches b
				ON b.id=l.branch_id
			INNER JOIN users u
				ON u.id=l.user_id
		WHERE b.branch=?
        ORDER BY l.datetime ASC
        ;
	`)
	if err != nil {
		branch.Error = ErrorNotice{Description: ErrorDBQuery}
		return &branch
	}

	rows, err := stmt.Query(branchKey)
	if err != nil {
		branch.Error = ErrorNotice{Description: ErrorDBQuery}
		return &branch
	}
	defer rows.Close()

	for rows.Next() {
		var leaf Leaf

		err = rows.Scan(&branch.Name, &leaf.Id, &leaf.Body, &leaf.Username, &leaf.Datetime)
		if err != nil {
			branch.Error = ErrorNotice{Description: ErrorDBQuery}
			return &branch
		}

		branch.Leaves = append(branch.Leaves, leaf)
	}

	fmt.Println("Queried branch: " + branchKey + ".")

	return &branch
}

func updateLastBranch(token string, branchKey string) {
	var userId int
    var branchId int

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		fmt.Println(ErrorDBConnection)
		return
	}
	defer db.Close()

    // Get the user ID for the token.
	stmt, err := db.Prepare(`
		SELECT id
		FROM users
		WHERE user_token = ?
		;
	`)
	if err != nil {
		fmt.Println(ErrorDBOther)
		return
	}

	row := stmt.QueryRow(token)
	err = row.Scan(&userId)
	if err != nil {
		fmt.Println(ErrorDBQuery)
		return
	}

    // Get the branch ID for the branch key.
    stmt, err = db.Prepare(`
		SELECT b.id
		FROM branch b
		WHERE b.branch = ?
		;
	`)
	if err != nil {
		fmt.Println(ErrorDBOther)
		return
	}

	row = stmt.QueryRow(branchKey)
	err = row.Scan(&branchId)
	if err != nil {
		fmt.Println(ErrorDBQuery)
		return
	}

    // Update the user's last branch.
	stmt, err = db.Prepare(`
		UPDATE users
        SET last_branch_id = ?
        WHERE id = ?
		;
	`)
	if err != nil {
		fmt.Println(ErrorDBOther)
		return
	}

	_, err = stmt.Exec(branchId, userId)
	if err != nil {
		fmt.Println(ErrorDBUpdate)
		return
	}

	fmt.Println("Updated last visited branch for " + strconv.Itoa(userId) + " to " + branchKey + ".")
}

func insertLeaf(token string, leaf Leaf) *Leaf {
	var userId int

	db, err := sql.Open("sqlite3", "../database/main.db")
	if err != nil {
		leaf.Error = ErrorNotice{Description: ErrorDBConnection}
		return &leaf
	}
	defer db.Close()

	stmt, err := db.Prepare(`
		SELECT id, username
		FROM users
		WHERE user_token=?
		;
	`)
	if err != nil {
		leaf.Error = ErrorNotice{Description: ErrorDBOther}
		return &leaf
	}

	row := stmt.QueryRow(token)
	err = row.Scan(&userId, &leaf.Username)
	if err != nil {
		leaf.Error = ErrorNotice{Description: ErrorDBQuery}
		return &leaf
	}

	stmt2, err := db.Prepare(`
		INSERT INTO leaves
            (branch_id, body, user_id, datetime)
        SELECT b.id, ?, ?, ?
        FROM branches b
        WHERE b.branch = ?
		;
	`)
	if err != nil {
		leaf.Error = ErrorNotice{Description: ErrorDBOther}
		return &leaf
	}

	_, err = stmt2.Exec(leaf.Body, userId, leaf.Datetime, leaf.BranchKey)
	if err != nil {
		leaf.Error = ErrorNotice{Description: ErrorDBUpdate}
		return &leaf
	}

	fmt.Println("Inserted leaf for: " + leaf.BranchKey + ".")

	return &leaf
}
