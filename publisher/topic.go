/**
* Message sent to a `topic` exchange can't have an arbitary `routing key`, it must be a list of words, delimited by dots `.`
* It is like `direct` exchange, but, a message with a particular `routing key` would be delivered to all queues with a matching `binding key`.
* `*` star can substitute for exactly one word
* `#` hash can substitute for zero or more words
* When a queue is bound with `#`, it will receive all the messages regardless of routing key, like in `fanout` exchange
* When a special characters '*' or '#' aren't used, it would behave as `direct` exchange
 */
package main

// When a new consumer is spawned, why the previous message is not consumed?
// If the messages are routed to exchanges that are not bind by any queue, are those lost ?
// Or, if exact binding of queue is done, why aren't they recovered ?
