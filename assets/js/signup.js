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

document.getElementById("addOrganization").onclick = function() {
	addOrganization()
};

function addOrganization() {
	console.log("print");
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