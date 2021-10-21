- Start Date: 2021-10-18
- Target Version: 0.21.0

# Summary
This PR proposes adding a new attribute named `tags` to applications to allow more flexible filtering.

There are primarily two possible filtering methods: labels in the form of `key=value` and tags, but tags are simpler for both the user and the design.
The labels work a bit too much, as this feature is only for users to distinguish, not for the system to filter.
Therefore, this RFC covers only about tags.

# Motivation
Currently, it is able to filter with the embedded attributes `Environment` and `Kind`.
It works fine for a relatively tiny project with a few developers, but for a huge project like with so many microservices and multiple teams sharing the responsibility, it gets difficult to find.
It makes easier to find out applications they'd like that it allows their own attributes.

It would be nice to be able to share links among team members like:

```
https://control-plane.dev/applications?tags=searvice-a,team-1
```

# Difficulty
**One line**: array exact matching is difficult.

What we'd really like to do is to search for Applications/Deployments that have all the tags we specified.

Like for instance, let's say we have an application with tags: `["service-a", "team-1"]`

Then you requested: `["service-a", "team-1"]`

In this case, that application would be returned since all tags you requested is in the application.

But if you requested: `["service-a", "team-1", "payment"]`

In this case, the application would not be returned since the application we have doesn't have the "payment" tag.

That is, what we'd like is kind of like:

```sql
WHERE tags CONTAINS service-a
AND
WHERE tags CONTAINS team-1
```

RDBs allow you to relate them using a junction table, but we also support Firestore, a document-oriented database.
Firestore has a filter with the name `array-contains-any` which looks for matches on any of the specified elements.
But none of queries can find out documents that have all specified tags.

Besides, our MySQL also uses the JSON data type now, hence it is hard to relate between tables even in MySQL.

# Detailed design
To deal with the above difficulty, this RFC proposes inserting the application id and tag into [Elasticsearch](https://www.elastic.co/elasticsearch), which supports exact match for arrays.

```
  │
PipeCD
  │
  │
  │                ┌───────────────────┐
  ├──[AppID: Tags]──>  Elasticsearch
  │                └───────────────────┘
  │
  │                ┌───────────────────┐
  └──[Application]──>    Datastore
                   └───────────────────┘
```

Elasticsearch offers [the helm chart](https://github.com/elastic/helm-charts) so we can use it as child one. Values will be like:

```yaml
elasticsearch:
  sysctlInitContainer:
    enabled: false
```

Let's drill down to see how to implement it.

**When creating/updating**

First, we need to make changes when we create/update an application where the tag was created/updated.
Insert the tags and applicationID into Elasticsearch as well as into Datastore.

```
PUT /applications/_doc/1?refresh
{
  "id": "xxxx-yyyy",
  "tags": [ "service-a", "team-1" ],
  "required_matches": 2
}
```

**When fInding**

In the case of `ListApplication` with tags specified, first, get the IDs of the applications that have all the specified tags from Elasticsearch, and then pass those IDs to Datastore's IN query.
It's the same in case of `ListDeployment`.

```go
if len(o.Tags) > 0 {
	appIDs := elasticsearch.Search("Applications", o.Tags) // This is fuzzy code
	filters = append(filters, datastore.ListFilter{
		Field:    "Id",
		Operator: datastore.OperatorIn,
		Value:    appIDs,
	})
}
```

To retrieve documents that have all specified tags, we're looking to use [tems_set query](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-terms-set-query.html).
Like this:

```
GET applications/_search
{
  "query": {
    "terms_set": {
      "tags": {
        "terms": ["service-a", "team-1"],
        "minimum_should_match_script": {
          "source": "params.num_terms"
        }
      }
    }
  }
}
```


BTW:
- The license got changed lately but it works fine for us. See: https://www.elastic.co/pricing/faq/licensing
- Elasticseach can also be used to filter events by labels when we support Eventwatcher's event list in the future.

# Alternatives

## Idea1: Use a map (not practical)
To query documents that have all of the given tags from not an RDB, we can use a map instead of an array.

```go
tags := []string{
  "service-a",
  "team-1",
}
// Instead ->
tags := map[string]bool{
  "service-a": true,
  "team-a": true,
}
```

Then we can query like:

```
db.collection('Applications')
  .where('Tags.service-a', '==', true)
  .where('Tags.team-a', '==', true)
```

At first glance, this looks fine, but unfortunately, this is not practical because if we want to sort, we need composite indexes for all combinations.
We must always sort by UpdatedAt before paging, so this doesn't work well.

## Idea2: Make Firestore deprecated (worth considering)
If we can stop supporting Firestore as a datastore and only support RDBs, we can support the application tagging feature without any additional component.

In addition to the `applications` table, we can do that by creating the `tags` table, and the `applications_tags` table, which is a junction table.

```sql
SELECT applications.id
FROM applications
LEFT JOIN applications_tags ON applications_tags.application_id = applications.id
  AND applications_tags.tag_id = '?'
  AND applications_tags.tag_id = '?'
  AND applications_tags.tag_id = '?'
...
```
However, Firestore has already been Beta, hence there are more likely to many users who are using it in production environments.

Not just do we need to think about Firestore migration, but we also need to re-review our MySQL schema from scratch.

That being said, there is no reason why it has to be Firestore at this stage so it's worth considering.

# Unresolved questions
There are a couple of unresolved issues in using Elasticsearch.

## Needs to be the privilege mode
Elasticsearch doesn't work without increasing the limits on mmap counts: https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html

In order to set this up permanently, there is still an unresolved issue that requires the sysctl container to be run in privileged mode, so it's likely not to be usable by many users.
(Of cource, unavailable on GKE Autopilot as well because it's disallowed to change node-level settings): https://github.com/elastic/helm-charts/issues/1126#issuecomment-820189682

## Gets complicated
Currently we have two persistence methods, Datastore and Filestore.
This time we will add a new one, Tagstore (SearchEngine?), and need to be sure to query it when only manipulating tags.
It's a little bit too huge change for just adding tags, and it would cause two communications to external storage in one query. Possibly this can be solved with in-memory cache though.
