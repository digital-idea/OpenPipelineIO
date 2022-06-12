var webpage = "/pipelinesteps" // webpage
var endpoint = "/api/pipelinestep" // restAPI Endpoint
var uxprefix = "pipelinestep-" // UX prefix string

function UxToObject(obj) {
    obj.name = document.getElementById(uxprefix+'name').value
    obj.shortcode = document.getElementById(uxprefix+'shortcode').value
    obj.description = document.getElementById(uxprefix+'description').value
    // obj.textcolor = document.getElementById(uxprefix+'textcolor').value
    // obj.bgcolor = document.getElementById(uxprefix+'bgcolor').value
    // obj.bordercolor = document.getElementById(uxprefix+'bordercolor').value
    return obj
}

function ObjectToUx(obj) {
    document.getElementById(uxprefix+'id').value = obj.id
    document.getElementById(uxprefix+'name').value = obj.name
    document.getElementById(uxprefix+'shortcode').value = obj.shortcode
    document.getElementById(uxprefix+'description').value = obj.description
    // document.getElementById(uxprefix+'textcolor').value = obj.textcolor
    // document.getElementById(uxprefix+'bgcolor').value = obj.bgcolor
    // document.getElementById(uxprefix+'bordercolor').value = obj.bordercolor
}

function AddMode() {
    document.getElementById(uxprefix+'postbutton').hidden = false
    document.getElementById(uxprefix+'deletebutton').hidden = true
    document.getElementById(uxprefix+'putbutton').hidden = true
    InitModal()
}

function EditMode() {
    document.getElementById(uxprefix+'postbutton').hidden = true
    document.getElementById(uxprefix+'deletebutton').hidden = true
    document.getElementById(uxprefix+'putbutton').hidden = false
}

function DeleteMode() {
    document.getElementById(uxprefix+'postbutton').hidden = true
    document.getElementById(uxprefix+'deletebutton').hidden = false
    document.getElementById(uxprefix+'putbutton').hidden = true
}

function string2array(str) {
    let newArr = [];
    if (str === "") {
        return newArr
    }
    let arr = str.split(",");
    for (let i = 0; i < arr.length; i += 1) {
        newArr.push(arr[i].trim())
    }
    return newArr;
}

function InitModal() {
    let inputs = document.querySelectorAll("[id^='"+uxprefix+"']")
    for (let i = 0; i < inputs.length; i += 1) {
        if (inputs[i].type === "checkbox") {
            inputs[i].checked = false
        } else {
            inputs[i].value = ""
        }
    }
}

function SetModal(id) {
    EditMode()
    fetch(endpoint+'/'+id, {
        method: 'GET',
        headers: {"Authorization": "Basic "+ document.getElementById("token").value},
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {position:'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        ObjectToUx(obj)
    })
    .catch((err) => {
        console.log(err)
    });
}

function Post() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.name === "") {
        tata.error('Error', "Need name.",{position: 'tr',duration: 5000,onClose: null})
        return
    }
    fetch(endpoint, {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(obj),
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error',text,{position:'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        tata.success('Add', obj.name + "가 추가되었습니다.", {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function Put() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.name === "") {
        tata.error('Error',"Need name.",{position: 'tr',duration: 5000,onClose: null})
        return
    }
    fetch(endpoint+'/'+document.getElementById(uxprefix+'id').value, {
        method: 'PUT',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(obj),
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        tata.success('Edit', obj.name + "가 편집되었습니다.", {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function Delete() {
    fetch(endpoint+'/'+document.getElementById(uxprefix+'id').value, {
        method: 'DELETE',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {position: 'tr',duration: 5000,onClose: null})
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        tata.success('Delete', obj.name + "가 삭제되었습니다.", {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function tataLink() {
    window.location.replace(webpage)
}