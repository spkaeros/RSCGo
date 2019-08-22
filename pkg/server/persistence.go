/**
 * @Author: Zachariah Knight <zach>
 * @Date:   08-22-2019
 * @Email:  aeros.storkpk@gmail.com
 * @Project: RSCGo
 * @Last modified by:   zach
 * @Last modified time: 08-22-2019
 * @License: Use of this source code is governed by the MIT license that can be found in the LICENSE file.
 * @Copyright: Copyright (c) 2019 Zachariah Knight <aeros.storkpk@gmail.com>
 */

package server

import (
	"database/sql"
	"os"

	// Necessary for sqlite3 driver
	"bitbucket.org/zlacki/rscgo/pkg/entity"
	"bitbucket.org/zlacki/rscgo/pkg/list"
	_ "github.com/mattn/go-sqlite3"
)

//Objects List of the game objects in the world
var Objects = list.New(16384)

//LoadObjects Loads the game objects into memory from the SQLite3 database.
func LoadObjects() int {
	database, err := sql.Open("sqlite3", DataDirectory+string(os.PathSeparator)+"world.db")
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return 0
	}
	rows, err := database.Query("SELECT `id`, `direction`, `type`, `x`, `y` FROM `game_object_locations`")
	if err != nil {
		LogError.Println("Couldn't load SQLite3 database:", err)
		return 0
	}
	var id, direction, kind, x, y int
	counter := 0
	for rows.Next() {
		rows.Scan(&id, &direction, &kind, &x, &y)
		o := entity.NewObject(id, direction, x, y, kind != 0)
		o.Index = Objects.Add(o)
		entity.GetRegion(x, y).AddObject(o)
	}
	return counter
}
