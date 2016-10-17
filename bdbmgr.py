#!/usr/bin/python3

from bsddb3 import db

import argparse, sys, csv


class BDBMgr(object):
    rdb = None
    args = None

    def __init__(self):
        parser = argparse.ArgumentParser(
            prog='bdbmgr',
            description='manage berkeley db files',
            usage='''bdbmgr.py <command> [<args>]

The following subcommands are available:
   importdb     import from csv to bdb
   exportdb     export from bdb to csv
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

    def importdb(self):
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
        self.__import2db(args.db, args.csv)

    def exportdb(self):
        parser = argparse.ArgumentParser(
            description='export a bdb to csv')
        # NOT prefixing the argument with -- means it's not optional
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
        parser.add_argument('--csv', required=True )
        parser.add_argument('--db', required=True )
        args = parser.parse_args(sys.argv[2:])
        print('bdbmgr exportdb, csv={}, db={}'.format( args.csv, args.db ))
        self.rdbfile = args.db
        self.__export2csv(args.db, args.csv)

    def __import2db(self, dbfile, csvfile):

        try:
        	self.rdb.open(dbfile,None,db.DB_HASH, db.DB_CREATE)
        except Exception as e:
        	print('db error: unable to open db stack: {}'.format( dbfile ))
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


    def __export2csv(self, dbfile, csvfile):

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
                print('{:s},{:s}'.format(key.decode('ascii'), value.decode('utf-8')), file=f )
                rec = cursor.next()

        self.rdb.close()


if __name__ == '__main__':
    BDBMgr()

