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
