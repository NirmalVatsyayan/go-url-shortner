package handlers

import (
	db "github.com/NirmalVatsyayan/UrlShortnerRepo/Database"
	models "github.com/NirmalVatsyayan/UrlShortnerRepo/Models"
	"time"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"
	"gopkg.in/gin-gonic/gin.v1"
	"math/rand"
	"math"
	"log"
)


func stringInSlice(a string) bool {
	for _, b := range db.Data {
		if b == a {
			return true
		}
	}
	return false
}

func HelloHandler(c *gin.Context) {

	currentUser := c.MustGet("user").(models.User)

	currentTime := time.Now()
	currentTime.Format(time.RFC3339)
	c.JSON(200, gin.H{
		"current_time": currentTime,
		"text":"Hi " + currentUser.Username + ", You are login now.",
	})
}

func EncodedUrlRetriveHandler(c *gin.Context){
	encodedUrl := c.Param("encodedUrl")
	var urlObj models.UrlShortner

	count, err := db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{"urlencoded":encodedUrl}).Count()

	if err != nil  || count == 0{
		AbortWithError(c, http.StatusBadRequest, "Invalid URL")
		return
	}

	change := mgo.Change{
		Update: bson.M{"$inc": bson.M{"viewcount": 1}},
		ReturnNew: true,
	}

	_, err = db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{"urlencoded":encodedUrl}).Apply(change, &urlObj)
	if err != nil  {
		AbortWithError(c, http.StatusInternalServerError, "Some error occured, Kindly try later. Regrets for Inconvenience")
		return
	}

	c.JSON(200, gin.H{
		"original_url": urlObj.Url,
	})

}

func GetUrlInfoHandler(c *gin.Context){
	encodedUrl := c.Param("encodedUrl")
	currentUser := c.MustGet("user").(models.User)
	var urlObj models.UrlShortner

	err := db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{"urlencoded":encodedUrl,"userid":currentUser.ID}).One(&urlObj)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, "Given URL is not associated with you.")
		return
	}
	c.JSON(200, gin.H{
		"original_url": urlObj.Url,
		"view_count": urlObj.ViewCount,
		"encoded_url": urlObj.UrlEncoded,
		"created_on": urlObj.CreatedOn,
	})
}

func GetUserUrlsHandler(c *gin.Context){
	currentUser := c.MustGet("user").(models.User)

	pageNumber := c.DefaultQuery("page", "1")
	queryCount := c.DefaultQuery("count", "10")

	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, "Invalid page query param")
		return
	}

	count, err := strconv.Atoi(queryCount)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, "Invalid count query param")
		return
	}

	PrevUrl := ""
	NextUrl := ""

	user_urls_count, err := db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{"userid":currentUser.ID}).Count()
	if err != nil {
		AbortWithError(c, http.StatusInternalServerError, "Some Error occured, kindly try later")
		return
	}

	if user_urls_count == 0 {
		AbortWithError(c, http.StatusBadRequest, "No encoded urls associated with user")
		return
	}

	if user_urls_count > ((page)*count) {

		if page > 1{
			url_page := strconv.Itoa(page-1)
			PrevUrl = "http://localhost:8000/user/urls?count="+queryCount+"&page="+url_page

			url_page = strconv.Itoa(page+1)
			NextUrl = "http://localhost:8000/user/urls?count="+queryCount+"&page="+url_page
		}else{
			PrevUrl = ""
			url_page := strconv.Itoa(page+1)
			NextUrl = "http://localhost:8000/user/urls?count="+queryCount+"&page="+url_page
		}

	}else if user_urls_count < ((page)*count) {

		if page > 1{
			url_page := strconv.Itoa(page-1)
			PrevUrl = "http://localhost:8000/user/urls?count="+queryCount+"&page="+url_page

			url_page = strconv.Itoa(page+1)
			NextUrl = ""
		}else{
			PrevUrl = ""
			NextUrl = ""
		}
	}

	page_obj := models.Pagination{PrevUrl: PrevUrl, NextUrl:NextUrl, Count:user_urls_count}

	url_obj := []models.UrlShortner{}

	// Fetch user
	if err := db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{}).Sort("-createdon").Limit(count).Skip((page-1)*count).All(&url_obj); err != nil {
		AbortWithError(c, http.StatusInternalServerError, "Some Error occured, kindly try later")
		return
	}

	return_obj := models.UrlWrapper{Pagination:page_obj, Urls:url_obj}
	c.JSON(200, return_obj)
}


func indexOf(word string) (int) {
	for k, v := range db.Data {
		if word == v {
			return k
		}
	}
	return -1
}

func reverse(input string) string{

	ss := make([]string, 0)
	for _, data := range(input){
		ss = append(ss, string(data))
	}

	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
	return strings.Join(ss, "")
}



func UrlPostHandler(c *gin.Context) {

	currentUser := c.MustGet("user").(models.User)
	url := c.PostForm("url")
	encodedUrl := c.PostForm("urlEncoded")
	var url_value int64 = int64(rand.Intn(1000000))
	url_encoding_key := url_value
	found := false
	max_itr := 10
	itr_count := 0

	if url == "" {
		AbortWithError(c, http.StatusBadRequest, "Invaild URL")
		return
	}

	if encodedUrl != "" {
		for _, char := range encodedUrl {
			if !stringInSlice(string(char)){
				AbortWithError(c, http.StatusBadRequest, "Invaild Hash Url Provided")
				return
			}
		}

		encoding_url_count, _ := db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{"urlencoded":encodedUrl}).Count()
		if encoding_url_count > 0 {
			AbortWithError(c, http.StatusBadRequest, "Encoded URL already taken, Kindly request for a new One.")
			return
		}

		if len(encodedUrl) > 23 {
			AbortWithError(c, http.StatusBadRequest, "Encoded URL length too long.")
			return
		}

		var sum int64
		sum = 0
		index    := 0.0

		charLength := float64(len(db.Data))
		tokenLength := float64(len(encodedUrl))

		for _, c := range []byte(encodedUrl) {
			power := tokenLength - (index + 1)
			index := indexOf(string(c))
			sum += int64(index * int(math.Pow(charLength, power)))
			index++
		}

		log.Println(sum)

		url_obj := models.UrlShortner{
			UserID:       currentUser.ID,
			UrlEncodingKey: sum,
			Url: url,
			UrlEncoded: encodedUrl,
			CreatedOn: time.Now(),
			ViewCount: 0,
		}

		if err := db.MongoConn.DB(db.DbName).C("url_shortner").Insert(url_obj); err != nil {
			AbortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, url_obj)

	}else {

		for itr_count <= max_itr {
			encoding_key_count, _ := db.MongoConn.DB(db.DbName).C("url_shortner").Find(bson.M{"urlencodingkey":url_value}).Count()

			if encoding_key_count > 0 {
				url_value = url_value + db.MinVal
				itr_count = itr_count + 1
				continue
			}
			url_encoding_key = url_value
			found = true
			break
		}

		if found == false {
			AbortWithError(c, http.StatusBadRequest, "Some Error occured at server, Please try later. Regrets for inconvenience.")
			return
		}

		var hashSlice []int64
		var remainder int64 = 0
		index := 0
		hashString := ""

		for url_value > 0 {
			remainder = url_value % 62
			url_value = url_value / 62
			hashSlice = append(hashSlice, remainder)
		}

		hashSliceLength := len(hashSlice)

		if hashSliceLength > 23 {
			AbortWithError(c, http.StatusBadRequest, "Some Error occured at server, Please try later. Regrets for inconvenience.")
			return
		}

		for hashSliceLength > index {
			hashString = hashString + db.Data[hashSlice[index]]
			index++
		}


		hashString = reverse(hashString)

		url_obj := models.UrlShortner{
			UserID:       currentUser.ID,
			UrlEncodingKey: url_encoding_key,
			Url: url,
			UrlEncoded: hashString,
			CreatedOn: time.Now(),
			ViewCount: 0,
		}

		if err := db.MongoConn.DB(db.DbName).C("url_shortner").Insert(url_obj); err != nil {
			AbortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, url_obj)
	}
}
