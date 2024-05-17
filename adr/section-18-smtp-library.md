# SMTP Library

## Context

When constructing an email with SMTP, it will have a format similar to this:

```yaml
MIME-Version: 1.0
Date: Sun, 22 Jan 2023 11:54:48 -0500
To: jon@calhoun.io
From: test@lenslocked.com
Subject: This is a test email
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset=UTF-8

This is the body of the email
```

Modern emails tend to be more complicated and use HTML tags as well.

```yaml
MIME-Version: 1.0
Date: Sun, 22 Jan 2023 11:55:59 -0500
Subject: This is a test email
To: jon@calhoun.io
From: test@lenslocked.com
Content-Type: multipart/alternative;
 boundary=2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f

--2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f
Content-Transfer-Encoding: quoted-printable
Content-Type: text/plain; charset=UTF-8

This is the body of the email
--2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f
Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=UTF-8

<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>
--2df315b0bd754b2cea495f617b327626853ede3bcbf4608725384a95937f--
```

The standard SMTP package provided by Go is very limited as by today and there
is no plans to add features at the moment. So for SMTP, third party libraries
tend to be a better option in Go.

## Conclusion

We will be using `go-mail/mail` in this course.

```bash
go get github.com/go-mail/mail/v2
```

Link to [Go-mail v2 docs](https://pkg.go.dev/github.com/go-mail/mail/v2).
