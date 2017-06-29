package handlers

import (

	"net/http"
	"golang.org/x/crypto/bcrypt"
	"time"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"github.com/satori/go.uuid"
	db "github.com/NirmalVatsyayan/UrlShortnerRepo/Database"
	input "github.com/NirmalVatsyayan/UrlShortnerRepo/Input"
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


func ProfileHandler(c *gin.Context){
	currentUser := c.MustGet("user").(models.User)
	c.JSON(http.StatusOK, gin.H{
		"user_id": currentUser.ID ,
		"name": currentUser.Name,
		"username": currentUser.Username,
	})

}

func RefreshHandler(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)

	expire := time.Now().Add(db.ExpireTime)

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = currentUser.ID
	claims["exp"] = expire.Unix()
	token.Claims = claims
	// Set some claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(db.JWTSigningKey))

	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func LoginHandler(c *gin.Context) {
	var form input.Login
	var user models.User

	if c.BindJSON(&form) != nil {
		AbortWithError(c, http.StatusBadRequest, "Missing usename or password")
		return
	}

	err := db.MongoConn.DB(db.DbName).C("user_profile").Find(bson.M{"username":form.Username}).One(&user)

	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "DB Query Error")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
		AbortWithError(c, http.StatusUnauthorized, "Incorrect Username / Password")
		return
	}

	expire := time.Now().Add(db.ExpireTime)

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = expire.Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string

	tokenString, err := token.SignedString([]byte(db.JWTSigningKey))
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

func RegisterHandler(c *gin.Context) {
	var form input.Register

	if c.BindJSON(&form) != nil {
		AbortWithError(c, http.StatusBadRequest, "Missing name or usename or password")
		return
	}

	total_count, _ := db.MongoConn.DB(db.DbName).C("user_profile").Find(bson.M{"username":form.Username}).Count()

	if total_count > 0 {
		AbortWithError(c, http.StatusBadRequest, "Username is already exist")
		return
	}

	userId := uuid.NewV4().String()

	if digest, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost); err != nil {
		AbortWithError(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		form.Password = string(digest)
	}

	user_obj := models.User{
		ID:       userId,
		Name: form.Name,
		Username: form.Username,
		Password: form.Password,
	}

	if err := db.MongoConn.DB(db.DbName).C("user_profile").Insert(user_obj); err != nil {
		AbortWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	expire := time.Now().Add(db.ExpireTime)

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["id"] = user_obj.ID
	claims["exp"] = expire.Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString([]byte(db.JWTSigningKey))

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"token": tokenString,
	})
}

