require("../middleware/getDB")

export async function makeReview(req, res, db) {
    // let db;
    let bathroomID = req.params.bathroomID;

    try {
        // db = await getDB()
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
    } catch {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}