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