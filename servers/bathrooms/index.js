"use strict"

const express = require("express");
const mysql = require("promise-mysql");

const app = express();
app.use(express.json())

const addr = process.env.MESSAGESPORT || ":80";
const [host, port] = addr.split(":")

const dbHost = process.env.DBHOST
const dbPort = process.env.DBPORT
const dbUser = process.env.DBUSER
const dbPass = process.env.MYSQL_ROOT_PASSWORD
const dbName = process.env.DBNAME

//get all the bathrooms in the database
app.get("/bathroom", getAllBathrooms(req, res));

//gets a specific bathroom and all reviews
app.get("/bathroom/:bathroomID", getBathroom(req, res));

// create a review for a specific bathroom
app.post("/bathroom/:bathroomID/review", makeReview(req, res));

app.patch("/user/:userID/review/:reviewID", editReview(req, res));
app.delete("/user/:userID/review/:reviewID", deleteReview(req, res));

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

function checkAuth(req, res) {
    if (req.get("X-User") == undefined) 
      return res.status(401).json({ "message": "User must be authenticated" })
} 

async function getAllBathrooms(req, res) {
    checkAuth(req, res)
    let db; 
    
    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let bathrooms = await db.query(`
            SELECT * FROM tblBathroom
        `)

        if (db) db.end();
        return res.status(200).json(bathrooms)
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}

async function getBathroom(req, res) {
    checkAuth(req, res)
    let db;
    let bathroomID = req.params.bathroomID;

    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let bathroom = await db.query(`
            SELECT * FROM tblBathroom
            WHERE ID = ${bathroomID}
        `)
        if (bathroom.length != 1) {
            return res.status(403).send("Bathroom does not exist")
        }

        let reviews = await db.query(`
            SELECT * FROM tblReview
            WHERE BathroomID = ${bathroomID}
        `)
       

        if (db) db.end();
        return res.status(200).json(bathroom.concat(reviews))
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}

async function makeReview(req, res) {
    checkAuth(req, res)
    let db;
    let bathroomID = req.params.bathroomID;

    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let reviewJSON = req.body;
        let rows = await db.query(`
            INSERT INTO tblReview (UserID, BathroomID, Score, Content, CreatedAt, EditedAt)
            VALUES (${user.id}, ${bathroomID}, ${reviewJSON.Score} ${reviewJSON.Content}, NOW(), NOW())
        `)

        if (db) db.end();
        return res.status(201).json(reviewJSON);
    } catch {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}


async function editReview(req, res) {
    checkAuth(req, res)
    let db;
    let reviewID = req.params.reviewID;
    let userID = req.params.userID;

    try {
        db = await getDB();
        let user = JSON.parse(req.get("X-User"))

        // check review exists
        let result = await db.query(`
            SELECT * FROM tblReview
            WHERE ID = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Review does not exist")
        }

        // check to see if user is author of review
        if (userID != result[0].UserID) {
            return res.status(403).send("Your are not the creator of this review!")
        }

        if (req.body.content) {
            await db.query(`
                UPDATE tblReview SET Score = '${req.body.Score}'
                WHERE ID = ${reviewID}
            `)
            await db.query(`
                UPDATE Review SET Content = '${req.body.contnet}'
                WHERE ID = ${reviewID}
            `)
            await db.query(`
                UPDATE Review SET EditedAt = NOW()
                WHERE ID = ${reviewID}
            `)
        }
        if (db) db.end();
        return res.status(201).json(result[0])
    } catch {
        if (db) db.end();
        return res.status(500).json({"error": err.message})
    }
}

async function deleteReview(req, res) {
    checkAuth(req,res)
    let db;
    let reviewID = req.params.reviewID;
    let userID = req.params.userID;

    try {
        db = await getDB();
        let user = JSON.parse(req.get("X-User"))

         // check review exists
         let result = await db.query(`
            SELECT * FROM tblReview
            WHERE ID = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Review does not exist")
        }

         // check to see if user is author of review
         if (userID != result[0].UserID) {
            return res.status(403).send("Your are not the creator of this review!")
        }

        await db.query(`
            DELETE FROM tblReview
            WHERE ID = ${reviewID}
        `)

        if (db) db.end();
        return res.status(200).send("Successfully deleted!")
    } catch {
        if (db) db.end();
        return res.status(500).send("unexpected error: " + err.message)
    }
}

// /*
//   RESOURCE PATH: /bathroom/:bathroomID
//   SUPPORTED METHODS:
//   - GET
// */
// app.get("/bathroom/:bathroomID", async (req, res) => {
//     checkAuth(req, res)
//     let db;
//     let bathroomID = req.params.bathroomID;

//     try {
//         db = await getDB()
//         let user = JSON.parse(req.get("X-User"))

//         let result = await db.query(`
//             SELECT * FROM Bathroom
//             WHERE bathroom_id = ${bathroomID}
//         `)
//         if (result.length != 1) {
//             return res.status(403).send("Bathroom does not exist")
//         }

//         if (db) db.end();
//         return res.status(200).json(result)
//     } catch (err) {
//         if (db) db.end();
//         return res.status(500).json( {"error" : err.message })
//     }
// })