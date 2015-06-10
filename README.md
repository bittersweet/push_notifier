# push_notifier [![Build Status](https://drone.io/github.com/bittersweet/push_notifier/status.png)](https://drone.io/github.com/bittersweet/push_notifier/latest)

## What is it?

This is a small Go app that listens to incoming
[`push`](https://developer.github.com/v3/activity/events/types/#pushevent)
webhook events from Github. On receiving this it will post a message to Slack
with a message and the content of the file pushed.

At [Springest](github.com/springest/) we use this to update everyone on new
files to a repo we use to share scripts and snippets with eachother.
This encourages inline discussion and sharing!

## Requirements

This app requires the following environment variables to be present:

* `HOOK_SECRET` to make sure only GH can use the endpoint
* `GITHUB_TOKEN` to download file contents
* `SLACK_TOKEN` personal authentication token
* `SLACK_CHANNEL` channel to send the notification to

## Setup

Go to your project on Github and add a webhook. The defaults are fine, so
`application/json` Content-Type plus the `push` event.

### `HOOK_SECRET`

We only need to set up the URL (will know that later once deployed to Heroku)
plus a shared secret for Github and our app, to verify that messages can only
come from GH.

Something like the following is
[advised](https://developer.github.com/webhooks/securing/)
to generate a secret token:

`ruby -rsecurerandom -e 'puts SecureRandom.hex(20)'`

### `GITHUB_TOKEN`

A [personal access token](https://github.com/settings/tokens) so we have access
to the files commited, to display these inline in Slack. Only the `repo` scope
is needed for this token.

### `SLACK_TOKEN`

We need a token to be able to use the `postMessage` in Slack, you can set that
up [here](https://api.slack.com/web).

### `SLACK_CHANNEL`

The id of the channel. Go to
[https://api.slack.com/methods/channels.list/test](https://api.slack.com/methods/channels.list/test)
to get a list of channels via the Slack API Tester.
