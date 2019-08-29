# RSCGo

This project is being designed with security, simplicity, portability, and performance in mind.

The other RSClassic server implementations that I've seen that were written in Java are plagued by a number of issues, not the least 
of which being feature creep, performance problems due to reliance on the JVM, and memory leaks.

A large portion of the codebase they're based on was written over a decade ago by a team of hobbyists in an attempt to breathe life 
into a then-dying game, and only a few developers have bothered to try to really improve the code in between then and now, with 
varying degrees of success.

Bottom line, the codebases I've seen being used are sub-par and usually large, complex monstrosities.

RSCGo was designed to be simple, lightweight, and fast.  It was designed to leverage modern technologies to provide a simple yet
highly performant RSClassic server.  It should be able to handle large player loads with no issues, and it should be hardened against
securiy issues we've seen in practice on other RSClassic servers.  It should eventually behave just like Jagex's RSClassic
implementation had.

Currently, RSCGo supports logging in, supports multiple players, chat works, walking works, objects are where they should be, etc.

Here is a progress update picture: https://i.imgur.com/ZjzgBcE.png