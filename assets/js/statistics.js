
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
        let e = document.getElementById("this-month-deadline")
        e.innerHTML = `이번달 ${obj.total}개`
        e.style.width = (obj.total * 100) + "%"
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
        let e = document.getElementById("next-month-deadline")
        e.innerHTML = `다음달 ${obj.total}개`
        e.style.width = (obj.total * 100) + "%"
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


function TotalShotProgress() {
    fetch("/api1/statistics/shot", {
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
        let progress = document.getElementById("total-shot-progress")
        let total = 0
        let none = 0
        for (let [key, value] of Object.entries(obj.total)) {
            total += value
        }
        let order = ["none","assign","ready","wip","confirm","done","out","omit","hold"]
        for (let i = 0; i < order.length; i++) {
            for (let [key, value] of Object.entries(obj.total)) {
                if (value === 0) {
                    continue
                }
                if (!(order[i] == key)) {
                    continue
                }
                if (key === "none") {
                    // none status는 그래프를 그리지 않는다.
                    none = value
                    document.getElementById("NoneStatusNum").innerHTML = value
                    continue
                }
                
                let opt = document.createElement('div');
                opt.classList.add("progress-bar")
                opt.classList.add("bg-"+key)
                opt.role = "progressbar"
                let percent = (value / (total - none)) * 100
                opt.style = "width: " + percent.toFixed(1) + "%"
                opt.setAttribute("aria-valuenow",value)
                opt.setAttribute("aria-valuemin","0")
                opt.setAttribute("aria-valuemax",(total - none))
                opt.setAttribute("data-bs-toggle","tooltip")
                opt.setAttribute("data-bs-placement","top")
                opt.setAttribute("title", key + ":" + value)
                opt.innerHTML = `${key}<br>${percent.toFixed(1)}%(${value})`
                progress.appendChild(opt);
            }
        }
    })
    .catch((err) => {
        console.log(err)
    });
}
