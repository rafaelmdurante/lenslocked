# Hash Session Tokens

## Context

Our SessionService can generate tokens, but we don’t plan to store those 
directly in our database. Instead, we plan to store the hash of a token in our 
database. This is why we defined the sessions table in our database with a 
`token_hash` column.

```sql
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INT UNIQUE,
  token_hash TEXT UNIQUE NOT NULL
);
```

We do this to make our application more secure. If we store a session token in 
our database, an attacker who accesses our database will have access to valid 
session tokens and can impersonate users. If we only store a hash of the session
token, attackers won’t be able to reconstruct the original session token to 
impersonate them.

In order to validate a users session, we will: 
1. Read the session token from the request cookies.
2. If a session token is present, we will hash it. If not, we know the user 
   isn’t logged in.
3. Once we have a hashed session token we will search our database for a session
   with the same token hash, which will help us determine who the current user is.

This all means that we need to determine which hashing algorithm we will be
using to store our session tokens, then we can proceed with inserting sessions
into our database.

While it might be tempting to use `bcrypt` like we did with our passwords, this 
approach isn’t a good fit for session token hashes because it adds a unique salt
to every hash. If we were to hash the same session token three times in a row 
using bcrypt, we would get a different result every time because the salt would 
be different. The only way to validate a `bcrypt` hash is to use have access to 
both the hashed value, and the value we want to test.

```go
bcrypt.CompareHashAndPassword(hashedToken, plaintextToken)
```

Technically we could make this work by storing the user’s ID in a cookie, and 
using that to look up their sessions. This would give us access to the hashed 
tokens, but it would be slow, add unnecessary complexity, and make it more 
challenging to support multiple sessions per user in the future. What is even 
worse, we would be doing all of this extra work to use a hashing function that 
isn’t well suited to our needs. After all, our session tokens are:

1. Generated randomly.
2. Replaced every login.
3. Not reused across websites like passwords may be.

The end result is that bcrypt just doesn’t make sense, and we will instead be 
exploring other hashing functions.

## Conclusion

There area few options that we could explore that fit our needs, but 
`crypto/sha256` is probably the most common and best suited hashing algorithm 
for our needs.

### Why not HMAC?

HMAC (`crypto/hmac`) works by using a secret key, which is similar to the salt 
we apply to our passwords before we salt them, except the same secret key is 
applied to all values that are being hashed. Once the secret key is applied, a 
hashing algorithm like `SHA256` is used to do the actual hashing.

The primary benefit of HMAC is that if an attacker gets access to our database 
but does NOT have access to our secret key, they cannot even attempt rainbow 
table attacks. While this sounds good, in reality our session tokens are 
randomly generated and replaced at login, so a rainbow table attack is virtually
impossible to succeed.

On top of HMAC not really adding any extra security, chances are an attacker who
has access to our database will also have access to our secret key, making the 
use of HMAC pointless.




