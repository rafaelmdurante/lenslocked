# Why Do We Use 32 Bytes for Session Tokens?

To understand the reasoning behind why we use 32 bytes, we first need to learn a
bit about data types. A byte is a data type that stores 8 bits of data. A bit of
data is either a 0 or a 1. This means a byte can have 256 possible values. 
Putting that all together, we have two values per bit, and eight bits of data, 
so we have `2^8`, or 256, possible values that can be stored in a byte.

Now that we know how many possible values a byte has, lets explore how that 
affects an attackers ability to guess a value. Imagine that we generated session
tokens with only one byte of random data. This would mean that we only have 256
possible session tokens that we can create, and the math here is `256^1 = 256`.

If we used two bytes, we could create 65536 tokens. The math here is 
`256^2 = 65536`. Adding only one more byte vastly increases the total number of 
session tokens we can create.

If we jump to 32 bytes we have `256^32`, or `1e77` possible tokens. If we were 
to write that out, that is a 1 followed by 77 zeros. That is a massive number!

Our app’s session tokens are similar. The more users we have, the easier it is 
to guess a valid session token. Computers can also guess random session tokens 
very quickly. As we saw earlier, one of the ways to decrease the odds of 
someone guessing a value in a set of possible values is to increase the size of
that set. In other words, if our session tokens are composed of 32 bytes and 
there are many possible values, guessing a session token that is in use becomes
incredibly challenging, if not impossible. For this reason, we opt to use 32
bytes for our session tokens.

According to the OWASP Foundation, an organisation focused on security 
practices, a session should be at least 128 bits, or 16 bytes. We opt to use 32
as this is a bit better than the minimum, while not being so large that it 
hinders the performance of our application.

## Explaining the Odds

With 2 guesses there are three possible outcomes:

1. The guesser guesses right on the first try. There is a 30% chance of this
   happening, as 3 of 10 numbers are picked.
2. The guesses guesses wrong on the first try (7⁄10), then correct on the 
   second try (3⁄9). Multiplying these together, there is a 23.3% chance of this 
   happening.
3. The guesser guesses wrong on first try (7⁄10), then wrong again on the 
   second try (6⁄9). There is a 46.6% chance of this happening.
 
While the third case is by far the most likely case individually, we need to 
remember that either of the first two cases result in a win for the guesser, so
we have to add these odds together giving us a 53.3% chance of winning within 
two guesses.

The odds for three guesses work in a very similar manner, except they add up to 
a 70% chance of the guesser winning within three guesses.
