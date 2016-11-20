package redirektor

import (
	"gopkg.in/redis.v5"
	"net/url"
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
		r.dbs[key] = client
		i--
		DatabaseNames[i] = key
	}
	return &r, nil

}
func (r *Redirektor) NumDBs() int {
	sz := len(r.dbs)
	return sz
}
func (r *Redirektor) GetDBs(c *gin.Context) {
	c.JSON(http.StatusOK, r.DatabaseNames)
}
func (r *Redirektor) GetAll(c *gin.Context) {
	var response []RedirektRequestResponse
	for i := range r.DatabaseNames {
		dbname := r.DatabaseNames[i]
		dbclient := r.dbs[dbname]
		sz, err := dbclient.DbSize().Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"db size error": err})
		}
		if len(response) == 0 {
			response = make([]RedirektRequestResponse, sz, sz)
		} else {
			// make response fit the new records
			response = copy([]RedirektRequestResponse, sz)
			t := make([]byte, len(response), (cap(response) + sz))
			copy(t, response)
			response = t
		}

		pipe := dbclient.Pipeline()
		defer pipe.Close()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = pipe.Scan(cursor, "", 10).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"db scan error": err})
			}
			n += len(keys)
			for x := range keys {
				//HERE

			}
			if cursor == 0 {
				break
			}
		}
	}
}

func (r *Redirektor) Get(c *gin.Context) {

}
func (r *Redirektor) Set(c *gin.Context) {

}
