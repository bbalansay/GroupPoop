"use strict"

const express = require("express");
const mysql = require("promise-mysql");

const auth = require("./middleware/checkAuth")
const db = require("./middleware/getDB")
const allBath = require("./handlers/getAllBathrooms")
const getBath = require("./handlers/getBathroom")
const makeRev = require("./handlers/makeReview")
const editRev = require("./handlers/editReview")
const delRev = require("./handlers/deleteReview")

const app = express();
app.use(express.json())

const addr = process.env.BATHROOMPORT || ":80";
const [host, port] = addr.split(":")

//get all the bathrooms in the database
app.get("/bathroom", checkAuth(req, res, next), getDB(req, res, next), getAllBathrooms(req, res));

//gets a specific bathroom and all reviews
app.get("/bathroom/:bathroomID", checkAuth(req, res, next), getDB(req, res, next), getBathroom(req, res));

// create a review for a specific bathroom
app.post("/bathroom/:bathroomID/review", checkAuth(req, res, next), getDB(req, res, next), makeReview(req, res));

app.patch("/user/:userID/review/:reviewID", checkAuth(req, res, next), getDB(req, res, next), editReview(req, res));
app.delete("/user/:userID/review/:reviewID", checkAuth(req, res, next), getDB(req, res, next), deleteReview(req, res));

app.listen(port, host, () => {
    console.log(`server is listening at http://${addr}...`)
})