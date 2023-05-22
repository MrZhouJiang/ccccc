function setCookie(cname, cvalue, exOur) {
    var d = new Date();
    d.setTime(d.getTime() + (exOur * 1 * 60 * 60 * 1000));
    var expires = "expires=" + d.toGMTString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";

}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function SetUser(userName, roleId) {
    setCookie("username", userName, 6)
    setCookie("role", roleId, 6)
}

function checkCookie() {
    //先向服务器请求是否登陆
    const user = getCookie("username");
    //  const user = $.cookie("username");
    if (user != "") {
        return true
    } else {
        return false
    }
}

function getUserRole() {
    let role = getCookie("role");
    return role
}

function getUserName() {
    let user = getCookie("username");
    return user


}


function logout() {

    setCookie("username", "", -1)
    setCookie("role", "", -1)
}