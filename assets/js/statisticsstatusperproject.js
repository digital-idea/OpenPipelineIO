var data = new Object() // 프로젝트 통계를 넣는다.

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
            checkoption.setAttribute("status", projectStatusToText(obj[i].status))
            checkoption.setAttribute("project", project)
            checkoption.addEventListener("click", RenderPerChart);
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

function RenderPerChart() {
    data.checked = [] // 체크된 프로젝트 리스트 초기화
    
    // 켜져있는 프로젝트를 구한다.
    let inputs = document.querySelectorAll('input[id^="toggle-"]')
    
    for (let i = 0; i < inputs.length; i++) {
        let e = document.getElementById(inputs[i].value)
        if (inputs[i].checked === true) {
            data.checked.push(inputs[i].getAttribute("project")) // data에 체크된 프로젝트를 저장한다.
            e.style.display = "block" // 켜져있다면 해당 프로젝트를 켠다.
        } else {
            e.style.display = "none" // 꺼져있다면 해당 프로젝트를 끈다.
        }
    }

    // 선택된 프로젝트만 챠트를 렌더링한다.
    RenderSumChart()
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
        data = obj // data 글로벌 변수에 데이터에 넣습니다.
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
                    document.getElementById("TotalNoneStatusNum").innerHTML = value
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

function RenderSumChart() {
    // 초기화
    data.sum = new Object()
    data.sum.none = 0
    data.sum.assign = 0
    data.sum.ready = 0
    data.sum.wip = 0
    data.sum.confirm = 0
    data.sum.done = 0
    data.sum.out = 0
    data.sum.omit = 0
    data.sum.hold = 0
    data.sum.total = 0
    data.sum.onhold = 0
    data.sum.inprogress = 0
    data.sum.finalapproved = 0
    
    // 그래프를 그리기 전에 선택된 프로젝트의 더한 샷 갯수를 구해야 한다.
    for (let n = 0; n < data.checked.length; n++) {
        let p = data.checked[n] // 선택된 프로젝트 이름을 구한다.
        // 상태 챠트를 그릴 때 사용할 값
        data.sum.none += data.projects[p].none
        data.sum.assign += data.projects[p].assign
        data.sum.ready += data.projects[p].ready
        data.sum.wip += data.projects[p].wip
        data.sum.confirm += data.projects[p].confirm
        data.sum.done += data.projects[p].done
        data.sum.out += data.projects[p].out
        data.sum.omit += data.projects[p].omit
        data.sum.hold += data.projects[p].hold
        // Total값 추후 연산에 편하게 쓸 수 있다.
        data.sum.total += data.projects[p].none + data.projects[p].assign + data.projects[p].ready + data.projects[p].wip + data.projects[p].confirm + data.projects[p].done + data.projects[p].out + data.projects[p].omit + data.projects[p].hold
        // Progress 챠트 준비물
        data.sum.inprogress += data.projects[p].assign + data.projects[p].ready + data.projects[p].wip + data.projects[p].confirm + data.projects[p].out
        data.sum.finalapproved += data.projects[p].done
        data.sum.onhold += data.projects[p].omit + data.projects[p].hold
    }
    
    // Progress 합 챠트
    let progress = document.getElementById("sum-shot-progress")
    progress.innerHTML = "" // 기존 챠트를 초기화 한다.
    progresslist = [
        {"title":"In Progress", "style":"inprogress", "num":data.sum.inprogress},
        {"title":"Final Approved", "style":"finalapproved", "num":data.sum.finalapproved},
        {"title":"On Hold", "style":"onhold", "num":data.sum.onhold},
    ]
    for (let i in progresslist) {
        if (progresslist[i].num === 0) {
            continue
        }
        let opt = document.createElement('div');
        opt.classList.add("progress-bar")
        opt.classList.add("bg-"+progresslist[i].style)
        opt.role = "progressbar"
        let percent = (progresslist[i].num / (data.sum.total - data.sum.none)) * 100
        opt.style = "width: " + percent.toFixed(1) + "%"
        opt.setAttribute("aria-valuenow",progresslist[i].num)
        opt.setAttribute("aria-valuemin","0")
        opt.setAttribute("aria-valuemax",(data.sum.total - data.sum.none))
        opt.setAttribute("data-bs-toggle","tooltip")
        opt.setAttribute("data-bs-placement","top")
        let title = `${progresslist[i].title} ${percent.toFixed(1)}%(${progresslist[i].num})`
        opt.setAttribute("title", title)
        opt.innerHTML = title
        progress.appendChild(opt);
    }

    // Status 합 챠트
    let statusProgress = document.getElementById("sum-shot-status-progress")
    statusProgress.innerHTML = "" // 기존 챠트를 초기화 한다.
    
    for (let [key, value] of Object.entries(data.sum)) {
        if (value === 0 || key === "total" || key === "none"){ // 값이 0 이면 그리지 않는다.
            continue
        }
        if (key === "inprogress" || key === "onhold" || key === "finalapproved") {
            continue
        }
        let opt = document.createElement('div');
        opt.classList.add("progress-bar")
        opt.classList.add("bg-"+key)
        opt.role = "progressbar"
        let percent = (value / (data.sum.total - data.sum.none)) * 100
        opt.style = "width: " + percent.toFixed(1) + "%"
        opt.setAttribute("aria-valuenow",value)
        opt.setAttribute("aria-valuemin","0")
        opt.setAttribute("aria-valuemax",(data.sum.total - data.sum.none))
        opt.setAttribute("data-bs-toggle","tooltip")
        opt.setAttribute("data-bs-placement","top")
        opt.setAttribute("title", `${key} ${percent.toFixed(1)}%(${value})`)
        opt.innerHTML = `${key}<br>${percent.toFixed(1)}%(${value})`
        statusProgress.appendChild(opt);
    }

    //  None 값 설정
    document.getElementById("SumNoneStatusNum").innerHTML = data.sum.none
    
    initTooltip()
}

function initTooltip() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl)
    })
}

//전체 프리셋 버튼
function projectpresetAll(){
     // 전체 프리셋 버튼이 선택되면 전체 프로젝트 태그를 선택한다.
     let preset = document.getElementById("preset-all")
     let inputs = document.querySelectorAll('input')
     if (preset.checked === true) { //전체 프리셋 버튼을 켠다
         for (let i = 0; i < inputs.length; i++) {
            inputs[i].checked = true
         }
     
     } else { //전체 프리셋 버튼을 끈다
         for (let i = 0; i < inputs.length; i++) {
            inputs[i].checked = false
         }
     }
     RenderPerChart()
}

// 진행중 프리셋 버튼
function projectpresetPostProject(){
    // 진행중 프리셋 버튼이 선택되면 진행중인 프로젝트 태그를 선택한다.
    let preset = document.getElementById("preset-postproject")
    let inputs = document.querySelectorAll('input[status^="진행중"]')
    
    if (preset.checked === true) { //진행중 프리셋 버튼을 켠다
        for (let i = 0; i < inputs.length; i++) {
            inputs[i].checked = true
        }
    
    } else { //진행중 프리셋 버튼을 끈다
        for (let i = 0; i < inputs.length; i++) {
            inputs[i].checked = false
        }
    }
    RenderPerChart()
}

// 준비중 프리셋 버튼
function projectpresetPreProject(){
      // 준비중 프리셋 버튼이 선택되면 준비중인 프로젝트 태그를 선택한다.
      let preset = document.getElementById("preset-preproject")
      let inputs = document.querySelectorAll('input[status^="준비중"]')
      
      if (preset.checked === true) { //준비중 프리셋 버튼을 켠다
          for (let i = 0; i < inputs.length; i++) {
              inputs[i].checked = true
          }
      
      } else { //준비중 프리셋 버튼을 끈다
          for (let i = 0; i < inputs.length; i++) {
              inputs[i].checked = false
          }
      }
      RenderPerChart()
}

// 테스트 프리셋 버튼
function projectpresetTest(){
    // 테스트 프리셋 버튼이 선택되면 테스트인 프로젝트 태그를 선택한다.
    let preset = document.getElementById("preset-test")
    let inputs = document.querySelectorAll('input[status^="테스트"]')
    
    if (preset.checked === true) { //테스트 프리셋 버튼을 켠다
        for (let i = 0; i < inputs.length; i++) {
            inputs[i].checked = true
        }
    
    } else { //테스트 프리셋 버튼을 끈다
        for (let i = 0; i < inputs.length; i++) {
            inputs[i].checked = false
        }
    }
    RenderPerChart()
}

// 중단 프리셋 버튼
function projectpresetLayover(){
       // 중단 프리셋 버튼이 선택되면 중단인 프로젝트 태그를 선택한다.
       let preset = document.getElementById("preset-layover")
       let inputs = document.querySelectorAll('input[status^="중단"]')
       
       if (preset.checked === true) { //중단 프리셋 버튼을 켠다
           for (let i = 0; i < inputs.length; i++) {
               inputs[i].checked = true
           }
       
       } else { //중단 프리셋 버튼을 끈다
           for (let i = 0; i < inputs.length; i++) {
               inputs[i].checked = false
           }
       }
       RenderPerChart()
}

// 백업중 프리셋 버튼
function projectpresetBackup(){
       // 백업중 프리셋 버튼이 선택되면 백업중인 프로젝트 태그를 선택한다.
       let preset = document.getElementById("preset-backup")
       let inputs = document.querySelectorAll('input[status^="백업중"]')
       
       if (preset.checked === true) { //백업중 프리셋 버튼을 켠다
           for (let i = 0; i < inputs.length; i++) {
               inputs[i].checked = true
           }
       
       } else { //백업중 프리셋 버튼을 끈다
           for (let i = 0; i < inputs.length; i++) {
               inputs[i].checked = false
           }
       }
       RenderPerChart()
}

// 백업완료 프리셋 버튼
function projectpresetArchived(){
       // 백업완료 프리셋 버튼이 선택되면 백업완료인 프로젝트 태그를 선택한다.
       let preset = document.getElementById("preset-archive")
       let inputs = document.querySelectorAll('input[status^="백업완료"]')
       
       if (preset.checked === true) { //백업완료 프리셋 버튼을 켠다
           for (let i = 0; i < inputs.length; i++) {
               inputs[i].checked = true
           }
       
       } else { //백업완료 프리셋 버튼을 끈다
           for (let i = 0; i < inputs.length; i++) {
               inputs[i].checked = false
           }
       }
       RenderPerChart()
}

// 소송중 프리셋 버튼
function projectpresetLawsuit(){
   // 소송중 프리셋 버튼이 선택되면 소송중인 프로젝트 태그를 선택한다.
   let preset = document.getElementById("preset-lawsuit")
   let inputs = document.querySelectorAll('input[status^="소송중"]')
   
   if (preset.checked === true) { //소송중 프리셋 버튼을 켠다
       for (let i = 0; i < inputs.length; i++) {
           inputs[i].checked = true
       }
   
   } else { //소송중 프리셋 버튼을 끈다
       for (let i = 0; i < inputs.length; i++) {
           inputs[i].checked = false
       }
   }
   RenderPerChart()
}
