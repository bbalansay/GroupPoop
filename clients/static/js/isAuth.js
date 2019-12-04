if (!sessionStorage.auth && !sessionStorage.profile) {
  if (window.location.href.substring(window.location.href.lastIndexOf('/') + 1) != "signin.html" &&
      window.location.href.substring(window.location.href.lastIndexOf('/') + 1) != "register.html") {
    window.location.replace("/auth/signin.html");
  }
} else if (window.location.href.substring(window.location.href.lastIndexOf('/') + 1) == "signin.html" ||
           window.location.href.substring(window.location.href.lastIndexOf('/') + 1) == "register.html") {
  window.location.replace("/index.html");
}