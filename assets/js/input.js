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
    // 레거시 코드
    sleep(200);
    $.ajax({
        url: "/api/addnote",
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

function addNotes(project, text, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        name = cboxes[i].getAttribute("id")
        // 레거시 API
        sleep(200);
        $.ajax({
            url: "/api/addnote",
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

function addComment(project, name, text, token) {
    $.ajax({
        url: "/api/addcomment",
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
            console.info(data)
        }
    });
}

function addComments(project, text, token) {
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
            user: user,
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

function selectCheckboxAll() {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = true;
    }
}

function selectCheckboxNone() {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        cboxes[i].checked = false;
    }
}

function selectCheckboxInvert() {
    var cboxes = document.getElementsByName('selectID');
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