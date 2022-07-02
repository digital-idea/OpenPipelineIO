function initModalProjects() {
    fetch('/api2/projects', {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+ document.getElementById("token").value,
        },
    })
    .then((response) => {
        if (!response.ok) {
            throw Error(response.statusText + " - " + response.url);
        }
        return response.json()
    })
    .then((data) => {
        if (data === null) {
            return
        }
        let selectProject = document.getElementById('projectforpartner-projectname');
        if (selectProject) { // 해당 elements가 존재한다면 값을 추가한다.
            selectProject.innerHTML = "";
            for (let j = 0; j < data.length; j++){
                let opt = document.createElement('option');
                opt.value = data[j].id;
                opt.innerHTML = data[j].id;
                selectProject.appendChild(opt);
            }
        }
    })
    .catch((err) => {
        alert(err)
    });
}
initModalProjects() // 페이지가 로딩되면 먼저 실행한다.