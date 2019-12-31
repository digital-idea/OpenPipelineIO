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
    } else if (e.ctrlKey && e.altKey &&  e.shiftKey && e.which == 73) {
        selectCheckboxInvert()
    } else if (e.ctrlKey && e.altKey &&  e.shiftKey && e.which == 84) {
        scroll(0,0)
    } else if (e.ctrlKey && e.altKey &&  e.shiftKey && e.which == 77) {
        selectmode()
    } else if (e.ctrlKey && e.altKey &&  e.shiftKey && e.which == 69) {
        OpenEditfolder()
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
            document.getElementById('modal-edittask-id').value=id;
            document.getElementById("modal-edittask-status").value=data.task.status;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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


var multiInput = false;

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
    if (multiInput) {
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
                    let newItem = `<div class="row" id="${data.name}-task-${data.task}">
					<div id="${data.name}-task-${data.task}-status">
						<a class="mt-1 badge badge-assign statusbox">${data.task}</a>
					</div>
					<div id="${data.name}-task-${data.task}-predate"></div>
					<div id="${data.name}-task-${data.task}-date"></div>
					<div id="${data.name}-task-${data.task}-user"></div>
					<div id="${data.name}-task-${data.task}-playbutton"></div>
					<div class="ml-1">
						<span class="add" data-toggle="modal" data-target="#modal-edittask" onclick="
                        setEditTaskModal('${project}', '${data.id}', '${data.task}');
                        ">≡</span>
					</div>
                    </div>`;
                    document.getElementById(`${data.name}-tasks`).innerHTML = newItem + document.getElementById(`${data.name}-tasks`).innerHTML;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                let newItem = `<div class="row" id="${data.name}-task-${data.task}">
					<div id="${data.name}-task-${data.task}-status">
						<a class="mt-1 badge badge-assign statusbox">${data.task}</a>
					</div>
					<div id="${data.name}-task-${data.task}-predate"></div>
					<div id="${data.name}-task-${data.task}-date"></div>
					<div id="${data.name}-task-${data.task}-user"></div>
					<div id="${data.name}-task-${data.task}-playbutton"></div>
					<div class="ml-1">
						<span class="add" data-toggle="modal" data-target="#modal-edittask" onclick="
                        setEditTaskModal('${project}', '${data.id}', '${data.task}');
                        ">≡</span>
					</div>
				</div>`;
                document.getElementById(`${data.name}-tasks`).innerHTML = newItem + document.getElementById(`${data.name}-tasks`).innerHTML;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function rmTask(project, id, task) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    status: 'false',
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.name}-task-${data.task}`).remove();
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                status: 'false',
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.name}-task-${data.task}`).remove();
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubTask(project, id, task, userid) {
    let token = document.getElementById("token").value;
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraPubPath(project, id, path, userid) {
    let token = document.getElementById("token").value;
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setCameraProjection(project, id, userid) {
    let token = document.getElementById("token").value;
    let projection = document.getElementById("cameraoption-projection").checked;
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setNoteModal(project, id) {
    document.getElementById("modal-setnote-project").value = project;
    document.getElementById("modal-setnote-id").value = id;
    document.getElementById("modal-setnote-text").value = "";
    document.getElementById("modal-setnote-title").innerHTML = "Rm Tag" + multiInputTitle(id);
}

function setNote(project, id, text) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    name: id2name(currentID),
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
                        document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    } else {
                        // note-{{.Name}} 내부 내용에 추가한다.
                        document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.name).innerHTML;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                name: id2name(id),
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
                    document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                } else {
                    // note-{{.Name}} 내부 내용에 추가한다.
                    document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.name).innerHTML;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function addComment(project, id, text, media) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    let newComment = `<div id="comment-${data.name}-${data.date}">
                    <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
                    <span class="remove" data-toggle="modal" data-target="#rmcomment" onclick="setModal('rm-comment-name', '${data.name}');setModal('rm-comment-time', '${data.date}');setModal('rm-comment-text', '${data.text}');setModal('rm-comment-userid', '${data.userid}')">×</span>
                    <br><small class="text-white">${body}</small>`
                    if (data.media != "") {
                        if (data.media.includes("http")) {
                            newComment += `<br><a href="${data.media}" class="link">∞</a>`
                        } else {
                            newComment += `<br><a href="dilink://${data.media}" class="link">∞</a>`
                        }
                    }
                    newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
                    document.getElementById("comments-"+data.name).innerHTML = newComment + document.getElementById("comments-"+data.name).innerHTML;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                // comments-{{.Name}} 내부 내용에 추가한다.
                let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                let newComment = `<div id="comment-${data.name}-${data.date}">
                <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
                <span class="remove" data-toggle="modal" data-target="#rmcomment" onclick="setModal('rm-comment-name', '${data.name}');setModal('rm-comment-time', '${data.date}');setModal('rm-comment-text', '${data.text}');setModal('rm-comment-userid', '${data.userid}')">×</span>
                <br><small class="text-white">${body}</small>`
                if (data.media != "") {
                    if (data.media.includes("http")) {
                        newComment += `<br><a href="${data.media}" class="link">∞</a>`
                    } else {
                        newComment += `<br><a href="dilink://${data.media}" class="link">∞</a>`
                    }
                }
                newComment += `<hr class="my-1 p-0 m-0 divider"></hr></div>`
                document.getElementById("comments-"+data.name).innerHTML = newComment + document.getElementById("comments-"+data.name).innerHTML;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function rmComment(project, id, date, userid) {
    let token = document.getElementById("token").value;
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
            document.getElementById(`comment-${data.name}-${data.date}`).remove();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function addReference(project, id, title, path) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("token").value;
    if (multiInput) {
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
                        <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#rmreference" onclick="setModal('rm-reference-title', '' );setModal('rm-reference-name', '${data.name}')">－</span>
                        `
                    } else {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                    <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#rmreference" onclick="setModal('rm-reference-title', '' );setModal('rm-reference-name', '${data.name}')">－</span>
                    `
                } else {
                    document.getElementById("reference-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}


function rmReference(project, id, title) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    if (document.getElementById(`references-${data.ame}`).childElementCount > 0) {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#rmreference" onclick="setModal('rm-reference-title', '' );setModal('rm-reference-name', '${data.name}')">－</span>
                        `
                    } else {
                        document.getElementById("reference-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                    <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#rmreference" onclick="setModal('rm-reference-title', '' );setModal('rm-reference-name', '${data.name}')">－</span>
                    `
                } else {
                    document.getElementById("reference-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#addreference" onclick="setModal('add-reference-title', '' );setModal('add-reference-path', '' );setModal('add-reference-name', '${data.name}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            document.getElementById("button-thumbplay-"+data.name).innerHTML = `<a href="dilink://${data.path}" class="play">▶</a>`;
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTaskMov(project, id, task, mov) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    $.ajax({
        url: "/api/settaskmov",
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
                document.getElementById(`${data.name}-task-${data.task}-playbutton`).innerHTML = "";
            } else {
                document.getElementById(`${data.name}-task-${data.task}-playbutton`).innerHTML = `<a class="mt-1 ml-1 badge badge-light" href="dilink://${data.mov}">▶</a>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setTaskDue(project, id, task, due) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskUser(project, id, task, user) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    if (data.mov === "") {
                        document.getElementById(`${data.name}-task-${data.task}-user`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.name}-task-${data.task}-user`).innerHTML = `<span class="mt-1 ml-1 badge badge-light">${data.username}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                if (data.mov === "") {
                    document.getElementById(`${data.name}-task-${data.task}-user`).innerHTML = "";
                } else {
                    document.getElementById(`${data.name}-task-${data.task}-user`).innerHTML = `<span class="mt-1 ml-1 badge badge-light">${data.username}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskStatus(project, id, task, status) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    name: id2name(id),
                    task: task,
                    status: status,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`${data.name}-task-${data.task}-status`).innerHTML = `<a class="mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</a>`;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskstatus",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
                task: task,
                status: status,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`${data.name}-task-${data.task}-status`).innerHTML = `<a class="mt-1 badge badge-${data.status} statusbox" title="${data.status}">${data.task}</a>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskDate(project, id, task, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    name: id2name(id),
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
                        document.getElementById(`${data.name}-task-${data.task}-date`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.name}-task-${data.task}-date`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.shortdate}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskdate",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
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
                    document.getElementById(`${data.name}-task-${data.task}-date`).innerHTML = "";
                } else {
                    document.getElementById(`${data.name}-task-${data.task}-date`).innerHTML = `<span class="mt-1 ml-1 badge badge-darkmode">${data.shortdate}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskStartdate(project, id, task, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    name: id2name(id),
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskstartdate",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskUserNote(project, id, task, usernote) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}



function setTaskPredate(project, id, task, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    name: id2name(id),
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
                        document.getElementById(`${data.name}-task-${data.task}-predate`).innerHTML = "";
                    } else {
                        document.getElementById(`${data.name}-task-${data.task}-predate`).innerHTML = `<span class="mt-1 ml-1 badge badge-outline-darkmode">${data.shortdate}</span>`;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/settaskpredate",
            type: "post",
            data: {
                project: project,
                name: id2name(id),
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
                    document.getElementById(`${data.name}-task-${data.task}-predate`).innerHTML = "";
                } else {
                    document.getElementById(`${data.name}-task-${data.task}-predate`).innerHTML = `<span class="mt-1 ml-1 badge badge-outline-darkmode">${data.shortdate}</span>`;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setDeadline2D(project, id, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setDeadline3D(project, id, date) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    console.log(id);
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setShottype(project, id) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let e = document.getElementById("modal-shottype-type");
    let shottype = e.options[e.selectedIndex].value;
    
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAssettype(project, id) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    let types = document.getElementById("modal-assettype-type");
    let assettype = types.options[types.selectedIndex].value;
    console.log(id)
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setAddTagModal(project, id) {
    document.getElementById("modal-addtag-project").value = project;
    document.getElementById("modal-addtag-id").value = id;
    document.getElementById("modal-addtag-text").value = "";
    document.getElementById("modal-addtag-title").innerHTML = "Add Tag" + multiInputTitle(id);
}

function addTag(project, id, tag) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    name: id2name(id),
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
                    source = `<div id="tag-${data.name}-${data.tag}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.tag}</a></div>`;
                    document.getElementById("tags-"+data.name).innerHTML = document.getElementById("tags-"+data.name).innerHTML + source;
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`tags-${data.name}`).childElementCount > 0) {
                        document.getElementById("tag-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                        `
                    } else {
                        document.getElementById("tag-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        `
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/addtag",
            type: "post",
            
            data: {
                project: project,
                name: id2name(id),
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
                let source = `<div id="tag-${data.name}-${data.tag}"><a href="${url}" class="badge badge-outline-darkmode ml-1" alt="${data.userid}" title="${data.userid}">${data.tag}</a></div>`;
                document.getElementById("tags-"+data.name).innerHTML = document.getElementById("tags-"+data.name).innerHTML + source;
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`tags-${data.name}`).childElementCount > 0) {
                    document.getElementById("tag-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                    `
                } else {
                    document.getElementById("tag-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
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
    if (multiInput) {
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
                    name: id2name(id),
                    tag: tag,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    document.getElementById(`tag-${data.name}-${data.tag}`).remove();
                    // 요소갯수에 따라 버튼을 설정한다.
                    if (document.getElementById(`tags-${data.name}`).childElementCount > 0) {
                        document.getElementById("tag-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                        `;
                    } else {
                        document.getElementById("tag-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                        `;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/rmtag",
            type: "post",
            
            data: {
                project: project,
                name: id2name(id),
                tag: tag,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                document.getElementById(`tag-${data.name}-${data.tag}`).remove();
                // 요소갯수에 따라 버튼을 설정한다.
                if (document.getElementById(`tags-${data.name}`).childElementCount > 0) {
                    document.getElementById("tag-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#modal-rmtag" onclick="setRmTagModal('${data.project}','${data.id}')">－</span>
                    `;
                } else {
                    document.getElementById("tag-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#modal-addtag" onclick="setAddTagModal('${data.project}','${data.id}')">＋</span>
                    `;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function selectCheckbox() {
    let checknum = 0;
    var cboxes = document.getElementsByName('selectID');
    for (let i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === true) {
            checknum += 1
        }
    }
    if (checknum === 0) {
        multiInput = false;
        return
    }
    multiInput = true;
}

function selectCheckboxAll() {
    multiInput = true;
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = true;
    }
}

function selectCheckboxNone() {
    multiInput = false;
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = false;
    }
}

function selectCheckboxInvert() {
    var cboxes = document.getElementsByName('selectID');
    if (cboxes.length > 0) {
        multiInput = true;
    }
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            cboxes[i].checked = true;
        } else {
            cboxes[i].checked = false;
        }
    }
}

function setTaskLevel(project, id, task, level) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setObjectID(project, id, innum, outnum, userid) {
    let token = document.getElementById("token").value;
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}


function setPlatesize(project, id, size) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setUndistortionsize(project, id, size) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function setRendersize(project, id, size) {
    let token = document.getElementById("token").value;
    let userid = document.getElementById("userid").value;
    if (multiInput) {
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function CurrentProject() {
    let e = document.getElementById("Project");
    return e.options[e.selectedIndex].value;
}

function rmItem() {
    let project = CurrentProject()
    let token = document.getElementById("token").value;
    if (multiInput) {
        let cboxes = document.getElementsByName('selectID');
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
            alert("status:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
        /*DIV 하나를 생성한다.*/
        a = document.createElement("DIV");
        a.setAttribute("id", this.id + "autocomplete-list");
        a.setAttribute("class", "autocomplete-items");
        /*위에서 생성한 검색창을 부모에 붙힌다.*/
        this.parentNode.appendChild(a);
        /*각각의 아이템을 순환한다.*/
        for (i = 0; i < arr.length; i++) {
          /*검색어가 아이템에 포함되어 있다면, div를 생성한다.*/
          if (arr[i].includes(val)) {
            /*create a DIV element for each matching element:*/
            b = document.createElement("DIV");
            /*make the matching letters bold:*/
            b.innerHTML = "<strong>" + arr[i].substr(0, val.length) + "</strong>";
            b.innerHTML += arr[i].substr(val.length);
            /*insert a input field that will hold the current array item's value:*/
            b.innerHTML += "<input type='hidden' value='" + arr[i] + "'>";
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
                let tasks = data["Tasksettings"];
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
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                let tasks = data["Tasksettings"]
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
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                let tasks = data["Tasksettings"];
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
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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
                let tasks = data["Tasksettings"]
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
                alert("status:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
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