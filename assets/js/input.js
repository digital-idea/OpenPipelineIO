

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
            console.info(data)
        }
    });
}

function addNotes(project, text, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/addnote",
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

function addLink(project, name, text, token) {
    $.ajax({
        url: "/api/addlink",
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

function addLinks(project, text, token) {
    var cboxes = document.getElementsByName('selectID');
    for (var i = 0; i < cboxes.length; ++i) {
        if(cboxes[i].checked === false) {
            continue
        }
        $.ajax({
            url: "/api/addlink",
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
    console.log(project,task,user,token)
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
        console.log(name, tags);
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
        console.log(project, cboxes[i].getAttribute("id"), tag, token);
        $.ajax({
            url: "/api/rmtag",
            type: "post",
            
            data: {
                project: project,
                name: cboxes[i].getAttribute("id"),
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