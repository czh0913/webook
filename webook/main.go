package main

func main() {

	//server := initWebServer()
	//db := initDB()
	//rdb := initRedis()
	//u := initUser(db, rdb)
	//u.RegisterRoutes(server)
	server := InitWebServer()

	server.Run("0.0.0.0:8080")
}
