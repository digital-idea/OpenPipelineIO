function initModalPartners() {
    fetch('/api/partners', {
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
        let e = document.getElementById('projectforpartner-partnername');
        if (e === null) {
            return
        }
        e.innerHTML = "";
        for (let j = 0; j < data.length; j++){
            let opt = document.createElement('option');
            opt.value = data[j].name;
            opt.innerHTML = data[j].name;
            e.appendChild(opt);
        }
    })
    .catch((err) => {
        alert(err)
    });
}
initModalPartners() // 페이지가 로딩되면 먼저 실행한다.