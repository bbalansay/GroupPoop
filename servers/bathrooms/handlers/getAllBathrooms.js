async function getAllBathrooms(req, res) {
    // let db; 
    
    try {
        const db = req.db;

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

module.exports = {
    getAllBathrooms
}