
# qugo
QuGo is a queue manager in Go using redis.

This is based on my Queue library https://github.com/kavehmz/queue

---

## Approach

Focus of this design is mainly horizontal scalability via concurrency, partitioning and fault-detection.

My usual approach would be
- First **understanding** the field of business and its characteristics very well. How a business works matters in how we can scale it **cost effectively**.
- I would (work with my colleagues/search for/benchmark/study) several acceptable Paas/Saas solutions. Some scalability issues are already solved in other platforms. [1]
- If not able to use a ready solution for any reason like regulations, I would (work with my colleagues/search for/benchmark/study) for libre solutions.

But I didn't want to spend more than a day on this task or spend money. So I assumed a lot and fast forwarded many decisions that I would not do in a real product.


> **Note**: One assumption I make that affected code complexity was that **All Orders With The Same OrderID must be managed with a single worker and not distributed across different workers**. This I thought makes sense in a warehouse and QuGo library is able to handle it.

## Benchmarks and Picking "The Right Tools!"

This is a Queuing problems. We need a FIFO with a guarantee of atomic POPs.

### Queue
Redis is a DB with that ability. I also did benchmarks. Redis was able to handle around a 1M raw insert/reads per minute. But not depending on a limited benchmark **I added two reference benchmark to unit tests that can historically how relative performance of code**. As a reference you can see the results for AddTask and RemoveTask operations in https://travis-ci.org/kavehmz/qugo

![benchmarket](https://raw.githubusercontent.com/kavehmz/static/master/queue/benchmarket01.png)


> Note: As you see AddTask is a much lighter operation than RemoveTask. Because to make sure **We don't lose trace of running jobs in case of any crash** RemoveTask will add a log of running tasks in Redis. In case of crash that can be used to find those affected tasks.

### Language

This is a classical problem suited for those languages that support concurrency the best. Naming some Scala, Erlang, Go.

Sometimes there are restrictions about picking language. If I was forced to use a language because of restrictions I would then find the best tools for helping me as concurrent programming is all about communication (between Threads, Actors,... or what else is that language design).

Here **I picked Go** as it is a powerful concurrent language. But I didn't choose is over Scala. **I picked Go to play with this language a bit more** and finish the task in less than a day.

I would say for a production system based on my understanding **I would pick Scala** in general. Go for some usages. Erlang in some rare cases.

Regardless of language we chose, this is a concurrency problem (different than parallelism) and we need to use/apply the related best practices there.

### Design
I tried to design a horizontally scalable solution. 
For scalability we can separately add  more Redis instances and more Analysers (with as many workers as that analyser instance supports).

Partitioning factor to identify where each event must save is the OrderID.
RedisParition => OrderId % redisParitions
Queue => (OrderID/redisParitions) % queuePartitions

For example if we have redisParitions=2 and queuePartitions 3 we will save OrderID = 14 in 
Redis number 0
Queue number 1

**AddTask**
In problem there is a mentions of **The Device**. I assumed that is a single point and AddTask is normally a much faster operation than analyser. So I left it like that.

**AnalysePool**
This will use Go Routines (very light weight thread) and Channels (Go concept for communicating between Go Routines to spin up add many Analysers we want. When we run the app to as a Analyser we set two params

- -woker: number of concurrent workers which each analysers will have
- -id: Specfies the ID of analyser. This will set which redis and which queue this analyser will handle.

**Analyse Fucntions:** are defines in the app locally and passed to AnalysePool. This way AnalysePool and queue library remains **a general queue management library** and all the details of how warehouse analyses an even will be implmeneted separately. In languages like **Scala** or **Go** that functions are **first-class citizens** implementing this concept is easier.


> **Note:** AnalysePool will send all events of one OrderID to the same AnalyserWorker.

> **Note:** In Go number of workers can be specified by a simple concepts named Channel Buffering.

> **Pending Analyses and Crashes:** When removeTask get a new even from the queue it will set a flag in redis that show this task is pending. Later a waitforSuccess function will clear this flag. **If a crash happens for any reason, records of all unfinished analysis are saved in redis to for any inspection or crash handling**.

QuGo can handle multiple redis instance and in each instance can create multiple queue. This is to make sure Analysers can scale as we like. I assumed Analysers especially are the main bottleneck in scalability.

The whole design can be epxresses in one graph:

[![Diagram](https://raw.githubusercontent.com/kavehmz/static/master/queue/diagram.png)]()

I can still go on for a while about the design and the code but I think for those familiar with the concept this documents reveals enough about the ideas and process behind the choices.

### Very Important parts
This is just a practice to show my mind set. Otherwise this is not a complete solution.
In a complete solutions I would also work on,
- dev environment, I did a fast sample using **Vagrant** :https://github.com/kavehmz/queued
- metrics, monitoring and alerts (**datadog, newrelic, pagerduty**,..)
- automated setup and automated scaling (**Chef, amazon opsWorks**,..)
- I would use service discovery tools and services in my implementation (**Consul, ZooKeepr, etcd**,..)
- auto-scaling prediction and warm-up has predictable instant high-loads. For example in AWS both ELB and newly created volumes from a snapshot need **warm-up** for high volume loads. Alos **scale up-down** for costs and handling load.
- I might work on a concurrent event capturing devices and drop the assumption of one device that I interpreted from the word **The** in description. This is a very small and easy task anyway.
- I would use Scala instead of Go
- In Analysepool I would use Redis event instead of sleep to wait to queue to refill
.....

### References:
[1]http://queues.io/
