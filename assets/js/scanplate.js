document.getElementById("scanplateuploadzone").style.display = "none"

function checkUploadMethod() {
    if (document.getElementById("uploadmethod").checked) {
        document.getElementById("scanplateuploadzone").style.display = "block"
    } else {
        document.getElementById("scanplateuploadzone").style.display = "none"
    }
}