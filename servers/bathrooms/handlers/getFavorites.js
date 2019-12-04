async function getFavorites(req, res, {getDBConn}) {
    try {
        const db = getDBConn();
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

        res.set("Content-Type", "application/json")
        return res.status(200).json({"favorites": bathroomIDs});
    } catch (err) {
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    getFavorites
}