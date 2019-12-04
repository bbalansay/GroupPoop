"use strict"

const express = require("express");

const auth = require("./middleware/checkAuth")
const db = require("./middleware/getDB")
const allBath = require("./handlers/getAllBathrooms")
const getBath = require("./handlers/getBathroom")
const makeRev = require("./handlers/makeReview")
const editRev = require("./handlers/editReview")
const delRev = require("./handlers/deleteReview")
const addFav = require("./handlers/addFavorite")
const getFav = require("./handlers/getFavorites")

const app = express();
app.use(express.json())

const addr = process.env.BATHROOMPORT || ":80";
const [host, port] = addr.split(":")

//get all the bathrooms in the database
app.get("/bathroom", auth.checkAuth, db.getDB, allBath.getAllBathrooms);

//gets a specific bathroom and all reviews
app.get("/bathroom/:bathroomID", auth.checkAuth, db.getDB, getBath.getBathroom);

// create a review for a specific bathroom
app.post("/bathroom/:bathroomID/review", auth.checkAuth, db.getDB, makeRev.makeReview);

app.patch("/user/:userID/review/:reviewID", auth.checkAuth, db.getDB, editRev.editReview);
app.delete("/user/:userID/review/:reviewID", auth.checkAuth, db.getDB, delRev.deleteReview);

app.get("/user/:userID/favorites", auth.checkAuth, db.getDB, getFav.getFavorites);
app.post("/user/:userID/favorites", auth.checkAuth, db.getDB, addFav.addFavorite)

app.listen(port, host, () => {
    console.log(`server is listening at http://${addr}...`)
})