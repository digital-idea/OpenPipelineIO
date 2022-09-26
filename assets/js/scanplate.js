document.getElementById("scanplateuploadzone").style.display = "none"

function checkUploadMethod() {
    if (document.getElementById("uploadmethod").checked) {
        document.getElementById("scanplateuploadzone").style.display = "block"
    } else {
        document.getElementById("scanplateuploadzone").style.display = "none"
    }
}


function DeleteScanPlateTemp() {
    let dz = document.getElementById("scanplateuploaddropzone")
    dz.innerHTML = `<div class="fallback"><input name="file" type="file" hidden></div>`
    
    let endpoint = ""
    let token = ""
    if (document.getElementById("dns").value != "") {
        endpoint = document.getElementById("dns").value
    }
    if (document.getElementById("scanplatetoken").value != "") {
        token = document.getElementById("scanplatetoken").value
    } else {
        token = document.getElementById("token").value
    }
    fetch(endpoint+"/api/scanplatetemp", {
        method: 'DELETE',
        headers: {
            "Authorization": "Basic "+ token,
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
        
        // Dropzone.discover(); // deprecated 된 옵션. false로 해놓는걸 공식문서에서 명시
        // Dropzone.forElement('#scanplateuploaddropzone').removeAllFiles(true)
        tata.success('Delete', obj.path + "가 삭제되었습니다.", {position: 'tr',duration: 5000,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

