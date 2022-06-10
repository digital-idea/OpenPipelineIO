var webpage = "/partners" // webpage
var endpoint = "/api/partner" // restAPI Endpoint
var uxprefix = "partner-" // UX prefix string

function UxToObject(data) {
    data.name = document.getElementById(uxprefix + 'name').value
    data.symbol = document.getElementById(uxprefix + 'symbol').value
    data.domain = document.getElementById(uxprefix + 'domain').value
    data.size = document.getElementById(uxprefix + 'size').value
    data.homepage = document.getElementById(uxprefix + 'homepage').value
    data.address = document.getElementById(uxprefix + 'address').value
    data.phone = document.getElementById(uxprefix + 'phone').value
    data.email = document.getElementById(uxprefix + 'email').value
    data.location = document.getElementById(uxprefix + 'location').value
    data.timezone = document.getElementById(uxprefix + 'timezone').value
    data.description = document.getElementById(uxprefix + 'description').value
    data.businessregistrationnumber = document.getElementById(uxprefix + 'businessregistrationnumber').value
    data.manager = document.getElementById(uxprefix + 'manager').value
    data.managerphone = document.getElementById(uxprefix + 'managerphone').value
    data.manageremail = document.getElementById(uxprefix + 'manageremail').value
    data.ftp = document.getElementById(uxprefix + 'ftp').value
    data.ftpid = document.getElementById(uxprefix + 'ftpid').value
    data.ftppw = document.getElementById(uxprefix + 'ftppw').value
    data.opentime = document.getElementById(uxprefix + "opentime").options[document.getElementById(uxprefix + "opentime").selectedIndex].value,    
    data.closedtime = document.getElementById(uxprefix + "closedtime").options[document.getElementById(uxprefix + "closedtime").selectedIndex].value,
    data.paymentdate = document.getElementById(uxprefix + 'paymentdate').value
    data.bank = document.getElementById(uxprefix + 'bank').value
    data.bankaccount = document.getElementById(uxprefix + 'bankaccount').value
    data.monetaryunit = document.getElementById(uxprefix + "monetaryunit").options[document.getElementById(uxprefix + "monetaryunit").selectedIndex].value,
    data.projecthistory = document.getElementById(uxprefix + 'projecthistory').value
    data.reputation = document.getElementById(uxprefix + 'reputation').value
    data.status = document.getElementById(uxprefix + 'status').value
    data.codename = document.getElementById(uxprefix + 'codename').value
    data.companytype = document.getElementById(uxprefix + "companytype").options[document.getElementById(uxprefix + "companytype").selectedIndex].value,
    data.contactpoint = document.getElementById(uxprefix + 'contactpoint').value
    data.pmsurl = document.getElementById(uxprefix + 'pmsurl').value
    data.isabroad = document.getElementById(uxprefix + 'isabroad').checked
    data.isclient = document.getElementById(uxprefix + 'isclient').checked
    data.tags = string2array(document.getElementById(uxprefix + 'tags').value)
    return data
}

function ObjectToUx(data) {
    document.getElementById(uxprefix + 'id').value = data.id
    document.getElementById(uxprefix + 'name').value = data.name
    document.getElementById(uxprefix + 'symbol').value = data.symbol
    document.getElementById(uxprefix + 'domain').value = data.domain
    document.getElementById(uxprefix + 'size').value = data.size
    document.getElementById(uxprefix + 'homepage').value = data.homepage
    document.getElementById(uxprefix + 'address').value = data.address
    document.getElementById(uxprefix + 'phone').value = data.phone
    document.getElementById(uxprefix + 'email').value = data.email
    document.getElementById(uxprefix + 'location').value = data.location
    document.getElementById(uxprefix + 'timezone').value = data.timezone
    document.getElementById(uxprefix + 'description').value = data.description
    document.getElementById(uxprefix + 'businessregistrationnumber').value = data.businessregistrationnumber
    document.getElementById(uxprefix + 'manager').value = data.manager
    document.getElementById(uxprefix + 'managerphone').value = data.managerphone
    document.getElementById(uxprefix + 'manageremail').value = data.manageremail
    document.getElementById(uxprefix + 'ftp').value = data.ftp
    document.getElementById(uxprefix + 'ftpid').value = data.ftpid
    document.getElementById(uxprefix + 'ftppw').value = data.ftppw
    document.getElementById(uxprefix + "opentime").value = data.opentime
    document.getElementById(uxprefix + "closedtime").value = data.closedtime
    document.getElementById(uxprefix + 'paymentdate').value = data.paymentdate
    document.getElementById(uxprefix + 'bank').value = data.bank
    document.getElementById(uxprefix + 'bankaccount').value = data.bankaccount
    document.getElementById(uxprefix + "monetaryunit").value = data.monetaryunit
    document.getElementById(uxprefix + 'projecthistory').value = data.projecthistory
    document.getElementById(uxprefix + 'reputation').value = data.reputation
    document.getElementById(uxprefix + 'status').value = data.status
    document.getElementById(uxprefix + 'codename').value = data.codename
    document.getElementById(uxprefix + "companytype").value = data.companytype
    document.getElementById(uxprefix + 'contactpoint').value = data.contactpoint
    document.getElementById(uxprefix + 'pmsurl').value = data.pmsurl
    document.getElementById(uxprefix + 'isabroad').checked = data.isabroad
    document.getElementById(uxprefix + 'isclient').checked = data.isclient
    document.getElementById(uxprefix + 'tags').value = data.tags.join(",")
}

function AddMode() {
    document.getElementById(uxprefix + 'postbutton').hidden = false
    document.getElementById(uxprefix + 'deletebutton').hidden = true
    document.getElementById(uxprefix + 'putbutton').hidden = true
    InitModal()
}

function EditMode() {
    document.getElementById(uxprefix + 'postbutton').hidden = true
    document.getElementById(uxprefix + 'deletebutton').hidden = true
    document.getElementById(uxprefix + 'putbutton').hidden = false
}

function DeleteMode() {
    document.getElementById(uxprefix + 'postbutton').hidden = true
    document.getElementById(uxprefix + 'deletebutton').hidden = false
    document.getElementById(uxprefix + 'putbutton').hidden = true
}

function string2array(str) {
    var newArr = [];
    if (str === "") {
        return newArr
    }
    let arr = str.split(",");
    for (var i = 0; i < arr.length; i += 1) {
        newArr.push(arr[i].trim())
    }
    return newArr;
}

function InitModal() {
    let inputs = document.querySelectorAll("[id^='" + uxprefix + "']")
    for (var i = 0; i < inputs.length; i += 1) {
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
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                tata.error('Error', text, {
                    position: 'tr',
                    duration: 5000,
                    onClose: null,
                })
                return
            });
        }
        return response.json()
    })
    .then((data) => {
        ObjectToUx(data)
    })
    .catch((err) => {
        console.log(err)
    });
}

function Post() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.name === "") {
        tata.error('Error', "Need name.", {
            position: 'tr',
            duration: 5000,
            onClose: null,
        })
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
                tata.error('Error', text, {
                    position: 'tr',
                    duration: 5000,
                    onClose: null,
                })
                return
            });
        }
        return response.json()
    })
    .then((data) => {
        tata.success('Add', data.name + "가 추가되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClick: tataLink,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function Put() {
    let obj = new Object()
    obj = UxToObject(obj)
    if (obj.name === "") {
        tata.error('Error', "Need name.", {
            position: 'tr',
            duration: 5000,
            onClose: null,
        })
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
                tata.error('Error', text, {
                    position: 'tr',
                    duration: 5000,
                    onClose: null,
                })
                return
            });
        }
        return response.json()
    })
    .then((data) => {
        tata.success('Edit', data.name + "가 편집되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClick: tataLink,
            onClose: null,
        })
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
                tata.error('Error', text, {
                    position: 'tr',
                    duration: 5000,
                    onClose: null,
                })
                return
            });
        }
        return response.json()
    })
    .then((data) => {
        tata.success('Delete', data.name + "가 삭제되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClick: tataLink,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function tataLink() {
    window.location.replace(webpage)
}