async function getBathroom(req, res, {getDBConn}) {
    // let db;
    let bathroomID = req.params.bathroomID;

    try {
        const db = getDBConn();

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
       
        res.set("Content-Type", "application/json")
        return res.status(200).json(bathroom["reviews"] = reviews)
    } catch (err) {
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    getBathroom
}