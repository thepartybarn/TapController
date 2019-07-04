function ChangeRelay(Relay,State) {
    $.ajax({
    type: "GET",
    data: {"Relay":Relay,
    "State":State},
    url: "/api/setRelay"
    });
  }

/*
getHomeData();
var intervalID = setInterval(function(){getHomeData();}, 3000);

function getHomeData() {
  $.ajax({
  dataType: "json",
  url: "/api/gethomedata",
  success: updateHomeData
});
}

function updateHomeData(data){

  if (data["logs"] != null){
    document.getElementById("logs").innerHTML = data["logs"].join("<br />")
  }

//document.getElementById("TimeOfLastPictureUpdateFromN8").innerHTML = data["_lastN8PictureUpdateTime"]
}

function getClip(){

  document.getElementById("DVRModalButton").disabled = true;
  document.getElementById("DVRModalButton").innerHTML = "<span id=\"Spinner\" class=\"spinner-border spinner-border-sm\" role=\"status\" aria-hidden=\"true\"></span>";


  $.ajax({
    type: "GET",
  data: {"startTime":document.getElementById("startTime").value,
  "endTime":document.getElementById("endTime").value,
  "cameras":[document.getElementById("c1").checked,
document.getElementById("c2").checked,
document.getElementById("c3").checked,
document.getElementById("c4").checked,
document.getElementById("c5").checked,
document.getElementById("c6").checked,
document.getElementById("c7").checked,
document.getElementById("c8").checked]},
  url: "/api/getClip",
  success: clipDownload,
  error: clipFailure,
  complete: clipComplete
  });
}
function clipDownload(data){
    window.open(data["fileName"] ,"_blank");
}

function clipComplete(){
  document.getElementById("DVRModalButton").disabled = false;
  document.getElementById("DVRModalButton").innerHTML = "Download";

}

function clipFailure(){
alert("Failed")
}
*/
