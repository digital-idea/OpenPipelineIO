// modal이 뜨면 오토포커스가 되어야 한다.
$('#modal-addcomment').on('shown.bs.modal', function () {
    $('#modal-addcomment-text').trigger('focus')
})
$('#modal-editcomment').on('shown.bs.modal', function () {
    $('#modal-editcomment-text').trigger('focus')
})
$('#modal-setnote').on('shown.bs.modal', function () {
    $('#modal-setnote-text').trigger('focus')
})
$('#modal-addsource').on('shown.bs.modal', function () {
    $('#modal-addsource-subject').trigger('focus')
})
$('#modal-rmsource').on('shown.bs.modal', function () {
    $('#modal-rmsource-subject').trigger('focus')
})
$('#modal-setrnum').on('shown.bs.modal', function () {
    $('#modal-setrnum-text').trigger('focus')
})
$('#modal-addtag').on('shown.bs.modal', function () {
    $('#modal-addtag-text').trigger('focus')
})
$('#modal-rmtag').on('shown.bs.modal', function () {
    $('#modal-rmtag-text').trigger('focus')
})
$('#modal-deadline2d').on('shown.bs.modal', function () {
    $('#modal-deadline2d-date').trigger('focus')
})
$('#modal-deadline3d').on('shown.bs.modal', function () {
    $('#modal-deadline3d-date').trigger('focus')
})

// Hotkey: http://gcctech.org/csc/javascript/javascript_keycodes.htm
document.onkeyup = function(e) {
    if (e.ctrlKey && e.shiftKey && e.which == 65) {
        selectCheckboxAll()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 68) {
        selectCheckboxNone()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 73) {
        selectCheckboxInvert()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 84) {
        scroll(0,0)
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 77) {
        selectmode()
    } else if (e.ctrlKey && e.altKey && e.shiftKey && e.which == 69) {
        OpenEditfolder()
    } else if (e.which == 119) { // F8
        if ($('#modal-addcomment').hasClass('show')) {
            document.getElementById("modal-addcomment-addbutton").click();
        }
        if ($('#modal-editcomment').hasClass('show')) {
            document.getElementById("modal-editcomment-editbutton").click();
        }
        if ($('#modal-setnote').hasClass('show')) {
            document.getElementById("modal-setnote-editbutton").click();
        }
        if ($('#modal-editnote').hasClass('show')) {
            document.getElementById("modal-editnote-editbutton").click();
        }
    }
};

function OpenEditfolder() {
    let uri = document.getElementById("edit").href;
    window.location = uri;
}

function padNumber(number) {
    number = number.toString();
    while(number.length < 4) {
        number = "0" + number;
    }
    return number;
}

function id2name(id) {
    l = id.split("_");
    l.pop()
    return l.join("_")
}

// setModal 함수는 modalID와 value를 받아서 modal에 셋팅한다.
function setModal(modalID, value) {
    document.getElementById(modalID).value = value;
}

// setEditTaskModal 함수는 project, name, task 정보를 가지고 와서 Edit Task Modal에 값을 채운다.
function setEditTaskModal(project, id, task) {
    document.getElementById("modal-edittask-project").value = project;
    document.getElementById("modal-edittask-id").value = id;
    document.getElementById("modal-edittask-title").innerHTML = "Edit Task" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/task",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            task: task,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-edittask-startdate').value=data.task.startdate;
            document.getElementById('modal-edittask-predate').value=data.task.predate;
            document.getElementById('modal-edittask-date').value=data.task.date;
            document.getElementById('modal-edittask-expectday').value=data.task.expectday;
            document.getElementById('modal-edittask-resultday').value=data.task.resultday;
            document.getElementById('modal-edittask-level').value=data.task.tasklevel;
            document.getElementById('modal-edittask-task').value=data.task.title;
            document.getElementById('modal-edittask-path').value=data.task.mov;
            document.getElementById('modal-edittask-usernote').value=data.task.usernote;
            document.getElementById('modal-edittask-user').value=data.task.user;
            document.getElementById('modal-edittask-usercomment').value=data.task.usercomment;
            document.getElementById('modal-edittask-id').value=data.id;
            // ver2로 검색하면 modal-edittask-status 가 존재하지 않을 수 있다.
            try {
                document.getElementById("modal-edittask-status").value=data.task.status;
            }
            catch(err) {
                document.getElementById("modal-edittask-statusv2").value=data.task.statusv2;
            }
            
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setTimeModal 함수는 id 정보를 가지고 와서 Edit Time Modal 값을 채운다.
function setTimeModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-edittime-id").value = id;
    document.getElementById("modal-edittime-title").innerHTML = "Edit Time" + multiInputTitle(id);
    $.ajax({
        url: "/api/timeinfo",
        type: "post",
        data: {
            "project": project,
            "id": id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('scanin').value = data.scanin;
            document.getElementById('scanout').value = data.scanout;
            document.getElementById('scanframe').value = data.scanframe;
            document.getElementById('scantimecodein').value = data.scantimecodein;
            document.getElementById('scantimecodeout').value = data.scantimecodeout;
            document.getElementById('platein').value = data.platein;
            document.getElementById('plateout').value = data.plateout;
            document.getElementById('handlein').value = data.handlein;
            document.getElementById('handleout').value = data.handleout;
            document.getElementById('justin').value = data.justin;
            document.getElementById('justout').value = data.justout;
            document.getElementById('justtimecodein').value = data.justtimecodein;
            document.getElementById('justtimecodeout').value = data.justtimecodeout;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setShottypeModal 함수는 item id 정보를 이용해 Edit Shottype Modal에 값을 채운다.
function setShottypeModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-shottype-project").value = project
    document.getElementById("modal-shottype-title").innerHTML = "Shot Type" + multiInputTitle(id);
    $.ajax({
        url: "/api/shottype",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-shottype-id').value=id;
            document.getElementById("modal-shottype-type").value=data.shottype;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}



// setUsetypeModal 함수는 project, name을 받아서 Usetype Modal을 설정한다.
function setUsetypeModal(project, id, selectedType) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-usetype-project").value = project
    document.getElementById("modal-usetype-id").value = id
    document.getElementById("modal-usetype-title").innerHTML = "Use Type: " + id2name(id);
    $.ajax({
        url: `/api/usetypes?project=${project}&name=${id2name(id)}`,
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let types = data.types;
            // selectbox를 채운다.
            let sel = document.getElementById('modal-usetype-type');
            sel.innerHTML = "";
            for (let i = 0; i < types.length; i++) {
                let opt = document.createElement('option');
                opt.appendChild( document.createTextNode(types[i]) );
                opt.value = types[i]; 
                sel.appendChild(opt); 
            }
            // 이미 선택된 옵션을 selectbox에서 선택한다.
            document.getElementById("modal-usetype-type").value=selectedType;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setModalCheckbox(modalID, value) {
    if (value === "true") {
        document.getElementById(modalID).checked = true;
    } else {
        document.getElementById(modalID).checked = false;
    }
}

// *문자를 x문자로 바꾼다.
// X를 x문자로 바꾼다.
// 공백을 제거한다.
// 렌즈디스토션값을 입력시 2048*1280 -> 2048x1280 형태로 바꾸기 위함이다.
// 숫자와 x를 제외한 영문입력시 삭제됩니다.
function widthxHeight(event) {
	event = event || window.event;
	event.target.value = event.target.value.replace("*","x");
	event.target.value = event.target.value.replace("X","x");
	event.target.value = event.target.value.replace(/[^\d\x]/gi,"");
}

function sleep( millisecondsToWait ) {
    var now = new Date().getTime();
    while ( new Date().getTime() < now + millisecondsToWait ) {
        /* do nothing; this will exit once it reaches the time limit */
        /* if you want you could do something and exit */
    }
}

function multiInputTitle(id) {
    let checknum = 0;
    let cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            checknum += 1
        }
    }
    if (checknum === 0) {
        return ": " + id2name(id)
    } else if (checknum === 1) {
        return ": " + id2name(id)        
    } else {
        let name = id2name(id);
        let num = checknum - 1
        return `: ${name}외 ${num}건`
    }
}


function addTask(project, id, task) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/addtask",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    let newItem = `<div class="row" id="${data.id}-task-${data.task}">
					<div id="${data.id}-task-${data.task}-status">
						<span class="finger mt-1 badge badge-${data.status} statusbox">${data.task}</span>
					</div>
					<div id="${data.id}-task-${data.task}-predate"></div>
					<div id="${data.id}-task-${data.task}-date"></div>
					<div id="${data.id}-task-${data.task}-user"></div>
					<div id="${data.id}-task-${data.task}-playbutton"></div>
					<div class="ml-1">
						<span class="add" data-toggle="modal" data-target="#modal-edittask" onclick="
                        setEditTaskModal('${project}', '${data.id}', '${data.task}');
                        ">≡</span>
					</div>
                    </div>`;
                    document.getElementById(`${data.id}-tasks`).innerHTML = newItem + document.getElementById(`${data.id}-tasks`).innerHTML;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/addtask",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let newItem = `<div class="row" id="${data.id}-task-${data.task}">
					<div id="${data.id}-task-${data.task}-status">
						<span class="finger mt-1 badge badge-${data.status} statusbox">${data.task}</span>
					</div>
					<div id="${data.id}-task-${data.task}-predate"></div>
					<div id="${data.id}-task-${data.task}-date"></div>
					<div id="${data.id}-task-${data.task}-user"></div>
					<div id="${data.id}-task-${data.task}-playbutton"></div>
					<div class="ml-1">
						<span class="add" data-toggle="modal" data-target="#modal-edittask" onclick="
                        setEditTaskModal('${project}', '${data.id}', '${data.task}');
                        ">≡</span>
					</div>
				</div>`;
                document.getElementById(`${data.id}-tasks`).innerHTML = newItem + document.getElementById(`${data.id}-tasks`).innerHTML;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function rmTask(project, id, task) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/rmtask",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.id}-task-${data.task}`).remove();
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/rmtask",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-task-${data.task}`).remove();
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setFrame(mode, id, frame) {
    let token = document.getElementById("token").value;
    let project = CurrentProject()
    $.ajax({
        url: "/api/" + mode,
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            frame: frame,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (mode === "setscanin") {
                document.getElementById("scanin-"+data.name).innerHTML = `<span class="text-badge ml-1" title="scanin">${data.frame}</span>`;
            }
            if (mode === "setscanout") {
                document.getElementById("scanout-"+data.name).innerHTML = `<span class="text-badge ml-1" title="scanout">${data.frame}</span>`;
            }
            if (mode === "setscanframe") {
                document.getElementById("scanframe-"+data.name).innerHTML = `<span class="text-badge ml-1" title="scanframe">(${data.frame})</span>`;
            }
            if (mode === "sethandlein") {
                document.getElementById("handlein-"+data.name).innerHTML = data.frame;
            }
            if (mode === "sethandleout") {
                document.getElementById("handleout-"+data.name).innerHTML = data.frame;
            }
            if (mode === "setplatein") {
                document.getElementById("platein-"+data.name).innerHTML = `<span class="text-white black-opbg" title="platein">${data.frame}</span>`;
            }
            if (mode === "setplateout") {
                document.getElementById("plateout-"+data.name).innerHTML = `<span class="text-white black-opbg" title="plateout">${data.frame}</span>`;
            }
            if (mode === "setjustin") {
                document.getElementById("justin-"+data.name).innerHTML = `<span class="text-warning black-opbg" title="justin">${data.frame}</span>`;
            }
            if (mode === "setjustout") {
                document.getElementById("justout-"+data.name).innerHTML = `<span class="text-warning black-opbg" title="justout">${data.frame}</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setScanTimecodeIn(project, id, timecode, userid) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setscantimecodein",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            timecode: timecode,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("scantimecodein-"+data.name).innerHTML = `<span class="text-badge ml-1">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubTask() {
    $.ajax({
        url: "/api/setcamerapubtask",
        type: "post",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            task: document.getElementById('modal-cameraoption-pubtask').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("campubtask-"+data.id).innerHTML = `<span class="text-badge ml-1">${data.task}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraLensmm() {
    $.ajax({
        url: "/api/setcameralensmm",
        type: "post",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            lensmm: document.getElementById('modal-cameraoption-lensmm').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("camlensmm-"+data.id).innerHTML = `<span class="text-badge ml-1">${data.lensmm}mm</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubPath() {
    $.ajax({
        url: "/api/setcamerapubpath",
        type: "post",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            path: document.getElementById('modal-cameraoption-pubpath').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("campubpath-"+data.id).innerHTML = `<a href="dilink://${data.path}" class="text-badge ml-1">${data.path}</a>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraOptionModal(project, id) {
    document.getElementById("modal-cameraoption-project").value = project;
    document.getElementById("modal-cameraoption-id").value = id;
    document.getElementById("modal-cameraoption-title").innerHTML = "Camera Option" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-cameraoption-pubtask').value = data.productioncam.pubtask;
            document.getElementById('modal-cameraoption-lensmm').value = data.productioncam.lensmm;
            document.getElementById('modal-cameraoption-pubpath').value = data.productioncam.pubpath;
            if (data.productioncam.projection) {
                document.getElementById("modal-cameraoption-projection").checked = true;
            } else {
                document.getElementById("modal-cameraoption-projection").checked = false;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraProjection() {
    $.ajax({
        url: "/api/setcameraprojection",
        type: "post",
        data: {
            project: document.getElementById('modal-cameraoption-project').value,
            id: document.getElementById('modal-cameraoption-id').value,
            projection: document.getElementById("modal-cameraoption-projection").checked,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            if (data.projection) {
                document.getElementById("camprojection-"+data.id).innerHTML = `<span class="text-badge ml-1">ProjectionCam</span>`;
            } else {
                document.getElementById("camprojection-"+data.id).innerHTML = "";
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setScanTimecodeOut(project, id, timecode, userid) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setscantimecodeout",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            timecode: timecode,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("scantimecodeout-"+data.name).innerHTML = `<span class="text-badge ml-1">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setJustTimecodeIn(project, id, timecode, userid) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setjusttimecodein",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            timecode: timecode,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("justtimecodein-"+data.name).innerHTML = `<span class="text-warning black-opbg">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setJustTimecodeOut(project, id, timecode, userid) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setjusttimecodeout",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            timecode: timecode,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("justtimecodeout-"+data.name).innerHTML = `<span class="text-warning black-opbg">${data.timecode}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setNoteModal(project, id) {
    document.getElementById("modal-setnote-project").value = project;
    document.getElementById("modal-setnote-id").value = id;
    document.getElementById("modal-setnote-title").innerHTML = "Set Note" + multiInputTitle(id);
    document.getElementById("modal-setnote-text").value = "";
}

function editNoteModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-editnote-project").value = project;
    document.getElementById("modal-editnote-id").value = id;
    document.getElementById("modal-editnote-title").innerHTML = "Set Note" + multiInputTitle(id);
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-editnote-text").value = data.note.text;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setNote(project, id, text) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let currentID = cboxes[i].getAttribute("id");
            let overwrite = document.getElementById("modal-setnote-overwrite").checked;
            sleep(200);
            $.ajax({
                url: "/api/setnote",
                type: "post",
                data: {
                    project: project,
                    id: currentID,
                    text: text,
                    userid: userid,
                    overwrite: overwrite,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (overwrite) {
                        // note-{{.Name}} 내부 내용을 교체한다.
                        document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    } else {
                        // note-{{.Name}} 내부 내용에 추가한다.
                        document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.id).innerHTML;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        let overwrite = document.getElementById("modal-setnote-overwrite").checked;
        $.ajax({
            url: "/api/setnote",
            type: "post",
            data: {
                project: project,
                id: id,
                text: text,
                userid: userid,
                overwrite: overwrite,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (overwrite) {
                    // note-{{.Name}} 내부 내용을 교체한다.
                    document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                } else {
                    // note-{{.Name}} 내부 내용에 추가한다.
                    document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.id).innerHTML;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function editNote(project, id, text) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setnote",
        type: "post",
        data: {
            project: project,
            id: id,
            text: text,
            userid: userid,
            overwrite: true,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("note-"+data.id).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// initPublishModal 함수는 modal-addpublish 를 초기화 한다.
function initPublishModal() {
    document.getElementById("modal-addpublish-project").value = ""
    document.getElementById("modal-addpublish-name").value = ""
    document.getElementById("modal-addpublish-task").value = ""
    document.getElementById("modal-addpublish-secondarykey").value = ""
    document.getElementById("modal-addpublish-path").value = ""
    document.getElementById("modal-addpublish-status").value= "usethis"
    document.getElementById("modal-addpublish-tasktouse").value= ""
    document.getElementById("modal-addpublish-subject").value= ""
    document.getElementById("modal-addpublish-mainversion").value= 1
    document.getElementById("modal-addpublish-subversion").value= 0
    document.getElementById("modal-addpublish-filetype").value= ""
    document.getElementById("modal-addpublish-kindofusd").value= ""
}

function setAddPublishModal(project, name, task) {
    initPublishModal() // AddPublishModal을 한번 초기화 한다.
    document.getElementById("modal-addpublish-project").value = project
    document.getElementById("modal-addpublish-name").value = name
    document.getElementById("modal-addpublish-task").value = task
    // publishkey를 셋팅한다.
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/publishkeys",
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(datas) {
            if (datas.length == 0) {
                alert("PublishKey 등록이 필요합니다.");
                document.getElementById('modal-addpublish-addbutton').disabled = true;
                return
            }
            let keys = document.getElementById('modal-addpublish-key');
            keys.innerHTML = "";
            for (let i = 0; i < datas.length; i++){
                let opt = document.createElement('option');
                opt.value = datas[i].id;
                opt.innerHTML = datas[i].id;
                keys.appendChild(opt);
            }
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setEditPublishModal(project, id, task, tasktouse, key, createtime, path) {
    document.getElementById("modal-editpublish-project").value = project
    document.getElementById("modal-editpublish-id").value = id
    document.getElementById("modal-editpublish-task").value = task
    document.getElementById("modal-editpublish-tasktouse").value = tasktouse
    document.getElementById("modal-editpublish-path").value = path
    let token = document.getElementById("token").value;
    // publishkey를 셋팅한다.
    $.ajax({
        url: "/api/publishkeys",
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(datas) {
            if (datas.length == 0) {
                alert("PublishKey 등록이 필요합니다.");
                document.getElementById('modal-editpublish-editbutton').disabled = true;
                return
            }
            let keys = document.getElementById('modal-editpublish-key');
            keys.innerHTML = "";
            for (let i = 0; i < datas.length; i++){
                let opt = document.createElement('option');
                opt.value = datas[i].id;
                opt.innerHTML = datas[i].id;
                keys.appendChild(opt);
            }
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });

    // 아이템 정보를 가지고 와서 modal-editpublish를 채운다.
    $.ajax({
        url: "/api/getpublish",
        type: "post",
        data: {
            "project": project,
            "id": id,
            "task": task,
            "key": key,
            "path": path,
            "createtime": createtime,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-editpublish-key').value = key;
            document.getElementById('modal-editpublish-secondarykey').value = data.secondarykey;
            document.getElementById('modal-editpublish-path').value = data.path;
            document.getElementById('modal-editpublish-status').value = data.status;
            document.getElementById('modal-editpublish-tasktouse').value = data.tasktouse;
            document.getElementById('modal-editpublish-subject').value = data.subject;
            document.getElementById('modal-editpublish-mainversion').value = data.mainversion;
            document.getElementById('modal-editpublish-subversion').value = data.subversion;
            document.getElementById('modal-editpublish-filetype').value = data.filetype;
            document.getElementById('modal-editpublish-filetype').value = data.filetype;
            document.getElementById('modal-editpublish-kindofusd').value = data.kindofusd;
            document.getElementById('modal-editpublish-createtime').value = data.createtime;
            document.getElementById('modal-editpublish-isoutput').checked = data.isoutput;
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function editPublish() {
    let token = document.getElementById("token").value
    let project = document.getElementById('modal-editpublish-project').value
    let id = document.getElementById('modal-editpublish-id').value
    let task = document.getElementById('modal-editpublish-task').value
    let key = document.getElementById('modal-editpublish-key').value
    let createtime = document.getElementById('modal-editpublish-createtime').value
    let secondarykey = document.getElementById('modal-editpublish-secondarykey').value
    let path = document.getElementById('modal-editpublish-path').value
    let status = document.getElementById('modal-editpublish-status').value
    let tasktouse = document.getElementById('modal-editpublish-tasktouse').value
    let subject = document.getElementById('modal-editpublish-subject').value
    let mainversion = document.getElementById('modal-editpublish-mainversion').value
    let subversion = document.getElementById('modal-editpublish-subversion').value
    let filetype = document.getElementById('modal-editpublish-filetype').value
    let kindofusd = document.getElementById('modal-editpublish-kindofusd').value
    let isoutput = false
    if (document.getElementById('modal-editpublish-isoutput').checked) {
        isoutput = true
    }
    // 기존 데이터를 삭제한다.
    $.ajax({
        url: "/api/rmpublish",
        type: "post",
        data: {
            project: project,
            id: id,
            task: task,
            key: key,
            path: path,
            createtime: createtime,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function() {
            
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });

    // 새로운 데이터를 추가한다.
    $.ajax({
        url: "/api/addpublish",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            task: task,
            key: key,
            createtime: createtime,
            secondarykey: secondarykey,
            path: path,
            status: status,
            tasktouse: tasktouse,
            subject: subject,
            mainversion: mainversion,
            subversion: subversion,
            filetype: filetype,
            kindofusd: kindofusd,
            isoutput: isoutput,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAddCommentModal(project, id) {
    document.getElementById("modal-addcomment-project").value = project;
    document.getElementById("modal-addcomment-id").value = id;
    document.getElementById("modal-addcomment-title").innerHTML = "Add Comment" + multiInputTitle(id);
    // init media, media title
    document.getElementById("modal-addcomment-media").value = "";
    document.getElementById("modal-addcomment-mediatitle").value = "";
}

function addComment() {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let project = document.getElementById('modal-addcomment-project').value;
    let id = document.getElementById('modal-addcomment-id').value;
    let text = document.getElementById('modal-addcomment-text').value
    let media = document.getElementById('modal-addcomment-media').value
    let mediatitle = document.getElementById('modal-addcomment-mediatitle').value
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let currentID = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/addcomment",
                type: "post",
                data: {
                    project: project,
                    name: id2name(currentID),
                    text: text,
                    mediatitle: mediatitle,
                    media: media,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // comments-{{.Name}} 내부 내용에 추가한다.
                    let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    let newComment = `<div id="comment-${data.id}-${data.date}">
                    <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
                    <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}', '${data.mediatitle}', '${data.media}')">≡</span>
                    <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}')">×</span>
                    <br><small class="text-warning">${body}</small>`
                    if (data.media != "") {
                        if (data.media.includes("http")) {
                            newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                        } else {
                            newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="dilink://${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                        }
                    }
                    newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
                    document.getElementById("comments-"+data.id).innerHTML = newComment + document.getElementById("comments-"+data.id).innerHTML;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/addcomment",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                text: text,
                mediatitle: mediatitle,
                media: media,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // comments-{{$id}} 내부 내용에 추가한다.
                let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                let newComment = `<div id="comment-${data.id}-${data.date}">
                <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
                <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}', '${data.mediatitle}', '${data.media}')">≡</span>
                <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}')">×</span>
                <br><div class="text-warning small">${body}</div>`
                if (data.media != "") {
                    if (data.media.includes("http")) {
                        newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                    } else {
                        newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="dilink://${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                    }
                }
                newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
                document.getElementById("comments-"+data.id).innerHTML = newComment + document.getElementById("comments-"+data.id).innerHTML;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function editComment() {
    let token = document.getElementById("token").value;
    let project = document.getElementById('modal-editcomment-project').value;
    let id = document.getElementById('modal-editcomment-id').value;
    let time = document.getElementById('modal-editcomment-time').value
    let text = document.getElementById('modal-editcomment-text').value
    let mediatitle = document.getElementById('modal-editcomment-mediatitle').value
    let media = document.getElementById('modal-editcomment-media').value
    $.ajax({
        url: "/api/editcomment",
        type: "post",
        data: {
            project: project,
            id: id,
            time: time,
            text: text,
            mediatitle: mediatitle,
            media: media,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            // comments-${data.id}}-${data.time} 내부 내용을 업데이트 한다.
            let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
            let newComment = `<span class="text-badge">${data.time} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
            <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.project}', '${data.id}', '${data.time}', '${data.text}', '${data.mediatitle}', '${data.media}')">≡</span>
            <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.project}', '${data.id}', '${data.time}', '${data.text}')">×</span>
            <br><div class="text-warning small">${body}</div>`
            if (data.media != "") {
                if (data.media.includes("http")) {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                } else {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
								<a href="dilink://${data.media}" onclick="copyClipboard('${data.media}')">
									<img src="/assets/img/link.svg" class="finger">
								</a>
								<span class="text-white pl-2 small">${data.mediatitle}</span>
							</div>`
                }
            }
            document.getElementById(`comment-${data.id}-${data.time}`).innerHTML = newComment
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRmCommentModal(project, id, time, text) {
    document.getElementById("modal-rmcomment-project").value = project;
    document.getElementById("modal-rmcomment-id").value = id;
    document.getElementById("modal-rmcomment-time").value = time;
    document.getElementById("modal-rmcomment-text").value = text;
    document.getElementById("modal-rmcomment-title").innerHTML = "Rm Comment" + multiInputTitle(id);
}

function setEditReviewModal(id) {
    document.getElementById("modal-editreview-id").value = id;
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-editreview-project").value = data.project;
            document.getElementById("modal-editreview-task").value = data.task;
            document.getElementById("modal-editreview-name").value = data.name;
            document.getElementById("modal-editreview-stage").value = data.stage;
            document.getElementById("modal-editreview-createtime").value = data.createtime;
            document.getElementById("modal-editreview-path").value = data.path;
            document.getElementById("modal-editreview-mainversion").value = data.mainversion;
            document.getElementById("modal-editreview-subversion").value = data.subversion;
            document.getElementById("modal-editreview-fps").value = data.fps;
            document.getElementById("modal-editreview-description").value = data.description;
            document.getElementById("modal-editreview-camerainfo").value = data.camerainfo;
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setRmReviewCommentModal(id, time) {
    document.getElementById("modal-rmreviewcomment-id").value = id;
    document.getElementById("modal-rmreviewcomment-time").value = time;
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-rmreviewcomment-project").value = data.project;
            document.getElementById("modal-rmreviewcomment-name").value = data.name;
            for (let i = 0; i < data.comments.length; i++) {
                if (data.comments[i].date == time) {
                    document.getElementById("modal-rmreviewcomment-text").value = data.comments[i].text;
                    break
                }
            }
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setEditReviewCommentModal(id, time) {
    document.getElementById("modal-editreviewcomment-id").value = id;
    document.getElementById("modal-editreviewcomment-time").value = time;
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            for (let i = 0; i < data.comments.length; i++) {
                if (data.comments[i].date == time) {
                    document.getElementById("modal-editreviewcomment-text").value = data.comments[i].text;
                    document.getElementById("modal-editreviewcomment-media").value = data.comments[i].media;
                    document.getElementById("modal-editreviewcomment-frame").value = data.comments[i].frame;
                    break
                }
            }
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setRmReviewModal(id) {
    // review id의 데이터를 가지고 와서 모달을 설정한다.
    $.ajax({
        url: "/api/review",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-rmreview-id").innerHTML = "ID: " + id;
            document.getElementById("modal-rmreview-id").value = id;
            document.getElementById("modal-rmreview-project").innerHTML = "Project: " + data.project;
            document.getElementById("modal-rmreview-name").innerHTML = "Name: " + data.name;
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function setRmPublishKeyModal(project, id, task, key) {
    document.getElementById("modal-rmpublishkey-project").value = project;
    document.getElementById("modal-rmpublishkey-id").value = id;
    document.getElementById("modal-rmpublishkey-task").value = task;
    document.getElementById("modal-rmpublishkey-key").value = key;
    document.getElementById("modal-rmpublishkey-title").innerHTML = "Rm Publish Key" + multiInputTitle(id);
}

function setRmPublishModal(project, id, task, key, createtime, path) {
    document.getElementById("modal-rmpublish-project").value = project
    document.getElementById("modal-rmpublish-id").value = id
    document.getElementById("modal-rmpublish-task").value = task
    document.getElementById("modal-rmpublish-key").value = key
    document.getElementById("modal-rmpublish-createtime").value = createtime
    document.getElementById("modal-rmpublish-path").value = path
}

function setPublishModal(project, id, task, key, path, createtime, status) {
    document.getElementById("modal-setpublish-project").value = project;
    document.getElementById("modal-setpublish-id").value = id;
    document.getElementById("modal-setpublish-task").value = task;
    document.getElementById("modal-setpublish-key").value = key;
    document.getElementById("modal-setpublish-path").value = path;
    document.getElementById("modal-setpublish-createtime").value = createtime;
    document.getElementById("modal-setpublish-status").value = status;
    document.getElementById("modal-setpublish-status").innerHTML = status;
}

function setEditCommentModal(project, id, time, text, mediatitle, media) {
    document.getElementById("modal-editcomment-project").value = project;
    document.getElementById("modal-editcomment-id").value = id;
    document.getElementById("modal-editcomment-time").value = time;
    document.getElementById("modal-editcomment-text").value = text;
    document.getElementById("modal-editcomment-mediatitle").value = mediatitle;
    document.getElementById("modal-editcomment-media").value = media;
    document.getElementById("modal-editcomment-title").innerHTML = "Edit Comment" + multiInputTitle(id);
}

function setDetailCommentsModal(project, id) {
    document.getElementById("modal-detailcomments-title").innerHTML = "Detail Comments" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            // 기존 디테일을 지운다.
            document.getElementById('modal-detailcomments-body').innerHTML = "";
            // 코멘트를 추가한다.
            let comments = data.comments
            comments.reverse();
            for (var i = 0; i < comments.length; ++i) {
                let br = document.createElement("br")
                // elements way
                let cmt = document.createElement("div")
                cmt.setAttribute("id", "comment-"+data.id+"-"+comments[i].date)
                let info = document.createElement("span")
                info.setAttribute("class","text-badge")
                info.innerHTML = comments[i].date + " / "
                let userinfo = document.createElement("a")
                userinfo.setAttribute("href", "/user?id="+data.comments[i].author)
                userinfo.setAttribute("class","text-darkmode")
                userinfo.innerHTML = comments[i].author
                info.append(userinfo)
                cmt.append(info)
                cmt.append(br)
                let text = document.createElement("div")
                text.setAttribute("class","text-darkmode small")
                text.innerHTML = "<br />" + comments[i].text.replace(/\n/g, "<br />")
                cmt.append(text)
                cmt.append(br)
                if (comments[i].media !== "") {
                    let link = document.createElement("a")
                    let protocol = "dilink://"
                    if (comments[i].media.startsWith("http")) {
                        protocol = ""
                    }
                    link.setAttribute("href", protocol + comments[i].media)
                    link.innerHTML = `<img src="/assets/img/link.svg" class="finger">`
                    cmt.append(link)
                    let span = document.createElement("span")
                    span.setAttribute("class","text-darkmode small pl-2")
                    span.innerHTML = comments[i].mediatitle
                    cmt.append(span)
                }
                let line = document.createElement("hr")
                line.setAttribute("class","my-1 p-0 m-0 divider")
                cmt.append(line)
                let parents= document.getElementById('modal-detailcomments-body')
                parents.append(cmt)
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRmUserModal(id) {
    document.getElementById("modal-rmuser-id").value = id;
}

function rmComment(project, id, date) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/rmcomment",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            date: date
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`comment-${data.id}-${data.date}`).remove();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmReview() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/rmreview",
        type: "post",
        data: {
            id: document.getElementById("modal-rmreview-id").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`review-${data.id}`).remove();
            initCanvas();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function editReviewComment() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/editreviewcomment",
        type: "post",
        data: {
            id: document.getElementById("modal-editreviewcomment-id").value,
            time: document.getElementById("modal-editreviewcomment-time").value,
            text: document.getElementById("modal-editreviewcomment-text").value,
            media: document.getElementById("modal-editreviewcomment-media").value,
            frame: document.getElementById("modal-editreviewcomment-frame").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`reviewcomment-${data.id}-${data.time}-text`).innerText = data.text
            if (data.frame == 0) {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-frame`).remove();
            } else {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-frame`).innerText = data.frame + "/" + (data.frame + data.productionstartframe - 1)
            }
            
            if (data.media.startsWith("http") || data.media.startsWith("rvlink")) {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-media`).setAttribute("href", data.media)
            } else {
                document.getElementById(`reviewcomment-${data.id}-${data.time}-media`).setAttribute("href", "dilink://" + data.media)
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmReviewComment() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/rmreviewcomment",
        type: "post",
        data: {
            id: document.getElementById("modal-rmreviewcomment-id").value,
            time: document.getElementById("modal-rmreviewcomment-time").value,
            project: document.getElementById("modal-rmreviewcomment-project").value,
            name: document.getElementById("modal-rmreviewcomment-name").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`reviewcomment-${data.id}-${data.time}`).remove();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmPublishKey() {
    $.ajax({
        url: "/api/rmpublishkey",
        type: "post",
        data: {
            project: document.getElementById('modal-rmpublishkey-project').value,
            id: document.getElementById('modal-rmpublishkey-id').value,
            task: document.getElementById('modal-rmpublishkey-task').value,
            key: document.getElementById('modal-rmpublishkey-key').value
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function rmPublish() {
    $.ajax({
        url: "/api/rmpublish",
        type: "post",
        data: {
            project: document.getElementById('modal-rmpublish-project').value,
            id: document.getElementById('modal-rmpublish-id').value,
            task: document.getElementById('modal-rmpublish-task').value,
            key: document.getElementById('modal-rmpublish-key').value,
            createtime: document.getElementById('modal-rmpublish-createtime').value,
            path: document.getElementById('modal-rmpublish-path').value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function() {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAddSourceModal(project, id) {
    document.getElementById("modal-addsource-project").value = project;
    document.getElementById("modal-addsource-id").value = id;
    document.getElementById("modal-addsource-subject").value = "";
    document.getElementById("modal-addsource-path").value = "";
    document.getElementById("modal-addsource-title").innerHTML = "Add Source" + multiInputTitle(id);
}

function addSource(project, id, title, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            let currentID = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/addsource",
                type: "post",
                data: {
                    project: project,
                    name: id2name(currentID),
                    title: title,
                    path: path,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // 기존 Sources 추가된다.
                    let source = "";
                    if (path.startsWith("http")) {
                        source = `<div id="source-${data.name}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                    } else {
                        source = `<div id="source-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                    }
                    document.getElementById("sources-"+data.name).innerHTML = document.getElementById("sources-"+data.name).innerHTML + source;
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`sources-${data.name}`).childElementCount > 0) {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.project}','${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/addsource",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                title: title,
                path: path,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // 기존 Sources 추가된다.
                let source = "";
                if (path.startsWith("http")) {
                    source = `<div id="source-${data.name}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                } else {
                    source = `<div id="source-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                }
                document.getElementById("sources-"+data.name).innerHTML = document.getElementById("sources-"+data.name).innerHTML + source;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`sources-${data.name}`).childElementCount > 0) {
                    document.getElementById("source-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.project}','${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("source-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRmSourceModal(project, id) {
    document.getElementById("modal-rmsource-project").value = project;
    document.getElementById("modal-rmsource-id").value = id;
    document.getElementById("modal-rmsource-subject").value = "";
    document.getElementById("modal-rmsource-title").innerHTML = "Rm Source" + multiInputTitle(id);
}

function rmSource(project, id, title) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            currentID = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/rmsource",
                type: "post",
                data: {
                    project: project,
                    name: id2name(currentID),
                    title: title,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`source-${data.name}-${data.title}`).remove();
                    if (document.getElementById(`sources-${data.name}`).childElementCount > 0) {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.project}','${data.id}')>－</span>
                        `
                    } else {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/rmsource",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                title: title,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`source-${data.name}-${data.title}`).remove();
                if (document.getElementById(`sources-${data.name}`).childElementCount > 0) {
                    document.getElementById("source-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmsource" onclick="setRmSourceModal('${data.project}','${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("source-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addsource" onclick="setAddSourceModal('${data.project}','${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setAddReferenceModal(project, id) {
    document.getElementById("modal-addreference-project").value = project;
    document.getElementById("modal-addreference-id").value = id;
    document.getElementById("modal-addreference-subject").value = "";
    document.getElementById("modal-addreference-path").value = "";
    document.getElementById("modal-addreference-title").innerHTML = "Add Reference" + multiInputTitle(id);
}

function addReference(project, id, title, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            let currentID = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/addreference",
                type: "post",
                data: {
                    project: project,
                    name: id2name(currentID),
                    title: title,
                    path: path,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // 기존 References 추가된다.
                    let ref = "";
                    if (path.startsWith("http")) {
                        ref = `<div id="reference-${data.name}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                    } else {
                        ref = `<div id="reference-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                    }
                    document.getElementById("references-"+data.name).innerHTML = document.getElementById("references-"+data.name).innerHTML + ref;
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`references-${data.name}`).childElementCount > 0) {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference" onclick="setRmReferenceModal('${data.project}','${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/addreference",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                title: title,
                path: path,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // 기존 References 추가된다.
                let ref = "";
                if (path.startsWith("http")) {
                    ref = `<div id="reference-${data.name}-${data.title}"><a href="${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
                } else {
                    ref = `<div id="reference-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}" onclick="copyClipboard('${data.path}')">${data.title}</a></div>`;
                }
                document.getElementById("references-"+data.name).innerHTML = document.getElementById("references-"+data.name).innerHTML + ref;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`references-${data.name}`).childElementCount > 0) {
                    document.getElementById("reference-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference"  onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference"  onclick="setRmReferenceModal('${data.project}','${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("reference-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRmReferenceModal(project, id) {
    document.getElementById("modal-rmreference-project").value = project;
    document.getElementById("modal-rmreference-id").value = id;
    document.getElementById("modal-rmreference-subject").value = "";
    document.getElementById("modal-rmreference-title").innerHTML = "Rm Source" + multiInputTitle(id);
}

function rmReference(project, id, title) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            currentID = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/rmreference",
                type: "post",
                data: {
                    project: project,
                    name: id2name(currentID),
                    title: title,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`reference-${data.name}-${data.title}`).remove();
                    if (document.getElementById(`references-${data.name}`).childElementCount > 0) {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference" onclick="setRmReferenceModal('${data.project}','${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/rmreference",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                title: title,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`reference-${data.name}-${data.title}`).remove();
                if (document.getElementById(`references-${data.name}`).childElementCount > 0) {
                    document.getElementById("reference-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmreference" onclick="setRmReferenceModal('${data.project}','${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("reference-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addreference" onclick="setAddReferenceModal('${data.project}','${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setSeq(seq) {
    $.ajax({
        url: "/api/setseq",
        type: "post",
        data: {
            project: document.getElementById('modal-iteminfo-project').value,
            id: document.getElementById('modal-iteminfo-id').value,
            seq: seq,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setPlatePath(path) {
    $.ajax({
        url: "/api/setplatepath",
        type: "post",
        data: {
            project: document.getElementById('modal-iteminfo-project').value,
            id: document.getElementById('modal-iteminfo-id').value,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
        dataType: "json",
        success: function(data) {
            return data;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setThummov(path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let project = document.getElementById('modal-iteminfo-project').value;
    let id = document.getElementById('modal-iteminfo-id').value;
    $.ajax({
        url: "/api/setthummov",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            path: path,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("button-thumbplay-"+data.name).innerHTML = `<a href="dilink://${data.path}" class="play">PLAY</a>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setBeforemov(path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let project = document.getElementById('modal-iteminfo-project').value;
    let id = document.getElementById('modal-iteminfo-id').value;
    $.ajax({
        url: "/api/setbeforemov",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            path: path,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAftermov(path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let project = document.getElementById('modal-iteminfo-project').value;
    let id = document.getElementById('modal-iteminfo-id').value;
    $.ajax({
        url: "/api/setaftermov",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            path: path,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setEditmov(path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let project = document.getElementById('modal-iteminfo-project').value;
    let id = document.getElementById('modal-iteminfo-id').value;
    $.ajax({
        url: "/api/seteditmov",
        type: "post",
        data: {
            project: project,
            id: id,
            path: path,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setRetimeplate(project, id, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setretimeplate",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            path: path,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.path === "") {
                document.getElementById("button-retime-"+data.name).innerHTML = "";
            } else {
                document.getElementById("button-retime-"+data.name).innerHTML = `<a href="dilink://${data.path}" class="badge badge-danger">R</a>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setOCIOcc(project, id, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setociocc",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            path: path,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.path === "") {
                document.getElementById("button-ociocc-"+data.name).innerHTML = "";
            } else {
                document.getElementById("button-ociocc-"+data.name).innerHTML = `<span class="badge badge-info mt-1">N</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRollmedia(project, id, rollmedia) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setrollmedia",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            rollmedia: rollmedia,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.rollmedia === "") {
                document.getElementById(data.name+"-onsetbutton").innerHTML = "";
            } else {
                document.getElementById(data.name+"-onsetbutton").innerHTML = `<a href="/setellite?project=${project}&searchword=${data.rollmedia}" class="badge badge-done statusbox text-dark" target="_blink">onset</a>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setScanname(project, id, scanname) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setscanname",
        type: "post",
        data: {
            project: project,
            id: id,
            scanname: scanname,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.project}-${data.id}-scanname`).innerHTML = data.scanname;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setIteminfoModal(project, id) {
    document.getElementById("modal-iteminfo-project").value = project;
    document.getElementById("modal-iteminfo-id").value = id;
    document.getElementById("modal-iteminfo-title").innerHTML = "Iteminfo" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-iteminfo-name').value = data.name;
            document.getElementById('modal-iteminfo-type').value = data.type;
            document.getElementById('modal-iteminfo-seq').value = data.seq;
            document.getElementById('modal-iteminfo-cut').value = data.cut;
            document.getElementById('modal-iteminfo-platepath').value = data.platepath;
            document.getElementById('modal-iteminfo-thummov').value = data.thummov;
            document.getElementById('modal-iteminfo-beforemov').value = data.beforemov;
            document.getElementById('modal-iteminfo-aftermov').value = data.aftermov;
            document.getElementById('modal-iteminfo-editmov').value = data.editmov;
            document.getElementById('modal-iteminfo-retimeplate').value = data.retimeplate;
            document.getElementById('modal-iteminfo-ociocc').value = data.ociocc;
            document.getElementById('modal-iteminfo-rollmedia').value = data.rollmedia;
            document.getElementById('modal-iteminfo-scanname').value = data.scanname;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTaskMov(project, id, task, mov) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api2/settaskmov",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            task: task,
            mov: mov,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.mov === "") {
                document.getElementById(`${data.id}-task-${data.task}-playbutton`).innerHTML = "";
            } else {
                document.getElementById(`${data.id}-task-${data.task}-playbutton`).innerHTML = `<a class="mt-1 ml-1 badge badge-light" href="dilink://${data.mov}">▶</a>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTaskExpectDay(expectday) {
    let token = document.getElementById("token").value;
    let project = document.getElementById('modal-edittask-project').value;
    let id = document.getElementById('modal-edittask-id').value;
    let task = document.getElementById('modal-edittask-task').value;

    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/settaskexpectday",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    expectday: expectday,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    console.info(data);
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskexpectday",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                expectday: expectday,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data);
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskResultDay(resultday) {
    let token = document.getElementById("token").value;
    let project = document.getElementById('modal-edittask-project').value;
    let id = document.getElementById('modal-edittask-id').value;
    let task = document.getElementById('modal-edittask-task').value;

    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/settaskresultday",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    resultday: resultday,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    console.info(data);
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskresultday",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                resultday: resultday,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data);
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskUser() {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let project = document.getElementById('modal-edittask-project').value
    let id = document.getElementById('modal-edittask-id').value
    let task = document.getElementById('modal-edittask-task').value
    let user = document.getElementById('modal-edittask-user').value
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/settaskuser",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    task: task,
                    user: user,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.username === "") {
                        document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = `<span class="mt-1 ml-1 badge badge-light">${data.username}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskuser",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                task: task,
                user: user,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.username === "") {
                    document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-user`).innerHTML = `<span class="mt-1 ml-1 badge badge-light">${data.username}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskUserComment() {
    let token = document.getElementById("token").value;
    let project = document.getElementById('modal-edittask-project').value
    let id = document.getElementById('modal-edittask-id').value
    let task = document.getElementById('modal-edittask-task').value
    let usercomment = document.getElementById('modal-edittask-usercomment').value
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/settaskusercomment",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    usercomment: usercomment,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.username === "") {
                        document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.usercomment}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskusercomment",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                usercomment: usercomment,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.username === "") {
                    document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-usercomment`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.usercomment}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskStatus(project, id, task, status) { // legacy
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/settaskstatus",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    status: status,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.id}-task-${data.task}-status`).innerHTML = `<span class="finger mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskstatus",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                status: status,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-task-${data.task}-status`).innerHTML = `<span class="finger mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskStatusV2(project, id, task, status) {
    let token = document.getElementById("token").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api2/settaskstatus",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    status: status,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.id}-task-${data.task}-status`).innerHTML = `<a class="mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</a>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api2/settaskstatus",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                status: status,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.id}-task-${data.task}-status`).innerHTML = `<a class="mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</a>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskDate(project, id, task, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200)
            let id = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/settaskdate",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    date: date,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.mov === "") {
                        document.getElementById(`${data.id}-task-${data.task}-date`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.id}-task-${data.task}-date`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.shortdate}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskdate",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                date: date,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.mov === "") {
                    document.getElementById(`${data.id}-task-${data.task}-date`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-date`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.shortdate}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskStartdate(project, id, task, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200)
            let id = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/settaskstartdate",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    date: date,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    console.log(data.date);
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskstartdate",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                date: date,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.log(data.date);
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskUserNote(project, id, task, usernote) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200)
            let id = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/settaskusernote",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    task: task,
                    usernote: usernote,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    console.info(data)
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskusernote",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                task: task,
                usernote: usernote,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}



function setTaskPredate(project, id, task, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200)
            $.ajax({
                url: "/api/settaskpredate",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    task: task,
                    date: date,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.shortdate === "") {
                        document.getElementById(`${data.id}-task-${data.task}-predate`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.id}-task-${data.task}-predate`).innerHTML = `<span class="mt-1 ml-1 badge badge-outline-darkmode">${data.shortdate}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskpredate",
            type: "post",
            data: {
                project: project,
                id: id,
                task: task,
                date: date,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.shortdate === "") {
                    document.getElementById(`${data.id}-task-${data.task}-predate`).innerHTML = "";
                } else {
                    document.getElementById(`${data.id}-task-${data.task}-predate`).innerHTML = `<span class="mt-1 ml-1 badge badge-outline-darkmode">${data.shortdate}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function rfc3339toNormaltime(t) {
    if (t.includes("T")) {
        return t.split("T")[0]
    }
    return t
}

// setDeadline2dModal 함수는 project, id 정보를 이용해서 Deadline2d Modal 값을 채운다.
function setDeadline2dModal(project, id) {
    document.getElementById("modal-deadline2d-project").value = project;
    document.getElementById("modal-deadline2d-id").value = id;
    document.getElementById("modal-deadline2d-title").innerHTML = "Set Deadline 2D" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-deadline2d-date').value = rfc3339toNormaltime(data.ddline2d);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setDeadline2D(project, id, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setdeadline2d",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    date: date,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("deadline2d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline2d" onclick="setDeadline2dModal('${data.project}','${data.id}')">2D:${data.shortdate}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setdeadline2d",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                date: date,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("deadline2d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline2d" onclick="setDeadline2dModal('${data.project}','${data.id}')">2D:${data.shortdate}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

// setDeadline3dModal 함수는 project, id 정보를 이용해서 Deadline3d Modal 값을 채운다.
function setDeadline3dModal(project, id) {
    document.getElementById("modal-deadline3d-project").value = project;
    document.getElementById("modal-deadline3d-id").value = id;
    document.getElementById("modal-deadline3d-title").innerHTML = "Set Deadline 3D" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic " + token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-deadline3d-date').value = rfc3339toNormaltime(data.ddline3d);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setDeadline3D(project, id, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setdeadline3d",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    date: date,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("deadline3d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline3d" onclick="setDeadline3dModal('${data.project}','${data.id}')">3D:${data.shortdate}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setdeadline3d",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                date: date,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("deadline3d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline3d" onclick="setDeadline3dModal('${data.project}','${data.id}')">3D:${data.shortdate}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setShottype(project, id) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let e = document.getElementById("modal-shottype-type");
    let shottype = e.options[e.selectedIndex].value;
    
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setshottype",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    shottype: shottype,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("shottype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-shottype" onclick="setShottypeModal('${project}','${data.id}')">${data.type}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setshottype",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                shottype: shottype,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("shottype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-shottype" onclick="setShottypeModal('${project}','${data.id}')">${data.type}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setUsetype(project, id) {
    let token = document.getElementById("token").value;
    let e = document.getElementById("modal-usetype-type");
    let type = e.options[e.selectedIndex].value;
    
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setusetype",
                type: "post",
                data: {
                    project: project,
                    id: id,
                    type: type,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.project}-${data.id}-usetype`).innerHTML = `<span class="badge badge-warning ml-1" data-toggle="modal" data-target="#modal-usetype" onclick="setUsetypeModal('${project}','${data.id}','${data.type}')">${data.type}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setusetype",
            type: "post",
            data: {
                project: project,
                id: id,
                type: type,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.project}-${data.id}-usetype`).innerHTML = `<span class="badge badge-warning ml-1" data-toggle="modal" data-target="#modal-usetype" onclick="setUsetypeModal('${project}','${data.id}','${data.type}')">${data.type}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setAssettypeModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-assettype-project").value = project
    document.getElementById("modal-assettype-title").innerHTML = "Assettype Type" + multiInputTitle(id);
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-assettype-id').value=id;
            document.getElementById("modal-assettype-type").value=data.assettype;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAssettype(project, id) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let types = document.getElementById("modal-assettype-type");
    let assettype = types.options[types.selectedIndex].value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setassettype",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    assettype: assettype,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // assettype button update
                    document.getElementById("assettype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-assettype" onclick="setAssettypeModal('${project}', '${data.id}')">${data.type}</span>`;
                    // remove old assettype tag
                    document.getElementById(`assettag-${data.name}-${data.oldtype}`).remove();
                    // add new assettype tag
                    let url = `/inputmode?project=${data.project}&searchword=assettags:${data.type}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`;
                    source = `<div id="tag-${data.name}-${data.type}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.type}</a></div>`;
                    document.getElementById("assettags-"+data.name).innerHTML = document.getElementById("assettags-"+data.name).innerHTML + source;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setassettype",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                assettype: assettype,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // assettype button update
                document.getElementById("assettype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#modal-assettype" onclick="setAssettypeModal('${project}', '${data.id}')">${data.type}</span>`;
                // remove old assettype tag
                document.getElementById(`assettag-${data.name}-${data.oldtype}`).remove();
                // add new assettype tag
                let url = `/inputmode?project=${data.project}&searchword=assettags:${data.type}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`;
                source = `<div id="tag-${data.name}-${data.type}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.type}</a></div>`;
                document.getElementById("assettags-"+data.name).innerHTML = document.getElementById("assettags-"+data.name).innerHTML + source;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

// setRnumModal 함수는 project, id 정보를 이용해서 Edit Rnum Modal 값을 채운다.
function setRnumModal(project, id) {
    document.getElementById("modal-setrnum-project").value = project;
    document.getElementById("modal-setrnum-id").value = id;
    document.getElementById("modal-setrnum-title").innerHTML = "Set Rnum number" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-setrnum-text').value = data.rnum;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setRnum(project, id, rnum) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setrnum",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            rnum: rnum,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.rnum !== "") {
                document.getElementById("rnum-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-setrnum" onclick="setModal('modal-setrnum-text', '${data.rnum}' );setModal('modal-setrnum-id', '${data.id}')"{{end}}>${data.rnum}</span>`;
            } else {
                document.getElementById("rnum-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-setrnum" onclick="setModal('modal-setrnum-text', '${data.rnum}' );setModal('modal-setrnum-id', '${data.id}')"{{end}}>no rnum</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAddTagModal(project, id) {
    document.getElementById("modal-addtag-project").value = project;
    document.getElementById("modal-addtag-id").value = id;
    document.getElementById("modal-addtag-text").value = "";
    document.getElementById("modal-addtag-title").innerHTML = "Add Tag" + multiInputTitle(id);
}

function setRenameTagModal(project) {
    document.getElementById("modal-renametag-project").value = project;
    document.getElementById("modal-renametag-beforetag").value = "";
    document.getElementById("modal-renametag-aftertag").value = "";
}

function addTag(project, id, tag) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/addtag",
                type: "post",
                
                data: {
                    project: project,
                    id: id,
                    tag: tag,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    // 기존 Tags에 추가된다.
                    let url = `/inputmode?project=${data.project}&searchword=tag:${data.tag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`
                    source = `<div id="tag-${data.id}-${data.tag}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.tag}</a></div>`;
                    document.getElementById("tags-"+data.id).innerHTML = document.getElementById("tags-"+data.id).innerHTML + source;
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                        document.getElementById("tag-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("tag-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/addtag",
            type: "post",
            
            data: {
                project: project,
                id: id,
                tag: tag,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // 기존 Tags에 추가된다.
                let url = `/inputmode?project=${data.project}&searchword=tag:${data.tag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&task=`
                let source = `<div id="tag-${data.id}-${data.tag}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.tag}</a></div>`;
                document.getElementById("tags-"+data.id).innerHTML = document.getElementById("tags-"+data.id).innerHTML + source;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function renameTag(project, before, after) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/renametag",
        type: "post",
        
        data: {
            project: project,
            before: before,
            after: after,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            alert(`${data.project}프로젝트의 ${data.before} Tag가 ${data.after} Tag로 변경되었습니다.\n새로고침 해주세요.`);
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRmTagModal(project, id) {
    document.getElementById("modal-rmtag-project").value = project;
    document.getElementById("modal-rmtag-id").value = id;
    document.getElementById("modal-rmtag-text").value = "";
    document.getElementById("modal-rmtag-title").innerHTML = "Rm Tag" + multiInputTitle(id);
}

function rmTag(project, id, tag) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/rmtag",
                type: "post",
                
                data: {
                    project: project,
                    id: id,
                    tag: tag,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`tag-${data.id}-${data.tag}`).remove();
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                        document.getElementById("tag-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                        `;
                    } else {
                        document.getElementById("tag-button-"+data.id).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        `;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/rmtag",
            type: "post",
            
            data: {
                project: project,
                id: id,
                tag: tag,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`tag-${data.id}-${data.tag}`).remove();
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`tags-${data.id}`).childElementCount > 0) {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                    `;
                } else {
                    document.getElementById("tag-button-"+data.id).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    `;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function isMultiInput() {
    var cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            return true;
        }
    }
    return false;
}


function selectCheckbox(id) {
    let beforeNum = document.querySelectorAll('input[name=selectID]:checked').length;
    if (document.getElementById(id).checked) {
        if ((beforeNum - 1) === 0) {
            document.getElementById("topbtn").innerHTML = "Top"
        } else {
            document.getElementById("topbtn").innerHTML = "Top<br>" + (beforeNum - 1)
        }
    } else {
        if ((beforeNum + 1) === 0) {
            document.getElementById("topbtn").innerHTML = "Top"
        } else {
            document.getElementById("topbtn").innerHTML = "Top<br>" + (beforeNum + 1)
        }
        
    }
    // 선택된 아이템의 숫자를 갱신한다.
}

function selectCheckboxAll() {
    let cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = true;
    }
    document.getElementById("topbtn").innerHTML = "Top<br>" + (document.querySelectorAll('input[name=selectID]:checked').length)
}

function selectCheckboxNone() {
    let cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = false;
    }
    document.getElementById("topbtn").innerHTML = "Top"
}

function selectCheckboxInvert() {
    let cboxes = document.getElementsByName('selectID');
    let invertNum = 0
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            cboxes[i].checked = true;
            invertNum += 1
        } else {
            cboxes[i].checked = false;
        }
    }
    document.getElementById("topbtn").innerHTML = "Top<br>" + (invertNum)
}

function setTaskLevel(project, id, task, level) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/settasklevel",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    task: task,
                    level: level,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    console.info(data)
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settasklevel",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                task: task,
                level: level,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setObjectIDModal(project, id) {
    document.getElementById("modal-objectid-project").value = project;
    document.getElementById("modal-objectid-id").value = id;
    document.getElementById("modal-objectid-title").innerHTML = "Object ID" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-objectid-in').value = data.objectidin;
            document.getElementById('modal-objectid-out').value = data.objectidout;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setObjectID(project, id, innum, outnum) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setobjectid",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            in: innum,
            out: outnum,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("objectidnum-"+data.name).innerHTML = `<span class="text-badge ml-1">${data.in}-${data.out}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCrowdAsset(project, id) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setcrowdasset",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.crowdasset) {
                document.getElementById("crowdasset-"+data.id).innerHTML = `<span class="badge badge-warning finger" onclick="setCrowdAsset('${data.project}', '${data.id}')">Crowd</span>`;
            } else {
                document.getElementById("crowdasset-"+data.id).innerHTML = `<span class="badge badge-light fade finger" onclick="setCrowdAsset('${data.project}', '${data.id}')">Crowd</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setPlatesizeModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-platesize-project").value = project
    document.getElementById("modal-platesize-title").innerHTML = "Platesize" + multiInputTitle(id);
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-platesize-id').value = id;
            document.getElementById("modal-platesize-size").value = data.platesize;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setPlatesize(project, id, size) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setplatesize",
                type: "post",
                
                data: {
                    project: project,
                    name: id2name(id),
                    size: size,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("platesize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-platesize" onclick="setPlatesizeModal('${project}', '${data.id}')">S: ${data.size}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setplatesize",
            type: "post",
            
            data: {
                project: project,
                name: id2name(id),
                size: size,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("platesize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-platesize" onclick="setPlatesizeModal('${project}', '${data.id}')">S: ${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setUndistortionsizeModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-undistortionsize-project").value = project
    document.getElementById("modal-undistortionsize-title").innerHTML = "Undistortionsize" + multiInputTitle(id);
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-undistortionsize-id').value = id;
            document.getElementById("modal-undistortionsize-size").value = data.dsize;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setUndistortionsize(project, id, size) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setundistortionsize",
                type: "post",
                
                data: {
                    project: project,
                    name: id2name(id),
                    size: size,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("undistortionsize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-undistortionsize" onclick="setUndistortionsizeModal('${project}', '${data.id}', '${data.size}')">U: ${data.size}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setundistortionsize",
            type: "post",
            
            data: {
                project: project,
                name: id2name(id),
                size: size,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("undistortionsize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-undistortionsize" onclick="setUndistortionsizeModal('${project}', '${data.id}', '${data.size}')">U: ${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRendersizeModal(project, id) {
    let token = document.getElementById("token").value;
    document.getElementById("modal-rendersize-project").value = project
    document.getElementById("modal-rendersize-title").innerHTML = "Rendersize" + multiInputTitle(id);
    $.ajax({
        url: `/api/item?project=${project}&id=${id}`,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById('modal-rendersize-id').value = id;
            document.getElementById("modal-rendersize-size").value = data.rendersize;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRendersize(project, id, size) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/setrendersize",
                type: "post",
                
                data: {
                    project: project,
                    name: id2name(id),
                    size: size,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById("rendersize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setRendersizeModal('${project}', '${data.id}')">R: ${data.size}</span>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setrendersize",
            type: "post",
            
            data: {
                project: project,
                name: id2name(id),
                size: size,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById("rendersize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setRendersizeModal('${project}', '${data.id}')">R: ${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function CurrentProject() {
    let e = document.getElementById("searchbox-project");
    return e.options[e.selectedIndex].value;
}

function rmItem() {
    let project = CurrentProject()
    let token = document.getElementById("token").value;
    let cboxes = document.getElementsByName('selectID');
    let selectNum = 0;
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            selectNum += 1
        }
    }
    if (selectNum > 0) {
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/rmitemid",
                type: "post",
                data: {
                    project: project,
                    id: id,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`item-${data.id}`).remove();
                    document.getElementById("topbtn").innerHTML = "Top" // 삭제가되면 Top버튼에서 선택된 아이템 갯수를 리셋한다.
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        alert("삭제할 아이템을 선택해주세요.");
    }
}

function autocomplete(inp) {
    let arr
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/autocompliteusers",
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            arr = data.users;
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
    /*the autocomplete function takes two arguments,
    the text field element and an array of possible autocompleted values:*/
    var currentFocus;
    /*execute a function when someone writes in the text field:*/
    inp.addEventListener("input", function(e) {
        var a, b, i, val = this.value; // a:유저를 표기하는 박스, b:유저를 표기하는 박스 내부의 한 요소, i: for문에 필요한 인수, val: input창에 입력한 값
        /*이미 열려있는 리스트를 닫는다.*/
        closeAllLists();
        if (!val) { return false;}
        currentFocus = -1;
        // DIV 하나를 생성한다.
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        /*위에서 생성한 검색창을 부모에 붙힌다.*/
        this.parentNode.appendChild(a);
        /*각각의 아이템을 순환한다.*/
        for (i = 0; i < arr.length; i++) {
          /*검색어가 아이템에 포함되어 있다면, div를 생성한다.*/
          if (arr[i].searchword.includes(val)) {
            /*create a DIV element for each matching element:*/
            b = document.createElement("DIV");
            let userInfo = arr[i].id + "(" + arr[i].name + "," + arr[i].team + ")";
            /*make the matching letters bold:*/
            let otherList = userInfo.split(val)
            b.innerHTML = otherList[0] + "<span class='text-warning'>" + val + "</span>" + otherList[1];
            /*insert a input field that will hold the current array item's value:*/
            b.innerHTML += `<input type='hidden' value='${userInfo}'>`;
            /*execute a function when someone clicks on the item value (DIV element):*/
            b.addEventListener("click", function(e) {
                /*insert the value for the autocomplete text field:*/
                inp.value = this.getElementsByTagName("input")[0].value;
                /*close the list of autocompleted values,
                (or any other open lists of autocompleted values:*/
                closeAllLists();
            });
            a.appendChild(b);
          }
        }
    });
    /*execute a function presses a key on the keyboard:*/
    inp.addEventListener("keydown", function(e) {
        var x = document.getElementById(this.id + "autocomplete-list");
        if (x) x = x.getElementsByTagName("div");
        if (e.keyCode == 40) {
          /*If the arrow DOWN key is pressed,
          increase the currentFocus variable:*/
          currentFocus++;
          /*and and make the current item more visible:*/
          addActive(x);
        } else if (e.keyCode == 38) { //up
          /*If the arrow UP key is pressed,
          decrease the currentFocus variable:*/
          currentFocus--;
          /*and and make the current item more visible:*/
          addActive(x);
        } else if (e.keyCode == 13) {
          /*If the ENTER key is pressed, prevent the form from being submitted,*/
          e.preventDefault();
          if (currentFocus > -1) {
            /*and simulate a click on the "active" item:*/
            if (x) x[currentFocus].click();
          }
        }
    });
    function addActive(x) {
      /*a function to classify an item as "active":*/
      if (!x) return false;
      /*start by removing the "active" class on all items:*/
      removeActive(x);
      if (currentFocus >= x.length) currentFocus = 0;
      if (currentFocus < 0) currentFocus = (x.length - 1);
      /*add class "autocomplete-active":*/
      x[currentFocus].classList.add("autocomplete-active");
    }
    function removeActive(x) {
      /*a function to remove the "active" class from all autocomplete items:*/
      for (var i = 0; i < x.length; i++) {
        x[i].classList.remove("autocomplete-active");
      }
    }
    function closeAllLists(elmnt) {
      /*close all autocomplete lists in the document,
      except the one passed as an argument:*/
      var x = document.getElementsByClassName("autocomplete-items");
      for (var i = 0; i < x.length; i++) {
        if (elmnt != x[i] && elmnt != inp) {
        x[i].parentNode.removeChild(x[i]);
      }
    }
  }
  
  /*execute a function when someone clicks in the document:*/
  document.addEventListener("click", function (e) {
      closeAllLists(e.target);
  });
}

let input = document.getElementsByClassName("searchuser");
for (var i = 0; i < input.length; i++) {
    autocomplete(input[i])
}

function setAddTaskModal(project, id, type) {
    document.getElementById("modal-addtask-project").value = project;
    document.getElementById("modal-addtask-id").value = id;
    document.getElementById("modal-addtask-title").innerHTML = "Add Task" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    if (type === "org" || type === "left") {
        $.ajax({
            url: "/api/shottasksetting",
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data["tasksettings"];
                let addtasks = document.getElementById('modal-addtask-taskname');
                addtasks.innerHTML = "";
                for (let i = 0; i < tasks.length; i++){
                    let opt = document.createElement('option');
                    opt.value = tasks[i].name;
                    opt.innerHTML = tasks[i].name;
                    addtasks.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
    if (type === "asset") {
        $.ajax({
            url: "/api/assettasksetting",
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data["tasksettings"]
                let addtasks = document.getElementById('modal-addtask-taskname');
                addtasks.innerHTML = "";
                for (let i = 0; i < tasks.length; i++){
                    let opt = document.createElement('option');
                    opt.value = tasks[i].name;
                    opt.innerHTML = tasks[i].name;
                    addtasks.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRmTaskModal(project, id, type) {
    document.getElementById("modal-rmtask-project").value = project;
    document.getElementById("modal-rmtask-id").value = id;
    document.getElementById("modal-rmtask-title").innerHTML = "Rm Task" + multiInputTitle(id);
    let token = document.getElementById("token").value;
    if (type === "org" || type === "left") {
        $.ajax({
            url: "/api/shottasksetting",
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data["tasksettings"];
                let rmtasks = document.getElementById('modal-rmtask-taskname');
                rmtasks.innerHTML = "";
                for (let i = 0; i < tasks.length; i++){
                    let opt = document.createElement('option');
                    opt.value = tasks[i].name;
                    opt.innerHTML = tasks[i].name;
                    rmtasks.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
    if (type === "asset") {
        $.ajax({
            url: "/api/assettasksetting",
            type: "get",
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let tasks = data["tasksettings"]
                let rmtasks = document.getElementById('modal-rmtask-taskname');
                rmtasks.innerHTML = "";
                for (let i = 0; i < tasks.length; i++){
                    let opt = document.createElement('option');
                    opt.value = tasks[i].name;
                    opt.innerHTML = tasks[i].name;
                    rmtasks.appendChild(opt);
                }
            },
            error: function(request,status,error){
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function changeStatusURI(status) {
    let tags = document.getElementsByClassName("statusuri");
    for ( var i = 0; i < tags.length; i++) {
        let c = document.getElementById("searchbox-checkbox-" + status);
        if (tags[i].href.includes(status + "=true")) {
            tags[i].href = tags[i].href.replace(status + "=true", status + "=" + c.checked)
        } else {
            tags[i].href = tags[i].href.replace(status + "=false", status + "=" + c.checked)
        }
    }
}

function mailInfo(project, id) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/mailinfo",
        type: "post",
        data: {
            "project": project,
            "id": id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let mailString = "mailto:"
            if (data.mails) {
                mailString += data.mails.join(",")
            }
            mailString += `?subject=[${data.header}] ${data.title}&`;
            // 메일을 참조할 사람을 추가한다.
            if (data.cc) {
                mailString += `cc=${data.cc.join(",")}`
            }
            window.location.href = mailString;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function foldingmenu() {
    let searchbox = document.getElementById("searchbox")
    if(searchbox.style.display=="none") {
        // 펼치기
        searchbox.style.display='block';
        // item 위에 빈영역을 만들어 주어야 한다.
        document.getElementById("blinkspace").style.height = "550px";
        // 메뉴글씨 바꾸기
        document.getElementById("foldoption").innerText = "Collapse Searchbox ▲"
    } else {
        // 접기
        searchbox.style.display='none';
        // item 위에 빈영역을 만들어 주어야 한다.
        document.getElementById("blinkspace").style.height = "100px";
        // 메뉴글씨 바꾸기
        document.getElementById("foldoption").innerText = "Expand Searchbox ▼"
    }
}

// TopClick 함수는 스크롤시 보여지는 Top 버튼을 누를 때 발생하는 이벤트이다.
function TopClick() {
    document.body.scrollTop = 0;
    document.documentElement.scrollTop = 0;
    let searchbox = document.getElementById("searchbox")
    if (searchbox !== null) {
        searchbox.style.display = "block";
        document.getElementById("blinkspace").style.height = "550px";
        document.getElementById("foldoption").innerText = "Collapse Searchbox ▲"
    }
}

function setPublish() {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setpublishstatus",
        type: "post",
        data: {
            project: document.getElementById('modal-setpublish-project').value,
            id: document.getElementById('modal-setpublish-id').value,
            task: document.getElementById('modal-setpublish-task').value,
            key: document.getElementById('modal-setpublish-key').value,
            path: document.getElementById('modal-setpublish-path').value,
            createtime: document.getElementById('modal-setpublish-createtime').value,
            status: document.getElementById("modal-setpublish-status").value,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function() {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function addPublish() {
    let token = document.getElementById("token").value
    let project = document.getElementById('modal-addpublish-project').value
    let name = document.getElementById('modal-addpublish-name').value
    let task = document.getElementById('modal-addpublish-task').value
    let key = document.getElementById('modal-addpublish-key').value
    let secondarykey = document.getElementById('modal-addpublish-secondarykey').value
    let path = document.getElementById('modal-addpublish-path').value
    let status = document.getElementById('modal-addpublish-status').value
    let tasktouse = document.getElementById('modal-addpublish-tasktouse').value
    let subject = document.getElementById('modal-addpublish-subject').value
    let mainversion = document.getElementById('modal-addpublish-mainversion').value
    let subversion = document.getElementById('modal-addpublish-subversion').value
    let filetype = document.getElementById('modal-addpublish-filetype').value
    let kindofusd = document.getElementById('modal-addpublish-kindofusd').value
    let isoutput = false
    if (document.getElementById('modal-addpublish-isoutput').checked) {
        isoutput = true
    }
    $.ajax({
        url: "/api/addpublish",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            key: key,
            secondarykey: secondarykey,
            path: path,
            status: status,
            tasktouse: tasktouse,
            subject: subject,
            mainversion: mainversion,
            subversion: subversion,
            filetype: filetype,
            kindofusd: kindofusd,
            createtime: "",
            isoutput: isoutput,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function() {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function addReview() {
    let token = document.getElementById("token").value
    let reviewFps = document.getElementById("modal-addreview-fps")
    $.ajax({
        url: "/api/addreview",
        type: "post",
        data: {
            project: document.getElementById("modal-addreview-project").value,
            name: document.getElementById("modal-addreview-name").value,
            stage: document.getElementById("modal-addreview-stage").value,
            task: document.getElementById("modal-addreview-task").value,
            author: document.getElementById("modal-addreview-author").value,
            path: document.getElementById("modal-addreview-path").value,
            description: document.getElementById("modal-addreview-description").value,
            camerainfo: document.getElementById("modal-addreview-camerainfo").value,
            fps: reviewFps.options[reviewFps.selectedIndex].value,
            mainversion: document.getElementById("modal-addreview-mainversion").value,
            subversion: document.getElementById("modal-addreview-subversion").value,
            removeafterprocess: document.getElementById("modal-addreview-removeafterprocess").checked,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function() {
            alert("리뷰가 정상적으로 등록되었습니다.");
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function clickCommentButton() {
    if (document.getElementById("review-comment").value === "") {
        alert("comment가 비어있습니다")
        return
    }
    setReviewStatus('comment')
    addReviewComment()
}

function setReviewStatus(status) {
    $.ajax({
        url: "/api/setreviewstatus",
        type: "post",
        data: {
            status: status,
            id: document.getElementById("current-review-id").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            let item = document.getElementById("reviewstatus-"+data.id)
            // 상태 내부 글씨를 바꾼다.
            item.innerHTML = data.status
            // 상태의 색상을 바꾼다.
            if (data.status === "approve") {
                item.setAttribute("class","ml-1 badge badge-success")
                addReviewCommentText("Approved " + data.stage + " Stage.") // comment를 남긴다.
                setReviewNextStatus(data.id) // 다음 Status를 설정한다.
                setReviewNextStage(data.id) // 다음 Stage를 설정한다.
            } else if (data.status === "comment") {
                item.setAttribute("class","ml-1 badge badge-warning")
            } else {
                item.setAttribute("class","ml-1 badge badge-secondary")
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setReviewNextStatus(id) {
    $.ajax({
        url: "/api/setreviewnextstatus",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            let item = document.getElementById("reviewstatus-"+data.id)
            // 상태 내부 글씨를 바꾼다.
            item.innerHTML = data.status
            // 상태의 색상을 바꾼다.
            if (data.status === "approve") {
                item.setAttribute("class","ml-1 badge badge-success")
            } else if (data.status === "comment") {
                item.setAttribute("class","ml-1 badge badge-warning")
            } else {
                item.setAttribute("class","ml-1 badge badge-secondary")
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setReviewNextStage(id) {
    $.ajax({
        url: "/api/setreviewnextstage",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            let item = document.getElementById("review-stage-"+data.id)
            // 상태 내부 글씨를 바꾼다.
            item.innerHTML = data.stage
            // 상태의 색상을 바꾼다.
            item.setAttribute("class","ml-1 badge badge-stage-"+data.stage)
            // 현재 띄워진 화면의 우측하단의 Stage 상태를 변경한다.
            document.getElementById("current-review-stage").value = data.stage
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setReviewStage(stage) {
    $.ajax({
        url: "/api/setreviewstage",
        type: "post",
        data: {
            stage: stage,
            id: document.getElementById("current-review-id").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            // 해당 id의 stage 글씨와 색상을 바꾼다.
            let itemStage = document.getElementById("review-stage-"+data.id)
            itemStage.innerHTML = data.stage
            itemStage.setAttribute("class","ml-1 badge badge-stage-"+data.stage)
            // 현재 띄워진 화면의 우측하단의 Stage 상태를 변경한다.
            document.getElementById("current-review-stage").value = data.stage
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setReviewProcessStatus(id, status) {
    $.ajax({
        url: "/api/setreviewprocessstatus",
        type: "post",
        data: {
            status: status,
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            let item = document.getElementById("reviewstatus-"+data.id)
            // 상태 내부 글씨를 바꾼다.
            item.innerHTML = data.status
            // 상태의 색상을 바꾼다.
            if (data.processstatus === "wait") {
                item.setAttribute("class","ml-1 badge badge-danger")
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewProject 함수는 리뷰데이터의 Project를 변경한다.
function setReviewProject() {
    $.ajax({
        url: "/api/setreviewproject",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            project: document.getElementById("modal-editreview-project").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-project`).innerText = data.project
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewTask 함수는 리뷰데이터의 Task을 변경한다.
function setReviewTask() {
    $.ajax({
        url: "/api/setreviewtask",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            task: document.getElementById("modal-editreview-task").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-task`).innerText = data.task
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewPath 함수는 리뷰데이터의 Path를 변경한다.
function setReviewPath() {
    $.ajax({
        url: "/api/setreviewpath",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            path: document.getElementById("modal-editreview-path").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            setReviewProcessStatus(data.id, "wait") // 만약 Path가 수정되면 Status가 wait으로 바뀌고 동영상이 다시 연산이 되어야 한다.
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewCreatetime 함수는 리뷰데이터의 Createtime을 변경한다.
function setReviewCreatetime() {
    $.ajax({
        url: "/api/setreviewcreatetime",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            createtime: document.getElementById("modal-editreview-createtime").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// rfc3339 함수는 date 값을 받아서 RFC3339 포멧으로 반환한다.
function rfc3339(d) {
    function pad(n) {
        return n < 10 ? "0" + n : n;
    }
    function timezoneOffset(offset) {
        var sign;
        if (offset === 0) {
            return "Z";
        }
        sign = (offset > 0) ? "-" : "+";
        offset = Math.abs(offset);
        return sign + pad(Math.floor(offset / 60)) + ":" + pad(offset % 60);
    }
    return d.getFullYear() + "-" +
        pad(d.getMonth() + 1) + "-" +
        pad(d.getDate()) + "T" +
        pad(d.getHours()) + ":" +
        pad(d.getMinutes()) + ":" +
        pad(d.getSeconds()) + 
        timezoneOffset(d.getTimezoneOffset());
}

// setReviewCreatetimeNow 함수는 리뷰데이터의 시간을 현재시간으로 설정한다.
function setReviewCreatetimeNow() {
    let date = new Date()
    let time = rfc3339(date)
    $.ajax({
        url: "/api/setreviewcreatetime",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            createtime: time,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("modal-editreview-createtime").value = time
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
    
}

// setReviewMainVersion 함수는 리뷰데이터의 MainVersion을 변경한다.
function setReviewMainVersion() {
    $.ajax({
        url: "/api/setreviewmainversion",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            mainversion: document.getElementById("modal-editreview-mainversion").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-mainversion`).innerText = "v" + data.mainversion
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewSubVersion 함수는 리뷰데이터의 SubVersion을 변경한다.
function setReviewSubVersion() {
    $.ajax({
        url: "/api/setreviewsubversion",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            subversion: document.getElementById("modal-editreview-subversion").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-subversion`).innerText = "v" + data.subversion
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewFps 함수는 리뷰데이터의 Fps를 변경한다.
function setReviewFps() {
    $.ajax({
        url: "/api/setreviewfps",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            fps: document.getElementById("modal-editreview-fps").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewDescription함수는 리뷰데이터의 Description을 변경한다.
function setReviewDescription() {
    $.ajax({
        url: "/api/setreviewdescription",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            description: document.getElementById("modal-editreview-description").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("description").innerText = data.description
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewCameraInfo 함수는 리뷰데이터의 CameraInfo를 변경한다.
function setReviewCameraInfo() {
    $.ajax({
        url: "/api/setreviewcamerainfo",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            camerainfo: document.getElementById("modal-editreview-camerainfo").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            return
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

// setReviewName 함수는 리뷰데이터의 Name을 변경한다.
function setReviewName() {
    $.ajax({
        url: "/api/setreviewname",
        type: "post",
        data: {
            id: document.getElementById("modal-editreview-id").value,
            name: document.getElementById("modal-editreview-name").value,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            document.getElementById(`${data.id}-name`).innerText = data.name
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function addReviewComment() {
    $.ajax({
        url: "/api/addreviewcomment",
        type: "post",
        data: {
            id: document.getElementById("current-review-id").value,
            text: document.getElementById("review-comment").value,
            media: document.getElementById("review-media").value,
            stage: document.getElementById("current-review-stage").value,
            frame: document.getElementById("currentframe").innerHTML,
            framecomment: document.getElementById("review-framecomment").checked,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            // 데이터가 잘 들어가면 review-comments 에 들어간 데이터를 드로잉한다.
            let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
            let newComment = `<div id="reviewcomment-${data.id}-${data.date}" class="p-1">
            <span class="text-badge">${data.date} / <a href="/user?id=${data.author}" class="text-darkmode">${data.author}</a></span>
            <span class="edit" data-toggle="modal" data-target="#modal-editreviewcomment" onclick="setEditReviewCommentModal('${data.id}', '${data.date}')">≡</span>
            <span class="remove" data-toggle="modal" data-target="#modal-rmreviewcomment" onclick="setRmReviewCommentModal('${data.id}','${data.date}')">×</span>
            <br>
            <span class="badge badge-stage-${data.stage}">${data.stage}</span>`
            if (data.framecomment) {
                newComment += `<span class="badge badge-secondary m-1 finger" id="reviewcomment-${data.id}-${data.date}-frame" data-toggle="modal" data-target="#modal-gotoframe" onclick="setModalGotoFrame()">${data.frame}f / ${data.frame+data.productionstartframe-1}f</span>`
            }
            newComment += `<small class="text-white">${body}</small>`
            if (data.media != "") {
                if (data.media.includes("http")) {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
                        <a href="${data.media}" onclick="copyClipboard('${data.media}')">
                            <img src="/assets/img/link.svg" class="finger">
                        </a>
                        <span class="text-white pl-2 small">${data.mediatitle}</span>
                    </div>`
                } else {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
                        <a href="dilink://${data.media}" onclick="copyClipboard('${data.media}')">
                            <img src="/assets/img/link.svg" class="finger">
                        </a>
                        <span class="text-white pl-2 small">${data.mediatitle}</span>
                    </div>`
                }
            }
            newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
            document.getElementById("review-comments").innerHTML = newComment + document.getElementById("review-comments").innerHTML;
            // 입력한 값을 초기화 한다.
            document.getElementById("review-comment").value = ""; 
            document.getElementById("review-media").value = "";
            document.getElementById("review-framecomment").checked = false;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function addReviewCommentText(text) {
    $.ajax({
        url: "/api/addreviewcomment",
        type: "post",
        data: {
            id: document.getElementById("current-review-id").value,
            text: text,
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            // 데이터가 잘 들어가면 review-comments 에 들어간 데이터를 드로잉한다.
            let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
            let newComment = `<div id="reviewcomment-${data.id}-${data.date}" class="p-1">
            <span class="text-badge">${data.date} / <a href="/user?id=${data.author}" class="text-darkmode">${data.author}</a></span>
            <span class="edit" data-toggle="modal" data-target="#modal-editreviewcomment" onclick="setEditReviewCommentModal('${data.id}', '${data.date}')">≡</span>
            <span class="remove" data-toggle="modal" data-target="#modal-rmreviewcomment" onclick="setRmReviewCommentModal('${data.id}','${data.date}')">×</span>
            <br><small class="text-white">${body}</small>`
            if (data.media != "") {
                if (data.media.includes("http")) {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
                        <a href="${data.media}" onclick="copyClipboard('${data.media}')">
                            <img src="/assets/img/link.svg" class="finger">
                        </a>
                        <span class="text-white pl-2 small">${data.mediatitle}</span>
                    </div>`
                } else {
                    newComment += `<div class="row pl-3 pt-3 pb-1">
                        <a href="dilink://${data.media}" onclick="copyClipboard('${data.media}')">
                            <img src="/assets/img/link.svg" class="finger">
                        </a>
                        <span class="text-white pl-2 small">${data.mediatitle}</span>
                    </div>`
                }
            }
            newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
            document.getElementById("review-comments").innerHTML = newComment + document.getElementById("review-comments").innerHTML;
            document.getElementById("review-comment").value = ""; // 입력한 값을 초기화 한다.
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTypeAddShot(type, readOnly) {
    document.getElementById("addshot-type").value = type
    document.getElementById('addshot-type').readOnly = readOnly
}

// 리뷰를 위해서 동영상에 그림을 그리기 위해 필요한 글로벌 변수를 셋팅한다.
var drawCanvas, drawCtx;
var mouseStartX=0
var mouseStartY=0
var drawing = false;
var globalClientWidth = 0;
var globalClientHeight = 0;
var globalReviewRenderWidth = 0;
var globalReviewRenderHeight = 0;
var globalReviewRenderWidthOffset = 0;
var globalReviewRenderHeightOffset = 0;
var framelineOffset = 0;
var frameLineMarkHeight = 12; // 프레임 표시라인 높이

function initCanvas() {
    let playerbox = document.getElementById("playerbox"); // player 캔버스를담을 div를 가지고 온다.
    globalClientWidth = playerbox.clientWidth // 클라이언트 사용자의 가로 사이즈를 구한다.
    globalClientHeight = playerbox.clientHeight // 클라이언트 사용자의 세로 사이즈를 구한다.
    // Player 캔버스를 초기화 한다.
    let playerCanvas = document.getElementById("player");
    playerCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    playerCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    playerCanvas.setAttribute("width", globalClientWidth) // 캔버스를 클라이언트 사용자의 가로사이즈로 설정한다.
    playerCanvas.setAttribute("height", globalClientHeight) // 캔버스를 클라이언트 사용자의 세로사이즈로 설정한다.
    // Draw 캔버스를 초기화 한다.
    let drawCanvas = document.getElementById("drawcanvas");
    drawCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    drawCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    drawCanvas.setAttribute("width", globalClientWidth) // 그림을 그리는 캔버스 가로 사이즈를 설정한다.
    drawCanvas.setAttribute("height", globalClientHeight) // 그림을 그리는 캔버스 세로 사이즈를 설정한다.
    // UX 캔버스를 초기화 한다.
    let uxCanvas = document.getElementById("uxcanvas");
    uxCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    uxCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    uxCanvas.setAttribute("width", globalClientWidth) // UX 캔버스 가로 사이즈를 설정한다.
    uxCanvas.setAttribute("height", globalClientHeight) // UX 캔버스 세로 사이즈를 설정한다.
    // Animation UX 캔버스를 초기화 한다.
    let aniuxCanvas = document.getElementById("aniuxcanvas");
    aniuxCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    aniuxCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    aniuxCanvas.setAttribute("width", globalClientWidth) // Animation UX 캔버스 가로 사이즈를 설정한다.
    aniuxCanvas.setAttribute("height", globalClientHeight) // Animation UX 캔버스 세로 사이즈를 설정한다.
    // Screenshot 캔버스를 초기화 한다.
    let screenshotCanvas = document.getElementById("screenshot");
    screenshotCanvas.setAttribute("width", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    screenshotCanvas.setAttribute("height", 0) // 이 줄이 없으면 아이템을 클릭할 때 마다 캔버스가 계속 커진다.
    screenshotCanvas.setAttribute("width", globalClientWidth) // 스크린샷 캔버스 가로 사이즈를 설정한다.
    screenshotCanvas.setAttribute("height", globalClientHeight) // 스크린샷 캔버스 세로 사이즈를 설정한다.
}

function selectReviewItem(id, project, fps) {
    // 입력받은 프로젝트로 웹페이지의 Review Title을 변경한다.
    document.title = "Review: " + project;
    let playerbox = document.getElementById("playerbox"); // player 캔버스를담을 div를 가지고 온다.
    let clientWidth = playerbox.clientWidth // 클라이언트 사용자의 가로 사이즈를 구한다.
    let clientHeight = playerbox.clientHeight // 클라이언트 사용자의 세로 사이즈를 구한다.
    initCanvas();
    let playerCanvas = document.getElementById("player");
    let playerCtx = playerCanvas.getContext("2d");
    let drawCanvas = document.getElementById("drawcanvas");
    drawCtx = drawCanvas.getContext("2d")
    let uxCanvas = document.getElementById("uxcanvas");
    uxCtx = uxCanvas.getContext("2d")
    let aniuxCanvas = document.getElementById("aniuxcanvas");
    aniuxCtx = aniuxCanvas.getContext("2d")
    let screenshotCanvas = document.getElementById("screenshot");
    screenshotCtx = screenshotCanvas.getContext("2d")

    // 브러쉬 설정
    drawCtx.lineWidth = 4; // 브러시 사이즈
    drawCtx.strokeStyle = "#EFEAD6" // 브러시 컬러

    // 마우스 이벤트 처리
    drawCanvas.addEventListener("mousemove", function (e) {move(e)},false);
    drawCanvas.addEventListener("mousedown", function (e) {down(e)}, false);
    drawCanvas.addEventListener("mouseup", function (e) {up(e)}, false);
    drawCanvas.addEventListener("mouseout", function (e) {out(e)}, false);

    // 버튼설정 및 버튼 이벤트
    let playButton = document.getElementById("player-play");
    let pauseButton = document.getElementById("player-pause");
    let playAndPauseButton = document.getElementById("player-playandpause");
    let loopAndLoofOffButton = document.getElementById("player-loopandloopoff");
    let startButton = document.getElementById("player-start");
    let endButton = document.getElementById("player-end");
    let beforeFrameButton = document.getElementById("player-left");
    let afterFrameButton = document.getElementById("player-right");
    let gotoFrameInput = document.getElementById("modal-gotoframe-frame");
    let prevDrawing = document.getElementById("drawing-prev");
    let nextDrawing = document.getElementById("drawing-next");

    // fouseReview 함수는 리뷰 ID를 받아서 해당 ID를 가진 리뷰아이템의 스크롤을 포커싱 한다.
    if (id !== "") {
        document.getElementById('review-' + id).scrollIntoView(); // 해당아이템을 포커스한다.
        window.scrollTo({top:0, left:0, behavior:'auto'}); // 위 아이템을 포커싱 하면서 windows 포커싱이 틀어질 수 있다. windows 스크롤을 리셋한다.
    }
    

    // GotoFrame 모달창에서 프레임이 변경되면 해당 프레임으로 이동한다.
    gotoFrameInput.addEventListener("change", function() {
        targetFrame = document.getElementById("modal-gotoframe-frame").value
        video.currentTime = gotoFrame(targetFrame, fps) // video.currentTime이 바뀌기 때문에 video.addEventListener('timeupdate', function () {}) 이벤트가 발생해서 드로잉이 띄워진다.
    });

    // prev Drawing 버튼을 클릭할 때 이벤트
    prevDrawing.addEventListener("click", function() {
        $.ajax({
            url: "/api/reviewdrawingframe",
            type: "post",
            data: {
                id: id,
                frame: document.getElementById("currentframe").innerHTML,
                mode: "prev",
            },
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value
            },
            dataType: "json",
            success: function(data) {
                video.currentTime = gotoFrame(data.resultframe, fps)
            },
            error: function(){
                return
            }
        })
    });

    // next Drawing 버튼을 클릭할 때 이벤트
    nextDrawing.addEventListener("click", function() {
        $.ajax({
            url: "/api/reviewdrawingframe",
            type: "post",
            data: {
                id: id,
                frame: document.getElementById("currentframe").innerHTML,
                mode: "next",
            },
            headers: {
                "Authorization": "Basic "+ document.getElementById("token").value
            },
            dataType: "json",
            success: function(data) {
                video.currentTime = gotoFrame(data.resultframe, fps)
            },
            error: function(){
                return
            }
        })
    });
    // 플레이 버튼을 클릭할 때 이벤트
    playButton.addEventListener("click", function() {
        playAndPauseButton.className = "player-pause"
        video.play();
    });
    // 일시정지 버튼을 클릭할 때 이벤트
    pauseButton.addEventListener("click", function() {
        playAndPauseButton.className = "player-play"
        video.pause();
    });

    // 재생과 정지가 같이 진행되는 버튼
    playAndPauseButton.addEventListener("click", function() {
        if (!video.paused) {
            playAndPauseButton.className = "player-play"
            video.pause();    
        } else {
            playAndPauseButton.className = "player-pause"
            video.play();
        }
    });
    // Loop버튼 클릭시
    loopAndLoofOffButton.addEventListener("click", function() {
        if (loopAndLoofOffButton.className == "player-loop") {
            loopAndLoofOffButton.className = "player-loopoff"
            video.loop = false
        } else {
            loopAndLoofOffButton.className = "player-loop"
            video.loop = true
        }
    });
    startButton.addEventListener("click", function() {
        if (fps == 25) {
            video.currentTime = video.seekable.start(0)
        } else {
            video.currentTime = video.seekable.start(0) + (1/parseFloat(fps))
        }
    }); // 처음으로 이동하는 버튼을 클릭할 때 이벤트. 맨 앞에서 1/fps만큼 이후로 이동해야 함.
    endButton.addEventListener("click", function() {
        if (fps == 60 || fps == 23.976) {
            video.currentTime = video.seekable.end(0)
        } else {
            // 끝으로 이동하는 버튼을 클릭할 때 이벤트. 맨 뒤에서 1/fps만큼 이전으로 이동해야 함.
            video.currentTime = video.seekable.end(0) - (1/parseFloat(fps))
        }
    });
    // 이전 프레임으로 이동하는 버튼을 클릭할 때 이벤트
    beforeFrameButton.addEventListener("click", function() {
        if (fps == 25) {
            // 25fps를 가지고 있는 미디어는 시작 프레임이 2프레임에서 시작된다.
            if (video.currentTime > 0.0) {
                video.currentTime -= (1/parseFloat(fps));
            } else {
                video.currentTime = video.seekable.start(0)
            }
        } else {
            if (video.currentTime > 0.0) {
                video.currentTime -= (1/parseFloat(fps));
            } else {
                video.currentTime = video.seekable.start(0) + (1/parseFloat(fps))
            }
        }
    }); 
    // 다음 프레임으로 이동하는 버튼을 클릭할 때 이벤트
    afterFrameButton.addEventListener("click", function() {
        if (video.currentTime < video.seekable.end(0)) {
            video.currentTime += (1/parseFloat(fps));
        } else {
            video.currentTime = video.seekable.end(0) - (1/parseFloat(fps))
        }
    });

    // video 객체를 생성한다.
    var video = document.createElement('video');
    video.src = "/reviewdata?id=" + id;
    video.autoplay = true;
    video.loop = true;
    video.setAttribute("id", "currentvideo");
    
    // 플레이어창의 배경을 검정으로 한번 채운다.
    playerCtx.fillStyle = "#000000";
    playerCtx.fillRect(0, 0, clientWidth, clientHeight);
    
    // 비디오객체의 메타데이터를 로딩하면 실행할 함수를 설정한다.
    let totalFrame = 0
    let sketchesFrame = [];
    // 기존에 드로잉 되어 있는 데이터를 가지고 온다.
    $.ajax({
        url: "/api/review",
        type: "post",
        data: {
            id: id,
        },
        async: false,
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            for (let i = 0; i < data.sketches.length; i++) {
                sketchesFrame.push(data.sketches[i].frame)
            }
        },
        error: function(){
            return
        }
    })
    
    // 비디오가 로딩되면 메타데이터로 처리할 수 있는 과정을 처리한다.
    video.onloadedmetadata = function() {
        // Draw 캔버스에 프레임 표기 그림을 그린다.
        totalFrame = Math.round(this.duration * parseFloat(fps)) // round로 해야 23.976fps에서 frame 에러가 발생하지 않는다.
        // totalFrame을 표기한다.
        document.getElementById("totalframe").innerHTML = padNumber(totalFrame);
        // 프레임 표기바의 간격을 구하고 global 변수에 저장한다.
        framelineOffset = clientWidth / totalFrame
        
        // 프레임 위치에 해당하는 곳에 회색바로 박스를 그린다.
        for (let i = 0; i < totalFrame; i++) {
            uxCtx.beginPath();
            if (sketchesFrame.includes(i+1)) {                
                uxCtx.strokeStyle = '#FFCD31';
            } else {
                uxCtx.strokeStyle = '#333333';
            }
            uxCtx.lineWidth = 2;
            uxCtx.moveTo(i*framelineOffset + (framelineOffset / 2) , clientHeight - frameLineMarkHeight);
            uxCtx.lineTo(i*framelineOffset + (framelineOffset / 2), clientHeight);
            uxCtx.stroke();
            uxCtx.closePath();
        }
        
        // 재생에 필요한 모든 설정이 완료되면 리뷰 데이터를 플레이시킨다.
        playAndPauseButton.className = "player-pause"
        video.play();
    };
    
    video.addEventListener('play', function () {
        let $this = this; //cache
        (function loop() {
            if (!$this.paused && !$this.ended) {
                renderWidth = ($this.videoWidth * clientHeight) / $this.videoHeight // 실제로 렌더링되는 너비
                renderHeight = ($this.videoHeight * clientWidth) / $this.videoWidth // 실제로 렌더링되는 높이
                if (clientWidth <= renderWidth && renderHeight < clientHeight) {
                    // 가로형: 가로비율이 맞고, 높이가 적을 때
                    let hOffset = (clientHeight - renderHeight) / 2
                    globalReviewRenderWidth = clientWidth
                    globalReviewRenderHeight = renderHeight
                    globalReviewRenderHeightOffset = hOffset
                    playerCtx.drawImage($this, 0, hOffset, clientWidth, renderHeight);
                } else {
                    // 세로형: 가로비율이 작고 높이가 맞을 때
                    let wOffset = (clientWidth - renderWidth) / 2
                    globalReviewRenderWidth = renderWidth
                    globalReviewRenderHeight = clientHeight
                    globalReviewRenderWidthOffset = wOffset
                    playerCtx.drawImage($this, wOffset, 0, renderWidth, clientHeight);
                }
                // fps에 맞게 currentFrame을 드로잉한다.
                let currentFrame = Math.floor($this.currentTime * parseFloat(fps))
                if (currentFrame < totalFrame) {
                    document.getElementById("currentframe").innerHTML = padNumber(currentFrame + 1)
                } else {
                    // 재생이 멈추면 표기는 totalFrame이 되어야 하지만 실제 재생시점은 영상의 마지막이 되어야 한다.
                    document.getElementById("currentframe").innerHTML = padNumber(totalFrame)
                }
                // 커서의 위치를 드로잉 한다.
                aniuxCtx.clearRect(0, 0, clientWidth, clientHeight);
                aniuxCtx.strokeStyle = "#FF0000";
                aniuxCtx.lineWidth = 4;
                aniuxCtx.beginPath();
                aniuxCtx.moveTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight - frameLineMarkHeight);
                aniuxCtx.lineTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight);
                aniuxCtx.stroke();

                // 다음화면 갱신
                setTimeout(loop, 1000 / parseFloat(fps));

            }
        })();
    }, 0);

    video.addEventListener('timeupdate', function () {
        let $this = this; //cache 화
        if (clientWidth <= globalReviewRenderWidth && globalReviewRenderHeight < clientHeight) {
            // 가로형: 가로비율이 맞고, 높이가 적을 때
            let hOffset = (clientHeight - globalReviewRenderHeight) / 2
            playerCtx.drawImage($this, 0, hOffset, clientWidth, globalReviewRenderHeight);
        } else {
            // 세로형: 가로비율이 작고 높이을 꽉 채울 때
            let wOffset = (clientWidth - globalReviewRenderWidth) / 2
            playerCtx.drawImage($this, wOffset, 0, globalReviewRenderWidth, clientHeight);
        }
        // fps에 맞게 currentFrame을 드로잉한다.
        let currentFrame = Math.floor($this.currentTime * parseFloat(fps))
        if (currentFrame < totalFrame) {
            document.getElementById("currentframe").innerHTML = padNumber(currentFrame + 1)
        } else {
            // 재생이 멈추면 표기는 totalFrame이 되어야 하지만 실제 재생시점은 영상의 마지막이 되어야 한다.
            document.getElementById("currentframe").innerHTML = padNumber(totalFrame)
        }
        // 빨간 커서의 위치를 드로잉 한다.
        aniuxCtx.clearRect(0, 0, clientWidth, clientHeight);
        aniuxCtx.strokeStyle = "#FF0000";
        aniuxCtx.lineWidth = 4;
        aniuxCtx.beginPath();
        aniuxCtx.moveTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight - frameLineMarkHeight);
        aniuxCtx.lineTo(currentFrame * framelineOffset + (framelineOffset/2), clientHeight);
        aniuxCtx.stroke();
        
        // 드로잉 프레임은 비디오가 정지될 때만 보여야한다.
        if (video.paused) {
            // 프레임을 이동하면 기존 드로잉이 지워져야 한다.
            removeDrawing()
            // 드로잉이 존재하면 fg 캔버스에 그린다.
            let drawing = new Image()
            let id = document.getElementById("current-review-id").value
            let frame = document.getElementById("currentframe").innerHTML
            let url = `/reviewdrawingdata?id=${id}&frame=${frame}&time=${new Date().getTime()}`
            let http = new XMLHttpRequest();
            http.open("HEAD", url, false)
            http.send()
            if (http.status === 200) {
                let fg = document.getElementById("drawcanvas")
                let fgctx = fg.getContext("2d")
                drawing.src = url
                drawing.onload = function() {
                    fgctx.drawImage(drawing,
                        0, 0, drawing.width, drawing.height,
                        globalReviewRenderWidthOffset, globalReviewRenderHeightOffset, globalReviewRenderWidth, globalReviewRenderHeight
                    );
                };
            } else {
                return
            }
        }
        
    }, 0);
}

function gotoFrame(frame, fps) {
    return ((parseFloat(frame) - (1.0/parseFloat(fps))) / parseFloat(fps))
}

// draw 함수는 x,y 좌표를 받아 그림을 그린다.
function draw(curX, curY) {
    drawCtx.beginPath();
    drawCtx.moveTo(mouseStartX, mouseStartY);
    drawCtx.lineTo(curX, curY);
    drawCtx.stroke();
}

// down 함수는 마우스를 클릭할 때 현재 위치를 마우스의 시작좌표로 설정하고 그림을 그린다고 설정한다.
function down(e) {
    mouseStartX = e.offsetX;
    mouseStartY = e.offsetY;
    drawing = true;
}

// up 함수는 마우스 버튼을 땔 때 그림그리는 모드를 종료한다.
function up(e) {
    drawing = false;
    saveDrawing()
}

// move 함수는 그림을 그리는 상태이고 마우스가 이동할 때 현재 위치에 그림을 그리고 현재위치를 다시 마우스의 시작위치로 바꾼다.
function move(e) {
    if (!drawing) {
        return; // 마우스가 눌러지지 않았으면 리턴
    }
    var curX = e.offsetX, curY = e.offsetY;
    draw(curX, curY);
    mouseStartX = curX;
    mouseStartY = curY;
}

// out 함수는 화면 밖으로 커서가 나가면 그림그리는 모드를 끈다.
function out(e) {
    drawing = false;
}

// changeYellowDrawingFrame 함수는 프레임을 받아서 노란색 드로잉 마커를 체크한다.
function changeYellowDrawingFrame() {
    let uxCanvas = document.getElementById("uxcanvas");
    uxCtx = uxCanvas.getContext("2d")
    currentFrame = parseInt(document.getElementById("currentframe").innerHTML) - 1
    uxCtx.beginPath();
    uxCtx.strokeStyle = '#FFCD31';
    uxCtx.lineWidth = 2;
    uxCtx.moveTo(currentFrame*framelineOffset + (framelineOffset / 2) , globalClientHeight - frameLineMarkHeight);
    uxCtx.lineTo(currentFrame*framelineOffset + (framelineOffset / 2), globalClientHeight);
    uxCtx.stroke();
    uxCtx.closePath();
}

// changeDrawingGrayFrame 함수는 프레임을 받아서 노란색 드로잉 마커를 체크한다.
function changeDrawingGrayFrame() {
    let uxCanvas = document.getElementById("uxcanvas");
    uxCtx = uxCanvas.getContext("2d")
    currentFrame = parseInt(document.getElementById("currentframe").innerHTML) - 1
    uxCtx.beginPath();
    uxCtx.strokeStyle = '#333333';
    uxCtx.lineWidth = 2;
    uxCtx.moveTo(currentFrame*framelineOffset + (framelineOffset / 2) , globalClientHeight - frameLineMarkHeight);
    uxCtx.lineTo(currentFrame*framelineOffset + (framelineOffset / 2), globalClientHeight);
    uxCtx.stroke();
    uxCtx.closePath();
}


// screenshot 함수는 리뷰중인 스크린을 스크린샷 합니다.
function screenshot(filename) {
    let screenshot = document.getElementById("screenshot");
    let fg = document.getElementById("drawcanvas");
    let bg = document.getElementById("player");
    let screenshotctx = screenshot.getContext('2d')
    screenshotctx.drawImage(bg, 0, 0) // 배경에 bg를 그린다.
    screenshotctx.drawImage(fg, 0, 0) // 배경에 fg를 그린다.
    let dataURL = screenshot.toDataURL("image/png")
    let link = document.createElement('a');
    link.href = dataURL;
    let today = new Date();
    let y = today.getFullYear().toString();
    let m = ("0"+(today.getMonth()+1)).slice(-2);
    let d = ("0"+(today.getDate())).slice(-2);
    let hour = ("0"+(today.getHours())).slice(-2);
    let min = ("0"+(today.getMinutes())).slice(-2);
    let sec = ("0"+(today.getSeconds())).slice(-2);
    let timestamp = y + m + d + 'T' + hour + min + sec; // ISO 8601 format 
    let currentFrame = document.getElementById("currentframe").innerHTML + 'f';
    link.download = filename + '_'+ currentFrame + '_' + timestamp + '.png';
    link.setAttribute("type","hidden") // firefox에서는 꼭 DOM구조를 지켜야 한다.
    document.body.appendChild(link); // firefox에서는 꼭 DOM구조를 지켜야 한다.
    link.click();
    link.remove();
    // 스크린샷이 저장되면 기존에 캔버스에 합성된 이미지를 제거한다.
    let playerbox = document.getElementById("playerbox");
    let clientWidth = playerbox.clientWidth
    let clientHeight = playerbox.clientHeight
    screenshotctx.clearRect(0, 0, clientWidth, clientHeight);
    changeYellowDrawingFrame()
}

// saveDrawing 함수는 리뷰스크린에 드로잉된 이미지를 서버에 저장합니다.
function saveDrawing() {
    let id = document.getElementById("current-review-id").value;
    let token = document.getElementById("token").value;
    // Crop Canvas를 생성한다.
    let cropCanvas = document.createElement("canvas");
    cropCanvas.id = "cropCanvas";
    cropCanvas.width = globalReviewRenderWidth;
    cropCanvas.height = globalReviewRenderHeight;
    let ctx = cropCanvas.getContext('2d');
    let drawingCanvas = document.getElementById("drawcanvas")
    ctx.drawImage(drawingCanvas, globalReviewRenderWidthOffset, globalReviewRenderHeightOffset, globalReviewRenderWidth, globalReviewRenderHeight, 0, 0, globalReviewRenderWidth, globalReviewRenderHeight)

    // canvas의 드로잉을 .png 파일로 파일화 한다.
    let fg = cropCanvas.toDataURL("image/png");
    let blobBin = atob(fg.split(',')[1]); // base64 데이터를 바이너리로 변경한다.
    let array = [];
    for (let i = 0; i < blobBin.length; i++) {
        array.push(blobBin.charCodeAt(i));
    }
    let file = new Blob([new Uint8Array(array)], {type: 'image/png'}); // Blob 생성

    let formData = new FormData();
    formData.append("file", file); // .png 추가
    formData.append("id", id)
    formData.append("frame", document.getElementById("currentframe").innerHTML)
    $.ajax({
        url: "/api/uploadreviewdrawing",
        type: "POST",
        enctype: "multipart/form-data",
        processData: false,
        contentType : false,
        cache: false,
        data: formData,
        headers: {
            "Authorization": "Basic "+ token
        },
        success: function(data) {
            changeYellowDrawingFrame()
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

// removeDrawing 함수는 리뷰 스케치를 제거합니다. 프레임을 갱신할 때 사용합니다.
function removeDrawing() {
    // drawcanvas를 지운다.
    let fg = document.getElementById("drawcanvas");
    let fgctx = fg.getContext("2d");
    fgctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
    // screenshot를 지운다.
    let screenshot = document.getElementById("screenshot");
    let screenshotctx = screenshot.getContext("2d");
    screenshotctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
}

// removeDrawingAndData 함수는 리뷰 스케치를 제거하고 서버의 이미지도 제거합니다.
function removeDrawingAndData() {
    // drawcanvas를 지운다.
    let fg = document.getElementById("drawcanvas");
    let fgctx = fg.getContext("2d");
    fgctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
    // screenshot를 지운다.
    let screenshot = document.getElementById("screenshot");
    let screenshotctx = screenshot.getContext("2d");
    screenshotctx.clearRect(0, 0, globalClientWidth, globalClientHeight);
    // 서버에 파일이 존재하면 삭제한다.
    $.ajax({
        url: "/api/rmreviewdrawing",
        type: "post",
        data: {
            id: document.getElementById("current-review-id").value,
            frame: parseInt(document.getElementById("currentframe").innerHTML),
        },
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            // 파일이 잘 삭제되면, 그림이 그려진 프레임의 노란바를 회색바로 변경한다.
            changeDrawingGrayFrame()
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

// copyClipboard 는 value 값을 받아서, 클립보드로 복사하는 기능이다.
function copyClipboard(value) {
    let id = document.createElement("input");   // input요소를 만듬
    id.setAttribute("value", value);            // input요소에 값을 추가
    document.body.appendChild(id);              // body에 요소 추가
    id.select();                                // input요소를 선택
    document.execCommand("copy");               // 복사기능 실행
    document.body.removeChild(id);              // body에 요소 삭제
}

// copyClipboardAndMessage 는 value 값을 받아서, 클립보드로 복사하는 기능이다. 이후 메시지를 출력한다.
function copyClipboardAndMessage(value) {
    let id = document.createElement("input");   // input요소를 만듬
    id.setAttribute("value", value);            // input요소에 값을 추가
    document.body.appendChild(id);              // body에 요소 추가
    id.select();                                // input요소를 선택
    document.execCommand("copy");               // 복사기능 실행
    document.body.removeChild(id);              // body에 요소 삭제
    alert(value + "\n값이 클립보드에 복사되었습니다.");
}

// 리뷰페이지 핫키
// Hotkey: http://gcctech.org/csc/javascript/javascript_keycodes.htm
document.onkeydown = function(e) {
    // 인풋창에서는 화살표를 움직였을 때 페이지가 이동되면 안된다.
    if (event.target.tagName === "INPUT") {
        return
    }
    if (event.target.tagName === "TEXTAREA") {
        return
    }
    if (e.which == 37) { // arrow left
        document.getElementById("player-pause").click();
        document.getElementById("player-left").click();
    } else if (e.which == 39) { // arrow right
        document.getElementById("player-pause").click();
        document.getElementById("player-right").click();
    } else if (e.which == 80 || e.which == 83 || e.which == 32) { // p, s, space
        document.getElementById("player-playandpause").click();
    } else if (e.which == 219) { // [
        document.getElementById("player-pause").click();
        document.getElementById("player-start").click();
    } else if (e.which == 221) { // ]
        document.getElementById("player-pause").click();
        document.getElementById("player-end").click();
    } else if (e.which == 84) { // t
        document.getElementById("player-trash").click();
    } else if (e.which == 67) { // c
        document.getElementById("player-screenshot").click();
    } else if (e.which == 76) { // l
        document.getElementById("player-loopandloopoff").click();
    } else if (e.which == 190) { // .
        document.getElementById("drawing-next").click();
    } else if (e.which == 188) { // ,
        document.getElementById("drawing-prev").click();
    }
};

function rvplay(id) {
    // review id의 데이터를 가지고 path값을 구하고 dilink를 통해 rv player에 연결한다.
    $.ajax({
        url: "/api/review",
        type: "post",
        data: {
            id: id,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let link = document.createElement('a');
            link.href = "dilink://" + data.path;
            document.body.appendChild(link); // firefox에서는 꼭 DOM구조를 지켜야 한다.
            link.click();
            link.remove();
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

// playOriginal 함수는 project와 name을 받아서 해당 아이템의 썸네일동영상을 재생한다. 원본 플레이트 영상을 재생하기 위해서 사용한다.
function playOriginal(project, name) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: `/api2/item?project=${project}&name=${name}`,
        type: "get",
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            let link = document.createElement('a');
            link.href = "dilink://" + data.thummov; // 썸네일 동영상이 보통은 썸네일을 플레이할 때의 Original Plate 이다.
            document.body.appendChild(link); // firefox에서는 꼭 DOM구조를 지켜야 한다.
            link.click();
            link.remove();
        },
        error: function(request,status,error){
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    })
}

function rmUser() {
    let id = document.getElementById("modal-rmuser-id").value
    $.ajax({
        url: "/api2/user?id="+id,
        type: "DELETE",
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value
        },
        dataType: "json",
        success: function(data) {
            location.reload()
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setModalGotoFrame() {
    document.getElementById("modal-gotoframe-frame").value = parseInt(document.getElementById("currentframe").innerHTML);
}

// redirectPage 함수는 page를 받아서 해당 페이지로 리다이렉트 한다.
function redirectPage(page) {
    let href = new URL(window.location.href);
    href.searchParams.set('page', page);
    let url = href.toString();
    window.location.href = url
}

function selectUserID(id) {
    let yellow = "rgb(255, 196, 35)" // 선택된 색상
    let grey = "rgb(167, 165, 157)" // 기본 색상
    if (document.getElementById(id).style.borderColor === yellow) {
        document.getElementById(id).style.borderColor = grey
    } else {
        document.getElementById(id).style.borderColor = yellow
    }
    
}