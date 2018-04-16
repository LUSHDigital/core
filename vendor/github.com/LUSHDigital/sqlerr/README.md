# SQLErr

When `mysql` is built, it runs a tool called [`comp_err`](https://dev.mysql.com/doc/refman/5.7/en/comp-err.html).

*From MySQL's documentation:*

> comp_err creates the errmsg.sys file that is used by mysqld to determine the error messages to display for different error codes. comp_err normally is run automatically when MySQL is built. It compiles the errmsg.sys file from the text file located at sql/share/errmsg-utf8.txt in MySQL source distributions.

In essence, the internal package contains an extremely stripped down version of `comp_err`, which only handles the error definitions and error codes, as we need them, but ignores the translations and error messages.

The `Stringer` implementation is provided by the Go team's fantastic [stringer](https://godoc.org/golang.org/x/tools/cmd/stringer) tool.