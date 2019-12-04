async function getFavorites(req, res) {
    try {
        const db = req.db;
        let user = JSON.parse(req.get("X-User"))

        let bathroomsIDs = []
        let results = await db.query(`
            SELECT BathroomID
            FROM tblFavorites
            WHERE UserID = ${user.id}
        `)

        for (let result of results) {
            bathroomsIDs.push(result)
        }

        if (db) db.end();
        return res.status(200).json({"favorites": bathroomIDs});
    } catch {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    getFavorites
}