BEGIN TRANSACTION;
DROP TABLE IF EXISTS "SourceFile";
CREATE TABLE IF NOT EXISTS "SourceFile" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Name"	TEXT,
	"ObjectID" TEXT,
	"VersionList"	TEXT,
	"Checksum"	TEXT,
	"Filename"	TEXT,
	"FileModTime"	INTEGER,
	"FileSize"	INTEGER
);
DROP TABLE IF EXISTS "ServerFile";
CREATE TABLE IF NOT EXISTS "ServerFile" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Name"	TEXT,
	"ObjectID" TEXT,
	"VersionList"	TEXT,
	"Checksum"	TEXT,
	"Filename"	TEXT,
	"FileModTime"	INTEGER,
	"FileSize"	INTEGER
);
COMMIT;
