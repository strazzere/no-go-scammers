# no-go-scammers

Simple `golang` script which will stand up a simplistic web server and all you to call out to a specified scammer number, play the text + mp3 on a loop while recording.

The current default mp3 is Homer Simpson stating `Wait a minute. This could be some kind of scam... or possibly scamola!`.

# What? Why?

I got a few calls from "DEA" / "Federal Agents" / "Boarder and Customs" that was a recording, stating I must call a different number to follow up on an active case. Upon calling the second number, they then state you need to wire them money due to an outstanding warrent (IRS modified scam?).

So I set up a fake twilio number which would "catch" these recordings and also send a recording giving out that number. Once the inbound calls are processed, it was fed to a script which performed this call logic. Basically it appears the scammers could not process more than ~3 inbound calls and didn't understand how to block specific numbers (or restricted ones). This resulted in their lines going down for approximately ~2 weeks until they disconnected and got a new one. The hopeful "theory" I was working on is that all their voicemails sent out, will still point to the old number which now no longer works. The honeypots should catch updated numbers, rinse, repeat.

# How?

Just fix the `const` values and run the "server" via `go rung no.go`. This is not a secure script, it likely exposes things in a poor manner and I would not expose this to the internet without looking at the code yourself.

## Warnings

Please note:
 * This is very likely against Twilio's TOS, so don't abuse it against folks who aren't scammers
 * Recording these calls is highly dependent on the juristiction of where you are calling from and whom you are recording, IANAL so look it up and do not rely on me!

TODO:
 - [ ] Use throw away numbers vs just one
 - [ ] Track which throw away numbers are used with what outgoing scammer numbers
 - [ ] Commit the "honeypot" code, to catch the scammer number changes
