// changeReportExcelURI 함수는 프로젝트를 선택하면 reportexcel url을 설정한다.
function changeReportExcelURI() {
    let project = CurrentProject();
    document.getElementById("reportexcelURI").href = "/reportexcel?project=" + project;
}