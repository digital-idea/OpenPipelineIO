
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
    partner.opentime = document.getElementById('partner-opentime').value
    partner.closedtime = document.getElementById('partner-closedtime').value
    partner.paymentdate = document.getElementById('partner-paymentdate').value
    partner.bank = document.getElementById('partner-bank').value
    partner.bankaccount = document.getElementById('partner-bankaccount').value
    partner.monetaryunit = document.getElementById('partner-monetaryunit').value
    partner.projecthistory = document.getElementById('partner-projecthistory').value
    partner.reputation = document.getElementById('partner-reputation').value
    partner.status = document.getElementById('partner-status').value
    partner.codename = document.getElementById('partner-codename').value
    partner.companytype = document.getElementById('partner-companytype').value
    partner.contactpoint = document.getElementById('partner-contactpoint').value
    partner.isabroad = document.getElementById('partner-isabroad').checked
    partner.isclient = document.getElementById('partner-isclient').checked
    partner.tags = string2array(document.getElementById('partner-tags').value)
    console.log(partner)
    fetch('/api/partner', {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: JSON.stringify(partner),
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        let html = `
        <div class="col-lg-4 col-md-6 col-sm-12">
            <div class="card m-2 bg-darkmode">
                <h6 class="card-header">
                    ${data.name}
                </h6>
                <div class="card-body">
                    <h6 class="card-title">
                        은행: ${data.bank}
                    </h6>
                    <p class="card-text">
                        홈페이지: ${data.homepage}
                    <p>
                </div>
            </div>
        </div>`
        // process
        document.getElementById('partners').innerHTML = html + document.getElementById('partners').innerHTML

        // 메시지 띄우기
        tata.success('Add Partner', data.name + "가 추가되었습니다.", {
            position: 'tr',
            duration: 5000,
            onClose: null,
        })
    })
    .catch((err) => {
        alert(err)
    });
}



