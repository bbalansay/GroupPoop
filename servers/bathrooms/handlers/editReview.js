require("../middleware/getDB")

async function editReview(req, res) {
    let db;
    let reviewID = req.params.reviewID;
    let userID = req.params.userID;

    try {
        db = await getDB();
        let user = JSON.parse(req.get("X-User"))

        // check review exists
        let result = await db.query(`
            SELECT * FROM tblReview
            WHERE ID = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).send("Review does not exist")
        }

        // check to see if user is author of review
        if (userID != result[0].UserID) {
            return res.status(403).send("Your are not the creator of this review!")
        }

        if (req.body.content) {
            await db.query(`
                UPDATE tblReview SET Score = '${req.body.Score}'
                WHERE ID = ${reviewID}
            `)
            await db.query(`
                UPDATE Review SET Content = '${req.body.contnet}'
                WHERE ID = ${reviewID}
            `)
            await db.query(`
                UPDATE Review SET EditedAt = NOW()
                WHERE ID = ${reviewID}
            `)
        }
        if (db) db.end();
        return res.status(201).json(result[0])
    } catch {
        if (db) db.end();
        return res.status(500).json({"error": err.message})
    }
}