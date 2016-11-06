# redirektor

## WIP

Apache RewriteMap management - work in progress

Looks like python is best to manage berkeleydb files. Tried golang, wanted golang but no bindings for older berkleydb versions - I tried em all. :)

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

## bdbmgr

A command line utility to manage berkeley db files for use by Apache RewriteMaps.

The following subcommands are available:

```
$ ./bdbmgr.py 
usage: bdbmgr.py <command> [<args>]

The following subcommands are available:
   csv2db     import from csv to bdb
   db2csv     export from bdb to csv
   redis2db   import from redis to bdb
   redis2csv  export from redis to csv
   csv2redis  import from csv to redis
   db2redis   import from bdb to redis
   poll       subscribe to redis and update bdb with changes
bdbmgr: error: the following arguments are required: subcommand
```

Use the **poll** command to subscribe to a redis DB and update the local RewriteMap db file with changes.
Use supervisord to start and monitor the polling process.






## Notes to self

* enable redis event notifications

```
$ grep AKE /etc/redis/redis.conf
#  A     Alias for g$lshzxe, so that the "AKE" string means all the events.
notify-keyspace-events "AKE"
```

* or enable redis key events via cli

```
redis-cli config set notify-keyspace-events KEA
```


* delete all keys in one hit
```
redis keys '*'|cargs redis-cli del 
```
