
## merge-api

This is a restful API which returns at most max sorted posts 
combined from two sources. It has two GET endpoints, /health and /posts.

### To build and run (go 1.11+ required):
  
`go build`

`./merge-api`

This starts the api locally on (by default) port 8080.

### Try it with:

curl http://localhost:8080/health

curl http://localhost:8080/posts?max=5 

curl 'http://localhost:8080/posts?max=10&order=asc'

### Unit tests

`go test`

### Some issues and assumptions:

* https://fakerql.com seems to be broken (the DNS doesn't resolve).
So I just made the client for fakerql return some fake data generated locally.

* Error handling: I'm not sure about the error http response codes/messaging 
expected by clients so I put a few placeholders in. A team working on 
this would want to figure this out together.

* I've added a few validation checks to ensure the queries
are reasonable, but I don't know what sort of validation
should be done.

* I added some wiring for Cloud Datastore, if only 
to connect it to something, but I'm not sure the choice of database 
is very clear given the amount of information given. The data
looks relational-- many posts probably have the same author.
The makes a sql db look reasonable. On the other hand, that
would mean doing a join on every query for this api which
may not be wanted for performance reasons.

* Assuming that simple is better,
the algorithm is simple and not optimized: it retrieves max
items from each source, joins the two lists, sorts them and
returns max items. This might be improved. For example
one could optimize the go memory allocations or have 
the two sources do the sorting and and merge in the api.

* Setting max to an extremely large number will cause this 
algorithm problems (memory use, sorting time, time to query remote sources). 
To avoid this, I added an arbitrary hard limit of 10000 posts
guessing that clients would never need to ask for more than that.
On the other hand, if max can be very large sometimes, I guess 
we'd need to figure out how to paginate or stream the results
effectively.

* I was not sure whether duplicates can appear in the combined
list from both sources, so I assumed not. But if so, and they 
should be removed, duplicate removal would have to be added in.

* I was not sure whether, if clients ask again, they are expecting
to get a completely different list of postings, or the most recent
postings whether or not that has changed--- even if this overlaps
with what was returned on a previous query.



### Some possible further things to do:

* Caching: presumably caching by the api (or by a cdn like fastly) would be useful.
* Add authentication method for clients. Currently the endpoints are completely open. 
* Add secrets/authentication between api and database.
* Use structured/leveled logging. I've put in some very basic logging.
* Add metrics/tracing for observability.
* Implement clean shutdown. 
* Unit testing is thin. Should be evaluated for reasonable, judicious coverage.
* Integration testing 


### Deployment

I'm not sure how issues like price, performance, scalability,
speed of development and ease of maintenance balance in this use case, 
but Google App Engine (standard) is a likely place to run this;
for the app itself this would require not much more than setting up an app.yaml
file and changing the health check url to the place gcloud expects.
The database would also need to be configured.
GKE (Kubernetes in google cloud) is another possibility 
if GAE is too limiting but is a bit more involved to set up.


