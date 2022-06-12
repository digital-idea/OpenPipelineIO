var webpage = "/endpoints" // webpage
var endpoint = "/api/endpoint" // restAPI Endpoint
var uxprefix = "endpoint-" // UX prefix string

function UxToObject(obj) {
    obj.dns = document.getElementById(uxprefix+'dns').value
    obj.endpoint = document.getElementById(uxprefix+'endpoint').value
    obj.description = document.getElementById(uxprefix+'description').value
    obj.method = string2array(document.getElementById(uxprefix+'method').value)
    obj.parameter = string2array(document.getElementById(uxprefix+'parameter').value)
    obj.cors = document.getElementById(uxprefix+'cors').value
	obj.type = document.getElementById(uxprefix+'type').value
    obj.storagetype = document.getElementById(uxprefix+'storagetype').value
	obj.db = document.getElementById(uxprefix+'db').value
	obj.isserverless = document.getElementById(uxprefix+'isserverless').checked
    obj.user = document.getElementById(uxprefix+'user').checked
    obj.developer = document.getElementById(uxprefix+'developer').checked
    obj.admin = document.getElementById(uxprefix+'admin').checked
    obj.security = document.getElementById(uxprefix+'security').checked
    obj.isasset = document.getElementById(uxprefix+'isasset').checked
	obj.ispatent = document.getElementById(uxprefix+'ispatent').checked
    obj.isupload = document.getElementById(uxprefix+'isupload').checked
	obj.token = document.getElementById(uxprefix+'token').value
	obj.tags = string2array(document.getElementById(uxprefix+'tags').value)
    obj.supportformat = string2array(document.getElementById(uxprefix+'supportformat').value)
	obj.curl = document.getElementById(uxprefix+'curl').value
    obj.category = document.getElementById(uxprefix+'category').value
    obj.partner = document.getElementById(uxprefix+'partner').value
    obj.progress = document.getElementById(uxprefix+'progress').value
    return obj
}

function ObjectToUx(obj) {
    document.getElementById(uxprefix+'id').value = obj.id
    document.getElementById(uxprefix+'dns').value = obj.dns
    document.getElementById(uxprefix+'endpoint').value = obj.endpoint
    document.getElementById(uxprefix+'description').value = obj.description
    document.getElementById(uxprefix+'method').value = obj.method.join(",")
    document.getElementById(uxprefix+'parameter').value = obj.parameter.join(",")
    document.getElementById(uxprefix+'cors').value = obj.cors
	document.getElementById(uxprefix+'type').value = obj.type
    document.getElementById(uxprefix+'storagetype').value = obj.storagetype
	document.getElementById(uxprefix+'db').value = obj.db
	document.getElementById(uxprefix+'isserverless').checked = obj.isserverless
    document.getElementById(uxprefix+'user').checked = obj.user
    document.getElementById(uxprefix+'developer').checked = obj.developer
    document.getElementById(uxprefix+'admin').checked = obj.admin
    document.getElementById(uxprefix+'security').checked = obj.security
    document.getElementById(uxprefix+'isasset').checked = obj.isasset
	document.getElementById(uxprefix+'ispatent').checked = obj.ispatent
    document.getElementById(uxprefix+'isupload').checked = obj.isupload
	document.getElementById(uxprefix+'token').value = obj.token
    document.getElementById(uxprefix+'tags').value = obj.tags.join(",")
    document.getElementById(uxprefix+'supportformat').value = obj.supportformat.join(",")
	document.getElementById(uxprefix+'curl').value = obj.curl
    document.getElementById(uxprefix+'category').value = obj.category
    document.getElementById(uxprefix+'partner').value = obj.partner
    document.getElementById(uxprefix+'progress').value = obj.progress
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