package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Redirekt struct {
	Prefix string `json:"prefix"`
	InUrl  string `json:"incoming"`
	OutUrl string `json:"outgoing"`
}

type RedirektsResponse struct {
	Prefixes  []string   `json:"prefixes"`
	Redirekts []Redirekt `json:"redirekts"`
}

type DbsResponse struct {
	Dbs []DbResponse `json:"dbs"`
}
type DbResponse struct {
	Name string `json:"name"`
}
type Db struct {
	Client   *redis.Client
	Prefixes []string
}

type RedirektService interface {
	GetDBs(*gin.Context)
	Get(*gin.Context)
	Set(*gin.Context)
	Delete(*gin.Context)
}

// Client represents a client to the RedirektService
type Redirektor struct {

	// Returns the current time.
	Now func() time.Time

	config *Config
	// map of db objects
	Dbs       map[string]Db
	Prefix2Db map[string]Db

	// all available prefixes
	Prefixes []string
}

func NewRedirektor(c *Config) (*Redirektor, error) {

	r := Redirektor{config: c}
	i := len(c.Db)
	psz := 0 // store max size of slice for all db prefixes
	r.Dbs = make(map[string]Db, i)
	for key, dbprofile := range c.Db {
		fmt.Println("creating db connection for profile:", key)
		client := redis.NewClient(&redis.Options{
			Addr:     dbprofile.DbAddr,
			Password: dbprofile.DbPassword,
			DB:       dbprofile.DbID,
			//              ReadOnly: true,
		})
		db := Db{
			Client:   client,
			Prefixes: dbprofile.Prefix,
		}
		r.Dbs[key] = db
		psz += len(dbprofile.Prefix)
	}
	r.Prefixes = make([]string, 0, psz)
	r.Prefix2Db = make(map[string]Db, psz)
	for key, dbprofile := range c.Db {
		r.Prefixes = append(r.Prefixes, dbprofile.Prefix...)
		for _, pfx := range dbprofile.Prefix {
			r.Prefix2Db[pfx] = r.Dbs[key]
		}
	}
	return &r, nil

}
func (r *Redirektor) NumDBs() int {
	sz := len(r.Dbs)
	return sz
}

func (r *Redirektor) GetDBs(c *gin.Context) {
	sz := len(r.Dbs)
	resp := make([]string, sz)
	for key, _ := range r.Dbs {
		sz--
		resp[sz] = key
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (r *Redirektor) getDB(dbname string) ([]Redirekt, error) {

	db, ok := r.Dbs[dbname]
	if !ok {
		return nil, fmt.Errorf("error: unknown db name: %s\n", dbname)
	}

	sz, err := db.Client.DbSize().Result()
	if err != nil {
		return nil, fmt.Errorf("error: could not get db size: %s\n", err)
	}
	redirekts := make([]Redirekt, 0, sz)

	var cursor uint64
	for {
		var rkeys []string
		var err error
		rkeys, cursor, err = db.Client.Scan(cursor, "", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("error: could not scan db: %s\n", err)
		}
		for x := range rkeys {
			//HERE
			k := rkeys[x]
			outurl, err := db.Client.Get(k).Result()
			if err != nil {
				return nil, fmt.Errorf("error: could not get from db: %s\n", err)
			}
			s := strings.SplitN(k, ":", 2)
			var prefix, inurl string
			switch n := len(s); n {
			case 2:
				prefix, inurl = s[0], s[1]
			case 1:
				prefix = ""
				inurl = s[0]
			default:
				fmt.Printf("getDB drop invalid record: %q\n", s)
				continue
			}
			row := Redirekt{
				Prefix: prefix,
				InUrl:  inurl,
				OutUrl: outurl,
			}
			redirekts = append(redirekts, row)
			//				fmt.Printf("redirekts: %s:%s %s\n", dbname, prefix, inurl, outurl)

		}
		if cursor == 0 {
			break
		}
	}

	return redirekts, nil
}

func (r *Redirektor) GetAll(c *gin.Context) {

	var redirekts []Redirekt

	for dbname, _ := range r.Dbs {

		dbredirekts, err := r.getDB(dbname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sz := len(dbredirekts)
		if len(redirekts) > 0 {
			// make redirekts fit the new records
			t := make([]Redirekt, len(redirekts), (cap(redirekts) + sz))
			copy(t, redirekts)
			redirekts = append(t, dbredirekts...)
		} else {
			redirekts = dbredirekts
		}

	}
	response := RedirektsResponse{
		r.Prefixes,
		redirekts,
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
	return
}

func (r *Redirektor) Get(c *gin.Context) {
	dbname := c.Param("dbname")

	redirekts, err := r.getDB(dbname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirekts": redirekts})
	return
}

func (r *Redirektor) Set(c *gin.Context) {

	var json Redirekt
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid redirekt format"})
		return
	}

	inurl := json.InUrl
	outurl := json.OutUrl

	_, err = url.Parse(inurl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	_, err = url.Parse(outurl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	valid := govalidator.IsPrintableASCII(json.Prefix)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid prefix string"})
		return
	}

	fmt.Printf("set: %s:%s\n", json.Prefix, json.InUrl)

	db, validpfx := r.Prefix2Db[json.Prefix]
	if !validpfx {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid prefix string"})
		return
	}
	key := json.Prefix + ":" + inurl

	err = db.Client.Set(key, outurl, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not set redis key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "AOK"})
	return

}

func (r *Redirektor) Delete(c *gin.Context) {

	var json Redirekt
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid redirekt format"})
		return
	}

	inurl := json.InUrl

	_, err = url.Parse(inurl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	valid := govalidator.IsPrintableASCII(json.Prefix)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid prefix string"})
		return
	}

	fmt.Printf("delete: %s:%s\n", json.Prefix, json.InUrl)

	db, validpfx := r.Prefix2Db[json.Prefix]
	if !validpfx {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid prefix string"})
		return
	}
	key := json.Prefix + ":" + inurl

	err = db.Client.Del(key).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete redis key"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "AOK"})
	return

}
