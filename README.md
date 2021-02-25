# No Coin

Attempt to learn more about crypto by building my own implementation

# Question

Is it possible to build a chain that allows users to run dapps and stored data in a decentralized manner without
the headache of speculators making coin prices surge and therefore the cost of executing a dapp sore.

One possible idea is to use a central authority to distribute tokens for a stable price. This central authority would
sell and buy tokens at this same stable price and start by premining all tokens on the network. You could have the same
transfer of ownership models as other cryptos but by having this central authority you would be able to control the price
of the tokens.

I have no idea how feasable this idea even is. Or even if it is a good idea or if there is some fundemental flaw that I
haven't understood yet. I probably will fail at even building the simplest implementation but this is an opportunity to
learn about these types of systems and have a bit of fun at the same time!

The idea is to try and copy as many ideas from Bitcoin and Ethereum. It would be interesting to see if it would be possible
to run dapps in a WASM runtime? These dapps would basically be state reducers in that they take an input, modify the dapp state
which becomes the new output.

## TXs

3 types of transactions

TX credit transfer: 
UTXOs in
UTXOs out

TX function deploy: 
code in
UTXOs in <- TX fee
UTXOs out

TX function call: 
args in
previous return in
credit limit <- gas limit
UTXOs in <- TX fee
UTXOs out
return out

UXTO pool
Function return pool

# References

- Bitcoin original code :: [link](https://github.com/bitcoin/bitcoin/tree/4405b78d6059e536c36974088a8ed4d9f0f29898)
- Proof of Work :: [link](https://medium.com/blockchaintechnologies/blockchain-mechanics-proof-of-work-75f5df8c1c35)
- Digital signing :: [link](https://en.wikipedia.org/wiki/Digital_signature)
- What a block is made up of :: [link](https://learnmeabitcoin.com/technical/blkdat)
- Transaction data in the block :: [link](https://learnmeabitcoin.com/technical/transaction-data)
- Peer discovery in bitcoin network (this article is actually brilliant!) :: [link](http://sebastianappelt.com/understanding-blockchain-peer-discovery-and-establishing-a-connection-with-python/)
- Where is ethereum byte code stored :: [link](<https://stackoverflow.com/questions/52374352/where-bytecode-is-stored#:~:text=1%20Answer&text=Contracts%20live%20on%20the%20blockchain,Ethereum%20Virtual%20Machine%20(EVM).&text=Contract%20addresses%20have%20bytecode%20associated,private%20keys%20behind%20the%20contract>)
