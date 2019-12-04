async function makeReview(req, res) {
    // let db;
    let bathroomID = req.params.bathroomID;

    try {
        const db = req.db;
        let user = JSON.parse(req.get("X-User"))

        let reviewJSON = req.body;
        if (length(reviewJSON.Content) > 512) {
            return res.status(415).send("Review is too long")
        }
        let rows = await db.query(`
            INSERT INTO tblReview (UserID, BathroomID, Score, Content, CreatedAt, EditedAt)
            VALUES (${user.id}, ${bathroomID}, ${reviewJSON.Score} ${reviewJSON.Content}, NOW(), NOW())
        `)

        if (db) db.end();
        return res.status(201).json(reviewJSON);
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    makeReview
}