
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

function PostPartner() {
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

function goPartnerPage() {
    window.location.replace("partners")
}