package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Redirekt struct {
	Prefix string
	Db     string
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
func (r *Redirektor) GetAll(c *gin.Context) {

	var response []RedirektRequestResponse
	for dbname, dbclient := range r.Dbs {

		sz, err := dbclient.DbSize().Result()
		if err != nil {
			fmt.Printf("db size err: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"db size error": err})
			return
		}
		if len(response) == 0 {
			response = make([]RedirektRequestResponse, 0, sz)
		} else {
			// make response fit the new records
			t := make([]RedirektRequestResponse, len(response), (int64(cap(response)) + sz))
			copy(t, response)
			response = t
		}

		var cursor uint64
		for {
			var rkeys []string
			var err error
			rkeys, cursor, err = dbclient.Scan(cursor, "", 100).Result()
			if err != nil {
				fmt.Printf("db scan err: %s\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"db scan error": err})
				return
			}
			for x := range rkeys {
				//HERE
				k := rkeys[x]
				outurl, err := dbclient.Get(k).Result()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"db getall error": err})
					return
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
					fmt.Printf("unexpected record length: %q\n", s)
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
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
	return
}

func (r *Redirektor) Get(c *gin.Context) {

}
func (r *Redirektor) Set(c *gin.Context) {

}
