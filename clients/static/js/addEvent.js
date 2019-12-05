window.onload = () => {
  setTimeout(() => {
    $("#fav").click((e) => {
      e.preventDefault()

      fetch("https://api.grouppoop.icu/favorites/" + params.get("id"), {
        method: 'POST',
        headers: {
          'Authorization': sessionStorage.auth
        }
      })
        .then(checkStatus)
        .catch(console.log)
    })

    $("#btnRev").click((e) => {
      e.preventDefault()

      fetch("https://api.grouppoop.icu/bathroom/" + params.get("id") + "/review", {
        method: 'POST',
        body: JSON.stringify({
          Score: $("#score").val(),
          Content: $("#content").val()
        }),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': sessionStorage.auth
        }
      })
        .then(checkStatus)
        .then(resp => resp.json())
        .catch(console.log)
      })
  }, 1000)
}