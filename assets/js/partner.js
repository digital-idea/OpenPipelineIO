
var struct = "partner"
var prefix = struct + "-" // ux id prefix string
var endpoint = "/api/" + struct // restAPI endpoint
var page = struct + "s" // webpage

function UxToObject(data) {
    data.name = document.getElementById(prefix + 'name').value
    data.symbol = document.getElementById(prefix + 'symbol').value
    data.domain = document.getElementById(prefix + 'domain').value
    data.size = document.getElementById(prefix + 'size').value
    data.homepage = document.getElementById(prefix + 'homepage').value
    data.address = document.getElementById(prefix + 'address').value
    data.phone = document.getElementById(prefix + 'phone').value
    data.email = document.getElementById(prefix + 'email').value
    data.location = document.getElementById(prefix + 'location').value
    data.timezone = document.getElementById(prefix + 'timezone').value
    data.description = document.getElementById(prefix + 'description').value
    data.businessregistrationnumber = document.getElementById(prefix + 'businessregistrationnumber').value
    data.manager = document.getElementById(prefix + 'manager').value
    data.managerphone = document.getElementById(prefix + 'managerphone').value
    data.manageremail = document.getElementById(prefix + 'manageremail').value
    data.ftp = document.getElementById(prefix + 'ftp').value
    data.ftpid = document.getElementById(prefix + 'ftpid').value
    data.ftppw = document.getElementById(prefix + 'ftppw').value
    data.opentime = document.getElementById(prefix + "opentime").options[document.getElementById(prefix + "opentime").selectedIndex].value,    
    data.closedtime = document.getElementById(prefix + "closedtime").options[document.getElementById(prefix + "closedtime").selectedIndex].value,
    data.paymentdate = document.getElementById(prefix + 'paymentdate').value
    data.bank = document.getElementById(prefix + 'bank').value
    data.bankaccount = document.getElementById(prefix + 'bankaccount').value
    data.monetaryunit = document.getElementById(prefix + "monetaryunit").options[document.getElementById(prefix + "monetaryunit").selectedIndex].value,
    data.projecthistory = document.getElementById(prefix + 'projecthistory').value
    data.reputation = document.getElementById(prefix + 'reputation').value
    data.status = document.getElementById(prefix + 'status').value
    data.codename = document.getElementById(prefix + 'codename').value
    data.companytype = document.getElementById(prefix + "companytype").options[document.getElementById(prefix + "companytype").selectedIndex].value,
    data.contactpoint = document.getElementById(prefix + 'contactpoint').value
    data.pmsurl = document.getElementById(prefix + 'pmsurl').value
    data.isabroad = document.getElementById(prefix + 'isabroad').checked
    data.isclient = document.getElementById(prefix + 'isclient').checked
    data.tags = string2array(document.getElementById(prefix + 'tags').value)
    return data
}

function ObjectToUx(data) {
    document.getElementById(prefix + 'id').value = data.id
    document.getElementById(prefix + 'name').value = data.name
    document.getElementById(prefix + 'symbol').value = data.symbol
    document.getElementById(prefix + 'domain').value = data.domain
    document.getElementById(prefix + 'size').value = data.size
    document.getElementById(prefix + 'homepage').value = data.homepage
    document.getElementById(prefix + 'address').value = data.address
    document.getElementById(prefix + 'phone').value = data.phone
    document.getElementById(prefix + 'email').value = data.email
    document.getElementById(prefix + 'location').value = data.location
    document.getElementById(prefix + 'timezone').value = data.timezone
    document.getElementById(prefix + 'description').value = data.description
    document.getElementById(prefix + 'businessregistrationnumber').value = data.businessregistrationnumber
    document.getElementById(prefix + 'manager').value = data.manager
    document.getElementById(prefix + 'managerphone').value = data.managerphone
    document.getElementById(prefix + 'manageremail').value = data.manageremail
    document.getElementById(prefix + 'ftp').value = data.ftp
    document.getElementById(prefix + 'ftpid').value = data.ftpid
    document.getElementById(prefix + 'ftppw').value = data.ftppw
    document.getElementById(prefix + "opentime").value = data.opentime
    document.getElementById(prefix + "closedtime").value = data.closedtime
    document.getElementById(prefix + 'paymentdate').value = data.paymentdate
    document.getElementById(prefix + 'bank').value = data.bank
    document.getElementById(prefix + 'bankaccount').value = data.bankaccount
    document.getElementById(prefix + "monetaryunit").value = data.monetaryunit
    document.getElementById(prefix + 'projecthistory').value = data.projecthistory
    document.getElementById(prefix + 'reputation').value = data.reputation
    document.getElementById(prefix + 'status').value = data.status
    document.getElementById(prefix + 'codename').value = data.codename
    document.getElementById(prefix + "companytype").value = data.companytype
    document.getElementById(prefix + 'contactpoint').value = data.contactpoint
    document.getElementById(prefix + 'pmsurl').value = data.pmsurl
    document.getElementById(prefix + 'isabroad').checked = data.isabroad
    document.getElementById(prefix + 'isclient').checked = data.isclient
    document.getElementById(prefix + 'tags').value = data.tags.join(",")
}

function AddMode() {
    document.getElementById(prefix + 'postbutton').hidden = false
    document.getElementById(prefix + 'deletebutton').hidden = true
    document.getElementById(prefix + 'putbutton').hidden = true
    InitModal()
}

function EditMode() {
    document.getElementById(prefix + 'postbutton').hidden = true
    document.getElementById(prefix + 'deletebutton').hidden = true
    document.getElementById(prefix + 'putbutton').hidden = false
}

function DeleteMode() {
    document.getElementById(prefix + 'postbutton').hidden = true
    document.getElementById(prefix + 'deletebutton').hidden = false
    document.getElementById(prefix + 'putbutton').hidden = true
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
    // id가 "partner-" 로 시작하는 input > text, input > checkbox, select를 초기화 한다.
    let inputs = document.querySelectorAll("[id^='" + prefix + "']")
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
            onClick: goPage,
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
    fetch(endpoint+'/'+document.getElementById(prefix+'id').value, {
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
            onClick: goPage,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function Delete() {
    fetch(endpoint+'/'+document.getElementById(prefix+'id').value, {
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
            onClick: goPage,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function goPage() {
    window.location.replace(page)
}