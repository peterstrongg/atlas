import sys
import sqlite3

if len(sys.argv) < 2:
    exit(-1)

# Create database and table5
conn = sqlite3.connect("mac_vendors.sqlite")
curs = conn.cursor()
curs.execute('''
    CREATE TABLE IF NOT EXISTS mac_vendors (
        oui text,
        org text
    )
''')
conn.commit()

with open(sys.argv[1], "r", encoding="utf-8") as f:
    for line in f:
        oui = line[0:6]
        org = line[7:-1]
        curs.execute("INSERT INTO mac_vendors VALUES(?,?)", (oui, org))

conn.commit()
conn.close()
        

    

    