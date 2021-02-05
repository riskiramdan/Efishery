// Code generated by rice embed-go; DO NOT EDIT.
package databases

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "20200205205811_create_table.down.sql",
		FileModTime: time.Unix(1612533169, 0),

		Content: string("drop table if exists \"user\";\ndrop table if exists \"role\";"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "20200205205811_create_table.up.sql",
		FileModTime: time.Unix(1612533176, 0),

		Content: string("-- Table Definition ----------------------------------------------\nCREATE TABLE \"user\" (\n  \"id\" SERIAL PRIMARY KEY NOT NULL,\n  \"roleId\" int NOT NULL,\n  \"name\" varchar(80) NOT NULL,\n  \"phone\" varchar(80) NOT NULL,\n  \"password\" varchar NOT NULL,\n  \"token\" varchar,\n  \"tokenExpiredAt\" timestamp,\n  \"createdAt\" timestamp NOT NULL DEFAULT (now()),\n  \"createdBy\" varchar(20) DEFAULT 'admin',\n  \"updatedAt\" timestamp NOT NULL DEFAULT (now()),\n  \"updatedBy\" varchar(20) DEFAULT 'admin',\n  \"deletedAt\" timestamp,\n  \"deletedBy\" varchar(20)\n);\n\nCREATE TABLE \"role\" (\n  \"id\" SERIAL PRIMARY KEY NOT NULL,\n  \"name\" varchar(80) NOT NULL,\n  \"createdAt\" timestamp NOT NULL DEFAULT (now()),\n  \"createdBy\" varchar(20) DEFAULT 'admin',\n  \"updatedAt\" timestamp NOT NULL DEFAULT (now()),\n  \"updatedBy\" varchar(20) DEFAULT 'admin',\n  \"deletedAt\" timestamp,\n  \"deletedBy\" varchar(20)\n);\n\nALTER TABLE \"user\" ADD FOREIGN KEY (\"roleId\") REFERENCES \"role\" (\"id\");\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1612533176, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "20200205205811_create_table.down.sql"
			file3, // "20200205205811_create_table.up.sql"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`./migrations`, &embedded.EmbeddedBox{
		Name: `./migrations`,
		Time: time.Unix(1612533176, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"20200205205811_create_table.down.sql": file2,
			"20200205205811_create_table.up.sql":   file3,
		},
	})
}
