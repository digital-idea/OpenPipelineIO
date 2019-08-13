function setTags(project, name, tags, token) {
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

function setTaskStatus(project, name, task, status, token) {
    $.ajax({
        url: "/api/status",
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