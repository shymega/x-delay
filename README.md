milter-x-delay
==============

## What is it?

This milter will process outgoing emails from a compatible MTA,
such as Postfix, or Exim. This would process emails sent from
SMTP/Submission/SMTPS **and** Sendmail.

If the headers of an email include this header: `X-Delay`, then it
will read the aforementioned header, which should have a ISO-8601
compliant timestamp.

The milter will store the email in an internal queue, and wait for the
system time (with accordance to timezones) to reach the `X-Delay`
timestamp specified in the message.

Once this condition is met, the milter will allow the email to be sent.

Once the milter reaches the timestamp specified, it will release the
message to the MTA.

## Development status

Working locally on the milter, expect commits soon!

## License

This project is licensed under the Apache 2.0 license.
