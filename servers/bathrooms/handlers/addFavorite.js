async function addFavorite(req, res, {getDBConn}) {
    // let db;
    let bathroomID = req.params.bathroomID;

    try {
        const db = getDBConn();
        let user = JSON.parse(req.get("X-User"))

        let results = await db.query(`
            SELECT BathroomID
            FROM tblFavorites
            WHERE UserID = ${user.id}
        `)

        for (let result of results) {
            if (result.BathroomID == bathroomID) {
                res.status(304).json({"error": "Favorite already added."})
            }
        }

        let rows = await db.query(`
            INSERT INTO tblFavorites (UserID, BathroomID)
            VALUES (${user.id}, ${bathroomID})
        `)

        if (db) db.end();
        return res.status(201).json({"message": "Favorite added."});
    } catch(err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    addFavorite
}