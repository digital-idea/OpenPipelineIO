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



var multiInput = false;

function removeWhiteSpace(event) {
	event.value = event.value.replace(/ /g, '');
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



function setAssignTask(project, name, task, status, token) {
    $.ajax({
        url: "/api/setassigntask",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            status: status,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setFrame(mode, project, name, frame, token) {
    $.ajax({
        url: "/api/" + mode,
        type: "post",
        data: {
            project: project,
            name: name,
            frame: frame,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}


function setScanTimecodeIn(project, name, timecode, token) {
    $.ajax({
        url: "/api/setscantimecodein",
        type: "post",
        data: {
            project: project,
            name: name,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setScanTimecodeOut(project, name, timecode, token) {
    $.ajax({
        url: "/api/setscantimecodeout",
        type: "post",
        data: {
            project: project,
            name: name,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setJustTimecodeIn(project, name, timecode, token) {
    $.ajax({
        url: "/api/setjusttimecodein",
        type: "post",
        data: {
            project: project,
            name: name,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setJustTimecodeOut(project, name, timecode, token) {
    $.ajax({
        url: "/api/setjusttimecodeout",
        type: "post",
        data: {
            project: project,
            name: name,
            timecode: timecode,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setNote(project, name, text, userid, overwrite, token) {
    if (multiInput) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            
            if(cboxes[i].checked === false) {
                continue
            }
            let currentName = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/setnote",
                type: "post",
                data: {
                    project: project,
                    name: currentName,
                    text: text,
                    userid: userid,
                    overwrite: overwrite,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.error !== "") {
                        alert(data.error)
                        return
                    }
                    if (document.getElementById("set-note-overwrite").checked) {
                        // note-{{.Name}} 내부 내용을 교체한다.
                        document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    } else {
                        // note-{{.Name}} 내부 내용에 추가한다.
                        document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.name).innerHTML;
                    }
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/setnote",
            type: "post",
            data: {
                project: project,
                name: name,
                text: text,
                userid: userid,
                overwrite: overwrite,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.error !== "") {
                    alert(data.error)
                    return
                }
                if (document.getElementById("set-note-overwrite").checked) {
                    // note-{{.Name}} 내부 내용을 교체한다.
                    document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                } else {
                    // note-{{.Name}} 내부 내용에 추가한다.
                    document.getElementById("note-"+data.name).innerHTML = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>') + "<br>" + document.getElementById("note-"+data.name).innerHTML;
                }
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function addComment(project, name, text, userid, token) {
    if (multiInput) {
        let cboxes = document.getElementsByName('selectID');
        for (let i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            let currentName = cboxes[i].getAttribute("id");
            sleep(200);
            $.ajax({
                url: "/api/addcomment",
                type: "post",
                data: {
                    project: project,
                    name: currentName,
                    text: text,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.error !== "") {
                        alert(data.error)
                        return
                    }
                    // comments-{{.Name}} 내부 내용에 추가한다.
                    let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                    let newComment = `<div id="comment-${data.name}-${data.date}">
                    <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
                    <span class="remove" data-toggle="modal" data-target="#rmcomment" onclick="setModal('rm-comment-name', '${data.name}');setModal('rm-comment-time', '${data.date}');setModal('rm-comment-text', '${data.text}');setModal('rm-comment-userid', '${data.userid}')">×</span>
                    <br><small class="text-white">${body}</small><hr class="my-1 p-0 m-0 divider"></hr></div>`
                    document.getElementById("comments-"+data.name).innerHTML = newComment + document.getElementById("comments-"+data.name).innerHTML;
                },
                error: function(request,status,error){
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
                }
            });
        }
    } else {
        $.ajax({
            url: "/api/addcomment",
            type: "post",
            data: {
                project: project,
                name: name,
                text: text,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.error !== "") {
                    alert(data.error)
                    return
                }
                // comments-{{.Name}} 내부 내용에 추가한다.
                let body = data.text.replace(/(?:\r\n|\r|\n)/g, '<br>');
                let newComment = `<div id="comment-${data.name}-${data.date}">
                <span class="text-badge">${data.date} / <a href="/user?id=${data.userid}" class="text-darkmode">${data.userid}</a></span>
                <span class="remove" data-toggle="modal" data-target="#rmcomment" onclick="setModal('rm-comment-name', '${data.name}');setModal('rm-comment-time', '${data.date}');setModal('rm-comment-text', '${data.text}');setModal('rm-comment-userid', '${data.userid}')">×</span>
                <br><small class="text-white">${body}</small><hr class="my-1 p-0 m-0 divider"></hr></div>`
                document.getElementById("comments-"+name).innerHTML = newComment + document.getElementById("comments-"+name).innerHTML;
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function rmComment(project, name, date, userid, token) {
    $.ajax({
        url: "/api/rmcomment",
        type: "post",
        data: {
            project: project,
            name: name,
            date: date,
            userid: userid
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.error !== "") {
                alert(data.error)
                return
            }
            document.getElementById(`comment-${data.name}-${data.date}`).remove();
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}


function addSource(project, name, title, path, userid, token) {
    if (multiInput) {
        let cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            let currentName = cboxes[i].getAttribute("id")
            $.ajax({
                url: "/api/addsource",
                type: "post",
                data: {
                    project: project,
                    name: currentName,
                    title: title,
                    path: path,
                    userid: userid,
                },
                headers: {
                    "Authorization": "Basic "+ token
                },
                dataType: "json",
                success: function(data) {
                    if (data.error !== "") {
                        alert(data.error)
                        return
                    }
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/addsource",
            type: "post",
            data: {
                project: project,
                name: name,
                title: title,
                path: path,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                if (data.error !== "") {
                    alert(data.error)
                    return
                }
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}


function rmSource(project, name, title, userid, token) {
    if (multiInput) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            if(cboxes[i].checked === false) {
                continue
            }
            sleep(200);
            currentName = cboxes[i].getAttribute("id");
            $.ajax({
                url: "/api/rmsource",
                type: "post",
                data: {
                    project: project,
                    name: currentName,
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
                    alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
                }
            });
            
        }
    } else {
        $.ajax({
            url: "/api/rmsource",
            type: "post",
            data: {
                project: project,
                name: name,
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
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}


function setTaskUser(project, name, task, user, token) {
    $.ajax({
        url: "/api/settaskuser",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            user: removeWhiteSpace(user),
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setTaskUsers(project, task, user, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/settaskuser",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                task: task,
                user, user,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}


function setTaskStatus(project, name, task, status, token) {
    $.ajax({
        url: "/api/settaskstatus",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            status: status,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setTaskStatuses(project, task, status, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/settaskstatus",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                task: task,
                status: status,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setTaskDate(project, name, task, date, token) {
    console.log(project, name, task, date, token);
    $.ajax({
        url: "/api/settaskdate",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            date: date,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setTaskDates(project, task, date, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/settaskdate",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                task: task,
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}


function setTaskPredate(project, name, task, date, token) {
    console.log(project, name, task, date, token);
    $.ajax({
        url: "/api/settaskpredate",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            date: date,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setTaskPredates(project, task, date, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/settaskpredate",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                task: task,
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setDeadline2D(project, name, date, token) {
    $.ajax({
        url: "/api/setdeadline2d",
        type: "post",
        data: {
            project: project,
            name: name,
            date: date,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setDeadline2Ds(project, date, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        console.log(project, cboxes[i].getAttribute("id"), date)
        $.ajax({
            url: "/api/setdeadline2d",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setDeadline3D(project, name, date, token) {
    $.ajax({
        url: "/api/setdeadline3d",
        type: "post",
        data: {
            project: project,
            name: name,
            date: date,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setDeadline3Ds(project, date, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        console.log(project, cboxes[i].getAttribute("id"), date)
        $.ajax({
            url: "/api/setdeadline3d",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                date: date,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setShottype(project, name, shottype, token) {
    $.ajax({
        url: "/api/setshottype",
        type: "post",
        data: {
            project: project,
            name: name,
            shottype: shottype,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setShottypes(project, shottype, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/setshottype",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                shottype: shottype,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setAssettype(project, name, assettype, token) {
    $.ajax({
        url: "/api/setassettype",
        type: "post",
        data: {
            project: project,
            name: name,
            assettype: assettype,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setAssettypes(project, assettype, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/setassettype",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                assettype: assettype,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setRnum(project, name, rnum, userid, token) {
    $.ajax({
        url: "/api/setrnum",
        type: "post",
        data: {
            project: project,
            name: name,
            rnum: rnum,
            userid: userid,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            if (data.error !== "") {
                alert(data.error)
                return
            }
            if (data.rnum !== "") {
                document.getElementById("rnum-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#setrnum" onclick="setModal('set-rnum-text', '${data.rnum}' );setModal('set-rnum-name', '${data.name}');setModal('set-rnum-userid', '${data.userid}')"{{end}}>${data.rnum}</span>`;
            } else {
                document.getElementById("rnum-"+data.name).innerHTML = `<span class="black-opbg" data-toggle="modal" data-target="#setrnum" onclick="setModal('set-rnum-text', '${data.rnum}' );setModal('set-rnum-name', '${data.name}');setModal('set-rnum-userid', '${data.userid}')"{{end}}>no rnum</span>`;
            }
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function addTags(project, tag, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        name = cboxes[i].getAttribute("id");
        // 기존 폼에 있는 값에도 태그를 추가한다.
        before = document.getElementById("input-tag-" + name).value;
        let splitString = ","
        if (before === "") {
            splitString = ""
        }
        document.getElementById("input-tag-" + name).value = before + splitString + tag;
        $.ajax({
            url: "/api/addtag",
            type: "post",
            
            data: {
                project: project,
                name: name,
                tag: tag,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setTags(project, tags, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        name = cboxes[i].getAttribute("id");
        document.getElementById("input-tag-"+name).value = tags;
        $.ajax({
            url: "/api/settags",
            type: "post",
            
            data: {
                project: project,
                name: name,
                tags: tags,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

function setTag(project, name, tags, token) {
    $.ajax({
        url: "/api/settags",
        type: "post",
        
        data: {
            project: project,
            name: name,
            tags: tags,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}


function rmTags(project, tag, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        name = cboxes[i].getAttribute("id");
        if (document.getElementById("input-tag-"+name).value.startsWith(tag+",")) {
            // 태그가 처음에 위치할 때
            document.getElementById("input-tag-"+name).value = document.getElementById("input-tag-"+name).value.replace(tag+",","");
        } else if (document.getElementById("input-tag-"+name).value.includes(","+tag+",")) {
            // 태그가 중간에 있을 때
            document.getElementById("input-tag-"+name).value = document.getElementById("input-tag-"+name).value.replace(","+tag,",");
        } else {
            // 태그가 끝에 있을 때
            document.getElementById("input-tag-"+name).value = document.getElementById("input-tag-"+name).value.replace(","+tag,"");
        }
        $.ajax({
            url: "/api/rmtag",
            type: "post",
            
            data: {
                project: project,
                name: name,
                tag: tag,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
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

function setTaskLevel(project, name, task, level, token) {
    $.ajax({
        url: "/api/settasklevel",
        type: "post",
        data: {
            project: project,
            name: name,
            task: task,
            level: level,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        },
        error: function(request,status,error){
            alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
        }
    });
}

function setTaskLevels(project, task, level, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        name = cboxes[i].getAttribute("id");
        $.ajax({
            url: "/api/settasklevel",
            type: "post",
            data: {
                project: project,
                name: name,
                task: task,
                level, level,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            },
            error: function(request,status,error){
                alert("code:"+request.status+"\n"+"status:"+status+"\n"+"Msg:"+JSON.parse(request.responseText).error+"\n"+"error:"+error);
            }
        });
    }
}

// setModal 함수는 modalID와 value를 받아서 modal에 셋팅한다.
function setModal(modalID, value) {
    document.getElementById(modalID).value=value;
}