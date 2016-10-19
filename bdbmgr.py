#!/usr/bin/python3

from bsddb3 import db

import argparse, sys, csv, redis


class BDBMgr(object):
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
   syncdb       poll redis master and update bdb with changes
''')
        parser.add_argument('subcommand', help='subcommand to run')
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
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
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
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
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
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
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
        parser.add_argument('--redis', required=None )
        parser.add_argument('--db', required=True )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr redis2db, db={}'.format( args.db ))
        if not args.redis:
            self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        self.rdbfile = args.db
        self.__redis2db(args.db)

    def csv2redis(self):
        parser = argparse.ArgumentParser(
            description='import from csv file to redis')
        # prefixing the argument with -- means it's optional
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
        parser.add_argument('--csv', required=True )
        parser.add_argument('--redis', required=None )
        # now that we're inside a subcommand, ignore the first
        # TWO argvs, ie the command (git) and the subcommand (commit)
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr csv2redis, csv={}'.format( args.csv ))
        if not args.redis:
            self.rredis = redis.StrictRedis(host='localhost', port=6379, db=0)
        self.__csv2redis(args.csv)


    def __csv2db(self, dbfile, csvfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('error: unable to open db: {}'.format( dbfile ))
        	print(e)
        	sys.exit(99)

        cursor = self.rdb.cursor()
        rec = cursor.first()

        with open(csvfile, newline='') as f:
            reader = csv.reader(f)
            try:
                for row in reader:
                    self.rdb.put(str.encode(row[0]), str.encode(row[1]))
            except Exception as e:
                sys.exit('csv error: file {}, line {}: {}'.format(csvfile, reader.line_num, e))

        self.rdb.close()


    def __db2csv(self, dbfile, csvfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH,db.DB_DIRTY_READ)
        except Exception as e:
        	print('error: unable to open db stack: {}'.format( dbfile ))
        	print(e)
        	sys.exit(99)

        with open(csvfile, 'w') as f:
            cursor = self.rdb.cursor()
            rec = cursor.first()

            while rec:
                key = rec[0]
                value = rec[1]
                print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')), file=f )
                rec = cursor.next()

        self.rdb.close()

    def __redis2db(self, dbfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('error: unable to open db: {}'.format( dbfile ))
        	print(e)
        	sys.exit(99)

        cursor = self.rdb.cursor()
        rec = cursor.first()

	# pull all redis entries
        try:
            for key in self.rredis.scan_iter():
                value = self.rredis.get(key)
                #print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')) )
                self.rdb.put(key, value)
        except Exception as e:
        	print('error scanning redis into db: {}'.format( dbfile ))
        	print(e)
        	sys.exit(99)

        self.rdb.close()

    def __csv2redis(self, csvfile):

        with open(csvfile, newline='') as f:
            reader = csv.reader(f)
            try:
                with self.rredis.pipeline() as pipe:
                    for row in reader:
                        pipe.set(str.encode(row[0]), str.encode(row[1]))
                    pipe.execute()
            except Exception as e:
                sys.exit('csv2redis error: file {}, line {}: {}'.format(csvfile, reader.line_num, e))


if __name__ == '__main__':
    BDBMgr()

