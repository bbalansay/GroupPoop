"use strict"

const express = require("express");
const mysql = require("promise-mysql");

import {checkAuth} from './middleware/checkAuth';
import {getDB} from './middleware/getDB';
import {getAllBathrooms} from './handlers/getAllBathrooms';
import {getBathroom} from './handlers/getBathroom';
import {makeReview} from './handlers/makeReview';
import {editReview} from './handlers/editReview';
import {editReview} from './handlers/deleteReview';

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

app.listen(port, host, () => {
    console.log(`server is listening at http://${addr}...`)
})