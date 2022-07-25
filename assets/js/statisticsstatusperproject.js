function SetProjectsTags() {
    fetch("/api2/projects", {
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
        let projectstags = document.getElementById("projectstags")
        
        for (let i = 0; i < obj.length; i++) {
            // 버튼 생성
            let button = document.createElement("span")
            button.classList.add("btn")
            button.classList.add("btn-sm")
            button.classList.add("btn-outline-darkmode")
            button.classList.add("m-1")
            // 체크박스 옵션
            let checkoption = document.createElement("input")
            let project = obj[i].id
            checkoption.classList.add("me-1")
            checkoption.id = "toggle-" + project
            checkoption.value = "chart-" + project
            checkoption.addEventListener("click", toggleCharts);
            checkoption.type = "checkbox"
            button.appendChild(checkoption)
            // 프로젝트 제목
            let title = document.createElement("span")
            title.classList.add("text-white")
            title.innerText = `${obj[i].id} (${projectStatusToText(obj[i].status)})`
            button.appendChild(title)

            // 추가
            projectstags.appendChild(button)
        }
    })
    .catch((err) => {
        console.log(err)
    });
}

function toggleCharts() {
    // 켜져있는 프로젝트를 구한다.
    let inputs = document.querySelectorAll('input[id^="toggle-"]')
    for (let i = 0; i < inputs.length; i++) {
        let e = document.getElementById(inputs[i].value)
        if (inputs[i].checked === true) {
            e.style.display = "block" // 켜져있다면 해당 프로젝트를 켠다.
        } else {
            e.style.display = "none" // 꺼져있다면 해당 프로젝트를 끈다.
        }
    }
    
    
}

function projectStatusToText(num) {
    switch (num) {
        case 0:
            return "테스트"
        case 1:
            return "준비중"
        case 2:
            return "진행중"
        case 3:
            return "중단"
        case 4:
            return "백업중"
        case 5:
            return "백업완료"
        case 6:
            return "소송중"
        default:
            return "Unknown"
    }
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
            let chart = document.createElement("div");
            chart.id = "chart-" + key

            // 줄긋기
            let line = document.createElement("hr")
            line.classList.add("p-0")
            line.classList.add("m-0")
            line.classList.add("my-2")
            chart.appendChild(line)

            // 프로젝트별 제목 그리기
            let projectRoot = document.createElement("div");
            let title  = document.createElement("span");
            title.classList.add('align-middle');
            title.innerText = key
            projectRoot.appendChild(title);
            chart.appendChild(projectRoot)

            // Progress바 추가
            let progressBar = document.createElement("div")
            progressBar.classList.add("progress")
            progressBar.classList.add("mt-2")
            progressBar.id = key + "-shot-progress"
            progressBar.style.height = "30px"
            chart.appendChild(progressBar)
            
            // Status바 추가
            let statusBar = document.createElement("div")
            statusBar.classList.add("progress")
            statusBar.classList.add("mt-2")
            statusBar.id = key + "-shot-status-progress"
            statusBar.style.height = "50px"
            chart.appendChild(statusBar)

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
                if (currentProgresslist[i].num === 0) {
                    continue
                }
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
            chart.appendChild(noneInfo)

            // 완성된 챠트 1개를 추가한다.
            chart.style.display = "none"
            projectList.appendChild(chart)
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
        // 전체 진행률 챠트 그리기
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

function initTooltip() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl)
    })
}

