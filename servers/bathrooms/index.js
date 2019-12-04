"use strict"

const express = require("express");
const mysql = require("promise-mysql");

require("./middleware/checkAuth")
require("./handlers/getAllBathrooms")
require("./handlers/getBathroom")
require("./handlers/makeReview")
require("./handlers/editReview")
require("./handlers/deleteReview")

const app = express();
app.use(express.json())

const addr = process.env.BATHROOMPORT || ":80";
const [host, port] = addr.split(":")

const dbHost = process.env.DBHOST
const dbPort = process.env.DBPORT
const dbUser = process.env.DBUSER
const dbPass = process.env.MYSQL_ROOT_PASSWORD
const dbName = process.env.DBNAME

async function getDB() {
    let db = await mysql.createConnection({
        host: dbHost,
        port: dbPort,
        user: dbUser,
        password: dbPass,
        database: dbName
    });
    return db;
}

//get all the bathrooms in the database
app.get("/bathroom", getAllBathrooms(req, res));

//gets a specific bathroom and all reviews
app.get("/bathroom/:bathroomID", getBathroom(req, res));

// create a review for a specific bathroom
<<<<<<< HEAD
app.post("/bathroom/:bathroomID/review", checkAuth(req, res), makeReview(req, res));

app.patch("/user/:userID/review/:reviewID", checkAuth(req, res), editReview(req, res));
app.delete("/user/:userID/review/:reviewID", checkAuth(req, res), deleteReview(req, res));


=======
app.post("/bathroom/:bathroomID/review", checkAuth(req, res, next), makeReview(req, res));

app.patch("/user/:userID/review/:reviewID", checkAuth(req, res, next), editReview(req, res));
app.delete("/user/:userID/review/:reviewID", checkAuth(req, res, next), deleteReview(req, res));
>>>>>>> 1638dafc850a295e7a8f7f2763d1f201a961451c

app.listen(port, host, () => {
    console.log(`server is listening at http://${addr}...`)
})