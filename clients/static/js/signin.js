$(document).ready(() => {
  $("#register").click(() => {
    window.location.href = "register.html";
  });

  $("#submit").click((e) => {
    e.preventDefault()

    let credentials = {
      email: $("#email").val(),
      password: $("#password").val()
    }

    fetch("https://api.grouppoop.icu/login", {
      method: 'POST',
      body: JSON.stringify(credentials),
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(checkStatus)
      .then(saveSession)
      .then((resp) => resp.json())
      .then(saveProfile)
      .then(redirect)
      .catch(() => {
        setTimeout(() => $("#alert").html(`<br><div class="alert alert-danger" role="alert">Could not sign in using these credentials.</div>`), 500);
        setTimeout(() => $("#alert").html(""), 5000);
      })
  })
})

const checkStatus = (response) => {
  if (response.status >= 200 && response.status < 300) {
    return response;
  } else {
    return Promise.reject(new Error());
  }
}

const saveSession = (response) => {
  for (let pair of response.headers.entries()) {
    if (pair[0] == "authorization") {
      sessionStorage.auth = pair[1];
      console.log("Set sessionStorage.auth to: " + pair[1])
      return response;
    }
  }

  return Promise.reject(new Error());
}

const saveProfile = (json) => {
  sessionStorage.profile = JSON.stringify(json)
  console.log("Set sessionStorage.profile to: " + JSON.stringify(json))
  return json
}

const redirect = () => {
  setTimeout(() => window.location.replace("../index.html"), 500);
}