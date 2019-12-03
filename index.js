"use strict"

// const express = require("express");
// const mysql = require("promise-mysql");

// const app = express();
// app.use(express.json())

// const addr = process.env.MESSAGESPORT || ":80";
// const [host, port] = addr.split(":")

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

function checkAuth(req, res) {
    if (req.get("X-User") == undefined) 
      return res.status(401).json({ "message": "User must be authenticated" })
} 

/*
    RESOURCE PATH: /review
    SUPPORTED METHODS:
    - POST
*/
app.post("/review/:reviewID", async (req, res) => {
    checkAuth(req, res)
    let db;
    let reviewID = req.params.reviewID;

    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let reviewJSON = req.body;
        let rows = await db.query(`
            INSERT INTO Review (user_id, bathroom_id, content, time)
            VALUES (${user.id}, ${reviewJSON.bathroom_id}, ${reviewJSON.content}, NOW())
        `)

        if (db) db.end();
        return res.status(200).json(channels);
    } catch {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
})

/*
    RESOURCE PATH: /review:reviewID
    SUPPORTED METHODS:
    - GET
    - PATCH
    - DELETE
*/
app.get("/review/:reviewID", async (req, res) => {
    checkAuth(req, res)
    let db;
    let reviewID = req.params.reviewID;

    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let result = await db.query(`
            SELECT * FROM Review
            WHERE review_id = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Review does not exist")
        }

        if (db) db.end();
        return res.status(200).json(result)
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
})

app.patch("/review/:reviewID", async (req, res) => {
    checkAuth(req, res)
    let db;
    let reviewID = req.params.reviewID;

    try {
        db = await getDB();
        let user = JSON.parse(req.get("X-User"))

        // check review exists
        let result = await db.query(`
            SELECT * FROM Review
            WHERE review_id = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Review does not exist")
        }

        // check to see if user is author of review
        if (user.id != result[0].user_id) {
            return res.status(40).send("Your are not the creator of this review!")
        }

        if (req.body.content) {
            await db.query(`
                UPDATE Review SET content = '${req.body.contnet}'
                WHERE review_id = ${reviewID}
            `)
            await db.query(`
                UPDATE Review SET time = NOW()
                WHERE id = ${reviewID}
            `)
        }
    } catch {
        if (db) db.end();
        return res.status(500).json({"error": err.message})
    }
})

app.delete("/review/:reviewID", async (req,res) => {
    checkAuth(req,res)
    let db;
    let reviewID = req.params.reviewID;

    try {
        db = await getDB();
        let user = JSON.parse(req.get("X-User"))

         // check review exists
         let result = await db.query(`
            SELECT * FROM Review
            WHERE review_id = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Review does not exist")
        }

         // check to see if user is author of review
         if (user.id != result[0].user_id) {
            return res.status(40).send("Your are not the creator of this review!")
        }

        await db.query(`
            DELETE FROM Review
            WHERE review_id = ${reviewID}
        `)

        if (db) db.end();
        return res.status(200).send("Successfully deleted!")
    } catch {
        if (db) db.end();
        return res.status(500).send("unexpected error: " + err.message)
    }
})

/*
  RESOURCE PATH: /bathroom/:bathroomID
  SUPPORTED METHODS:
  - GET
*/
app.get("/bathroom/:bathroomID", async (req, res) => {
    checkAuth(req, res)
    let db;
    let bathroomID = req.params.bathroomID;

    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let result = await db.query(`
            SELECT * FROM Bathroom
            WHERE bathroom_id = ${bathroomID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Bathroom does not exist")
        }

        if (db) db.end();
        return res.status(200).json(result)
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
})