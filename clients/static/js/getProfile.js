$(document).ready(() => {
  fetch("https://api.grouppoop.icu/user/me", {
    method: 'GET',
    headers: {
      'Authorization': sessionStorage.auth
    }
  })
  .then(checkStatus)
  .then((resp) => resp.json())
  .then(populateProfile)
  .catch((err) => {
    console.log(err)
  })
  
  fetch("https://api.grouppoop.icu/favorites", {
    method: 'GET',
    headers: {
      'Authorization': sessionStorage.auth
    }
  })
  .then(checkStatus)
  .then((resp) => resp.json())
  .then(populateFavorites)
  .catch((err) => {
    console.log(err)
  })
  
  $("#edit-profile-button").click(() => {
    let updates = {
      firstName: $("#firstname-update").val(),
      lastName: $("#lastname-update").val(),
    }
    
    fetch("https://api.grouppoop.icu/user/me", {
      method: 'PATCH',
      body: JSON.stringify(updates),
      headers: {
        'Content-Type': 'application/json',
        'Authorization': sessionStorage.auth
      }
    })
      .then(checkStatus)
      .then(window.location.reload())
      .catch(() => {
        setTimeout(() => $("#alert").html(`<br><div class="alert alert-danger" role="alert">Unable to register an account with these credentials.</div>`), 1000);
        setTimeout(() => $("#alert").html(""), 5000);
      })
  })
  
  
})

const populateFavorites = async (res) => {
  for (let fav in res.favorites) {
    let currFav = await getBathroom(fav.bathroomID);
    $("#favorites").append(currFav);
  }
}

const populateProfile = (profile) => {
  let user = profile.user;
  let userInfo = document.createElement("div");
  
  let title = document.createElement("h3");
  title.textContent = "Profile";
  
  let userName = document.createElement("h5");
  userName.textContent = "Username: " + user.userName;
  userInfo.appendChild(userName);
  
  let userFirstName = document.createElement("h5");
  userFirstName.textContent = "First Name: " + user.firstName;
  userInfo.appendChild(userFirstName);
  
  let userLastName = document.createElement("h5");
  userLastName.textContent = "Last Name: " + user.lastName;
  userInfo.appendChild(userLastName);
  
  
  let userReviews = document.createElement("div");
  let reviewTitle = document.createElement("h3");
  reviewTitle.textContent = "Your Reviews";
  userReviews.appendChild(reviewTitle);
  for (let review of profile.reviews) {
    let currReview = document.createElement("div");
    let bathID = document.createElement("p");
    bathID.textContent = review.bathroomID;
    currReview.appendChild(bathID);
    let score = document.createElement("p");
    score.textContent = "Score: " + review.score;
    currReview.appendChild(score);
    let content = document.createElement("p");
    content.textContent = review.content;
    currReview.appendChild(content);
    userReviews.appendChild(currReview);
  }
  
  
  $("#profile").append(userInfo);
  $("#reviews").append(userReviews);
}

async function getBathroom(id) {
  res = await fetch("https://api.grouppoop.icu/bathroom/" + id, {
    method: 'GET',
    headers: {
      'Authorization': sessionStorage.auth
    }
  })
  .catch((err) => {
    console.log(err);
  })
  console.log(res.body);
  let currFav = document.createElement("div");
  let bathName = document.createElement("h5");
  bathName.textContent = res.body.name;
  currFav.appendChild(bathName);
  return currFav;
}