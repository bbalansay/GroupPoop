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
        console.log("AHHHH")
        for (let result of results) {
            if (result.BathroomID == bathroomID) {
                res.status(304).json({"error": "Favorite already added."})
            }
        }
        console.log("about to add favorite")
        let rows = await db.query(`
            INSERT INTO tblFavorites (UserID, BathroomID)
            VALUES (${user.id}, ${bathroomID});
        `)

        return res.status(201).json({"message": "Favorite added."});
    } catch(err) {
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    addFavorite
}