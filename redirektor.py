#!/usr/bin/python3.4

from bsddb3 import db
from urllib.parse import urlparse
from time import sleep

import argparse, sys, csv, redis, datetime, os


class BDBMgr(object):
    version = '0.1.0'
    rdb = None
    rredis = None
    redis_url = 'redis://localhost:6379/0'
    args = None

    def __init__(self):
        parser = argparse.ArgumentParser(
            prog='redirektor',
            description='manage berkeley db files',
            usage='''redirektor.py <command> [<args>]

The following subcommands are available:
   csv2db     import from csv to bdb
   db2csv     export from bdb to csv
   redis2db   import from redis to bdb
   redis2csv  export from redis to csv
   csv2redis  import from csv to redis
   db2redis   import from bdb to redis
   poll       subscribe to redis and update bdb with changes
''')
        parser.add_argument('subcommand', help='subcommand to run')
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        # parse_args defaults to [1:] for args, but you need to
        # exclude the rest of the args too, or validation will fail
        args = parser.parse_args(sys.argv[1:2])
        if not hasattr(self, args.subcommand):
            print('unrecognized subcommand')
            parser.print_help()
            exit(1)
        if 'REDIREKTOR_BDB' in os.environ:
            self.rdbfile = os.getenv['REDIREKTOR_BDB']
        if 'REDIREKTOR_REDIS' in os.environ:
            self.redis_url = os.getenv['REDIREKTOR_REDIS']
        # use dispatch pattern to invoke method with same name
        getattr(self, args.subcommand)()

    def __del__(self):
        if self.rdb:
            self.rdb.sync()
            self.rdb.close()

    def csv2db(self):
        parser = argparse.ArgumentParser(
            description='import from csv file to bdb')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--csv', required=True )
        parser.add_argument('--db', required=False )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        self.rdb = db.DB()
        # args.db overrides REDIREKTOR_BDB
        if args.db:
            self.rdbfile = args.db
        if not self.rdbfile:
            print('error: no db file provided; use --db or env var REDIREKTOR_BDB')
            sys.exit(99)
        print('redirektor csv2db, csv={}, db={}'.format( args.csv, self.rdbfile ))
        self.__csv2db(args.csv)

    def db2csv(self):
        parser = argparse.ArgumentParser(
            description='export a bdb to csv')
        # NOT prefixing the argument with -- means it's not optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--csv', required=True )
        parser.add_argument('--db', required=False )
        args = parser.parse_args(sys.argv[2:])
        self.rdb = db.DB()
        # args.db overrides REDIREKTOR_BDB
        if args.db:
            self.rdbfile = args.db
        if not self.rdbfile:
            print('error: no db file provided; use --db or env var REDIREKTOR_BDB')
            sys.exit(99)
        print('redirektor db2csv, csv={}, db={}'.format( args.csv, self.rdbfile ))
        self.__db2csv(args.csv)

    def redis2db(self):
        parser = argparse.ArgumentParser(
            description='import from csv file to bdb')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--db', required=None )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        self.rdb = db.DB()
        # args.db overrides REDIREKTOR_BDB
        if args.db:
            self.rdbfile = args.db
        if not self.rdbfile:
            print('error: no db file provided; use --db or env var REDIREKTOR_BDB')
            sys.exit(99)
        if args.redis:
            self.redis_url = args.redis
        self.rredis = redis.from_url(self.redis_url)
        print('redirektor redis2db, redis={} db={}'.format( self.redis_url, self.rdbfile ))
        self.__redis2db()

    def redis2csv(self):
        parser = argparse.ArgumentParser(
            description='export from redis to csv')
        # NOT prefixing the argument with -- means it's not optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--csv', required=True )
        args = parser.parse_args(sys.argv[2:])
        if args.redis:
            self.redis_url = args.redis
        self.rredis = redis.from_url(self.redis_url)
        print('redirektor redis2csv, redis={} csv={}'.format( self.redis_url, args.csv ))
        self.__redis2csv(args.csv)

    def csv2redis(self):
        parser = argparse.ArgumentParser(
            description='import from csv file to redis')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--csv', required=True )
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        if args.redis:
            self.redis_url = args.redis
        self.rredis = redis.from_url(self.redis_url)
        print('redirektor csv2redis, csv={} redis={}'.format( args.csv, self.redis_url ))
        self.__csv2redis(args.csv)

    def db2redis(self):
        parser = argparse.ArgumentParser(
            description='import from bdb file to redis')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--db', required=None )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        self.rdb = db.DB()
        # args.db overrides REDIREKTOR_BDB
        if args.db:
            self.rdbfile = args.db
        if not self.rdbfile:
            print('error: no db file provided; use --db or env var REDIREKTOR_BDB')
            sys.exit(99)
        if args.redis:
            self.redis_url = args.redis
        self.rredis = redis.from_url(self.redis_url)
        print('redirektor db2redis, db={} redis={}'.format( self.rdbfile, self.redis_url ))
        self.__db2redis()

    def poll(self):
        parser = argparse.ArgumentParser(
            description='subscribe to redis events and update bdb')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--db', required=None )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])

        # default to localhost
        self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        self.rdb = db.DB()
        # args.redis overrides REDIREKTOR_REDIS
        if args.redis:
            self.redis_url = args.redis
        self.rredis = redis.from_url(self.redis_url)

        # args.db overrides REDIREKTOR_BDB
        if args.db:
            self.rdbfile = args.db
        if not self.rdbfile:
            print('error: no db file provided; use --db option or env var REDIREKTOR_BDB')
            sys.exit(99)

        try:
        	self.rdb.open(self.rdbfile,None,db.DB_HASH)
        except Exception as e:
        	print('error: unable to open db: {}'.format( self.rdbfile, e ))
        	sys.exit(99)
        print('redirektor poll, db={} redis={}'.format( self.rdbfile, self.redis_url ))
        self.__poll()

    def __csv2db(self, csvfile):

        try:
        	self.rdb.open(self.rdbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('error: unable to open db: {}'.format( self.rdbfile ))
        	print(e)
        	sys.exit(99)

        cursor = self.rdb.cursor()
        rec = cursor.first()
        count = 0

        with open(csvfile, encoding="latin-1") as f:
            reader = csv.reader(f)
            try:
                for row in reader:
                    self.rdb.put(str.encode(row[0]), str.encode(row[1]))
                    #print('csv: {}'.format(row))
                    count += 1
            except Exception as e:
                sys.exit('csv error: file {}, line {}: {}'.format(csvfile, reader.line_num, e))
        print('read {} records from csv to bdb: {}'.format( count, csvfile, self.rdbfile ))


    def __db2csv(self, csvfile):

        try:
        	self.rdb.open(self.rdbfile,None,db.DB_HASH,db.DB_DIRTY_READ)
        except Exception as e:
        	print('error: unable to open db stack: {}'.format( self.rdbfile ))
        	print(e)
        	sys.exit(99)

        count = 0
        with open(csvfile, 'w') as f:
            cursor = self.rdb.cursor()
            rec = cursor.first()

            while rec:
                key = rec[0]
                value = rec[1]
                print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')), file=f )
                rec = cursor.next()
                count += 1
        print('read {} records from bdb to csv: {}'.format( count, self.rdbfile, csvfile ))

    def __redis2db(self):

        try:
        	self.rredis.ping()
        except Exception as e:
            print('{:%Y-%b-%d %H:%M:%S} redis2db error: unable to ping redis: {}'.format( datetime.datetime.now(), e ))
            sys.exit(99)

        try:
        	self.rdb.open(self.rdbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('error: unable to open db: {}'.format( self.rdbfile ))
        	print(e)
        	sys.exit(99)

        cursor = self.rdb.cursor()
        rec = cursor.first()
        count = 0

	# pull all redis entries
        try:
            for key in self.rredis.scan_iter():
                value = self.rredis.get(key)
                #print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')) )
                self.rdb.put(key, value)
                count += 1
        except Exception as e:
        	print('error scanning redis into db: {}'.format( self.rdbfile ))
        	print(e)
        	sys.exit(99)
        print('read {} records from redis into db: {}'.format( count, self.rdbfile ))

    def __redis2csv(self, csvfile):

        try:
        	self.rredis.ping()
        except Exception as e:
            print('{:%Y-%b-%d %H:%M:%S} redis2csv error: unable to ping redis: {}'.format( datetime.datetime.now(), e ))
            sys.exit(99)

        count = 0
	# pull all redis entries
        try:
            with open(csvfile, 'w') as f:
                for key in self.rredis.scan_iter():
                    value = self.rredis.get(key)
                    #print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')) )
                    print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')), file=f )
                    count += 1
        except Exception as e:
        	print('error scanning redis into csv: {}'.format( csvfile ))
        	print(e)
        	sys.exit(99)

        print('read {} records from redis to csv: {}'.format( count, csvfile ))

    def __csv2redis(self, csvfile):

        try:
        	self.rredis.ping()
        except Exception as e:
            print('{:%Y-%b-%d %H:%M:%S} csv2redis error: unable to ping redis: {}'.format( datetime.datetime.now(), e ))
            sys.exit(99)

        count = 0
        with open(csvfile, encoding="latin-1") as f:
            reader = csv.reader(f)
            try:
                with self.rredis.pipeline() as pipe:
                    for row in reader:
                        pipe.set(str.encode(row[0]), str.encode(row[1]))
                        count +=1
                    pipe.execute()
            except Exception as e:
                sys.exit('csv2redis error: file {}, line {}: {}'.format(csvfile, reader.line_num, e))
        print('read {} records from csv {} to redis'.format( count, csvfile ))

    def __db2redis(self):

        try:
        	self.rredis.ping()
        except Exception as e:
            print('{:%Y-%b-%d %H:%M:%S} db2redis error: unable to ping redis: {}'.format( datetime.datetime.now(), e ))
            sys.exit(99)

        try:
        	self.rdb.open(self.rdbfile,None,db.DB_HASH,db.DB_DIRTY_READ)
        except Exception as e:
        	print('error: unable to open db stack: {}'.format( self.rdbfile ))
        	print(e)
        	sys.exit(99)

        count = 0
        cursor = self.rdb.cursor()
        rec = cursor.first()

        try:
            with self.rredis.pipeline() as pipe:
                while rec:
                    key = rec[0]
                    value = rec[1]
                    pipe.set(key.decode('utf-8'), value.decode('utf-8'))
                    rec = cursor.next()
                    count += 1
                pipe.execute()
        except Exception as e:
                sys.exit('db2redis error: db {}, counter {}: {}'.format(self.rdbfile, count, e))
        print('read {} records from bdb {} to redis'.format( count, self.rdbfile ))


    def __process_set(self, item):
        #print('set event received full: {}'.format( item ) )
        # lookup redis key/value
        # set dbm key/value
        key = item['data']
        value = self.rredis.get(key)
        #print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')) )
        rc = self.rdb.put(key, value)
        if rc:
            print('{:%Y-%b-%d %H:%M:%S} db put error: {}'.format( datetime.datetime.now(), rc ))
        rc = self.rdb.sync()
        print('{:%Y-%b-%d %H:%M:%S} set event received: {} {}'.format( datetime.datetime.now(), item['channel'], item['data'] ) )

    def __process_del(self, item):
        # delete dbm key
        key = item['data']
        rc = self.rdb.delete(key)
        self.rdb.sync()
        print('{:%Y-%b-%d %H:%M:%S} delete event received: {} {}'.format( datetime.datetime.now(), item['channel'], item['data'] ) )

    def __process_default(self, item):
        print('{:%Y-%b-%d %H:%M:%S} default event received: {} {}'.format( datetime.datetime.now(), item['channel'], item['data'] ) )
    
    def __poller(self):

        try:
            for item in self.pubsub.listen():
                if item['data'] == "KILL":
                    self.pubsub.unsubscribe()
                    print("unsubscribed and finished, outta here.")
                    break

                if item['type'] == "pmessage":
                    if item['channel'].decode().endswith(":set"):
                        self.__process_set(item)
                    if item['channel'].decode().endswith(":del"):
                        self.__process_del(item)
                    elif item['channel'].decode().endswith(":expired"):
                        self.__process_del(item)
                    else:
                        self.__process_default(item)
                else:
                    self.__process_default(item)
        except Exception as e:
                print('{:%Y-%b-%d %H:%M:%S} poller error caught: {}'.format( datetime.datetime.now(), e ))
                # try again in a while so that supervisord won't be restarting again and adain
                self.rdb.sync()
                sleep(120)
                self.__poll()

    def __poll(self):
        print('{:%Y-%b-%d %H:%M:%S} redirektor poll start'.format( datetime.datetime.now() ))
        try:
        	self.rredis.ping()
        except Exception as e:
            print('{:%Y-%b-%d %H:%M:%S} error: unable to ping redis: {}'.format( datetime.datetime.now(), e ))
            sys.exit(99)

        self.pubsub = self.rredis.pubsub()
        self.pubsub.psubscribe(['__keyevent*__:set','__keyevent*__:del','__keyevent*__:expired'])
        self.__poller()


if __name__ == '__main__':
    try:
    	BDBMgr()
    except KeyboardInterrupt:
        sys.exit(0)

