
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
        document.getElementById("this-month").innerHTML = obj.total
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
        document.getElementById("next-month").innerHTML = obj.total
    })
    .catch((err) => {
        console.log(err)
    });
}

