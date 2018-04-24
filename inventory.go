package main

func main() {
	db := GetDB()
	defer db.Close()

	SetupDB(db)
	AddLocation(db, Location{"closet", ""})
	AddLocation(db, Location{"shelf", "closet"})
	AddItem(db, Item{"bag", "shelf"})
}
