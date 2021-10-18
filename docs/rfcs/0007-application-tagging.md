- Start Date: 2021-10-18
- Target Version: 0.21.0

# Summary
This PR proposes adding a new attribute named Tags to applications to allow more flexible filtering.

# Motivation
Currently, it is able to filter with the embedded attributes `Environment` and `Kind`.
It works fine for a relatively tiny project with a few developers, but for a huge project like with so many microservices and multiple teams sharing the responsibility, it gets difficult to find.
It makes easier to find out applications they'd like that it allows their own attributes.

# Detailed design
There are primarily two possible filtering methods: labels in the form of `key=value` and tags, but tags are simpler for both the user and the design.

```proto
message Application {
    ...
    repeated string tags = 15;
}
```

Letting a Deployment have `application_tags` will bring us to filter deployment they want.

```proto
message Deployment {
    ...
    repeated string application_tags = 15;
}
```

It works well to embed comma-separated tags in the URL to share.

```
https://control-plane.dev/applications?tags=payment,teamA
```

# Alternatives
Using Labels.

```proto
message Application {
    ...
    map<string,string> labels = 33;
}
```

Query parameter:

```
https://control-plane.dev/applications?tags=service:payment,team:A
```

The web UI gets also slightly more complex.
The labels are over-specified, as this feature is only for users to distinguish, not for the system to filter.

# Unresolved questions

