// modal이 뜨면 오토포커스가 되어야 한다.
$('#addcomment').on('shown.bs.modal', function () {
    $('#add-comment-text').trigger('focus')
})
$('#setnote').on('shown.bs.modal', function () {
    $('#set-note-text').trigger('focus')
})
$('#addsource').on('shown.bs.modal', function () {
    $('#add-source-title').trigger('focus')
})
$('#setrnum').on('shown.bs.modal', function () {
    $('#set-rnum-text').trigger('focus')
})
$('#addtag').on('shown.bs.modal', function () {
    $('#add-tag-text').trigger('focus')
})
$('#rmtag').on('shown.bs.modal', function () {
    $('#rm-tag-text').trigger('focus')
})
$('#deadline2d').on('shown.bs.modal', function () {
    $('#deadline2d-date').trigger('focus')
})
$('#deadline3d').on('shown.bs.modal', function () {
    $('#deadline3d-date').trigger('focus')
})

// Hotkey
document.onkeyup = function(e) {
    if (e.ctrlKey && e.shiftKey && e.which == 65) {
        selectCheckboxAll()
    } else if (e.ctrlKey && e.shiftKey && e.which == 68) {
        selectCheckboxNone()
    } else if (e.ctrlKey && e.shiftKey && e.which == 73) {
        selectCheckboxInvert()
    } else if (e.ctrlKey && e.shiftKey && e.which == 84) {
        scroll(0,0)
    } else if (e.ctrlKey && e.shiftKey && e.which == 77) {
        selectmode()
    }
};

function id2name(id) {
    l = id.split("_");
    l.pop()
    return l.join("_")
}

// setModal 함수는 modalID와 value를 받아서 modal에 셋팅한다.
function setModal(modalID, value) {
    document.getElementById(modalID).value=value;
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
						<span class="add" data-toggle="modal" data-target="#edittask" onclick="
							setModal('edittask-startdate', '');
							setModal('edittask-due', '');
							setModal('edittask-level', '');
							setModal('edittask-predate', '');
							setModal('edittask-date', '');
							setModal('edittask-task', '${data.task}' );
							setModal('edittask-path', '' );
							setModal('edittask-usernote', '' );
							setModal('edittask-name', '${data.name}');
							setModal('edittask-user', '');
							setModal('edittask-userid', '${data.userid}')">≡</span>
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
						<span class="add" data-toggle="modal" data-target="#edittask" onclick="
							setModal('edittask-startdate', '');
							setModal('edittask-due', '');
							setModal('edittask-level', '');
							setModal('edittask-predate', '');
							setModal('edittask-date', '');
							setModal('edittask-task', '${data.task}' );
							setModal('edittask-path', '' );
							setModal('edittask-usernote', '' );
							setModal('edittask-name', '${data.name}');
							setModal('edittask-user', '');
							setModal('edittask-userid', '${data.userid}')">≡</span>
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

function setFrame(mode, project, id, frame, userid) {
    let token = document.getElementById("token").value;
    $.ajax({
        url: "/api/" + mode,
        type: "post",
        data: {
            project: project,
            name: id2name(id),
            frame: frame,
            userid: userid,
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

function setNote(project, id, text, userid) {
    let token = document.getElementById("token").value;
    if (multiInput) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let currentID = cboxes[i].getAttribute("id");
            let overwrite = document.getElementById("set-note-overwrite").checked;
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
        let overwrite = document.getElementById("set-note-overwrite").checked;
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


function addSource(project, id, title, path, userid) {
    let token = document.getElementById("token").value;
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
                        <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#rmsource" onclick="setModal('rm-source-title', '' );setModal('rm-source-name', '${data.name}');setModal('rm-source-userid', '${data.userid}')">－</span>
                        `
                    } else {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
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
                    <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#rmsource" onclick="setModal('rm-source-title', '' );setModal('rm-source-name', '${data.name}');setModal('rm-source-userid', '${data.userid}')">－</span>
                    `
                } else {
                    document.getElementById("source-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}


function rmSource(project, id, title, userid) {
    let token = document.getElementById("token").value;
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
                    if (document.getElementById(`sources-${data.ame}`).childElementCount > 0) {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#rmsource" onclick="setModal('rm-source-title', '' );setModal('rm-source-name', '${data.name}');setModal('rm-source-userid', '${data.userid}')">－</span>
                        `
                    } else {
                        document.getElementById("source-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
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
                    <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#rmsource" onclick="setModal('rm-source-title', '' );setModal('rm-source-name', '${data.name}');setModal('rm-source-userid', '${data.userid}')">－</span>
                    `
                } else {
                    document.getElementById("source-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#addsource" onclick="setModal('add-source-title', '' );setModal('add-source-path', '' );setModal('add-source-name', '${data.name}');setModal('add-source-userid', '${data.userid}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setThummov(project, id, path, userid) {
    let token = document.getElementById("token").value;
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

function setBeforemov(project, id, path, userid) {
    let token = document.getElementById("token").value;
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

function setAftermov(project, id, path, userid) {
    let token = document.getElementById("token").value;
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


function setRetimeplate(project, id, path, userid) {
    let token = document.getElementById("token").value;
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

function setOCIOcc(project, id, path, userid) {
    let token = document.getElementById("token").value;
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

function setRollmedia(project, id, rollmedia, userid, token) {
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

function setTaskMov(project, id, task, mov, userid) {
    let token = document.getElementById("token").value;
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

function setTaskDue(project, id, task, due, userid) {
    let token = document.getElementById("token").value;
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

function setTaskUser(project, id, task, user, userid) {
    let token = document.getElementById("token").value;
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

function setTaskStatus(project, id, task, status, userid) {
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

function setTaskDate(project, id, task, date, userid) {
    let token = document.getElementById("token").value;
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

function setTaskStartdate(project, id, task, date, userid) {
    let token = document.getElementById("token").value;
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

function setTaskUserNote(project, id, task, usernote, userid) {
    let token = document.getElementById("token").value;
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



function setTaskPredate(project, id, task, date, userid) {
    let token = document.getElementById("token").value;
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


function setDeadline2D(project, id, date, userid) {
    let token = document.getElementById("token").value;
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
                    document.getElementById("deadline2d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline2d" onclick="setModal('deadline2d-date', '${data.date}');setModal('deadline2d-name', '${data.name}');setModal('deadline2d-userid', '${data.userid}')">2D:${data.shortdate}</span>`;
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
                document.getElementById("deadline2d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline2d" onclick="setModal('deadline2d-date', '${data.date}');setModal('deadline2d-name', '${data.name}');setModal('deadline2d-userid', '${data.userid}')">2D:${data.shortdate}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setDeadline3D(project, id, date, userid) {
    let token = document.getElementById("token").value;
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
                    document.getElementById("deadline3d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline3d" onclick="setModal('deadline3d-date', '${data.date}');setModal('deadline3d-name', '${data.name}');setModal('deadline3d-userid', '${data.userid}')">3D:${data.shortdate}</span>`;
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
                document.getElementById("deadline3d-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#deadline3d" onclick="setModal('deadline3d-date', '${data.date}');setModal('deadline3d-name', '${data.name}');setModal('deadline3d-userid', '${data.userid}')">3D:${data.shortdate}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setShottype(project, id, userid) {
    let token = document.getElementById("token").value;
    let shottypes = document.getElementsByName('shottype');
    let shottype = "";
    for (var i = 0, length = shottypes.length; i < length; i++) {
        if (shottypes[i].checked) {
            shottype = shottypes[i].value
            break;
        }
    }
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
                    document.getElementById("shottype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#shottype" onclick="setModal('shottype-name', '${data.name}');setModal('shottype-userid', '${data.userid}')">${data.type}</span>`;
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
                document.getElementById("shottype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#shottype" onclick="setModal('shottype-name', '${data.name}');setModal('shottype-userid', '${data.userid}')">${data.type}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setAssettype(project, id, userid) {
    let token = document.getElementById("token").value;
    let types = document.getElementById("assettypes");
    let assettype = types.options[types.selectedIndex].value;
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
                    document.getElementById("assettype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#assettype" onclick="setModal('assettype-name', '${data.name}');setModal('assettype-userid', '${data.userid}')">${data.type}</span>`;
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
                document.getElementById("assettype-"+data.name).innerHTML = `<span class="badge badge-light ml-1" data-toggle="modal" data-target="#assettype" onclick="setModal('assettype-name', '${data.name}');setModal('assettype-userid', '${data.userid}')">${data.type}</span>`;
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

function setRnum(project, id, rnum, userid) {
    let token = document.getElementById("token").value;
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
                document.getElementById("rnum-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#setrnum" onclick="setModal('set-rnum-text', '${data.rnum}' );setModal('set-rnum-name', '${data.name}');setModal('set-rnum-userid', '${data.userid}')"{{end}}>${data.rnum}</span>`;
            } else {
                document.getElementById("rnum-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#setrnum" onclick="setModal('set-rnum-text', '${data.rnum}' );setModal('set-rnum-name', '${data.name}');setModal('set-rnum-userid', '${data.userid}')"{{end}}>no rnum</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
        }
    });
}

function addTag(project, id, tag, userid) {
    let token = document.getElementById("token").value;
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
                        <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#rmtag" onclick="setModal('rm-tag-text', '' );setModal('rm-tag-name', '${data.name}');setModal('rm-tag-userid', '${data.userid}')">－</span>
                        `
                    } else {
                        document.getElementById("tag-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
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
                    <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#rmtag" onclick="setModal('rm-tag-text', '' );setModal('rm-tag-name', '${data.name}');setModal('rm-tag-userid', '${data.userid}')">－</span>
                    `
                } else {
                    document.getElementById("tag-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
                    `
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function rmTag(project, id, tag, userid) {
    let token = document.getElementById("token").value;
    if (multiInput) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let name = cboxes[i].getAttribute("id");
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
                        <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
                        <span class="remove ml-0" data-toggle="modal" data-target="#rmtag" onclick="setModal('rm-tag-text', '' );setModal('rm-tag-name', '${data.name}');setModal('rm-tag-userid', '${data.userid}')">－</span>
                        `;
                    } else {
                        document.getElementById("tag-button-"+data.name).innerHTML = `
                        <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
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
                    <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
                    <span class="remove ml-0" data-toggle="modal" data-target="#rmtag" onclick="setModal('rm-tag-text', '' );setModal('rm-tag-name', '${data.name}');setModal('rm-tag-userid', '${data.userid}')">－</span>
                    `;
                } else {
                    document.getElementById("tag-button-"+data.name).innerHTML = `
                    <span class="add ml-1" data-toggle="modal" data-target="#addtag" onclick="setModal('add-tag-text', '' );setModal('add-tag-name', '${data.name}');setModal('add-tag-userid', '${data.userid}')">＋</span>
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
    var cboxes = document.getElementsByName('selectID');
    if (cboxes.length > 0) {
        multiInput = true;
    }
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

function setTaskLevel(project, id, task, level, userid) {
    let token = document.getElementById("token").value;
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



function setPlatesize(project, id, size, userid) {
    let token = document.getElementById("token").value;
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
                    document.getElementById("platesize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-platesize" onclick="setModal('platesize', '${data.size}');setModal('platesize-name', '${data.name}');setModal('platesize-userid', '${data.userid}')">S: ${data.size}</span>`;
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
                document.getElementById("platesize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-platesize" onclick="setModal('platesize', '${data.size}');setModal('platesize-name', '${data.name}');setModal('platesize-userid', '${data.userid}')">S: ${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setUndistortionsize(project, id, size, userid) {
    let token = document.getElementById("token").value;
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
                    document.getElementById("undistortionsize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-undistortionsize" onclick="setModal('undistortionsize', '${data.size}');setModal('undistortionsize-name', '${data.name}');setModal('undistortionsize-userid', '${data.userid}')">U: ${data.size}</span>`;
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
                document.getElementById("undistortionsize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-undistortionsize" onclick="setModal('undistortionsize', '${data.size}');setModal('undistortionsize-name', '${data.name}');setModal('undistortionsize-userid', '${data.userid}')">U: ${data.size}</span>`;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+request.responseText+"\n"+"error:"+error);
            }
        });
    }
}

function setRendersize(project, id, size, userid) {
    let token = document.getElementById("token").value;
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
                    document.getElementById("rendersize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setModal('rendersize', '${data.size}');setModal('rendersize-name', '${data.name}');setModal('rendersize-userid', '${data.userid}')">R: ${data.size}</span>`;
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
                document.getElementById("rendersize-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#modal-rendersize" onclick="setModal('rendersize', '${data.size}');setModal('rendersize-name', '${data.name}');setModal('rendersize-userid', '${data.userid}')">R: ${data.size}</span>`;
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

autocomplete(document.getElementById("edittask-user"));

function inputAddTasksetting(type) {
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
                let addtasks = document.getElementById('addtask-taskname');
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
                let addtasks = document.getElementById('addtask-taskname');
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

function inputRmTasksetting(type) {
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
                let rmtasks = document.getElementById('rmtask-taskname');
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
                let rmtasks = document.getElementById('rmtask-taskname');
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