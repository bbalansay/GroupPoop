async function makeReview(req, res, {getDBConn}) {
    // let db;
    let bathroomID = req.params.bathroomID;

    try {
        const db = getDBConn();
        let user = JSON.parse(req.get("X-User"))

        let reviewJSON = req.body;
        if (reviewJSON.Content.length > 512) {
            return res.status(415).send("Review is too long")
        }
        let rows = await db.query(`
            INSERT INTO tblReview (UserID, BathroomID, Score, Content, CreatedAt, EditedAt)
            VALUES (${user.id}, ${bathroomID}, ${reviewJSON.Score}, ${reviewJSON.Content}, SELECT NOW(), SELECT NOW())
        `)

        res.set("Content-Type", "application/json")
        return res.status(201).json(reviewJSON);
    } catch (err) {
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    makeReview
}