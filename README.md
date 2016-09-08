# redirektor

## WIP

An apache rewritemap program using BoltDB.

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
