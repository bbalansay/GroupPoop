$(document).ready(() => {
  $("#signin").click(() => {
    window.location.href = "signin.html";
  });

  $("#submit").click((e) => {
    e.preventDefault()

    let newUser = {
      email: $("#email").val(),
      userName: $("#userName").val(),
      firstName: $("#firstName").val(),
      lastName: $("#lastName").val(),
      password: $("#password").val(),
      passwordConf: $("#passwordConf").val()
    }

    fetch("https://api.grouppoop.icu/user", {
      method: 'POST',
      body: JSON.stringify(newUser),
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(checkStatus)
      .then(redirect)
      .catch(() => {
        setTimeout(() => $("#alert").html(`<br><div class="alert alert-danger" role="alert">Unable to register an account with these credentials.</div>`), 500);
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

const redirect = () => {
  setTimeout(() => window.location.replace("signin.html"), 500);
}