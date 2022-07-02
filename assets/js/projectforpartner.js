var webpage = "/projectsforpartner" // webpage
var endpoint = "/api/projectforpartner" // restAPI Endpoint
var uxprefix = "projectforpartner-" // UX prefix string

function UxToObject(obj) {
    obj.projectname = document.getElementById(uxprefix+'projectname').options[document.getElementById(uxprefix+"projectname").selectedIndex].value
    obj.partnername = document.getElementById(uxprefix+'partnername').options[document.getElementById(uxprefix+"partnername").selectedIndex].value
    obj.rnr = document.getElementById(uxprefix+'rnr').value
    obj.projecttype = document.getElementById(uxprefix+'projecttype').options[document.getElementById(uxprefix+"projecttype").selectedIndex].value
    obj.language = document.getElementById(uxprefix+'language').value
    obj.partnerinternalmanager = document.getElementById(uxprefix+'partnerinternalmanager').value
    obj.manageremail = document.getElementById(uxprefix+'manageremail').value
    obj.messenger = document.getElementById(uxprefix+'messenger').value
    obj.messengerid = document.getElementById(uxprefix+'messengerid').value
    obj.startdate = document.getElementById(uxprefix+'startdate').value
    obj.enddate = document.getElementById(uxprefix+'enddate').value
    obj.leftovershot = parseInt(document.getElementById(uxprefix+'leftovershot').value)
    obj.totalshot = parseInt(document.getElementById(uxprefix+'totalshot').value)
    obj.projectbudget = parseFloat(document.getElementById(uxprefix+'projectbudget').value)
    obj.percentageoftotalbudget = parseFloat(document.getElementById(uxprefix+'percentageoftotalbudget').value)
    obj.priceperframe = parseFloat(document.getElementById(uxprefix+'priceperframe').value)
    obj.paymentdateforvender = parseInt(document.getElementById(uxprefix+'paymentdateforvender').value)
    obj.paymentdateforclient = parseInt(document.getElementById(uxprefix+'paymentdateforclient').value)
    obj.pricepershot = parseFloat(document.getElementById(uxprefix+'pricepershot').value)
    obj.paymentcycle = document.getElementById(uxprefix+'paymentcycle').value
    obj.manday = parseInt(document.getElementById(uxprefix+"manday").value)
    obj.description = document.getElementById(uxprefix+"description").value
    return obj
}

function ObjectToUx(obj) {
    document.getElementById(uxprefix+'id').value = obj.id
    document.getElementById(uxprefix+'projectname').value = obj.projectname
    document.getElementById(uxprefix+'partnername').value = obj.partnername
    document.getElementById(uxprefix+'rnr').value = obj.rnr
    document.getElementById(uxprefix+'projecttype').value = obj.projecttype
    document.getElementById(uxprefix+'language').value = obj.language
    document.getElementById(uxprefix+'partnerinternalmanager').value = obj.partnerinternalmanager
    document.getElementById(uxprefix+'manageremail').value = obj.manageremail
    document.getElementById(uxprefix+'messenger').value = obj.messenger
    document.getElementById(uxprefix+'messengerid').value = obj.messengerid
    document.getElementById(uxprefix+'startdate').value = obj.startdate
    document.getElementById(uxprefix+'enddate').value = obj.enddate
    document.getElementById(uxprefix+'leftovershot').value = obj.leftovershot
    document.getElementById(uxprefix+'totalshot').value = obj.totalshot
    document.getElementById(uxprefix+'projectbudget').value = obj.projectbudget
    document.getElementById(uxprefix+'percentageoftotalbudget').value = obj.percentageoftotalbudget
    document.getElementById(uxprefix+'priceperframe').value = obj.priceperframe
    document.getElementById(uxprefix+'paymentdateforvender').value = obj.paymentdateforvender
    document.getElementById(uxprefix+'paymentdateforclient').value = obj.paymentdateforclient
    document.getElementById(uxprefix+'pricepershot').value = obj.pricepershot
    document.getElementById(uxprefix+'paymentcycle').value = obj.paymentcycle
    document.getElementById(uxprefix+'manday').value = obj.manday
    document.getElementById(uxprefix+'description').value = obj.description
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
        tata.success('Add', `${obj.projectname}에 ${obj.partnername} 값이 추가되었습니다.`, {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function Put() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.projectname === "") {
        tata.error('Error',"Need ProjectName.",{position: 'tr',duration: 5000,onClose: null})
        return
    }
    if (obj.partnername === "") {
        tata.error('Error',"Need PartnerName.",{position: 'tr',duration: 5000,onClose: null})
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
        tata.success('Edit', `${obj.projectname}에 ${obj.partnername} 값이 수정되었습니다.`, {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
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
        tata.success('Delete', `${obj.projectname}에 ${obj.partnername} 값이 삭제되었습니다.`, {position: 'tr',duration: 5000,onClick: tataLink,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function tataLink() {
    window.location.replace(webpage)
}