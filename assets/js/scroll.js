// When the user scrolls down 20px from the top of the document, show the button
window.onscroll = function() {scrollFunction()};
function scrollFunction() {
    let topbtn = document.getElementById("topbtn");
	if (document.body.scrollTop > 20 || document.documentElement.scrollTop > 20) {
        topbtn.style.display = "block";
	} else {
        topbtn.style.display = "none";
	}
}