"use strict"

const express = require("express");
const mysql = require("promise-mysql");
const auth = require("./middleware/checkAuth")
const allBath = require("./handlers/getAllBathrooms")
const getBath = require("./handlers/getBathroom")
const makeRev = require("./handlers/makeReview")
const editRev = require("./handlers/editReview")
const delRev = require("./handlers/deleteReview")
const addFav = require("./handlers/addFavorite")
const getFav = require("./handlers/getFavorites")

const dbHost = process.env.DBHOST
const dbPort = process.env.DBPORT
const dbUser = process.env.DBUSER
const dbPass = process.env.MYSQL_ROOT_PASSWORD
const dbName = process.env.DBNAME

const app = express();
app.use(express.json())

const addr = process.env.BATHROOMSPORT || ":80";
const [host, port] = addr.split(":")

let db;

const getDBConn = () => {
    return db
}

const RequestWrapper = (handler, x) => {
	return (req, res) => {
		handler(req, res, x);
	}
}


//get all the bathrooms in the database
app.get("/bathroom", auth.checkAuth, RequestWrapper(allBath.getAllBathrooms, {getDBConn}));

//gets a specific bathroom and all reviews
app.get("/bathroom/:bathroomID", auth.checkAuth, RequestWrapper(getBath.getBathroom, { getDBConn}));

// create a review for a specific bathroom
app.post("/bathroom/:bathroomID/review", auth.checkAuth, RequestWrapper(makeRev.makeReview, { getDBConn}));

app.patch("/review/:reviewID", auth.checkAuth, RequestWrapper(editRev.editReview, { getDBConn}));
app.delete("/review/:reviewID", auth.checkAuth, RequestWrapper(delRev.deleteReview, { getDBConn}));

app.get("/favorites", auth.checkAuth, RequestWrapper(getFav.getFavorites, { getDBConn}));
app.post("/favorites/:bathroomID", auth.checkAuth, RequestWrapper(addFav.addFavorite, { getDBConn}));

async function main() {
    try {
        db = await mysql.createConnection({
            host: dbHost,
            port: dbPort,
            user: dbUser,
            password: dbPass,
            database: dbName
        });
    
        app.listen(port, host, () => {
            console.log(`server is listening at http://${addr}...`)
        })
    } catch (err) {
        console.log(err)
        process.exit(1)
    }   
}

main()