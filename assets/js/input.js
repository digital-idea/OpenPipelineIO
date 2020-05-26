// modal이 뜨면 오토포커스가 되어야 한다.
$('#modal-addcomment').on('shown.bs.modal', function () {
    $('#modal-addcomment-text').trigger('focus')
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
            document.getElementById('modal-edittask-due').value=data.task.due;
            document.getElementById('modal-edittask-level').value=data.task.tasklevel;
            document.getElementById('modal-edittask-task').value=data.task.title;
            document.getElementById('modal-edittask-path').value=data.task.mov;
            document.getElementById('modal-edittask-usernote').value=data.task.usernote;
            document.getElementById('modal-edittask-user').value=data.task.user;
            document.getElementById('modal-edittask-id').value=data.id;
            // ver2로 검색하면 modal-edittask-status 가 존재하지 않을 수 있다.
            try {
                document.getElementById("modal-edittask-status").value=data.task.status;
            }
            catch(err) {
                console.log(err);
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
    let userid = document.getElementById("userid").value;
    if (isMultiInput()) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let id = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/setassigntask",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    task: task,
                    status: 'true',
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    let newItem = `<div class="row" id="${data.id}-task-${data.task}">
					<div id="${data.id}-task-${data.task}-status">
						<span class="finger mt-1 badge badge-assign statusbox">${data.task}</span>
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
            url: "/api/setassigntask",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                task: task,
                status: 'true',
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                let newItem = `<div class="row" id="${data.id}-task-${data.task}">
					<div id="${data.id}-task-${data.task}-status">
						<span class="finger mt-1 badge badge-assign statusbox">${data.task}</span>
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
    let userid = document.getElementById("userid").value;
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

function setCameraPubTask(project, id, task) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setcamerapubtask",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            task: task,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            document.getElementById("campubtask-"+data.name).innerHTML = `<span class="text-badge ml-1">Pub-${data.task}</span>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubPath(project, id, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/setcamerapubpath",
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
            document.getElementById("campubpath-"+data.name).innerHTML = `, <a href="dilink://${data.path}" class="text-badge ml-1">${data.path}</a>`;
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

function setCameraProjection(project, id) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let projection = document.getElementById("modal-cameraoption-projection").checked;
    $.ajax({
        url: "/api/setcameraprojection",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            projection: projection,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.projection === true) {
                document.getElementById("camprojection-"+data.name).innerHTML = `<span class="text-badge ml-1">, Projection</span>`;
            } else {
                document.getElementById("camprojection-"+data.name).innerHTML = "";
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

function setAddPublishModal(project, name, task) {
    document.getElementById("modal-addpublish-project").value = project
    document.getElementById("modal-addpublish-name").value = name
    document.getElementById("modal-addpublish-task").value = task
}

function setAddCommentModal(project, id) {
    document.getElementById("modal-addcomment-project").value = project;
    document.getElementById("modal-addcomment-id").value = id;
    document.getElementById("modal-addcomment-media").value = "";
    document.getElementById("modal-addcomment-title").innerHTML = "Add Comment" + multiInputTitle(id);
}

function addComment(project, id, text, media) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
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
                    <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}', '${data.media}')">≡</span>
                    <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}')">×</span>
                    <br><small class="text-warning">${body}</small>`
                    if (data.media != "") {
                        if (data.media.includes("http")) {
                            newComment += `<br><a href="${data.media}" class="link">∞</a>`
                        } else {
                            newComment += `<br><a href="dilink://${data.media}" class="link">∞</a>`
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
                <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}', '${data.media}')">≡</span>
                <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.project}', '${data.id}', '${data.date}', '${data.text}')">×</span>
                <br><small class="text-warning">${body}</small>`
                if (data.media != "") {
                    if (data.media.includes("http")) {
                        newComment += `<br><a href="${data.media}" class="link">∞</a>`
                    } else {
                        newComment += `<br><a href="dilink://${data.media}" class="link">∞</a>`
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

function editComment(project, id, time, text, media) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/editcomment",
        type: "post",
        data: {
            project: project,
            id: id,
            time: time,
            text: text,
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
            <span class="edit" data-toggle="modal" data-target="#modal-editcomment" onclick="setEditCommentModal('${data.project}', '${data.id}', '${data.time}', '${data.text}', '${data.media}')">≡</span>
            <span class="remove" data-toggle="modal" data-target="#modal-rmcomment" onclick="setRmCommentModal('${data.project}', '${data.id}', '${data.time}', '${data.text}')">×</span>
            <br><small class="text-warning">${body}</small>`
            if (data.media != "") {
                if (data.media.includes("http")) {
                    newComment += `<br><a href="${data.media}" class="link">∞</a>`
                } else {
                    newComment += `<br><a href="dilink://${data.media}" class="link">∞</a>`
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

function setRmPublishKeyModal(project, id, task, key) {
    document.getElementById("modal-rmpublishkey-project").value = project;
    document.getElementById("modal-rmpublishkey-id").value = id;
    document.getElementById("modal-rmpublishkey-task").value = task;
    document.getElementById("modal-rmpublishkey-key").value = key;
    document.getElementById("modal-rmpublishkey-title").innerHTML = "Rm Publish Key" + multiInputTitle(id);
}

function setRmPublishModal(project, id, task, key, time) {
    document.getElementById("modal-rmpublish-project").value = project
    document.getElementById("modal-rmpublish-id").value = id
    document.getElementById("modal-rmpublish-task").value = task
    document.getElementById("modal-rmpublish-key").value = key
    document.getElementById("modal-rmpublish-time").value = time
}

function setPublishModal(project, id, task, key, time) {
    document.getElementById("modal-setpublish-project").value = project;
    document.getElementById("modal-setpublish-id").value = id;
    document.getElementById("modal-setpublish-task").value = task;
    document.getElementById("modal-setpublish-key").value = key;
    document.getElementById("modal-setpublish-time").value = time;
    document.getElementById("modal-setpublish-message").innerHTML = `
        ${project} > ${id} > ${task} > ${key}<br>
        ${time} 에 Publish 된 데이터를<br>
        타팀의 아티스트가 사용해야 할 버전으로 설정하시겠습니까?
    `
}

function setEditCommentModal(project, id, time, text, media) {
    document.getElementById("modal-editcomment-project").value = project;
    document.getElementById("modal-editcomment-id").value = id;
    document.getElementById("modal-editcomment-time").value = time;
    document.getElementById("modal-editcomment-text").value = text;
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
                let text = document.createElement("small")
                text.setAttribute("class","text-white")
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
                    link.setAttribute("class","link")
                    link.innerHTML = "∞"
                    cmt.append(link)
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

function rmComment(project, id, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/rmcomment",
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            date: date,
            userid: userid
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
            time: document.getElementById('modal-rmpublish-time').value
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
                        source = `<div id="source-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
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
                    source = `<div id="source-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
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
                        ref = `<div id="reference-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
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
                    ref = `<div id="reference-${data.name}-${data.title}"><a href="dilink://${data.path}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.title}</a></div>`;
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

function setThummov(project, id, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
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

function setBeforemov(project, id, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
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

function setAftermov(project, id, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
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
            console.log(data);
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
            document.getElementById('modal-iteminfo-thummov').value = data.thummov;
            document.getElementById('modal-iteminfo-beforemov').value = data.beforemov;
            document.getElementById('modal-iteminfo-aftermov').value = data.aftermov;
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

function setTaskDue(project, id, task, due) {
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
                url: "/api/settaskdue",
                type: "post",
                data: {
                    project: project,
                    name: id2name(id),
                    task: task,
                    due: due,
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
    } else {
        $.ajax({
            url: "/api/settaskdue",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                task: task,
                due: due,
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
}

function setTaskUser(project, id, task, user) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
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
                    let url = `/inputmode?project=${data.project}&searchword=assettags:${data.type}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&template=index2&task=`;
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
                let url = `/inputmode?project=${data.project}&searchword=assettags:${data.type}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&template=index2&task=`;
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
                    let url = `/inputmode?project=${data.project}&searchword=tag:${data.tag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&template=index2&task=`
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
                let url = `/inputmode?project=${data.project}&searchword=tag:${data.tag}&sortkey=slug&sortkey=slug&assign=true&ready=true&wip=true&confirm=true&done=false&omit=false&hold=false&out=false&none=false&template=index2&task=`
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
                    console.log(data);
                    document.getElementById(`item-${data.id}`).remove();
                    
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
        var a, b, i, val = this.value;
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

autocomplete(document.getElementById("modal-edittask-user"));

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

// When the user scrolls down 20px from the top of the document, show the button
window.onscroll = function() {scrollFunction()};

function scrollFunction() {
    let topbtn = document.getElementById("topbtn");
    
	if (document.body.scrollTop > 20 || document.documentElement.scrollTop > 20) {
        topbtn.style.display = "block";
	} else {
        topbtn.style.display = "none";
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

function setPublish(project, id, task, key, createtime) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/setpublishstatus",
        type: "post",
        data: {
            project: project,
            id: id,
            task: task,
            key: key,
            createtime: createtime,
            status: "usethis",
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
    let path = document.getElementById('modal-addpublish-path').value
    let status = document.getElementById('modal-addpublish-status').value
    let subject = document.getElementById('modal-addpublish-subject').value
    let mainversion = document.getElementById('modal-addpublish-mainversion').value
    let subversion = document.getElementById('modal-addpublish-subversion').value
    let kindofusd = document.getElementById('modal-addpublish-kindofusd').value
    $.ajax({
        url: "/api/addpublish",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            key: key,
            path: path,
            status: status,
            subject: subject,
            mainversion: mainversion,
            subversion: subversion,
            kindofusd: kindofusd,
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