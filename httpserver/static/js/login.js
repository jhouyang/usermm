function login() {
    var username = document.getElementById("username")
    var password = document.getElementById("password")

    if (username.value == "") {
        alert("请输入用户名")
    } else if (password.value == "") {
        alert("请输入密码")
    }
    if (!username.value.match(/^\S{2,20}$/)) {
        console.log("get focus")
        username.className = 'userRed';
        username.focus();
        return;
    } 


    if (password.value.length<1 || password.value.length>20) {
        console.log("passwd get focus")
        password.className = 'userRed';
        password.focus();
        return;
    } 
    var xhr = new XMLHttpRequest();
    xhr.open('post', 'http://192.168.33.10:8080/login')
    xhr.setRequestHeader("Content-type","application/x-www-form-urlencoded")
    xhr.send('username=' + username.value + "&passwd=" + $.md5(password.value))
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            console.log(xhr.responseText)
            var json = eval("("+xhr.responseText+")");
            console.log(json.code)
            console.log(json.msg)
            console.log(json.data)
            if (json.code == 0) {
                window.location.href = "http://192.168.33.10:8080/static/index.html?name=" + username.value
                //window.location.href = "http://192.168.33.10:8080/static/index.html"
                window.event.returnValue = false
            } else {
                alert("账号或密码错误。")
            }
        }
    }
}
