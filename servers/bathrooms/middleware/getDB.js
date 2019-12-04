async function getDB() {
    let db = await mysql.createConnection({
        host: dbHost,
        port: dbPort,
        user: dbUser,
        password: dbPass,
        database: dbName
    });
    return db;
}