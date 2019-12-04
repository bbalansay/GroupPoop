fetch("https://api.grouppoop.icu/bathroom", {
  method: 'GET',
  headers: {
    'Authorization': sessionStorage.auth
  }
})
  .then(checkStatus)
  .then((resp) => resp.json())
  .then(console.log(JSON.stringify(data)))
  .catch((err) => {
    alert(err)
  })

const checkStatus = (response) => {
  if (response.status >= 400) {
    return Promise.reject(new Error());
  }

  return response;
}

const clearStorage = () => {
  sessionStorage.clear()
}
