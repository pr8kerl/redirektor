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
	Prefix string
	InUrl  url.URL
	OutUrl url.URL
}

type RedirektRequestResponse struct {
	Db     string `json:"db"`
	Prefix string `json:"prefix"`
	InUrl  string `json:"incoming"`
	OutUrl string `json:"outcoming"`
}
type DbsResponse struct {
	Dbs []DbResponse `json:"dbs"`
}
type DbResponse struct {
	Name string `json:"name"`
}

type RedirektService interface {
	GetDBs(*gin.Context)
	Get(*gin.Context)
	Set(*gin.Context)
}

// Client represents a client to the RedirektService
type Redirektor struct {

	// Returns the current time.
	Now func() time.Time

	// db options
	Dbs    map[string]redis.Client
	config *Config
}

func NewRedirektor(c *Config) (*Redirektor, error) {

	r := Redirektor{config: c}
	i := len(c.Db)
	r.Dbs = make(map[string]redis.Client, i)
	for key, dbprofile := range c.Db {
		fmt.Println("creating db connection for profile:", key)
		client := redis.NewClient(&redis.Options{
			Addr:     dbprofile.DbAddr,
			Password: dbprofile.DbPassword,
			DB:       dbprofile.DbID,
			//              ReadOnly: true,
		})
		r.Dbs[key] = *client
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

func (r *Redirektor) getDB(dbname string) ([]RedirektRequestResponse, error) {

	dbclient, ok := r.Dbs[dbname]
	if !ok {
		return nil, fmt.Errorf("error: unknown db name: %s\n", dbname)
	}

	sz, err := dbclient.DbSize().Result()
	if err != nil {
		return nil, fmt.Errorf("error: could not get db size: %s\n", err)
	}
	response := make([]RedirektRequestResponse, 0, sz)

	var cursor uint64
	for {
		var rkeys []string
		var err error
		rkeys, cursor, err = dbclient.Scan(cursor, "", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("error: could not scan db: %s\n", err)
		}
		for x := range rkeys {
			//HERE
			k := rkeys[x]
			outurl, err := dbclient.Get(k).Result()
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
			row := RedirektRequestResponse{
				Db:     dbname,
				Prefix: prefix,
				InUrl:  inurl,
				OutUrl: outurl,
			}
			response = append(response, row)
			//				fmt.Printf("response: %s:%s %s\n", dbname, prefix, inurl, outurl)

		}
		if cursor == 0 {
			break
		}
	}

	return response, nil
}

func (r *Redirektor) GetAll(c *gin.Context) {

	var response []RedirektRequestResponse
	for dbname, _ := range r.Dbs {

		dbresponse, err := r.getDB(dbname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error getting db": err})
			return
		}

		sz := len(dbresponse)
		if len(response) > 0 {
			// make response fit the new records
			t := make([]RedirektRequestResponse, len(response), (cap(response) + sz))
			copy(t, response)
			response = t
		} else {

			response = make([]RedirektRequestResponse, 0, sz)
		}

	}
	c.JSON(http.StatusOK, gin.H{"response": response})
	return
}

func (r *Redirektor) Get(c *gin.Context) {
	dbname := c.Param("dbname")

	dbresponse, err := r.getDB(dbname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getting db": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dbresponse})
	return

}
func (r *Redirektor) Set(c *gin.Context) {
	dbname := c.Param("dbname")

	dbclient, ok := r.Dbs[dbname]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error: unknown db name": dbname})
		return
	}

	var json Redirekt
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid redirekt format"})
		return
	}

	inurl := json.InUrl.String()
	outurl := json.OutUrl.String()

	valid := govalidator.IsURL(inurl)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error: invalid incoming url": inurl})
		return
	}
	valid = govalidator.IsURL(outurl)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error: invalid outgoing url": outurl})
		return
	}
	valid = govalidator.IsPrintableASCII(json.Prefix)
	if !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "invalid prefix string"})
		return
	}

	key := json.Prefix + ":" + inurl

	err = dbclient.Set(key, outurl, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error: setting key": key})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "AOK"})
	return

}
