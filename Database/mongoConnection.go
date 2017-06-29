package database

import (
	"gopkg.in/mgo.v2"
	"time"
	config "github.com/NirmalVatsyayan/go-url-shortner/Config"
)

var (
	MongoConn *mgo.Session
	DbName string
	Data = [...] string {"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t",
		"u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N",
		"O","P","Q","R","S","T","U","V","W","X","Y","Z","1","2","3","4","5","6","7","8",
		"9","0"}
)

const (
	JWTSigningKey string        = "nirmalvatsyayan"
	ExpireTime    time.Duration = time.Minute * 60 * 24 * 30
	Realm         string        = "jwt auth"
	MinVal        int64  = 10000
	MaxVal        int64 = -9223372036854775808
)

func InitDB() {
	configs, _ := config.ReadConfig("Config.json")

	conn, err := mgo.Dial("mongodb://"+configs.DB_HOST)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	MongoConn = conn
	DbName = configs.DB_NAME
}
