#!/usr/bin/python3

from bsddb3 import db


import argparse
import sys


class BDBMgr(object):
    db = None

    def __init__(self):
        parser = argparse.ArgumentParser(
            prog='bdbmgr',
            description='manage berkeley db files',
            usage='''bdbmgr.py <command> [<args>]

The following sub commands are available:
   importdb     import from csv to bdb
   exportdb     export from bdb to csv
   syncdb       poll redis master and update bdb with changes
''')
        parser.add_argument('command', help='subcommand to run')
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
        # parse_args defaults to [1:] for args, but you need to
        # exclude the rest of the args too, or validation will fail
        args = parser.parse_args(sys.argv[1:2])
        if not hasattr(self, args.command):
            print('unrecognized command')
            parser.print_help()
            exit(1)
        # use dispatch pattern to invoke method with same name
        getattr(self, args.command)()

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
        print('running bdbmgr importdb, csv={}, db={}'.format( args.csv, args.db ))

    def exportdb(self):
        parser = argparse.ArgumentParser(
            description='export a bdb to csv')
        # NOT prefixing the argument with -- means it's not optional
        parser.add_argument('--version', action='version', version='%(prog)s 0.1.0')
        parser.add_argument('--csv', required=False )
        parser.add_argument('--db', required=True )
        args = parser.parse_args(sys.argv[2:])
        print('running bdbmgr exportdb, csv={}, db={}'.format( args.csv, args.db ))
        export2csv(args.db)


def export2csv(dbfile):
    rdb = db.DB()
    try:
    	rdb.open(dbfile,None,db.DB_HASH,db.DB_DIRTY_READ)
    except Exception as e:
    	print('error: unable to open db stack: {}'.format( dbfile ))
    	print(e)
    	sys.exit(99)

    cursor = rdb.cursor()
    rec = cursor.first()

    while rec:
        key = rec[0]
        value = rec[1]
        print('{:s},{:s}'.format(key.decode('utf-8'), value.decode('utf-8')))
        rec = cursor.next()
    rdb.close()


if __name__ == '__main__':
    BDBMgr()

