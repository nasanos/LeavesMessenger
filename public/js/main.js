// To do:
//   Change `alert` messages to modal messages.
//   Consider including the served error messages in the modal.

var currentToken = "";
var currentBranchKey = "";
var leafBodyCount = 0;

getToken = (username, password) => {
    $.ajax({
        type: "POST",
        dataType: "json",
        crossDomain: true,
        url: "/logon/",
        data: {"username": username, "password": password},
        success: (token) => {
            if (token.Error.Description != "") {
                alert("Authentication failed.");
                checkToken();
            } else {
                currentToken = token.Token;
                checkToken();
                getUser(token.Token);
            }
        },
        error: (err) => {
            alert("Authentication failed.");
            checkToken();
        }
    })
}

getUser = (token) => {
    $.ajax({
        beforeSend: (request) => {
            request.setRequestHeader("Authentication", token);
        },
        dataType: "json",
        url: "/u/",
        success: (user) => {
            if (user.Error.Description != "") {
                alert("Unable to get user information.");
                currentToken = "";
                checkToken();
            } else {
                getBranch(token, user.LastBranchKey);
            }
        },
        error: (err) => {
            alert("Unable to get user information.");
            currentToken = "";
            checkToken();
        }
    });
}

getBranch = (token, branchKey) => {
    $.ajax({
        beforeSend: (request) => {
            request.setRequestHeader("Authentication", token);
        },
        dataType: "json",
        url: "/b/" + branchKey,
        success: (branch) => {
            if (branch.Error.Description != "") {
                alert("Unable to get branch infomration.");
                currentToken = "";
                checkToken();
            } else {
                currentBranchKey = branchKey;

                console.log(currentBranchKey)

                $("#branchName").text(branch.Name);

                $("#branchBody").html("");
                for (i=0; i<branch.Leaves.length; i++) {
                    appendLeaf(branch.Leaves[i]);
                }

                websocketConnection(branchKey, token);
            }
        },
        error: (err) => {
            alert("Unable to get branch infomration.");
            currentToken = "";
            checkToken();
        }
    });
}

appendLeaf = (leaf) => {
    leafBody = leaf.Body.replace(/(?:\r\n|\r|\n)/g, "<br/>");

    $("#branchBody").append(`
        <div class="row justify-content-center">
            <div id="leafBody-` + leafBodyCount + `" class="col-10 col-md-8">` + leafBody + `</div>
        </div>
        <div class="row justify-content-end">
            <div id="leafHeader-` + leafBodyCount + `" class="col-10 col-md-6 col-lg-5f">
                <p class="leafHeader">` + leaf.Username + ` | ` + leaf.Datetime + `</p>
            </div>
        </div>
        <hr class="leafDivider"/>
    `);
}

checkToken = () => {
    if (currentToken == "") {
        $("#loginModal").modal("show")
    } else {
        $("#loginModal").modal("hide")
    }
}

websocketConnection = (branchKey, token) => {
    var socket = new WebSocket("ws://" + window.location.host + "/ws/?branch=" + branchKey + "&token=" + token);
    socket.onopen = () => console.log("Websocket connection established.");
    socket.onclose = () => console.log("Websocket connection lost.");
    socket.onmessage = (response) => {
        leaf = JSON.parse(response.data)
        appendLeaf(leaf);
    }
}

$("#loginSubmit").click(() => {
    getToken($("#loginUsernameField").val(), $("#loginPasswordField").val());
});

$("#loginUsernameField").on("keypress", (e) => {
    if (e.keyCode == 13) {
        getToken($("#loginUsernameField").val(), $("#loginPasswordField").val());
    }
});

$("#loginPasswordField").on("keypress", (e) => {
    if (e.keyCode == 13) {
        getToken($("#loginUsernameField").val(), $("#loginPasswordField").val());
    }
});

$("#leafPostSubmit").click(() => {
    leafPostBody = $("#leafPostBodyField").val();
    postData = {"body": leafPostBody, "branchKey": currentBranchKey};

    $.ajax({
        type: "POST",
        dataType: "json",
        beforeSend: (request) => {
            request.setRequestHeader("Authentication", currentToken);
        },
        data: postData,
        url: "/l/" + currentBranchKey,
        success: (response) => {
            $("#leafPostBodyField").val("");

            if (response.Error.Description != "") {
                alert("An error occurred while posting the message.");
            }
        },
        error: (err) => {
            alert("An error occurred while posting the message.");
        }
    });
});

$(document).ready(function() {
    $(".modal").modal({
        backdrop: "static",
        keyboard: false
    });

    $(".modal").on("shown.bs.modal", () => {checkToken();});
    $(".modal").on("hidden.bs.modal", () => {checkToken();});

    // For testing:
    getToken("BabyBird", "gumbo");
});
