// changeReportExcelURI 함수는 프로젝트를 선택하면 reportexcel url을 설정한다.
function changeReportExcelURI() {
    let project = CurrentProject();
    document.getElementById("reportexcelURI").href = "/reportexcel?project=" + project;
}

// changeExportFormatType 함수는 ExportExcel의 포멧 옵션이 변경될 때 실행되는 함수이다.
function changeExportFormatType() {
    let token = document.getElementById("token").value;
    let format = document.getElementById("format").value;
    if (format === "all") {
        let sel = document.getElementById('task');
        sel.innerHTML = "";
        let allOpt = document.createElement('option');
        allOpt.appendChild(document.createTextNode("All"));
        allOpt.value = "all";
        sel.appendChild(allOpt);
    }
    if (format === "shot") {
        $.ajax({
            url: `/api/shottasksetting`,
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data.tasksettings
                // task를 갱신한다.
                let sel = document.getElementById('task');
                sel.innerHTML = "";
                // all 추가
                let allOpt = document.createElement('option');
                allOpt.appendChild(document.createTextNode("all"));
                allOpt.value = "all";
                sel.appendChild(allOpt);
                // shot Task 추가
                for (let i = 0; i < tasks.length; i++) {
                    let opt = document.createElement('option');
                    opt.appendChild( document.createTextNode(tasks[i].name));
                    opt.value = tasks[i].name;
                    sel.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
    if (format === "asset") {
        $.ajax({
            url: `/api/assettasksetting`,
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data.tasksettings
                // task를 갱신한다.
                let sel = document.getElementById('task');
                sel.innerHTML = "";
                // all 추가
                let allOpt = document.createElement('option');
                allOpt.appendChild(document.createTextNode("all"));
                allOpt.value = "all";
                sel.appendChild(allOpt);
                // shot Task 추가
                for (let i = 0; i < tasks.length; i++) {
                    let opt = document.createElement('option');
                    opt.appendChild( document.createTextNode(tasks[i].name));
                    opt.value = tasks[i].name;
                    sel.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}


// exportExcelCurrentPage는 현재 페이지를 엑셀로 뽑는다.
function exportExcelCurrentPage() {
    let project = document.getElementById("searchbox-project").value
    let task = document.getElementById("searchbox-task").value
    let searchword = document.getElementById("searchbox-searchword").value
    let sortkey = document.getElementById("searchbox-sortkey").value
    let searchbartemplate = document.getElementById("searchbox-searchbar-template").value
    let assign = document.getElementById("searchbox-checkbox-assign").checked
    let ready = document.getElementById("searchbox-checkbox-ready").checked
    let wip = document.getElementById("searchbox-checkbox-wip").checked
    let confirm = document.getElementById("searchbox-checkbox-confirm").checked
    let done = document.getElementById("searchbox-checkbox-done").checked
    let omit = document.getElementById("searchbox-checkbox-omit").checked
    let hold = document.getElementById("searchbox-checkbox-hold").checked
    let out = document.getElementById("searchbox-checkbox-out").checked
    let none = document.getElementById("searchbox-checkbox-none").checked
    let truestatusList = []
    let checkStatus = document.querySelectorAll('*[id^="searchbox-checkbox-"]');
    for (i=0;i<checkStatus.length;i++) {
        if (checkStatus[i].checked) {
            truestatusList.push(checkStatus[i].getAttribute("status"))
        }
    }
    truestatus = truestatusList.join(",")
    // 요청
    let url = `/download-excel-file?project=${project}&task=${task}&searchword=${searchword}&sortkey=${sortkey}&searchbartemplate=${searchbartemplate}&assign=${assign}&ready=${ready}&wip=${wip}&confirm=${confirm}&done=${done}&omit=${omit}&hold=${hold}&out=${out}&none=${none}&truestatus=${truestatus}`
    location.href = url
}