function setTags(project, name, tags) {
    $.post("/api/settags",
    {
        project: project,
        name: name,
        tags: tags,
    },
    function(data, status){
        console.log("Data: " + data + "\nStatus: " + status);
    });
}