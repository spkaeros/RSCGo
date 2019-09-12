package world

/**
 * Package world contains all of the code that has to do with the actual game world, and the entities within
 * said game world.  It may undergo a renaming in the near future to world, or something like that.
 *
 * It contains code that splits up the entire map into 48x48 chunks of 2304 tiles, as holding every tile in
 * memory would be too expensive on RAM, and would be generally less efficient.  Each region is initialized
 * on-demand and if no entity is ever in a region, then that region will stay unloaded until a mobile entity
 * travels into it.  Generally, when updating clients, we must grab their own region and the closest adjacent
 * regions, just in case they are close enough to the edge of their own region to see entities within other
 * regions within their view area.
 *
 * This package also contains all of the code for path finding.  Currently, the client generates the path, and
 * all the server does is traverse the clients generated path.  This may change in the future to the server
 * generating all of the path waypoints and such, or it may stay like this.  I must add clipping, since as a
 * general rule of thumb, we can't trust anything the client tells us, e.g if a modified client sends us a path
 * which walks through walls, currently the server code will honor that path and walk right through the wall,
 * since I haven't loaded the map data into memory on the server yet.  Once I load the map into the server,
 * this sort of behavior will be rectified most likely in the file pathway.go.
 *
 * All scene entities, e.g players, npcs, ground items, boundarys(doors, some fences), and game objects will all
 * be represented by data structures within this package.  Same with most of the logic for these entities, as well.
 */
