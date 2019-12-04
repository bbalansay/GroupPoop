const mysql = require("promise-mysql");

const dbHost = process.env.DBHOST
const dbPort = process.env.DBPORT
const dbUser = process.env.DBUSER
const dbPass = process.env.MYSQL_ROOT_PASSWORD
const dbName = process.env.DBNAME

async function getDB(req, res, next) {
    try {
        let db = await mysql.createConnection({
            host: dbHost,
            port: dbPort,
            user: dbUser,
            password: dbPass,
            database: dbName
        });

        if (db) {
            req.db = db;
            next()
        }
    } catch (err) {
        if (db) db.end();
        return res.status(500).send("unexpected error: " + err.message)
    }
}

module.exports = {
    getDB
}