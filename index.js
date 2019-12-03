"use strict"

const express = require("express");
const mysql = require("promise-mysql");

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

app.get("/bathroom/:bathroomID", async (req, res) => {
    checkAuth(req, res)
    let db;
    let bathroomID = req.params.bathroomID;

    try {
        db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let result = await db.query(`
            Select * FROM Bathroom
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