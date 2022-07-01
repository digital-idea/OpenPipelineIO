var webpage = "/projectsforpartner" // webpage
var endpoint = "/api/projectforpartner" // restAPI Endpoint
var uxprefix = "projectforpartner-" // UX prefix string

function UxToObject(obj) {
    obj.projectname = document.getElementById(uxprefix+'projectname').value
    obj.partnername = document.getElementById(uxprefix+'partnername').value
    obj.domain = document.getElementById(uxprefix+'domain').value
    obj.size = document.getElementById(uxprefix+'size').value
    obj.homepage = document.getElementById(uxprefix+'homepage').value
    obj.address = document.getElementById(uxprefix+'address').value
    obj.phone = document.getElementById(uxprefix+'phone').value
    obj.email = document.getElementById(uxprefix+'email').value
    obj.location = document.getElementById(uxprefix+'location').value
    obj.timezone = document.getElementById(uxprefix+'timezone').value
    obj.description = document.getElementById(uxprefix+'description').value
    obj.businessregistrationnumber = document.getElementById(uxprefix+'businessregistrationnumber').value
    obj.manager = document.getElementById(uxprefix+'manager').value
    obj.managerphone = document.getElementById(uxprefix+'managerphone').value
    obj.manageremail = document.getElementById(uxprefix+'manageremail').value
    obj.ftp = document.getElementById(uxprefix+'ftp').value
    obj.ftpid = document.getElementById(uxprefix+'ftpid').value
    obj.ftppw = document.getElementById(uxprefix+'ftppw').value
    obj.opentime = document.getElementById(uxprefix+"opentime").options[document.getElementById(uxprefix+"opentime").selectedIndex].value,    
    obj.closedtime = document.getElementById(uxprefix+"closedtime").options[document.getElementById(uxprefix+"closedtime").selectedIndex].value,
    obj.paymentdate = document.getElementById(uxprefix+'paymentdate').value
    obj.bank = document.getElementById(uxprefix+'bank').value
    obj.bankaccount = document.getElementById(uxprefix+'bankaccount').value
    obj.monetaryunit = document.getElementById(uxprefix+"monetaryunit").options[document.getElementById(uxprefix+"monetaryunit").selectedIndex].value,
    obj.projecthistory = document.getElementById(uxprefix+'projecthistory').value
    obj.reputation = document.getElementById(uxprefix+'reputation').value
    obj.status = document.getElementById(uxprefix+'status').value
    obj.codename = document.getElementById(uxprefix+'codename').value
    obj.companytype = document.getElementById(uxprefix+"companytype").options[document.getElementById(uxprefix+"companytype").selectedIndex].value,
    obj.contactpoint = document.getElementById(uxprefix+'contactpoint').value
    obj.pmsurl = document.getElementById(uxprefix+'pmsurl').value
    obj.isabroad = document.getElementById(uxprefix+'isabroad').checked
    obj.isclient = document.getElementById(uxprefix+'isclient').checked
    obj.tags = string2array(document.getElementById(uxprefix+'tags').value)
    return obj
}

function ObjectToUx(obj) {
    document.getElementById(uxprefix+'id').value = obj.id
    document.getElementById(uxprefix+'projectname').value = obj.projectname
    document.getElementById(uxprefix+'partnername').value = obj.partnername
    document.getElementById(uxprefix+'domain').value = obj.domain
    document.getElementById(uxprefix+'size').value = obj.size
    document.getElementById(uxprefix+'homepage').value = obj.homepage
    document.getElementById(uxprefix+'address').value = obj.address
    document.getElementById(uxprefix+'phone').value = obj.phone
    document.getElementById(uxprefix+'email').value = obj.email
    document.getElementById(uxprefix+'location').value = obj.location
    document.getElementById(uxprefix+'timezone').value = obj.timezone
    document.getElementById(uxprefix+'description').value = obj.description
    document.getElementById(uxprefix+'businessregistrationnumber').value = obj.businessregistrationnumber
    document.getElementById(uxprefix+'manager').value = obj.manager
    document.getElementById(uxprefix+'managerphone').value = obj.managerphone
    document.getElementById(uxprefix+'manageremail').value = obj.manageremail
    document.getElementById(uxprefix+'ftp').value = obj.ftp
    document.getElementById(uxprefix+'ftpid').value = obj.ftpid
    document.getElementById(uxprefix+'ftppw').value = obj.ftppw
    document.getElementById(uxprefix+"opentime").value = obj.opentime
    document.getElementById(uxprefix+"closedtime").value = obj.closedtime
    document.getElementById(uxprefix+'paymentdate').value = obj.paymentdate
    document.getElementById(uxprefix+'bank').value = obj.bank
    document.getElementById(uxprefix+'bankaccount').value = obj.bankaccount
    document.getElementById(uxprefix+"monetaryunit").value = obj.monetaryunit
    document.getElementById(uxprefix+'projecthistory').value = obj.projecthistory
    document.getElementById(uxprefix+'reputation').value = obj.reputation
    document.getElementById(uxprefix+'status').value = obj.status
    document.getElementById(uxprefix+'codename').value = obj.codename
    document.getElementById(uxprefix+"companytype").value = obj.companytype
    document.getElementById(uxprefix+'contactpoint').value = obj.contactpoint
    document.getElementById(uxprefix+'pmsurl').value = obj.pmsurl
    document.getElementById(uxprefix+'isabroad').checked = obj.isabroad
    document.getElementById(uxprefix+'isclient').checked = obj.isclient
    document.getElementById(uxprefix+'tags').value = obj.tags.join(",")
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