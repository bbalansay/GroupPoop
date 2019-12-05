let params = new URLSearchParams(window.location.search);

$(document).ready(() => {
  fetch("https://api.grouppoop.icu/bathroom/" + params.get("id"), {
    method: 'GET',
    headers: {
      'Authorization': sessionStorage.auth
    }
  })
  .then(checkStatus)
  .then((resp) => resp.json())
  .then(populateBathAndReviews)
  .catch(() => {
    let head1 = document.createElement("h1")
    head1.textContent = "Whoops!"
    $("#bathroom").append(head1);
    let head4 = document.createElement("h4")
    head4.textContent = "We ran into an error... maybe focus on your poop."
    $("#bathroom").append(head4);
  })
})

const populateBathAndReviews = (bathAndRev) => {
  let bath = bathAndRev.bathroom[0]
  let reviews = bathAndRev.reviews

  let header = document.createElement("h1")
  header.textContent = bath.Name
  $("#bathroom").append(header);

  let fav = document.createElement("input")
  fav.id = "fav"
  fav.type = "submit"
  fav.value = "Favorite"
  fav.style = "position: absolute; top: 12%; right: 7%"
  fav.className = "btn btn-primary"
  $("#bathroom").append(fav)


  let desc = document.createElement("h6")
  desc.textContent = bath.Description
  $("#bathroom").append(desc);

  let details = document.createElement("h6")
  details.textContent = bath.Gender + " | " + bath.Location
  $("#bathroom").append(details);

  let numerics = document.createElement("p")
  numerics.style = "text-align: left;"
  numerics.innerHTML = "Gender: " + bath.Gender + "<br>" +
                       "Location: " + bath.Location + "<br><br>" + 
                       "# of Sinks: " + bath.NumSinks + "<br>" +
                       "# of Toilets: " + bath.NumToilets + "<br>" +
                       "# of Urinals: " + bath.NumUrinals + "<br>" +
                       "# of Trash Cans: " + bath.NumTrashCans + "<br>" +
                       "# of Air Dryers: " + bath.NumAirDryers + "<br>" +
                       "# of Towel Dispensers: " + bath.NumTowelDispensers + "<br><br>"
  $("#bathroom").append(numerics);


  let makeRevHead = document.createElement("h4")
  makeRevHead.textContent = "Make a Review:"
  makeRevHead.style = "text-align: left"
  $("#bathroom").append(makeRevHead);

  let txt = document.createElement("textarea")
  txt.id = "rev"
  txt.name = "rev"
  txt.rows = "5"
  txt.cols = "100"
  $("#bathroom").append(txt);

  let btn = document.createElement("input")
  btn.id = "btnRev"
  btn.type = "submit"
  btn.value = "Submit"
  btn.className = "btn btn-primary"
  $("#bathroom").append(btn)


  let revHead = document.createElement("h4")
  revHead.textContent = "Reviews"
  revHead.style = "text-align: left"
  $("#bathroom").append(revHead);


  for (let review of reviews) {
    let card = document.createElement("div")
    card.style = "text-align: left"
    
    let body = document.createElement("div");
    body.className = "card-body";
    card.appendChild(body);
    
    let title = document.createElement("h6");
    title.className = "card-title";
    title.textContent = "Creator: User " + review.UserID
    body.appendChild(title)

    let content = document.createElement("p");
    content.className = "card-text";
    content.innerHTML = "Score: " + review.Score + "<br>" + 
                          "Review: "+ review.Content + "<br>" + 
                          "Created At: " + review.CreatedAt + " | " + "Edited At: " + review.EditedAt + "<br><br>"
    body.appendChild(content);


    $("#bathroom").append(card);
  }
}

// $(fav).click((e) => {
//   e.preventDefault()

  
// })

async function postFav() {
  fetch("https://api.grouppoop.icu/favorites/" + params.get("id"), {
    method: 'POST',
    headers: {
      'Authorization': sessionStorage.auth
    }
  })
  .then(checkStatus)
  .catch(console.log)
}