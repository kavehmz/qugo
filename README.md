# qugo
[![Build Status](https://travis-ci.org/kavehmz/qugo.svg)](https://travis-ci.org/kavehmz/qugo)

This is a queue manager using redis.

The main buzwords around its design are concurrency, paritioning and fault-detection [1]

I had some assumptions about the design.
- I just used my understading of problem to avoid lenghty discussions. For a real product I dont look for a solution until I understand the problem properly. It is vital to understand all aspects of a problem to be able to offer a good solution. There are different solutions in IT world for a reason. There is no sivler bullet.
- I wanted to finish the project in less than a day (So I did not do thorough search of all available options and languages). For a real solutions I will search, benchmark and talk to other team members.
- I didnt wanted to do it using local, easiliy setupable tools. Something that I wont do for real products. 
- I make the solution more Complex by assigning all events of the same Order to one Analyser as I found it more realistic. Without that solution would be far simpler.

The whole design can be epxresses in one graph:

[![Diagram](https://raw.githubusercontent.com/kavehmz/static/master/queue/diagram.png)]()
