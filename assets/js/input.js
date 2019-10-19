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
        }
    });
}

function addNote(project, name, text, token) {
    sleep(200);
    $.ajax({
        url: "/api/setnote",
        type: "post",
        data: {
            project: project,
            name: name,
            text: text,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            // console.info(data)
        }
    });
}

function addNotes(project, text, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        name = cboxes[i].getAttribute("id")
        // 다음 사용될 API
        sleep(200);
        $.ajax({
            url: "/api/setnote",
            type: "post",
            data: {
                project: project,
                name: name,
                text: text,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                // console.info(data)
            }
        });
    }
}

function addComment(project, name, text, userid, token) {
    if (multiInput) {
        var cboxes = document.getElementsByName('selectID');
        for (var i = 0; i < cboxes.length; ++i) {
            
            if(cboxes[i].checked === false) {
                continue
            }
            currentName = cboxes[i].getAttribute("id");
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
                    console.info(data)
                }
            });
            // comments-{{.Name}} 내부 내용에 추가한다.
            let body = text.replace("\n","<br>")
            let newComment = `<span class="text-badge">Now / <a href="/user?id=${userid}" class="text-darkmode">${userid}</a></span><br><small class="text-white">${body}</small><hr class="my-1 p-0 m-0 divider"></hr>`
            document.getElementById("comments-"+currentName).innerHTML = newComment + document.getElementById("comments-"+currentName).innerHTML;
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
                console.info(data)
            }
        });
        // comments-{{.Name}} 내부 내용에 추가한다.
        let body = text.replace("\n","<br>")
        let newComment = `<span class="text-badge">Now / <a href="/user?id=${userid}" class="text-darkmode">${userid}</a></span><br><small class="text-white">${body}</small><hr class="my-1 p-0 m-0 divider"></hr>`
        document.getElementById("comments-"+name).innerHTML = newComment + document.getElementById("comments-"+name).innerHTML;
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
            console.info(data)
        }
    });
    document.getElementById(`comments-${name}-${date}`).remove();
}

function addComments(project, text, userid, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/addcomment",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                text: text,
                userid: userid,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            }
        });
    }
}

function addSource(project, name, title, path, token) {
    $.ajax({
        url: "/api/addsource",
        type: "post",
        data: {
            project: project,
            name: name,
            title: title,
            path: path,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
        }
    });
}

function addSources(project, title, path, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        sleep(200);
        $.ajax({
            url: "/api/addsource",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                title: title,
                path: path,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
            }
        });
    }
}



function rmSources(project, title, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        sleep(200);
        $.ajax({
            url: "/api/rmsource",
            type: "post",
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
                title: title,
            },
            headers: {
                "Authorization": "Basic "+ token
            },
            dataType: "json",
            success: function(data) {
                console.info(data)
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
            }
        });
    }
}

function setRnum(project, name, rnum, token) {
    $.ajax({
        url: "/api/setrnum",
        type: "post",
        data: {
            project: project,
            name: name,
            rnum: rnum,
        },
        headers: {
            "Authorization": "Basic "+ token
        },
        dataType: "json",
        success: function(data) {
            console.info(data)
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
        console.log(project, cboxes[i].getAttribute("id"), tag, token);
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
        console.log(project,name,task,level);
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
            }
        });
    }
}

// setModal 함수는 modalid와 value를 받아서 modal에 셋팅한다.
function setModal(modalid, value) {
    document.getElementById(modalid).value=value;
}