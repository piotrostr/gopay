# gopay

Some random thought that I had is that since it is impossible to perform some
operations on-chain since solidity cannot really handle much
computational/storage load, it would be nice if there was something like a
stripe-like mixin that would enable communication with the chain and run code
on a server in Python or Go but only if there is a transaction confirmed, if it
is exisiting and not confirmed await its confirmation. Something similar to
what was done in SMPLverse, but more flexible. All of this sugarcoated nicely
for the user and with a neat sdk for devs.

## Setup

In order to run tests, the env variable of `PRIVATE_KEY` has to be set, as well
as the address corresponding to the private key in the `genesis.json` file with
non-zero balance

## Core Functionality

Core assumption is that there is some service offered on one-time-basis
could be something like iStats menus or a lifetime subscription fee.

Could be something like converting a bunch of pdfs one-time or
buying a coffee for someone.

The receipt is hashed and put in the smart contract during the transaction to
be checked.

The api waits for the transaction to finish and verifies the payment.

Essentially, all of the information goes to the blockchain as well as the server.

```graph
               /->blockchain
              /      =
             /       |
            / tx     | await the transaction, verify tx data
           /         |
          /          =
client  -->-------backend---------> accept payment or handle error
            tx                            -> store the information and make it
           data                              available in the private api
```

The server verifies the transaction and then unlocks the functionality.

All of the payments are stored on-chain, so if the same address comes
back it can be re-checked.

The tx data contains a hash of the functionality to be provided,
done both on the server side and the client side to ensure that the pipeline is secure.

Ideally, I would like to replicate the way that Stripe or other payment
processors work but for Crypto.

## Additional assumptions

The main api-route should be idempotent, meaning that

- there is a payment existing -> it is unlocked
- there is no payment for the `owner address` - payment is made -> it is unlocked

Smart contract is very simple, only handles the payment and stores the hashes
of content.
