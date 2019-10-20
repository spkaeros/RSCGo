package packets

/**
 * Package packets should contain any code related to building data packets to send to the client, or that
 * RSCGo is receiving from the client.  This package should also contain helper functions to easily craft
 * (and, for packets that always contain the same information regardless of situation, store in package-scoped
 * variables) packets to send to the client, so as to avoid code duplication throughout the code base.
 *
 * This packages scope and intended features may potentially be expanded later on to contain most or all socket
 * related code in general, but for now, this is mainly to provide a clean and easy to use interface for building
 * packet data structures.
 */