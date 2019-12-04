export async function getAllBathrooms(req, res, db) {
    // let db; 
    
    try {
        // db = await getDB()
        let user = JSON.parse(req.get("X-User"))

        let bathrooms = await db.query(`
            SELECT * FROM tblBathroom
        `)

        if (db) db.end();
        return res.status(200).json(bathrooms)
    } catch (err) {
        if (db) db.end();
        return res.status(500).json( {"error" : err.message })
    }
}