- Start Date: 2020-02-22
- Target Version: 1.0.0

# Summary

This RFC proposes adding SQL database as datastore for PipeCD Control-plane. Currently, we're considering between PostgreSQL and MySQL for those NoSQL support features.

# Motivation

The background of this proposal is we are going to add complex queries and more indexes to the control-plane datastore for some new features of PipeCD, the current in-use NoSQL datastore is blocking those requirements.

# Detailed design

As mention is the motivation section, the root cause of this change is: we're going to add complex queries with those indexes to the datastore of PipeCD control-plane. The use-case of those queries requires:
- easy way to pagination query result
- easy way to ordering query result
- easy way to add indexes for future adding queries

Besides, to keep `simplicity on installing` characteristic of PipeCD, the chosen must be:
- easy to be installed on cloud provider environment (GCP, AWS, etc)
- can use the updated version easily

We're using NoSQL for now due to its schemaless characteristic is good for the deployment process, so we focus on NoSQL support features on the chosen SQL database.

The sample table creates commands for the `project` model and `application` model (other models than project model is as same as the application) as below.

```sql
CREATE TABLE project (
	id VARCHAR ( 32 ) PRIMARY KEY,
	data JSON NOT NULL,
	created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE application (
	id VARCHAR ( 32 ) PRIMARY KEY,
	project_id VARCHAR ( 32 ) NOT NULL,
	data JSON NOT NULL,
	created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

For indexing issue:
- PostgreSQL indexing for the attributes of Jsonb data type is as follow: https://www.postgresql.org/docs/current/datatype-json.html#JSON-INDEXING
- MySQL also support indexing for attributes of Json data type as follow: https://dev.mysql.com/doc/refman/8.0/en/create-table-secondary-indexes.html

(note: both Jsonb and Json point to the type of storing Json data in those databases. The key difference is that JSON data is stored as an exact copy of the JSON input text, whereas jsonb stores data in a decomposed binary form. Jsonb seems a better choice since it reduces the cost of encode/decode on read/write operator)

# Alternatives

Currently, we consider between MySQL and PostgreSQL for those support for NoSQL features, PostgreSQL has a longer time in this field.

# Unresolved questions

Currently, we have 2 points which need to investigate more
1. The support of each cloud providers (GCP and AWS for now) for the chosen SQL database, since it's not the native service of those cloud providers (not as firestore of GCP and dynamodb of AWS for instance).
- AWS is going to support latest MySQL 8.0 in a near future (not currently).
- AWS is supporting PostgreSQL 12 for now (the latest version is 13).
- GCP is supporting PostgreSQL 13 for now (the latest version).
2. In case `PostgreSQL` is chosen, there is no official driver for golang currently, we have a list of candidates: [lib/pg](https://github.com/lib/pq), [go-pg/pg](https://github.com/go-pg/pg).
