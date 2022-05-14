var colorPicker = new iro.ColorPicker("#picker", {
  width: 250,
  layout: [
  { 
    component: iro.ui.Wheel,
  }
  ]
});

var client = new HttpClient();
getCurrentSettings()
document.querySelectorAll('.serverEscape').forEach(item => {
  item.addEventListener('click', event => {
    let confirmAction = confirm("This will change Vector's server environment to Escape Pod. This will NOT clear user data but will make parts of this web app inoperable and will restart onboarding. Vector's personality will remain intact. Would you like to continue?");
    if (confirmAction) {
      fetch("/api/server_escape")
      alert("Executing. After a while, Vector's screen should show 'Configuring'. Once Vector reaches the onboarding screen, the process has finished and you can start onboarding Vector. IMPORTANT: The Escape Pod's built-in settings implementation will not work. This is normal and is currently being worked on. Once he is set up, press OK then refresh this page. After that, press the 'Skip Onboarding' button (located after the server environment settings).");
    }
  })
})
document.querySelectorAll('.serverProd').forEach(item => {
  item.addEventListener('click', event => {
    let confirmAction = confirm("This will change Vector's server environment to Production (normal, stock). This will NOT clear user data but may affect the functionality of this web app and will restart onboarding. Vector's personality will remain intact. Would you like to continue?");
    if (confirmAction) {
      fetch("/api/server_prod")
      alert("Executing. After a while, Vector's screen should show 'Configuring'. Once Vector reaches the onboarding screen, the process has finished and you can start onboarding Vector. Once he is set up, press OK then refresh this page.");
    }
  })
})
function sendForm(formURL) {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", formURL);
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  xhr.send();
  xhr.onload = function() { 
    getCurrentSettings()
  };
}
function sendFormRainbowOff(formURL) {
  let xhr = new XMLHttpRequest();
  xhr.open("POST", formURL);
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  xhr.send();
  xhr.onload = function() { 
    getCurrentSettings()
  };
}
var as = document.getElementById('authStatus');
const asP = document.createElement('p');
asP.textContent =  "Checking for info...";
as.innerHTML = '';
as.appendChild(asP);
fetch("/api/get_auth_status")
.then(response => response.text())
.then((response) => {
 res = response.replace(/\s/g,'');
 as.innerHTML = '';
 if (`${res}` == "notauthorized1") {
  asP.textContent = "You must sign into your Anki account for many of these functions to work. This can be done at the following link. If Vector needs to be onboarded, that can be done here too:"
  as.appendChild(asP);
  const asA = document.createElement('a');
  var asAtext = document.createTextNode("Click here to authorize")
  asA.appendChild(asAtext);
  asA.title = "Click here to authorize";
  asA.href = "/auth.html";
  as.appendChild(asA);
} else if (`${res}` == "notauthorized2") {
  asP.textContent = "Signing in was attempted, but failed. Authorization is required for many functions here to work. This can be done at this link:" 
  as.appendChild(asP);
  const asA = document.createElement('a');
  var asAtext = document.createTextNode("Click here to authorize")
  asA.appendChild(asAtext);
  asA.title = "Click here to authorize";
  asA.href = "/auth.html";
  as.appendChild(asA);
} else if (`${res}` == "authorized") {
  asP.textContent = "App is authorized and everything should be working!"
  as.appendChild(asP);
} else if (`${res}` == "escapepod") {
  asP.textContent = "Bot is using Escape Pod so many of these functions will not work."
  as.appendChild(asP);
} else {
  asP.textContent = "An unknown error has occured. This app is likely not authenticated with your Anki account. This is required for many functions here to work. Authentication can be done at this link:"
  as.appendChild(asP);
  const asA = document.createElement('a');
  var asAtext = document.createTextNode("Click here to authorize")
  asA.appendChild(asAtext);
  asA.title = "Click here to authorize";
  asA.href = "/auth.html";
  as.appendChild(asA);
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
  xhr.onload = function() { 
    getCurrentSettings()
  };
};

function getCurrentSettings() {
    setTimeout(function(){
  let xhr = new XMLHttpRequest();
  xhr.open("POST", "/api/get_current_settings");
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  xhr.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
  xhr.responseType = 'json';
  xhr.send();
  xhr.onload = function() {
    //console.log(data)
    console.log(xhr.response)
    var jdocSettings = xhr.response
    let xhr2 = new XMLHttpRequest();
    xhr2.open("POST", "/api/rainbow_status");
    xhr2.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr2.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr2.send();
    xhr2.onload = function() {
     var rainbowStatus = xhr2.response
     res = rainbowStatus.replace(/\s/g,'');
     if (`${res}` == "on") {
      rainbowEye = "on"
    } else {
      rainbowEye = "off"
    }
    let xhr3 = new XMLHttpRequest();
    xhr3.open("POST", "/api/server_status");
    xhr3.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr3.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr3.send();
    xhr3.onload = function() {
     var serverResponse = xhr3.response
     res = serverResponse.replace(/\s/g,'');
     if (`${res}` == "escape") {
      serverStatus = "Escape Pod"
    } else if (`${res}` == "prod") {
      serverStatus = "Production"
    } else {
      serverStatus = "Unknown"
    }
    let xhr4 = new XMLHttpRequest();
    xhr4.open("POST", "/api/snore_status");
    xhr4.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr4.setRequestHeader("Cache-Control", "no-cache, no-store, max-age=0");
    xhr4.send();
    xhr4.onload = function() {
      var snoreResponse = xhr4.response
      res = snoreResponse.replace(/\s/g,'');
      if (`${res}` == "on") {
        snoreStatus = "Enabled"
      } else {
        snoreStatus = "Disabled"
      }
      if ( jdocSettings["jdoc"]["custom_eye_color"]) {
        var customECE = jdocSettings["jdoc"]["custom_eye_color"]["enabled"]
        var customECH = jdocSettings["jdoc"]["custom_eye_color"]["hue"]
        var customECS = jdocSettings["jdoc"]["custom_eye_color"]["saturation"]
      }
      var eyeColorS = jdocSettings["jdoc"]["eye_color"]
      var volumeS = jdocSettings["jdoc"]["master_volume"]
      var localeS = jdocSettings["jdoc"]["locale"]
      var timeSetS = jdocSettings["jdoc"]["clock_24_hour"]
      var tempFormatS = jdocSettings["jdoc"]["temp_is_fahrenheit"]
      var buttonS = jdocSettings["jdoc"]["button_wakeword"]
      if (`${rainbowEye}` == "on") {
       var eyeColorT = "Rainbow"
     } else if ( jdocSettings["jdoc"]["custom_eye_color"]) {
       if (`${customECE}` == "true") {
         var setHue = customECH * 360
         var setHue = setHue.toFixed(3)
         var setSat = customECS * 100
         var setSat = setSat.toFixed(3)
         colorPicker.color.hsl = { h: setHue, s: setSat, l: 50 };     
         var eyeColorT = "Custom"
       } else { 
         if (`${eyeColorS}` == 0) {
          var eyeColorT = "Teal"
        } else if  (`${eyeColorS}` == 1) {
          var eyeColorT = "Orange"
        } else if  (`${eyeColorS}` == 2) {
          var eyeColorT = "Yellow"
        } else if  (`${eyeColorS}` == 3) {
          var eyeColorT = "Lime Green"
        } else if  (`${eyeColorS}` == 4) {
          var eyeColorT = "Azure Blue"
        } else if  (`${eyeColorS}` == 5) {
          var eyeColorT = "Purple"
        } else if  (`${eyeColorS}` == 6) {
          var eyeColorT = "Matrix Green"
        } else {
          var eyeColorT = "none"
        }  
      } } else { if (`${eyeColorS}` == 0) {
        var eyeColorT = "Teal"
      } else if  (`${eyeColorS}` == 1) {
        var eyeColorT = "Orange"
      } else if  (`${eyeColorS}` == 2) {
        var eyeColorT = "Yellow"
      } else if  (`${eyeColorS}` == 3) {
        var eyeColorT = "Lime Green"
      } else if  (`${eyeColorS}` == 4) {
        var eyeColorT = "Azure Blue"
      } else if  (`${eyeColorS}` == 5) {
        var eyeColorT = "Purple"
      } else if  (`${eyeColorS}` == 6) {
        var eyeColorT = "Matrix Green"
      } else {
        var eyeColorT = "none"
      } }
      function pEyeColorST() {
       if (`${eyeColorS}` == 0) {
        var eyeColorT = "Teal"
      } else if (`${eyeColorS}` == 1) {
        var eyeColorT = "Orange"
      } else if (`${eyeColorS}` == 2) {
        var eyeColorT = "Yellow"
      } else if (`${eyeColorS}` == 3) {
        var eyeColorT = "Lime Green"
      } else if (`${eyeColorS}` == 4) {
        var eyeColorT = "Azure Blue"
      } else if (`${eyeColorS}` == 5) {
        var eyeColorT = "Purple"
      } else if (`${eyeColorS}` == 6) {
        var eyeColorT = "Matrix Green"
      } else {
        var eyeColorT = "none"
      }
    }
    if (`${volumeS}` == 0) {
      var volumeT = "Mute"
    } else if (`${volumeS}` == 1) {
      var volumeT = "Low"
    } else if (`${volumeS}` == 2) {
      var volumeT = "Medium Low"
    } else if (`${volumeS}` == 3) {
      var volumeT = "Medium"
    } else if (`${volumeS}` == 4) {
      var volumeT = "Medium High"
    } else if (`${volumeS}` == 5) {
      var volumeT = "High"
    } else {
      var volumeT = "none"
    }
    if (`${timeSetS}` == "false") {
      var timeSetT = "12 Hour"
    } else {
      var timeSetT = "24 Hour"
    }
    if (`${tempFormatS}` == "true") {
      var tempFormatT = "Fahrenheit"
    } else {
      var tempFormatT = "Celcius"
    }
    if (`${buttonS}` == 0) {
      var buttonT = "Hey Vector"
    } else {
      var buttonT = "Alexa"
    }
  //s1 = volume, s2 = eye-color, s3 = locale
  var s1 = document.getElementById('currentVolume');
  const s1P = document.createElement('p');
  s1P.textContent = "Current Volume: " + volumeT
  s1.innerHTML= ''
  s1.appendChild(s1P);
  var s2 = document.getElementById('currentEyeColor');
  const s2P = document.createElement('p');
  s2P.textContent = "Current Eye Color: " + eyeColorT
  s2.innerHTML = ''
  s2.appendChild(s2P);
  var s3 = document.getElementById('currentLocale');
  const s3P = document.createElement('p');
  s3P.textContent = "Current Locale: " + localeS
  s3.innerHTML = ''
  s3.appendChild(s3P);
  var s4 = document.getElementById('currentTimeSet');
  const s4P = document.createElement('p');
  s4P.textContent = "Current Time Format: " + timeSetT
  s4.innerHTML = ''
  s4.appendChild(s4P);
  var s5 = document.getElementById('currentTempFormat');
  const s5P = document.createElement('p');
  s5P.textContent = "Current Temp Format: " + tempFormatT
  s5.innerHTML = ''
  s5.appendChild(s5P);
  var s6 = document.getElementById('currentButton');
  const s6P = document.createElement('p');
  s6P.textContent = "Current Button Setting: " + buttonT
  s6.innerHTML = ''
  s6.appendChild(s6P);
  var s7 = document.getElementById('currentServer');
  const s7P = document.createElement('p');
  s7P.textContent = "Current Server Environment: " + `${serverStatus}`
  s7.innerHTML = ''
  s7.appendChild(s7P);
  var s8 = document.getElementById('currentSnore');
  const s8P = document.createElement('p');
  s8P.textContent = "Current Snore Setting: " + `${snoreStatus}`
  s8.innerHTML = ''
  s8.appendChild(s8P);
};
};
};
};
}, 300);
};
