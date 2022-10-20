var scanplateUploadZone = document.getElementById("scanplateuploadzone")
if(scanplateUploadZone){
    scanplateUploadZone.style.display = "none"
}


function setRmScanPlateModal(id) {
    document.getElementById("modal-rmscanplate-id").value = id
}

function SetScanPlateProcessStatusModal(id, status) {
    document.getElementById("modal-scanplateprocessstatus-id").value = id
    document.getElementById("modal-scanplateprocessstatus-status").value = status
}

function rmscanplate() {
    let id = document.getElementById("modal-rmscanplate-id").value
    fetch("/api/scanplate/" + id, {
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
    .then((data) => {
        document.getElementById(data.id).remove();
        tata.success('Delete', data.id, {position: 'tr',duration: 5000,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

function setscanplateprocessstatus() {
    let id = document.getElementById("modal-scanplateprocessstatus-id").value
    let status = document.getElementById("modal-scanplateprocessstatus-status").value
    fetch("/api/scanplate/" + id, {
        method: 'PATCH',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        body: new URLSearchParams({
            processstatus: status,
        }),
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
    .then((data) => {
        // 상태를 바꾸기
        let statusbtn = document.getElementById("processstatus-"+data.id)
        let html = ""
        if (status == "processing") {
            html = `<span class="badge bg-success finger" data-bs-toggle="modal" data-bs-target="#modal-scanplateprocessstatus" onclick="SetScanPlateProcessStatusModal('${data.id}', 'processing')">processing</span>`
        } else if (status == "wait") {
            html = `<span class="badge bg-warning text-dark finger" data-bs-toggle="modal" data-bs-target="#modal-scanplateprocessstatus" onclick="SetScanPlateProcessStatusModal('${data.id}', 'wait')">wait</span>`
        } else if (status == "error") {
            html = `<span class="badge bg-danger finger" data-bs-toggle="modal" data-bs-target="#modal-scanplateprocessstatus" onclick="SetScanPlateProcessStatusModal('${data.id}', 'error')">error</span>`
        } else {
            html = `<span class="badge bg-secondary finger" data-bs-toggle="modal" data-bs-target="#modal-scanplateprocessstatus" onclick="SetScanPlateProcessStatusModal('${data.id}', 'done')">done</span>`
        }
        statusbtn.innerHTML = ""
        statusbtn.innerHTML = html
        tata.success('Edit Status', data.id, {position: 'tr',duration: 5000,onClose: null})
    })
    .catch((err) => {
        console.log(err)
    });
}

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

