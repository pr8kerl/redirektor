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
	Key    string
	Label  string
	InUrl  url.URL
	OutUrl url.URL
}

type RedirektRequestResponse struct {
	Label  string `json:"label"`
	InUrl  string `json:"incoming"`
	OutUrl string `json:"outcoming"`
}
type DbsResponse struct {
	Dbs []DbResponse `json:"dbs"`
}
type DbResponse struct {
	Name string `json:"name"`
}

type Labeller struct {
	Keys map[string]string
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
	dbs           map[string]redis.Client
	DatabaseNames []string
	config        *Config
}

func NewRedirektor(c *Config) (*Redirektor, error) {

	r := Redirektor{config: c}
	i := len(c.Db)
	r.dbs = make(map[string]redis.Client, i)
	r.DatabaseNames = make([]string, i)
	for key, dbprofile := range c.Db {
		fmt.Println("creating db connection for profile:", key)
		client := redis.NewClient(&redis.Options{
			Addr:     dbprofile.DbAddr,
			Password: dbprofile.DbPassword,
			DB:       dbprofile.DbID,
			//              ReadOnly: true,
		})
		r.dbs[key] = *client
		i--
		r.DatabaseNames[i] = key
	}
	return &r, nil

}
func (r *Redirektor) NumDBs() int {
	sz := len(r.dbs)
	return sz
}
func (r *Redirektor) GetDBs(c *gin.Context) {
	c.JSON(http.StatusOK, r.DatabaseNames)
	return
}
func (r *Redirektor) GetAll(c *gin.Context) {

	var response []RedirektRequestResponse
	for i := range r.DatabaseNames {
		dbname := r.DatabaseNames[i]
		dbclient := r.dbs[dbname]
		sz, err := dbclient.DbSize().Result()
		if err != nil {
			fmt.Printf("db size err: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"db size error": err})
			return
		}
		if len(response) == 0 {
			response = make([]RedirektRequestResponse, sz, sz)
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
				s := strings.Split(k, ":")
				label, inurl := s[0], s[1]
				row := RedirektRequestResponse{
					Label:  label,
					InUrl:  inurl,
					OutUrl: outurl,
				}
				response = append(response, row)
				fmt.Printf("response: %s:%s %s\n", label, inurl, outurl)

			}
			if cursor == 0 {
				break
			}
		}
		c.JSON(http.StatusOK, gin.H{"response": response})
		return
	}
}

func (r *Redirektor) Get(c *gin.Context) {

}
func (r *Redirektor) Set(c *gin.Context) {

}
