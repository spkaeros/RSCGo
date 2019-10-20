package log

/**
 * Package log is a simple package to allow RSCGo to re-use logger instances through out the entire
 * application code base.  It is a very simple package, facilitating an easy-to-use interface for the
 * entire RSCGo code base to use for logging any sort of data that it needs to.
 *
 * I may potentially expand the scope of this package to support logging to an external database of some sort,
 * but for now, it uses stdout, stderr, and ./log/*.log plain-text log files.
 */
