function id(id) {
    return document.getElementById(id);
}

function showAlert(message) {
    id("serverMessage").innerText = message;
    id("coolAlert").style.display = "block";
    setTimeout( function() {
        id("coolAlertInner").className = "coolAlertInner";
    }, 20);
}

function closeAlert() {
    id("coolAlert").style.display = "none";
    id("coolAlertInner").className = "coolAlertInner coolAlertInnerHidden"
}

function startForwarding() {
    send("", "action?id=startForwarding");
    return false;
}

function send(request, url) {
    let ajax = new XMLHttpRequest();
    ajax.open("POST", url);
    ajax.responseType = "text";
    ajax.onload = function () {
        showAlert(ajax.response.substr(0, 50));
    }
    ajax.send(request);
}