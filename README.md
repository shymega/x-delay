milter-x-delay
==============

## What is it?

This project will process outgoing messages from a compatible MTA,
such as Postfix.

If the headers _do_ include this header: `X-Delay-Timestamp`, then it
will read the aforementioned header, which should have a ISO 8601
compliant timestamp.

Once the milter reaches the timestamp specified, it will send the
message to the MTA to be send ASAP.

## Status

Currently working locally on the milter. Will be out when its out.

## Copyright

This program is licensed under the Apache 2.0 license, copyright Dom
Rodriguez (2019)
