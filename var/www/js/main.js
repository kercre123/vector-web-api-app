var colorPicker = new iro.ColorPicker("#picker", {
  width: 250,
  layout: [
  { 
    component: iro.ui.Wheel,
  }
  ]
});

var client = new HttpClient();
var notices = document.getElementById('notices');
const noticesP = document.createElement('p');
noticesP.textContent =  "Checking for info...";
notices.innerHTML = '';
notices.appendChild(noticesP);
fetch("/api/get_notices")
.then(response => response.text())
.then((response) => {
 res = response.replace(/\s/g,'');
 notices.innerHTML = '';
 if (`${res}` == "notauthorized1") {
  noticesP.textContent = "You must sign into your Anki account for many of these functions to work. This can be done at this link:"
  notices.appendChild(noticesP);
  const noticesA = document.createElement('a');
  var noticesAtext = document.createTextNode("Click here to authorize")
  noticesA.appendChild(noticesAtext);
  noticesA.title = "Click here to authorize";
  noticesA.href = "/auth.html";
  notices.appendChild(noticesA);
} else if (`${res}` == "notauthorized2") {
  noticesP.textContent = "Signing in was attempted, but failed. Authorization is required for many functions here to work. This can be done at this link:" 
  notices.appendChild(noticesP);
  const noticesA = document.createElement('a');
  var noticesAtext = document.createTextNode("Click here to authorize")
  noticesA.appendChild(noticesAtext);
  noticesA.title = "Click here to authorize";
  noticesA.href = "/auth.html";
  notices.appendChild(noticesA);
} else if (`${res}` == "authorized") {
  noticesP.textContent = "App is authorized and everything should be working!"
  notices.appendChild(noticesP);
} else {
  noticesP.textContent = "An unknown error has occured. This app is likely not authenticated with your Anki account. This is required for many functions here to work. Authentication can be done at this link:"
  notices.appendChild(noticesP);
  const noticesA = document.createElement('a');
  var noticesAtext = document.createTextNode("Click here to authorize")
  noticesA.appendChild(noticesAtext);
  noticesA.title = "Click here to authorize";
  noticesA.href = "/auth.html";
  notices.appendChild(noticesA);
}
});

function sendCustomColor() {
  var pickerHue = colorPicker.color.hue;
  var pickerSat = colorPicker.color.saturation;
  var sendHue = pickerHue / 360
  var sendHue = sendHue.toFixed(3)
  var sendSat = pickerSat / 100
  var sendSat = sendSat.toFixed(3)
  let data = "hue=" + sendHue + "&" + "sat=" + sendSat
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "/api/custom_eye_color");
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  xhr.send(data);
};
