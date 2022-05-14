function submitCreds() {
  const form = document.getElementById('authForm');
  event.preventDefault();
  var usernameForm = form.elements['username'];
  var passwordForm = form.elements['password'];
  let emailSend = usernameForm.value;
  let passSend = passwordForm.value;
  var data = "username=" + emailSend + "&password=" + passSend
    var client = new HttpClient();
    var result = document.getElementById('authResult');
    const resultP = document.createElement('p');
    resultP.textContent =  "Authenticating...";
    result.innerHTML = '';
    result.appendChild(resultP);
    fetch("/api/sdk_auth?" + data)
    .then(response => response.text())
    .then((response) => {
      res = response.replace(/\s/g,'');
      result.innerHTML = '';
      if (`${res}` == "success") {
        resultP.textContent = "Authentication successful! Now you can use the app." 
        result.appendChild(resultP);
        const resultA = document.createElement('a');
        var resultAtext = document.createTextNode("Click here to return to the app")
        resultA.appendChild(resultAtext);
        resultA.title = "Click here to authorize";
        resultA.href = "/";
        result.appendChild(resultA);     
      } else if (`${res}` == "error") {
        resultP.textContent = "Invalid username or password. Please try again."
        result.appendChild(resultP);
      } else {
        resultP.textContent = "An unknown error has occurred."
        result.appendChild(resultP);
      };
    });
  };
