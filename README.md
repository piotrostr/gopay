# gopay

Some random thought that I had is that since it is impossible to perform some
operations on-chain since solidity cannot really handle much
computational/storage load, it would be nice if there was something like a
stripe-like mixin that would enable communication with the chain and run code
on a server in Python or Go but only if there is a transaction confirmed, if it
is exisiting and not confirmed await its confirmation.

All of this sugarcoated nicely for the user and with a neat sdk for devs.

Something similar to what was done in SMPLverse, but more flexible.
