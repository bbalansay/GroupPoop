export async function getBathroom(req, res, db) {
    // let db;
    let bathroomID = req.params.bathroomID;

    try {
        // db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let bathroom = await db.query(`
            SELECT * FROM tblBathroom
            WHERE ID = ${bathroomID}
        `)
        if (bathroom.length != 1) {
            return res.status(403).send("Bathroom does not exist")
        }

        let reviews = await db.query(`
            SELECT * FROM tblReview
            WHERE BathroomID = ${bathroomID}
        `)
       

        if (db) db.end();
        return res.status(200).json(bathroom.concat(reviews))
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}