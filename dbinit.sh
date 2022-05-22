rm -f events.db
sqlite3 events.db 'create table events(uuid text primary key, type text, severity text, event text, eventid text, descr text, ts int)'