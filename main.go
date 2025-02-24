package main

func main() {
	app := App{}
	app.Initalise(DbUser, DbPassword, DbName)
	app.Run("localhost:10000")

}
