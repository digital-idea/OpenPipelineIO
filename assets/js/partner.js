
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
    document.getElementById('partner-id').value = ""
    document.getElementById('partner-name').value = ""
    document.getElementById('partner-symbol').value = ""
    document.getElementById('partner-domain').value = ""
    document.getElementById('partner-size').value = ""
    document.getElementById('partner-homepage').value = ""
    document.getElementById('partner-address').value = ""
    document.getElementById('partner-phone').value = ""
    document.getElementById('partner-email').value = ""
    document.getElementById('partner-location').value = ""
    document.getElementById('partner-timezone').value = ""
    document.getElementById('partner-description').value = ""
    document.getElementById('partner-businessregistrationnumber').value = ""
    document.getElementById('partner-manager').value = ""
    document.getElementById('partner-managerphone').value = ""
    document.getElementById('partner-manageremail').value = ""
    document.getElementById('partner-ftp').value = ""
    document.getElementById('partner-ftpid').value = ""
    document.getElementById('partner-ftppw').value = ""
    document.getElementById("partner-opentime").value = ""
    document.getElementById("partner-closedtime").value = ""
    document.getElementById('partner-paymentdate').value = ""
    document.getElementById('partner-bank').value = ""
    document.getElementById('partner-bankaccount').value = ""
    document.getElementById("partner-monetaryunit").value = ""
    document.getElementById('partner-projecthistory').value = ""
    document.getElementById('partner-reputation').value = ""
    document.getElementById('partner-status').value = ""
    document.getElementById('partner-codename').value = ""
    document.getElementById("partner-companytype").value = ""
    document.getElementById('partner-contactpoint').value = ""
    document.getElementById('partner-isabroad').checked = false
    document.getElementById('partner-isclient').checked = false
    document.getElementById('partner-tags').value = ""
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
        document.getElementById('partner-id').value = id
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
        document.getElementById('partner-isabroad').checked = data.isabroad
        document.getElementById('partner-isclient').checked = data.isclient
        document.getElementById('partner-tags').value = data.tags.join(",")
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
    InitPartnerModal()
    let partner = new Object()
    partner.name = document.getElementById('partner-name').value
    partner.symbol = document.getElementById('partner-symbol').value
    partner.domain = document.getElementById('partner-domain').value
    partner.size = document.getElementById('partner-size').value
    partner.homepage = document.getElementById('partner-homepage').value
    partner.address = document.getElementById('partner-address').value
    partner.phone = document.getElementById('partner-phone').value
    partner.email = document.getElementById('partner-email').value
    partner.location = document.getElementById('partner-location').value
    partner.timezone = document.getElementById('partner-timezone').value
    partner.description = document.getElementById('partner-description').value
    partner.businessregistrationnumber = document.getElementById('partner-businessregistrationnumber').value
    partner.manager = document.getElementById('partner-manager').value
    partner.managerphone = document.getElementById('partner-managerphone').value
    partner.manageremail = document.getElementById('partner-manageremail').value
    partner.ftp = document.getElementById('partner-ftp').value
    partner.ftpid = document.getElementById('partner-ftpid').value
    partner.ftppw = document.getElementById('partner-ftppw').value
    partner.opentime = document.getElementById("partner-opentime").options[document.getElementById("partner-opentime").selectedIndex].value,    
    partner.closedtime = document.getElementById("partner-closedtime").options[document.getElementById("partner-closedtime").selectedIndex].value,
    partner.paymentdate = document.getElementById('partner-paymentdate').value
    partner.bank = document.getElementById('partner-bank').value
    partner.bankaccount = document.getElementById('partner-bankaccount').value
    partner.monetaryunit = document.getElementById("partner-monetaryunit").options[document.getElementById("partner-monetaryunit").selectedIndex].value,
    partner.projecthistory = document.getElementById('partner-projecthistory').value
    partner.reputation = document.getElementById('partner-reputation').value
    partner.status = document.getElementById('partner-status').value
    partner.codename = document.getElementById('partner-codename').value
    partner.companytype = document.getElementById("partner-companytype").options[document.getElementById("partner-companytype").selectedIndex].value,
    partner.contactpoint = document.getElementById('partner-contactpoint').value
    partner.isabroad = document.getElementById('partner-isabroad').checked
    partner.isclient = document.getElementById('partner-isclient').checked
    partner.tags = string2array(document.getElementById('partner-tags').value)
    
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
        tata.success('Add Partner', data.name + "가 추가되었습니다.", {
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
    partner.name = document.getElementById('partner-name').value
    partner.symbol = document.getElementById('partner-symbol').value
    partner.domain = document.getElementById('partner-domain').value
    partner.size = document.getElementById('partner-size').value
    partner.homepage = document.getElementById('partner-homepage').value
    partner.address = document.getElementById('partner-address').value
    partner.phone = document.getElementById('partner-phone').value
    partner.email = document.getElementById('partner-email').value
    partner.location = document.getElementById('partner-location').value
    partner.timezone = document.getElementById('partner-timezone').value
    partner.description = document.getElementById('partner-description').value
    partner.businessregistrationnumber = document.getElementById('partner-businessregistrationnumber').value
    partner.manager = document.getElementById('partner-manager').value
    partner.managerphone = document.getElementById('partner-managerphone').value
    partner.manageremail = document.getElementById('partner-manageremail').value
    partner.ftp = document.getElementById('partner-ftp').value
    partner.ftpid = document.getElementById('partner-ftpid').value
    partner.ftppw = document.getElementById('partner-ftppw').value
    partner.opentime = document.getElementById("partner-opentime").options[document.getElementById("partner-opentime").selectedIndex].value,    
    partner.closedtime = document.getElementById("partner-closedtime").options[document.getElementById("partner-closedtime").selectedIndex].value,
    partner.paymentdate = document.getElementById('partner-paymentdate').value
    partner.bank = document.getElementById('partner-bank').value
    partner.bankaccount = document.getElementById('partner-bankaccount').value
    partner.monetaryunit = document.getElementById("partner-monetaryunit").options[document.getElementById("partner-monetaryunit").selectedIndex].value,
    partner.projecthistory = document.getElementById('partner-projecthistory').value
    partner.reputation = document.getElementById('partner-reputation').value
    partner.status = document.getElementById('partner-status').value
    partner.codename = document.getElementById('partner-codename').value
    partner.companytype = document.getElementById("partner-companytype").options[document.getElementById("partner-companytype").selectedIndex].value,
    partner.contactpoint = document.getElementById('partner-contactpoint').value
    partner.isabroad = document.getElementById('partner-isabroad').checked
    partner.isclient = document.getElementById('partner-isclient').checked
    partner.tags = string2array(document.getElementById('partner-tags').value)
    
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
        tata.success('Add Partner', data.name + "가 편집되었습니다.", {
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
        tata.success('Delete Partner', data.name + "가 삭제되었습니다.", {
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