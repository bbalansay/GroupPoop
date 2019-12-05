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
})

const populateProfile = (profile) => {
  
  console.log(profile);
  
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
    
  }
  
  
  $("#profile").append(userInfo);
  $("#profile").append(userReviews);
}