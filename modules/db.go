package modules

// 这个包是失败的，启用将会导致程序无法启动（假死、闪退）

// "log"
// "os"
// "gorm.io/driver/sqlite"
// "gorm.io/gorm"

func DBsave() {

	// f, err := os.OpenFile("db.txt", os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// // 查找文件末尾的偏移量
	// n, _ := f.Seek(0, os.SEEK_END)
	// // 从末尾的偏移量开始写入内容
	// _, err = f.WriteAt([]byte(","+data), n)
	// if err != nil {
	// 	log.Print(err)
	// }
	// defer f.Close()

	// type Msg struct {
	// 	gorm.Model
	// 	UserName string
	// 	Message  string
	// 	Time     string
	// }
	// db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	// if err != nil {
	// 	log.Print("连接数据库失败")
	// }

	// // Migrate the schema
	// db.AutoMigrate(&Msg{})

	// // 插入内容
	// db.Create(&Msg{
	// 	UserName: "0",
	// 	Message:  "0",
	// 	Time:     "0",
	// })

	// 读取内容
	// var Database Database
	// db.First(&Database, 1)                 // find product with integer primary key
	// db.First(&Database, "code = ?", "D42") // find product with code D42

	// 删除操作：
	// db.Delete(&Database, 1)
}
