# contravent

> This is very TODO level at the moment

Contract testing for events. Using an externally defined schema verify consumers
and producers.

How to use:

- define events using JSON schema somewhere
- write something that produces that event
- write a test to verify that the produced event matches the schema
- (somewhere else) consume that event
- (and) write a test to verify it handles that shape of event

See [the example](./example).

## Difference to Pact

This uses an externally defined schema to match against, whereas Pact produces
the schema that will be matched from the tests themselves.
