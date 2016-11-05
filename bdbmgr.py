#!/usr/bin/python3

from bsddb3 import db
from urllib.parse import urlparse

import argparse, sys, csv, redis


class BDBMgr(object):
    version = '0.1.0'
    rdb = None
    rredis = None
    args = None

    def __init__(self):
        parser = argparse.ArgumentParser(
            prog='bdbmgr',
            description='manage berkeley db files',
            usage='''bdbmgr.py <command> [<args>]

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
        self.rdb = db.DB()
        # use dispatch pattern to invoke method with same name
        getattr(self, args.subcommand)()

    def csv2db(self):
        parser = argparse.ArgumentParser(
            description='import from csv file to bdb')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--csv', required=True )
        parser.add_argument('--db', required=True )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr importdb, csv={}, db={}'.format( args.csv, args.db ))
        self.rdbfile = args.db
        self.__csv2db(args.db, args.csv)

    def db2csv(self):
        parser = argparse.ArgumentParser(
            description='export a bdb to csv')
        # NOT prefixing the argument with -- means it's not optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--csv', required=True )
        parser.add_argument('--db', required=True )
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr exportdb, csv={}, db={}'.format( args.csv, args.db ))
        self.rdbfile = args.db
        self.__db2csv(args.db, args.csv)

    def redis2db(self):
        parser = argparse.ArgumentParser(
            description='import from csv file to bdb')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--db', required=True )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr redis2db, db={}'.format( args.db ))

        self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        if args.redis:
            self.rredis = redis.from_url(args.redis)

        self.rdbfile = args.db
        self.__redis2db(args.db)

    def redis2csv(self):
        parser = argparse.ArgumentParser(
            description='export from redis to csv')
        # NOT prefixing the argument with -- means it's not optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--csv', required=True )
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr redis2csv, csv={}'.format( args.csv ))
        self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        if args.redis:
            self.rredis = redis.from_url(args.redis)

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
        print('bdbmgr csv2redis, csv={}'.format( args.csv ))
        self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        if args.redis:
            self.rredis = redis.from_url(args.redis)

        self.__csv2redis(args.csv)

    def db2redis(self):
        parser = argparse.ArgumentParser(
            description='import from bdb file to redis')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--db', required=True )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr db2redis, db={}'.format( args.db ))
        self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        if args.redis:
            self.rredis = redis.from_url(args.redis)

        self.rdbfile = args.db
        self.__db2redis(args.db)

    def poll(self):
        parser = argparse.ArgumentParser(
            description='subscribe to redis events and update bdb')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s {version}'.format(version=self.version))
        parser.add_argument('--redis', required=None, help='redis connection string: redis://localhost:6379/0' )
        parser.add_argument('--db', required=True )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr poll, db={}'.format( args.db ))

        self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        if args.redis:
            self.rredis = redis.from_url(args.redis)

        self.rdbfile = args.db
        self.__poll(args.db)

    def __csv2db(self, dbfile, csvfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('error: unable to open db: {}'.format( dbfile ))
        	print(e)
        	sys.exit(99)

        cursor = self.rdb.cursor()
        rec = cursor.first()
        count = 0

        with open(csvfile, newline='') as f:
            reader = csv.reader(f)
            try:
                for row in reader:
                    self.rdb.put(str.encode(row[0]), str.encode(row[1]))
                    count += 1
            except Exception as e:
                sys.exit('csv error: file {}, line {}: {}'.format(csvfile, reader.line_num, e))

        self.rdb.close()
        print('read {} records from csv to bdb: {}'.format( count, csvfile, dbfile ))


    def __db2csv(self, dbfile, csvfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH,db.DB_DIRTY_READ)
        except Exception as e:
        	print('error: unable to open db stack: {}'.format( dbfile ))
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

        self.rdb.close()
        print('read {} records from bdb to csv: {}'.format( count, dbfile, csvfile ))

    def __redis2db(self, dbfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('error: unable to open db: {}'.format( dbfile ))
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
        	print('error scanning redis into db: {}'.format( dbfile ))
        	print(e)
        	sys.exit(99)

        self.rdb.close()
        print('read {} records from redis into db: {}'.format( count, dbfile ))

    def __redis2csv(self, csvfile):

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

        count = 0
        with open(csvfile, newline='') as f:
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

    def __db2redis(self, dbfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH,db.DB_DIRTY_READ)
        except Exception as e:
        	print('error: unable to open db stack: {}'.format( dbfile ))
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
                sys.exit('db2redis error: db {}, counter {}: {}'.format(dbfile, count, e))

        self.rdb.close()
        print('read {} records from bdb {} to redis'.format( count, dbfile ))


if __name__ == '__main__':
    BDBMgr()

