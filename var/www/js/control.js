var x = document.getElementById("sdkActions");
var keysKey = document.getElementById("keysKey");
keysKey.style.display = "none"
var useKeyboardControl = false
x.style.display = "none";

var isMovingForward = false
var isMovingLeft = false
var isMovingRight = false
var isMovingFL = false
var isMovingFR = false
var isMovingBack = false
var isMovingBL = false
var isMovingBR = false
var isStopped = false

function toggleKeyboard() {
    if (useKeyboardControl == false) {
        useKeyboardControl = true
        keysKey.style.display = "block"
    } else {
        useKeyboardControl = false
        keysKey.style.display = "none"
    }
}

function sdkInit() {
    sendForm('/api/initSDK')
    var x = document.getElementById("sdkActions");
    x.style.display = "block";
}

function sdkUnInit() {
    sendForm('/api/release_behavior_control')
    var x = document.getElementById("sdkActions");
    x.style.display = "none";
}

function sendForm(formURL) {
    let xhr = new XMLHttpRequest();
      xhr.open("POST", formURL);
      xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
      xhr.send();
      return
  }

let keysPressed = {};

keysPressed["w"] = false
keysPressed["a"] = false
keysPressed["s"] = false
keysPressed["d"] = false


document.addEventListener('keyup', (event) => {
    keysPressed[event.key] = false
    if (useKeyboardControl == true) {
    if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isStopped == false) {
            sendForm("/api/move_wheels?lw=0&rw=0")
        }
        isStopped = true
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == true && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isMovingForward == false) {
            sendForm("/api/move_wheels?lw=140&rw=140")
        }
        isStopped = false
        isMovingForward = true
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == true && keysPressed["a"] == true && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isMovingFL == false) {
            sendForm("/api/move_wheels?lw=100&rw=190")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = true
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == true && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == true) {
        if (isMovingFR == false) {
            sendForm("/api/move_wheels?lw=190&rw=100")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = true
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == true) {
        if (isMovingRight == false) {
            sendForm("/api/move_wheels?lw=150&rw=-150")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = true
    } else if (keysPressed["w"] == false && keysPressed["a"] == true && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isMovingLeft == false) {
            sendForm("/api/move_wheels?lw=-150&rw=150")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = true
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == true && keysPressed["d"] == false) {
        if (isMovingBack == false) {
            sendForm("/api/move_wheels?lw=-150&rw=-150")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = true
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == true && keysPressed["s"] == true && keysPressed["d"] == false) {
        if (isMovingBL == false) {
            sendForm("/api/move_wheels?lw=-100&rw=190")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = true
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == true && keysPressed["d"] == true) {
        if (isMovingBR == false) {
            sendForm("/api/move_wheels?lw=-190&rw=100")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = true
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    }
}
 });

document.addEventListener('keydown', function(event) {
    keysPressed[event.key] = true;
    if (useKeyboardControl == true) {
    if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isStopped == false) {
            sendForm("/api/move_wheels?lw=0&rw=0")
        }
        isStopped = true
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == true && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isMovingForward == false) {
            sendForm("/api/move_wheels?lw=140&rw=140")
        }
        isStopped = false
        isMovingForward = true
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == true && keysPressed["a"] == true && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isMovingFL == false) {
            sendForm("/api/move_wheels?lw=100&rw=190")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = true
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == true && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == true) {
        if (isMovingFR == false) {
            sendForm("/api/move_wheels?lw=190&rw=100")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = true
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == false && keysPressed["d"] == true) {
        if (isMovingRight == false) {
            sendForm("/api/move_wheels?lw=150&rw=-150")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = true
    } else if (keysPressed["w"] == false && keysPressed["a"] == true && keysPressed["s"] == false && keysPressed["d"] == false) {
        if (isMovingLeft == false) {
            sendForm("/api/move_wheels?lw=-150&rw=150")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = true
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == true && keysPressed["d"] == false) {
        if (isMovingBack == false) {
            sendForm("/api/move_wheels?lw=-150&rw=-150")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = false
        isMovingBack = true
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == true && keysPressed["s"] == true && keysPressed["d"] == false) {
        if (isMovingBL == false) {
            sendForm("/api/move_wheels?lw=-100&rw=190")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = true
        isMovingBR = false
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    } else if (keysPressed["w"] == false && keysPressed["a"] == false && keysPressed["s"] == true && keysPressed["d"] == true) {
        if (isMovingBR == false) {
            sendForm("/api/move_wheels?lw=-190&rw=100")
        }
        isStopped = false
        isMovingForward = false
        isMovingBL = false
        isMovingBR = true
        isMovingBack = false
        isMovingFL = false
        isMovingFR = false
        isMovingLeft = false
        isMovingRight = false
    }
}
});