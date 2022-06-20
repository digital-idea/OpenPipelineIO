
function DeadlinenumThismonth() {
    let date = new Date();
    let year = date.getFullYear();
    let month = '' + (date.getMonth() + 1); // 자바스크립트는 월을 0부터 센다.
    if (month.length < 2) {
        month = '0' + month;
    }
    thismonth = [year,month].join("-")

    fetch("/api/statistics/deadlinenum?date="+thismonth, {
        method: 'GET',
        headers: {"Authorization": "Basic "+ document.getElementById("token").value},
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                alert(text)
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        document.getElementById("this-month-deadline").innerHTML = obj.total
        let detail = document.getElementById("thisMonthDeadlineDetail")
        // 프로젝트별 정보 전달
        for (let [key, value] of Object.entries(obj.projects)) {
            let opt = document.createElement('div');
            opt.innerHTML = `${key} 프로젝트: ${value}개`
            detail.appendChild(opt);
        }
    })
    .catch((err) => {
        console.log(err)
    });
}


function DeadlinenumNextmonth() {
    let date = new Date();
    let year = date.getFullYear();
    let month = '' + (date.getMonth() + 2); // 자바스크립트는 월을 0부터 센다.
    if (month.length < 2) {
        month = '0' + month;
    }
    if (month === "12") {
        month = "01"
        year += 1
    }
    nextmonth = [year,month].join("-")

    fetch("/api/statistics/deadlinenum?date="+nextmonth, {
        method: 'GET',
        headers: {"Authorization": "Basic "+ document.getElementById("token").value},
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                alert(text)
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        document.getElementById("next-month-deadline").innerHTML = obj.total
        let detail = document.getElementById("nextMonthDeadlineDetail")
        // 프로젝트별 정보 전달
        for (let [key, value] of Object.entries(obj.projects)) {
            let opt = document.createElement('div');
            opt.innerHTML = `${key} 프로젝트: ${value}개`
            detail.appendChild(opt);
        }
    })
    .catch((err) => {
        console.log(err)
    });
}

function NeedDeadlinenum() {
    fetch("/api/statistics/needdeadlinenum", {
        method: 'GET',
        headers: {"Authorization": "Basic "+ document.getElementById("token").value},
    })
    .then((response) => {
        if (!response.ok) {
            response.text().then(function (text) {
                alert(text)
                return
            });
        }
        return response.json()
    })
    .then((obj) => {
        // 전체 갯수 전달
        document.getElementById("need-deadline").innerHTML = obj.total
        let detail = document.getElementById("needDeadlineDetail")
        // 프로젝트별 정보 전달
        for (let [key, value] of Object.entries(obj.projects)) {
            let opt = document.createElement('div');
            opt.innerHTML = `${key} 프로젝트: ${value}개`
            detail.appendChild(opt);
        }
    })
    .catch((err) => {
        console.log(err)
    });
}

