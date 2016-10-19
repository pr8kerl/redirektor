# redirektor

## WIP

Apache RewriteMap management - work in progress

Looks like python is best to manage berkeleydb files.

* download and install [BerkeleyDB 4.7](http://download.oracle.com/berkeley-db/db-4.7.25.tar.gz) (same version as used by Apache HTTPD in Amazon Linux)
* build it

```
tar -xvzf db-4.7.25.tar.gz
cd db-4.7.25/build_unix
../dist/configure --enable-static --enable-shared --enable-compat185 
make
sudo make install
export BERKELEYDB_DIR=/usr/local/BerkeleyDB.4.7
pip3 install bsddb3
```


And another apache rewritemap program using Redis.

Probably better to use Redis...
* use slave redis instance locally on each web server
* use central management redis server
* set max memory and LRU on local redis slave
* can expire keys

And a v basic csv importer too...
incoming_url,outgoing_url

```
$ csvimporter --help
usage: csvimporter [--version] [--help] <command> [<args>]

Available commands are:
    bolt     import csv to a BoltDB file
    redis    import csv to a RedisDB

$ csvimporter redis --help
csvimporter redis: import csv data to a Redis DB

		--database <string>	the redis DB connection string - default is localhost:6379
		--prefix <value>	the key prefix to use for import data
		--csv <filename>	the csv filename to import data from
		--db <int>	the redis DB number to use - defaults to 0


$ csvimporter bolt --help
csvimporter bolt: import csv data to a BoltDB file

		--database <filename>	the BoltDB filename to use
		--bucket <value>	the BoltDB bucket to import data into
		--csv <filename>	the csv filename to import data from
```

* delete all keys in one hit
```
redis keys '*'|cargs redis-cli del 
```
