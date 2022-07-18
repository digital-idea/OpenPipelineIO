
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
        let total = 0
        let none = 0
        let onHold = 0
        let inProgress = 0
        let finalApproved = 0

        // 전체 샷 진행챠트 그리기
        let statusProgress = document.getElementById("total-shot-status-progress")
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
                if (key === "assign" || key === "ready" || key === "wip" || key === "confirm" || key === "out") {
                    inProgress += value
                }
                if (key === "done") {
                    finalApproved += value
                }
                if (key === "omit" || key === "hold") {
                    onHold += value
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
                opt.setAttribute("title", `${key} ${percent.toFixed(1)}%(${value})`)
                opt.innerHTML = `${key}<br>${percent.toFixed(1)}%(${value})`
                statusProgress.appendChild(opt);
            }
        }
        // 진행률 챠트 그리기
        let progress = document.getElementById("total-shot-progress")
        progresslist = [
            {"title":"In Progress", "style":"inprogress", "num":inProgress},
            {"title":"Final Approved", "style":"finalapproved", "num":finalApproved},
            {"title":"On Hold", "style":"onhold", "num":onHold},
        ]
        for (let i in progresslist) {
            let opt = document.createElement('div');
            opt.classList.add("progress-bar")
            opt.classList.add("bg-"+progresslist[i].style)
            opt.role = "progressbar"
            let percent = (progresslist[i].num / (inProgress + finalApproved + onHold)) * 100
            opt.style = "width: " + percent.toFixed(1) + "%"
            opt.setAttribute("aria-valuenow",progresslist[i].num)
            opt.setAttribute("aria-valuemin","0")
            opt.setAttribute("aria-valuemax",(inProgress + finalApproved + onHold))
            opt.setAttribute("data-bs-toggle","tooltip")
            opt.setAttribute("data-bs-placement","top")
            let title = `${progresslist[i].title} ${percent.toFixed(1)}%(${progresslist[i].num})`
            opt.setAttribute("title", title)
            opt.innerHTML = title
            progress.appendChild(opt);
        }
        initTooltip()
    })
    .catch((err) => {
        console.log(err)
    });
}

function StatusPerProject() {
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
        let total = 0
        let none = 0
        let onHold = 0
        let inProgress = 0
        let finalApproved = 0
    
        let projectList = document.getElementById("statusperprojects")
        for (let [key, value] of Object.entries(obj.projects)) {
            // 줄긋기
            let line = document.createElement("hr")
            line.classList.add("p-0")
            line.classList.add("m-0")
            line.classList.add("my-2")
            projectList.appendChild(line)

            // 프로젝트별 제목 그리기
            let projectRoot = document.createElement("div");
            let title  = document.createElement("span");
            title.classList.add('align-middle');
            title.innerText = key
            projectRoot.appendChild(title);
            projectList.appendChild(projectRoot)

            // Progress바 추가
            let progressBar = document.createElement("div")
            progressBar.classList.add("progress")
            progressBar.classList.add("mt-2")
            progressBar.id = key + "-shot-progress"
            progressBar.style.height = "30px"
            projectList.appendChild(progressBar)
            
            // Status바 추가
            let statusBar = document.createElement("div")
            statusBar.classList.add("progress")
            statusBar.classList.add("mt-2")
            statusBar.id = key + "-shot-status-progress"
            statusBar.style.height = "50px"
            projectList.appendChild(statusBar)

            // Status바에 데이터 넣기
            let currentTotal = 0
            let currentNone = 0
            let currentOnHold = 0
            let currentInProgress = 0
            let currentFinalApproved = 0
            let currentStatus = obj.projects[key]
            
            // 그래프를 그리기 전에 전체 갯수를 구해야 한다.
            for (let [key, value] of Object.entries(currentStatus)) {
                currentTotal += value
            }

            let order = ["none","assign","ready","wip","confirm","done","out","omit","hold"]
            
            for (let i = 0; i < order.length; i++) {
                for (let [currentKey, currentValue] of Object.entries(currentStatus)) {
                    if (currentValue === 0) {
                        continue
                    }
                    if (!(order[i] == currentKey)) {
                        continue
                    }
                    if (currentKey === "none") {
                        currentNone = currentValue
                        continue
                    }
                    if (currentKey === "assign" || currentKey === "ready" || currentKey === "wip" || currentKey === "confirm" || currentKey === "out") {
                        currentInProgress += currentValue
                    }
                    if (currentKey === "done") {
                        currentFinalApproved += currentValue
                    }
                    if (currentKey === "omit" || currentKey === "hold") {
                        currentOnHold += currentValue
                    }
                    let opt = document.createElement('div');
                    opt.classList.add("progress-bar")
                    opt.classList.add("bg-"+currentKey)
                    opt.role = "progressbar"
                    let percent = (currentValue / (currentTotal - currentNone)) * 100
                    opt.style = "width: " + percent.toFixed(1) + "%"
                    opt.setAttribute("aria-valuenow",currentValue)
                    opt.setAttribute("aria-valuemin","0")
                    opt.setAttribute("aria-valuemax",(currentTotal - currentNone))
                    opt.setAttribute("data-bs-toggle","tooltip")
                    opt.setAttribute("data-bs-placement","top")
                    opt.setAttribute("title", `${currentKey} ${percent.toFixed(1)}%(${currentValue})`)
                    opt.innerHTML = `${currentKey}<br>${percent.toFixed(1)}%(${currentValue})`
                    statusBar.appendChild(opt);
                }
            }

            // 진행률 챠트 그리기
            
            currentProgresslist = [
                {"title":"In Progress", "style":"inprogress", "num":currentInProgress},
                {"title":"Final Approved", "style":"finalapproved", "num":currentFinalApproved},
                {"title":"On Hold", "style":"onhold", "num":currentOnHold},
            ]
            for (let i in currentProgresslist) {
                let opt = document.createElement('div');
                opt.classList.add("progress-bar")
                opt.classList.add("bg-"+currentProgresslist[i].style)
                opt.role = "progressbar"
                let percent = (currentProgresslist[i].num / (currentInProgress + currentFinalApproved + currentOnHold)) * 100
                opt.style = "width: " + percent.toFixed(1) + "%"
                opt.setAttribute("aria-valuenow",currentProgresslist[i].num)
                opt.setAttribute("aria-valuemin","0")
                opt.setAttribute("aria-valuemax",(currentInProgress + currentFinalApproved + currentOnHold))
                opt.setAttribute("data-bs-toggle","tooltip")
                opt.setAttribute("data-bs-placement","top")
                let title = `${currentProgresslist[i].title} ${percent.toFixed(1)}%(${currentProgresslist[i].num})`
                opt.setAttribute("title", title)
                opt.innerHTML = title
                progressBar.appendChild(opt);
            }
            

            // None 갯수 추가
            let noneInfo = document.createElement("small")
            noneInfo.innerText = "None Status 갯수: " + value["none"]
            projectList.appendChild(noneInfo)
            
            
        }

        // 전체 샷 진행챠트 그리기
        let statusProgress = document.getElementById("total-shot-status-progress")
        // 그래프를 그리기 전에 전체 갯수를 구해야 한다.
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
                if (key === "assign" || key === "ready" || key === "wip" || key === "confirm" || key === "out") {
                    inProgress += value
                }
                if (key === "done") {
                    finalApproved += value
                }
                if (key === "omit" || key === "hold") {
                    onHold += value
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
                opt.setAttribute("title", `${key} ${percent.toFixed(1)}%(${value})`)
                opt.innerHTML = `${key}<br>${percent.toFixed(1)}%(${value})`
                statusProgress.appendChild(opt);
            }
        }
        // 진행률 챠트 그리기
        let progress = document.getElementById("total-shot-progress")
        progresslist = [
            {"title":"In Progress", "style":"inprogress", "num":inProgress},
            {"title":"Final Approved", "style":"finalapproved", "num":finalApproved},
            {"title":"On Hold", "style":"onhold", "num":onHold},
        ]
        for (let i in progresslist) {
            let opt = document.createElement('div');
            opt.classList.add("progress-bar")
            opt.classList.add("bg-"+progresslist[i].style)
            opt.role = "progressbar"
            let percent = (progresslist[i].num / (inProgress + finalApproved + onHold)) * 100
            opt.style = "width: " + percent.toFixed(1) + "%"
            opt.setAttribute("aria-valuenow",progresslist[i].num)
            opt.setAttribute("aria-valuemin","0")
            opt.setAttribute("aria-valuemax",(inProgress + finalApproved + onHold))
            opt.setAttribute("data-bs-toggle","tooltip")
            opt.setAttribute("data-bs-placement","top")
            let title = `${progresslist[i].title} ${percent.toFixed(1)}%(${progresslist[i].num})`
            opt.setAttribute("title", title)
            opt.innerHTML = title
            progress.appendChild(opt);
        }
        initTooltip()
    })
    .catch((err) => {
        console.log(err)
    });
}

function Shottype() {
    fetch("/api/statistics/shottype", {
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
        let total = obj.totalnum.type2d + obj.totalnum.type3d + obj.totalnum.typenone
        let percent2d = (obj.totalnum.type2d / total) * 100
        let percent3d = (obj.totalnum.type3d / total) * 100
        let percentnone = (obj.totalnum.typenone / total) * 100
        // 2D갯수 셋팅
        let e1 = document.getElementById("shottype-2d")
        let title2d = `2D ${percent2d.toFixed(1)}% (${obj.totalnum.type2d})`
        e1.innerHTML = title2d
        e1.setAttribute("title", title2d)
        e1.style.width = percent2d + "%"
        // 3D갯수 셋팅
        let e2 = document.getElementById("shottype-3d")
        let title3d = `3D ${percent3d.toFixed(1)}% (${obj.totalnum.type3d})`
        e2.innerHTML = title3d
        e2.setAttribute("title", title3d)
        e2.style.width = percent3d + "%"
        // None 갯수 셋팅
        let e3 = document.getElementById("shottype-none")
        let titlenone = `None ${percentnone.toFixed(1)}% (${obj.totalnum.typenone})`
        e3.innerHTML = titlenone
        e3.setAttribute("title", titlenone)
        e3.style.width = percentnone + "%"
        // 툴팁 요소가 추가되었다. 툴팁셋팅을 초기화 한다.
        initTooltip()
    })
    .catch((err) => {
        console.log(err)
    });
}

function Itemtype() {
    fetch("/api/statistics/itemtype", {
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
        let total = obj.totalnum.shot + obj.totalnum.asset
        let percentshot = (obj.totalnum.shot / total) * 100
        let percentasset = (obj.totalnum.asset / total) * 100
        // Shot 갯수 셋팅
        let e1 = document.getElementById("itemtype-shot")
        let titleshot = `Shot ${percentshot.toFixed(1)}% (${obj.totalnum.shot})`
        e1.innerHTML = titleshot
        e1.setAttribute("title", titleshot)
        e1.style.width = percentshot + "%"
        // Asset 갯수 셋팅
        let e2 = document.getElementById("itemtype-asset")
        let titleasset = `Asset ${percentasset.toFixed(1)}% (${obj.totalnum.asset})`
        e2.innerHTML = titleasset
        e2.setAttribute("title", titleasset)
        e2.style.width = percentasset + "%"
        // 툴팁 요소가 추가되었다. 툴팁셋팅을 초기화 한다.
        initTooltip()
    })
    .catch((err) => {
        console.log(err)
    });
}

function initTooltip() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
    })
}

