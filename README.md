# mothra

A Slackbot that answers questions about [BugZilla](https://www.bugzilla.org/) issues.

## Configuration

## Slackbot Usage

Start `mothra` in Slackbot-mode:

```sh
$ mothra serve -p 8080
TODO(njhale): add output
```
TODO(njhale): add slack UX

## Command Line Usage

`mothra` comes with a command line interface that allows you to issue queries against the [configured](#configuration) BugZilla instance.

```sh
$ mothra queries
# show all available queries
$ mothra get [query]
# get a query
```

