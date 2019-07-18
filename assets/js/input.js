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