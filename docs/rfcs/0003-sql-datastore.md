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
# PostgreSQL
CREATE TABLE projects (
  id UUID PRIMARY KEY,
  data JSONB NOT NULL,
  disabled BOOL NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL
);
CREATE TABLE applications (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  data JSONB NOT NULL,
  disabled BOOL NOT NULL,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL
);

# MySQL
CREATE TABLE `projects` (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  disabled BOOL NOT NULL,
  created_at INT(11) NOT NULL,
  updated_at INT(11) NOT NULL
) ENGINE=InnoDB;

CREATE TABLE `applications` (
  id BINARY(16) PRIMARY KEY,
  project_id BINARY(16) NOT NULL,
  data JSON NOT NULL,
  disabled BOOL NOT NULL,
  created_at INT(11) NOT NULL,
  updated_at INT(11) NOT NULL
) ENGINE=InnoDB;
```

For indexing issue:
- PostgreSQL indexing for the attributes of Jsonb data type is as follow: [Json indexing](https://www.postgresql.org/docs/current/datatype-json.html#JSON-INDEXING)
- MySQL also support indexing for attributes of Json data type as follow: [Json indexing by secondary indexes](https://dev.mysql.com/doc/refman/8.0/en/create-table-secondary-indexes.html)

(note: both Jsonb and Json point to the type of storing JSON data in those databases. The key difference is that Json data is stored as an exact copy of the JSON input text, whereas Jsonb stores data in a decomposed binary form. Jsonb seems a better choice since it supports indexing with significant performance increases)

# Alternatives

Currently, we consider between MySQL and PostgreSQL for those support for NoSQL features, PostgreSQL has a longer time in this field. Some considering factors:
- ability to index a specific attribute of the JSON field
- supported operators on specific JSON field
- performance of query operation on indexed JSON field
- able to keep the advantage of schemaless pattern

## Ability to index a specific attribute of the JSON field

For `MySQL 8.0`, we have 2 points which have to be considered.

1. UUID data for `id` and `project_id` fields

Since we're using `application side UUID generate` pattern, (which mean to MySQL, those ids are true random UUID) store those ids under `VARCHAR(32)` data type is costly for both read and write operation due to the indexes does not work, using `UUID_TO_BIN` & `BIN_TO_UUID` with `swap_flag` and store data under `BINARY(16)` would help.

From MySQL docs
```
  o If swap_flag is 0, the two-argument form is equivalent to the
    one-argument form. The binary result is in the same order as the
    string argument.

  o If swap_flag is 1, the format of the return value differs: The
    time-low and time-high parts (the first and third groups of
    hexadecimal digits, respectively) are swapped. This moves the more
    rapidly varying part to the right and can improve indexing
    efficiency if the result is stored in an indexed column.

```

2. Create indexes for JSON attribute without affecting the schemaless advantage

We could use `CREATE INDEX` to create secondary indexes on generated columns which stored JSON attributes as an indirect way to indexes JSON attributes.
In order to keep the advantage of the schemaless pattern, we will use `virtual generated columns` instead of `stored generated columns` (which will physically store along with other columns of table). __The virtual generated columns wouldn't be generated on READ as long as we keep all generated columns as part of some secondary indexes__, which reduce the cost of recomputing from READ operation (note that computing virtual columns value cost on WRITE remains).

Sample secondary indexes creation as follow

```sql
mysql> CREATE INDEX idx ON applications ( ( CAST( data->>"$.name" AS CHAR(10) ) ) );
# OR
mysql> CREATE INDEX idy ON applications ( (JSON_VALUE(data, '$.name' RETURNING CHAR(10))) );
```

ref:
- https://dev.mysql.com/doc/refman/8.0/en/create-table-secondary-indexes.html
- https://dev.mysql.com/doc/refman/8.0/en/create-index.html

## Supported operators on specific JSON field

For our use-case, we plan to only focus on search function across attributes of JSON data and will always get back full JSON column data instead of just part (JSON object which contains only necessary keys or raw values) of it. Both MySQL and PostgreSQL search functions work with indexed attributes of JSON by using `->>` or `->` operator to specific attribute for `where` condition.

In case of using MySQL, though it had been noted on [the docs](https://dev.mysql.com/doc/refman/8.0/en/json-search-functions.html#function_json-value) that `JSON_VALUE()` equal to `CAST(JSON_UNQUOTE(JSON_EXTRACT(json_doc, path)))` or `CAST(json_doc->>path)`, we have to use exactly `where condition` which matches the `index expression` in order to make indexes work.

With the above defined `applications` table
```sql
mysql> SELECT data FROM applications;
+---------------------------------------------------------------+
| data                                                          |
+---------------------------------------------------------------+
| {"attr": "value_1", "name": "app-1", "tags": ["test"]}        |
| {"attr": "value_2", "name": "app-2", "tags": ["app"]}         |
| {"attr": "value_3", "name": "app-3", "tags": ["app", "test"]} |
+---------------------------------------------------------------+
mysql> CREATE INDEX idx ON applications ((CAST(data->>"$.name" AS CHAR(10))));
mysql> CREATE INDEX idy ON applications ( (JSON_VALUE(data, '$.name' RETURNING CHAR(10))) );
...
mysql> EXPLAIN SELECT data->>"$.name" FROM applications WHERE CAST(data->>'$.name' AS CHAR(10)) = 'app-1';
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
| id | select_type | table        | partitions | type | possible_keys | key  | key_len | ref   | rows | filtered | Extra |
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
|  1 | SIMPLE      | applications | NULL       | ref  | idx           | idx  | 13      | const |    1 |   100.00 | NULL  |
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
1 row in set, 1 warning (0.00 sec)

mysql> EXPLAIN SELECT data->>"$.name" FROM applications WHERE JSON_VALUE(data, '$.name' RETURNING CHAR(10)) = 'app-1';
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
| id | select_type | table        | partitions | type | possible_keys | key  | key_len | ref   | rows | filtered | Extra |
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
|  1 | SIMPLE      | applications | NULL       | ref  | idy           | idy  | 43      | const |    1 |   100.00 | NULL  |
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
1 row in set, 1 warning (0.00 sec)
```

note: indexing using `JSON_VALUE` costs more than `CAST` (key_len value is longer) and may cost more disk usage.

Another issue is that `LIKE` operator could not use the indexes for query on JSON attributes.

```sql
mysql> EXPLAIN SELECT data->>"$.name" FROM applications FORCE INDEX (idx) WHERE CAST(data->>'$.name' AS CHAR(10)) = 'app-1';
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
| id | select_type | table        | partitions | type | possible_keys | key  | key_len | ref   | rows | filtered | Extra |
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
|  1 | SIMPLE      | applications | NULL       | ref  | idx           | idx  | 13      | const |    1 |   100.00 | NULL  |
+----+-------------+--------------+------------+------+---------------+------+---------+-------+------+----------+-------+
1 row in set, 1 warning (0.00 sec)

mysql> EXPLAIN SELECT data->>"$.name" FROM applications FORCE INDEX (idx) WHERE CAST(data->>'$.name' AS CHAR(10)) LIKE 'app-%';
+----+-------------+--------------+------------+------+---------------+------+---------+------+------+----------+-------------+
| id | select_type | table        | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra       |
+----+-------------+--------------+------------+------+---------------+------+---------+------+------+----------+-------------+
|  1 | SIMPLE      | applications | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    6 |   100.00 | Using where |
+----+-------------+--------------+------------+------+---------------+------+---------+------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)

mysql> EXPLAIN SELECT data->>"$.name" FROM applications FORCE INDEX (idy) WHERE JSON_VALUE(data, '$.name' RETURNING CHAR(10)) LIKE 'app-%';
+----+-------------+--------------+------------+------+---------------+------+---------+------+------+----------+-------------+
| id | select_type | table        | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra       |
+----+-------------+--------------+------------+------+---------------+------+---------+------+------+----------+-------------+
|  1 | SIMPLE      | applications | NULL       | ALL  | NULL          | NULL | NULL    | NULL |    6 |   100.00 | Using where |
+----+-------------+--------------+------------+------+---------------+------+---------+------+------+----------+-------------+
1 row in set, 1 warning (0.00 sec)
```

## Performance of query operation on indexed JSON field

For queries which uses search function on indexed JSON fields and without using JOIN (in our use-case)

- For read queries, MySQL has a bit advantage due to its fast read-only characteristic. Besides, in case all virtual generated columns are secondary indexed columns, generated column values are materialized in the records of the index, which means MySQL will not recalculate virtual generated columns on query.
- For write queries, PostgreSQL has a bit advantage due to MySQL cost on calculating virtual generated columns on each writes.

ref: https://dev.mysql.com/doc/refman/8.0/en/create-table-secondary-indexes.html
## Able to keep advantage of schemaless pattern

Yes, for both 🎉

# Unresolved questions

Currently, we have 2 points which need to investigate more
1. The support of each cloud providers (GCP and AWS for now) for the chosen SQL database, since it's not the native service of those cloud providers (not as firestore of GCP and dynamodb of AWS for instance). Both AWS and GCP support latest versions of PostgreSQL(v13) and MySQL(v8).
2. In case `PostgreSQL` is chosen, there is no official driver for golang currently, we have a list of candidates: [pgx](https://github.com/jackc/pgx), [go-pg/pg](https://github.com/go-pg/pg).
