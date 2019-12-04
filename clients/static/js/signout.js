const signout = () => {
  fetch("https://api.grouppoop.icu/login/mine", {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': sessionStorage.auth
    }
  })
    .then(checkStatus)
    .then(clearStorage)
    .then(redirect)
    .catch((err) => {
      alert(err)
    })
}

const checkStatus = (response) => {
  if (response.status >= 400) {
    return Promise.reject(new Error());
  }

  return response;
}

const clearStorage = () => {
  sessionStorage.clear()
}

const redirect = () => {
  setTimeout(() => window.location.replace("/auth/signin.html"), 1000);
}