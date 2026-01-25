# Firestore index ensurer developers guide
This package automatically creates the needed composite indexes for Google Firestore.
It's based on well-defined JSON file named `indexes.json`, which is obtained by running the following, with the `__name__` field removed:

```
gcloud firestore indexes composite list --format=json --sort-by=name
```

For details for this command: https://cloud.google.com/sdk/gcloud/reference/firestore/indexes/composite/list
