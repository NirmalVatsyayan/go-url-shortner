package middlewares

import (
	"github.com/dgrijalva/jwt-go/request"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"log"
	"gopkg.in/mgo.v2/bson"
	db "github.com/NirmalVatsyayan/UrlShortnerRepo/Database"
	models "github.com/NirmalVatsyayan/UrlShortnerRepo/Models"
	"gopkg.in/gin-gonic/gin.v1"

)



func AbortWithError(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+db.Realm)
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
	c.Abort()
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		token, err := request.ParseFromRequest(c.Request,request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(db.JWTSigningKey))

			return b, nil
		})

		if err != nil {
			AbortWithError(c, http.StatusUnauthorized, "Invaild User Token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		log.Printf("Current user id: %s", claims["id"])

		err = db.MongoConn.DB(db.DbName).C("user_profile").Find(bson.M{"id":claims["id"]}).One(&user)

		if err != nil {
			AbortWithError(c, http.StatusInternalServerError, "DB Query Error")
			return
		}

		c.Set("user", user)
	}
}
