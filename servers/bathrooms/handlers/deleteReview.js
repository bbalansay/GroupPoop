async function deleteReview(req, res, {getDBConn}) {
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
            return res.status(403).json({"message": "Your are not the creator of this review!"})
        }

        await db.query(`
            DELETE FROM tblReview
            WHERE ID = ${reviewID}
        `)

        return res.status(200).json({"message": "Successfully deleted!"})
    } catch(err) {
        return res.status(500).json({"error: ": err.message})
    }
}

module.exports = {
    deleteReview
}