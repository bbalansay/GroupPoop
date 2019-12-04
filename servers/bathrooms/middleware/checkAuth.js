function checkAuth(req, res) {
    if (req.get("X-User") == undefined) 
      return res.status(401).json({ "message": "User must be authenticated" })
} 