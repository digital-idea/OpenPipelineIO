let userData = {
	"ID":"",
	"Password":"",
	"ConfirmPassword":"",
	"FirstNameKor":"",
	"LastNameKor":"",
	"FirstNameEng":"",
	"LastNameEng":"",
	"FirstNameChn":"",
	"LastNameChn":"",
	"Email":"",
	"EmailExternal":"",
	"Phone":"",
	"Hotline":"",
	"Organizations": [],
	"AdditionalTags": "",
	"Location": "",
	"Timezone": "",
	"CaptchaID": "",
	"CaptchaNum": "",
}

let organization = {
	"ID":"",
	"Division":"",
	"DivisionName":"",
	"Department":"",
	"DepartmentName":"",
	"Team":"",
	"TeamName":"",
	"Role":"",
	"RoleName":"",
	"Position":"",
	"PositionName":"",
}

// Restricts input for the given textbox to the given inputFilter.
function setInputFilter(textbox, inputFilter) {
	["input", "keydown", "keyup", "mousedown", "mouseup", "select", "contextmenu", "drop"].forEach(function(event) {
	  textbox.addEventListener(event, function() {
		if (inputFilter(this.value)) {
		  this.oldValue = this.value;
		  this.oldSelectionStart = this.selectionStart;
		  this.oldSelectionEnd = this.selectionEnd;
		} else if (this.hasOwnProperty("oldValue")) {
		  this.value = this.oldValue;
		  this.setSelectionRange(this.oldSelectionStart, this.oldSelectionEnd);
		}
	  });
    });
}


// addOrganization 버튼을 누르면 조직이 추가된다.
document.getElementById("addOrganization").onclick = function() {
	addOrganization();
	// add Tags in Organizations
	renderOrganizations();
};

function addOrganization() {
	let currentDate = new Date();
	orgnz = Object.create(organization)
	orgnz.ID = currentDate.getTime();
	orgnz.Division = document.getElementById("Division").value;
	orgnz.Department = document.getElementById("Department").value;
	orgnz.Team = document.getElementById("Team").value;
	orgnz.Role = document.getElementById("Role").value;
	orgnz.Position = document.getElementById("Position").value;
	// Name save
	var division = document.getElementById("Division");
	orgnz.DivisionName = division.options[division.selectedIndex].text;
	var department = document.getElementById("Department");
	orgnz.DepartmentName = department.options[department.selectedIndex].text;
	var team = document.getElementById("Team");
	orgnz.TeamName = team.options[team.selectedIndex].text;
	var role = document.getElementById("Role");
	orgnz.RoleName = role.options[role.selectedIndex].text;
	var position = document.getElementById("Position");
	orgnz.PositionName = position.options[position.selectedIndex].text;

	userData.Organizations.push(orgnz);
}

function renderOrganizations() {
	document.getElementById("Organizations").innerHTML = "";
	for (let i = 0; i < userData.Organizations.length; i++) {
		let div = document.createElement("div");
		div.setAttribute("id", userData.Organizations[i].ID);
		div.setAttribute("class", "alert alert-warning small");
		div.setAttribute("role", "alert");
		div.innerHTML += `${userData.Organizations[i].DivisionName}, `;
		div.innerHTML += `${userData.Organizations[i].DepartmentName}, `;
		div.innerHTML += `${userData.Organizations[i].TeamName}, `;
		div.innerHTML += `${userData.Organizations[i].RoleName}, `;
		div.innerHTML += `${userData.Organizations[i].PositionName} `;
		div.innerHTML += `<a href="#" class="alert-link">&bigotimes;</a>`;
		div.onclick = removeItem;
		document.getElementById("Organizations").appendChild(div);
	}
}

function removeItem(e) {
	id = e.target.parentElement.getAttribute("id");
	for (i = 0; i < userData.Organizations.length; i++) {
		if ( userData.Organizations[i].ID == id ) {
			// console.log(id);
			userData.Organizations.splice(i,1);
		}
	}
	renderOrganizations()
}

setInputFilter(document.getElementById("Hotline"), function(value) {
	return /^\d*$/.test(value);
});
setInputFilter(document.getElementById("Phone"), function(value) {
	return /^\d*$/.test(value);
});
setInputFilter(document.getElementById("CaptchaID"), function(value) {
	return /^\d*$/.test(value);
});

// SignUp 버튼을 누르면 가입이 된다.
document.getElementById("SignUp").onclick = function() {
	addUser()
};

function addUser() {
	// check Error
	if (document.getElementById("ID").value == "") {
		alert("ID는 빈 문자열이 될 수 없습니다.");
		return
	}
	if (document.getElementById("Password").value == "") {
		alert("Password는 빈 문자열이 될 수 없습니다.");
		return
	}
	if (document.getElementById("ConfirmPassword").value == "") {
		alert("ConfirmPassword는 빈 문자열이 될 수 없습니다.");
		return
	}
	if (document.getElementById("Password").value != document.getElementById("ConfirmPassword").value) {
		alert("Password가 서로 다릅니다.");
		return
	}
	if (document.getElementById("CaptchaNum").value == "") {
		alert("CaptchaNum는 빈 문자열이 될 수 없습니다.");
		return
	}

	if (userData.Organizations.length == 0) {
		alert("조직 정보를 등록하지 않았습니다. 그대로 가입이 진행됩니다.");
	}
	userData.ID = document.getElementById("ID").value;
	userData.Password = document.getElementById("Password").value;
	userData.ConfirmPassword = document.getElementById("ConfirmPassword").value;
	userData.FirstNameKor = document.getElementById("FirstNameKor").value;
	userData.LastNameKor = document.getElementById("LastNameKor").value;
	userData.FirstNameEng = document.getElementById("FirstNameEng").value;
	userData.LastNameEng = document.getElementById("LastNameEng").value;
	userData.FirstNameChn = document.getElementById("FirstNameChn").value;
	userData.LastNameChn = document.getElementById("LastNameChn").value;
	userData.Email = document.getElementById("Email").value;
	userData.EmailExternal = document.getElementById("EmailExternal").value;
	userData.Phone = document.getElementById("Phone").value;
	userData.Hotline = document.getElementById("Hotline").value;
	userData.AdditionalTags = document.getElementById("AdditionalTags").value;
	userData.Location = document.getElementById("Location").value;
	userData.Timezone = document.getElementById("Timezone").value;
	userData.CaptchaID = document.getElementById("CaptchaID").value;
	userData.CaptchaNum = document.getElementById("CaptchaNum").value;

    $.ajax({
        url: "/api/adduser",
		type: "post",
		data: userData,
        dataType: "json",
        success: function(data) {
            console.info(data)
        }
	});
	location.replace("/invalidaccess")
}