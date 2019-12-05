async function editReview(req, res, {getDBConn}) {
    // let db;
    let reviewID = req.params.reviewID;
    let userID = req.params.userID;

    try {
        const db = getDBConn();
        let user = JSON.parse(req.get("X-User"))

        // check review exists
        let result = await db.query(`
            SELECT * FROM tblReview
            WHERE ID = ${reviewID}
        `)
        if (result.length != 1) {
            return res.status(403).json({"message": "Review does not exist"})
        }

        // check to see if user is author of review
        if (user.id != result[0].UserID) {
            return res.status(403).json({"message": "You are not the creator of this review!"})
        }
        console.log("AH" + req.body.Content)
        if (req.body.Content) {
            console.log("GRAH" + req.body.Content)
            await db.query(`
                UPDATE tblReview SET Score = '${req.body.Score}'
                WHERE ID = ${reviewID}
            `)
            await db.query(`
                UPDATE tblReview SET Content = '${req.body.Content}'
                WHERE ID = ${reviewID}
            `)
            await db.query(`
                UPDATE tblReview SET EditedAt = NOW()
                WHERE ID = ${reviewID}
            `)
        }

        let newResult = await db.query(`
            SELECT * FROM tblReview
            WHERE ID = ${reviewID}
        `)

        res.set("Content-Type", "application/json")
        return res.status(201).json(newResult[0])
    } catch(err) {
        return res.status(500).json({"error": err.message})
    }
}

module.exports = {
    editReview
}