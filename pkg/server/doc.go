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

/**
 * Package server implements the core logic of the server.  It contains all of the networking code, all of the
 * core logic of the game, all of the packet handlers, all of the packet builders, etc.  It contains the main game
 * engine loop that processes the necessary synchronous actions every 650ms(standard jagex game tick rate).
 *
 * This package could use some seperation of concerns, honestly.  I want to isolate networking code to its own
 * package, I want to isolate packet handlers, packet builders, etc.  I intend to do some experimenting with
 * its structure once more of the main game features are implemented and the game is actually functioning properly.
 */
