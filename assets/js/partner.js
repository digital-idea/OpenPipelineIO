
function UxToPartnerObject(data) {
    data.name = document.getElementById('partner-name').value
    data.symbol = document.getElementById('partner-symbol').value
    data.domain = document.getElementById('partner-domain').value
    data.size = document.getElementById('partner-size').value
    data.homepage = document.getElementById('partner-homepage').value
    data.address = document.getElementById('partner-address').value
    data.phone = document.getElementById('partner-phone').value
    data.email = document.getElementById('partner-email').value
    data.location = document.getElementById('partner-location').value
    data.timezone = document.getElementById('partner-timezone').value
    data.description = document.getElementById('partner-description').value
    data.businessregistrationnumber = document.getElementById('partner-businessregistrationnumber').value
    data.manager = document.getElementById('partner-manager').value
    data.managerphone = document.getElementById('partner-managerphone').value
    data.manageremail = document.getElementById('partner-manageremail').value
    data.ftp = document.getElementById('partner-ftp').value
    data.ftpid = document.getElementById('partner-ftpid').value
    data.ftppw = document.getElementById('partner-ftppw').value
    data.opentime = document.getElementById("partner-opentime").options[document.getElementById("partner-opentime").selectedIndex].value,    
    data.closedtime = document.getElementById("partner-closedtime").options[document.getElementById("partner-closedtime").selectedIndex].value,
    data.paymentdate = document.getElementById('partner-paymentdate').value
    data.bank = document.getElementById('partner-bank').value
    data.bankaccount = document.getElementById('partner-bankaccount').value
    data.monetaryunit = document.getElementById("partner-monetaryunit").options[document.getElementById("partner-monetaryunit").selectedIndex].value,
    data.projecthistory = document.getElementById('partner-projecthistory').value
    data.reputation = document.getElementById('partner-reputation').value
    data.status = document.getElementById('partner-status').value
    data.codename = document.getElementById('partner-codename').value
    data.companytype = document.getElementById("partner-companytype").options[document.getElementById("partner-companytype").selectedIndex].value,
    data.contactpoint = document.getElementById('partner-contactpoint').value
    data.pmsurl = document.getElementById('partner-pmsurl').value
    data.isabroad = document.getElementById('partner-isabroad').checked
    data.isclient = document.getElementById('partner-isclient').checked
    data.tags = string2array(document.getElementById('partner-tags').value)
    return data
}

function partnerObjectToUx(data) {
    document.getElementById('partner-id').value = data.id
    document.getElementById('partner-name').value = data.name
    document.getElementById('partner-symbol').value = data.symbol
    document.getElementById('partner-domain').value = data.domain
    document.getElementById('partner-size').value = data.size
    document.getElementById('partner-homepage').value = data.homepage
    document.getElementById('partner-address').value = data.address
    document.getElementById('partner-phone').value = data.phone
    document.getElementById('partner-email').value = data.email
    document.getElementById('partner-location').value = data.location
    document.getElementById('partner-timezone').value = data.timezone
    document.getElementById('partner-description').value = data.description
    document.getElementById('partner-businessregistrationnumber').value = data.businessregistrationnumber
    document.getElementById('partner-manager').value = data.manager
    document.getElementById('partner-managerphone').value = data.managerphone
    document.getElementById('partner-manageremail').value = data.manageremail
    document.getElementById('partner-ftp').value = data.ftp
    document.getElementById('partner-ftpid').value = data.ftpid
    document.getElementById('partner-ftppw').value = data.ftppw
    document.getElementById("partner-opentime").value = data.opentime
    document.getElementById("partner-closedtime").value = data.closedtime
    document.getElementById('partner-paymentdate').value = data.paymentdate
    document.getElementById('partner-bank').value = data.bank
    document.getElementById('partner-bankaccount').value = data.bankaccount
    document.getElementById("partner-monetaryunit").value = data.monetaryunit
    document.getElementById('partner-projecthistory').value = data.projecthistory
    document.getElementById('partner-reputation').value = data.reputation
    document.getElementById('partner-status').value = data.status
    document.getElementById('partner-codename').value = data.codename
    document.getElementById("partner-companytype").value = data.companytype
    document.getElementById('partner-contactpoint').value = data.contactpoint
    document.getElementById('partner-pmsurl').value = data.pmsurl
    document.getElementById('partner-isabroad').checked = data.isabroad
    document.getElementById('partner-isclient').checked = data.isclient
    document.getElementById('partner-tags').value = data.tags.join(",")
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

function InitPartnerModal() {
    // id가 "partner-" 로 시작하는 input > text, input > checkbox, select를 초기화 한다.
    let inputs = document.querySelectorAll("[id^='partner-']")
    for (var i = 0; i < inputs.length; i += 1) {
        if (inputs[i].type === "checkbox") {
            inputs[i].checked = false
        } else {
            inputs[i].value = ""
        }
    }
}

function SetPartnerModal(id) {
    EditPartnerMode()
    fetch('/api/partner/'+id, {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    //return res.text().then(text => { throw new Error(text) })
    .then((response) => {
        if (!response.ok) {
            // response 값은 Promis 타입이다. then으로 처리한다.
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
        partnerObjectToUx(data)
    })
    .catch((err) => {
        console.log(err)
    });
}

function AddPartnerMode() {
    document.getElementById('postpartner-button').hidden = false
    document.getElementById('deletepartner-button').hidden = true
    document.getElementById('putpartner-button').hidden = true
    InitPartnerModal()
}

function EditPartnerMode() {
    document.getElementById('postpartner-button').hidden = true
    document.getElementById('deletepartner-button').hidden = true
    document.getElementById('putpartner-button').hidden = false
}

function DeletePartnerMode() {
    document.getElementById('postpartner-button').hidden = true
    document.getElementById('deletepartner-button').hidden = false
    document.getElementById('putpartner-button').hidden = true
}

function PostPartner() {
    let partner = new Object()
    partner = UxToPartnerObject(partner)
    if (partner.name === "") {
        tata.error('Error', "Need name.", {
            position: 'tr',
            duration: 5000,
            onClose: null,
        })
        return
    }
    fetch('/api/partner', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(partner),
    })
    //return res.text().then(text => { throw new Error(text) })
    .then((response) => {
        if (!response.ok) {
            // response 값은 Promis 타입이다. then으로 처리한다.
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
        // 성공 메시지 띄우기
        tata.success('Add', data.name + "가 추가되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClick: goPartnerPage,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function PutPartner() {
    let partner = new Object()
    partner = UxToPartnerObject(partner)
    if (partner.name === "") {
        tata.error('Error', "Need name.", {
            position: 'tr',
            duration: 5000,
            onClose: null,
        })
        return
    }
    fetch('/api/partner/'+document.getElementById('partner-id').value, {
        method: 'PUT',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(partner),
    })
    //return res.text().then(text => { throw new Error(text) })
    .then((response) => {
        if (!response.ok) {
            // response 값은 Promis 타입이다. then으로 처리한다.
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
        // 성공 메시지 띄우기
        tata.success('Edit', data.name + "가 편집되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClick: goPartnerPage,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function DeletePartner() {
    fetch('/api/partner/'+document.getElementById('partner-id').value, {
        method: 'DELETE',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            // response 값은 Promis 타입이다. then으로 처리한다.
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
        // 성공 메시지 띄우기
        tata.success('Delete', data.name + "가 삭제되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClick: goPartnerPage,
            onClose: null,
        })
    })
    .catch((err) => {
        console.log(err)
    });
}

function goPartnerPage() {
    window.location.replace("partners")
}