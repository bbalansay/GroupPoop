async function getAllBathrooms(req, res, {getDBConn}) {
    // let db; 
    
    try {
        const db = getDBConn();

        let bathrooms = await db.query(`
            SELECT * FROM tblBathroom
        `)

        res.set("Content-Type", "application/json")
        return res.status(200).json(bathrooms)
    } catch (err) {
        return res.status(500).json( {"error" : err.message })
    }
}

module.exports = {
    getAllBathrooms
}