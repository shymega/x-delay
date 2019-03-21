x-delay
=======

## What is it?

This SMTP server accepts outgoing emails from a (compatible) MTA, and when a
condition is met, releases them back to the MTA to be processed.

If the headers of an email include this header: `X-Delay-TS`, then it
will read the aforementioned header, which should have a ISO-8601
compliant timestamp.

The SMTP server will store the email in an internal queue, and wait for the
system time (with accordance to timezones) to reach the `X-Delay`
timestamp specified in the message.

Once this condition is met, the SMTP server will inject the message back into
the MTA's queuing system.

## Development status

Working locally on the SMTP server, expect commits soon!

## License

This project is licensed under the Apache 2.0 license.
