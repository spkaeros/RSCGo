package list

/**
 * Package list implements a generic slice-backed list API that can hold any types, and the types don't even have
 * to be uniform, though in practice they probably should be.  I may phase out this code in the near future,
 * as honestly it is nothing more than a wrapper around a standard go slice.  I am not sure, but I believe it is
 * VERY similar to the container/list package in the Go standard library.
 */
