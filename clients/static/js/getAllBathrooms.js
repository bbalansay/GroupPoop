$(document).ready(() => {
  fetch("https://api.grouppoop.icu/bathroom", {
    method: 'GET',
    headers: {
      'Authorization': sessionStorage.auth
    }
  })
  .then(checkStatus)
  .then((resp) => resp.json())
  .then(populateBathrooms)
  .catch((err) => {
    console.log(err)
  })
})

const checkStatus = (response) => {
  if (response.status >= 200 && response.status < 300) {
    return response;
  } else {
    return Promise.reject(new Error());
  }
}

const populateBathrooms = (bathrooms) => {
  for (let bathroom of bathrooms) {
    let card = document.createElement("div")
    card.id = bathroom.id
    
    let body = document.createElement("div");
    body.className = "card-body";
    card.appendChild(body);
    
    let title = document.createElement("h4");
    title.className = "card-title";
    title.textContent = bathroom.Name
    body.appendChild(title)
    
    let info = document.createElement("p");
    info.className = "card-text";
    info.textContent = bathroom.Gender + " | " + bathroom.Location
    body.appendChild(info);

    let desc = document.createElement("p");
    desc.className = "card-text";
    desc.textContent = bathroom.Description.length > 50 ? bathroom.Description.substring(0, 46) + "..." : bathroom.Description;
    body.appendChild(desc);

    $("#bathrooms").appendChild(card);
  }
}