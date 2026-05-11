# clonernews Audit and Requirements Code Map

This file answers each audit question and requirement with the file name, a small critical snippet, and a short explanation of how it is done.

## Audit Questions

### Functional

#### Does this post open without any errors?
File name: src/features/feed/feed.view.js, src/main.js, src/features/post-detail/post-detail.controller.js, src/features/post-detail/post-detail.view.js

Critical snippet:
```js
href: toDetailHref(item.id)
handleRoute(route)
load(id)
renderSanitizedText(textBody, viewModel.text)
```

How it is done: feed cards route into the item page, the shell mounts the detail view, and the detail view renders story and job content safely.

#### Does this post open without any errors?
File name: src/features/feed/feed.view.js, src/main.js, src/features/post-detail/post-detail.controller.js, src/features/post-detail/post-detail.view.js

Critical snippet:
```js
href: toDetailHref(item.id)
handleRoute(route)
load(id)
renderSanitizedText(textBody, viewModel.text)
```

How it is done: jobs use the same feed-to-detail route path as stories, so the same shell and view code handles them.

#### Does this post open without any errors?
File name: src/features/post-detail/post-detail.controller.js, src/features/polls/poll.view.js

Critical snippet:
```js
if (item.type === 'poll')
createPollView()
render(viewModel)
```

How it is done: poll items are detected in the controller and passed to the poll renderer, which draws the poll card and options.

#### Did the posts load without error and without spamming the user?
File name: src/features/feed/feed.view.js, src/features/feed/feed.controller.js

Critical snippet:
```js
handleSentinelIntersect(entries)
onStateChange(items, isLoading, error, hasMore)
if (isLoading || !hasMorePages)
```

How it is done: the feed uses an IntersectionObserver sentinel to load more only when the user reaches the bottom, and the controller blocks overlapping fetches.

#### Are the comments being displayed in the correct order (from newest to oldest)?
File name: src/core/use-cases/get-item.js, src/features/comments/comments-tree.view.js

Critical snippet:
```js
const depthLimit = Math.min(maxDepth, MAX_COMMENT_TREE_DEPTH)
const sortNodesNewestFirst = (nodes)
childListFragment.append(createCommentNode(...))
```

How it is done: the tree builder recurses with a depth cap, and the view sorts each level newest-first before appending children.

### General

#### Does the UI have at least stories, jobs and polls?
File name: src/features/feed/feed.view.js, src/features/polls/poll.view.js, src/features/post-detail/post-detail.controller.js

Critical snippet:
```js
FEED_TYPES = Object.freeze([...])
const createPollView = () => {
if (item.type === 'poll')
```

How it is done: the feed exposes the required tabs and the detail layer routes poll items into the poll renderer.

#### Are the posts displayed in the correct order (from newest to oldest)?
File name: src/core/use-cases/list-items.js, src/features/feed/feed.view.js

Critical snippet:
```js
const compareNewestFirst = (left, right)
const sortedItems = collectedItems.toSorted(compareNewestFirst)
```

How it is done: items are sorted before rendering, so the feed stays newest-first.

#### Does each comment present the right parent post?
File name: src/core/use-cases/get-item.js, src/features/comments/comments-tree.view.js

Critical snippet:
```js
const commentWithParent = { ...commentResult.data, parent: parentId }
data-testid: COMMENTS_TREE_TEST_IDS.parent
```

How it is done: the recursive fetch attaches the parent ID to each nested comment and the view prints that parent reference.

#### Does the UI notify the user when a certain post is updated?
File name: src/core/use-cases/poll-updates.js, src/features/live-banner/live-banner.controller.js, src/features/live-banner/live-banner.view.js, src/main.js, src/features/feed/feed.view.js

Critical snippet:
```js
const createPayload = (...)
const poll = async () => {
render(state)
role: 'status'
refreshActiveTab(updatedItemIds)
toItemsWithRefreshPriority(items, updatedItemIds)
```

How it is done: the updates use case diffs the live-data endpoint, the banner announces update counts through an accessible live region, and pressing Refresh reloads New feed with changed items pushed to the top and marked as new or updated on the cards.

#### Is the project using throttling to regulate the number of requests (every 5 seconds)?
File name: src/core/use-cases/poll-updates.js, src/infra/throttle.js

Critical snippet:
```js
export const MIN_POLL_INTERVAL_MS = 5000
if (currentTimeMs - lastPollAtMs < pollIntervalMs)
export const throttle = (fn, ms) => {
```

How it is done: the live-data poller refuses to run before five seconds pass, and the shared throttle utility provides the same timing guard.

### Bonus

#### Does the UI have more types of posts than stories, jobs and polls?
File name: src/features/feed/feed.view.js, src/main.js

Critical snippet:
```js
FEED_TYPES = Object.freeze([... 'ask', 'show'])
bindHeaderFeedTabs()
```

How it is done: the feed also exposes Ask HN and Show HN tabs, so the UI supports more than the minimum post types.

#### Have sub-comments (nested comments) been implemented?
File name: src/core/use-cases/get-item.js, src/features/comments/comments-tree.view.js

Critical snippet:
```js
buildCommentForest(...)
createCommentNode(...)
childListFragment.append(...)
```

How it is done: the comment tree is built recursively, so each comment can render its children underneath it.

## Requirements Questions

#### Handle stories, jobs, and polls
File name: src/core/use-cases/list-items.js, src/features/feed/feed.view.js, src/features/post-detail/post-detail.controller.js, src/features/polls/poll.view.js

Critical snippet:
```js
const FEED_TYPES = Object.freeze(['top', 'new', 'ask', 'show', 'job'])
const createPollView = () => {
```

How it is done: the feed accepts all required feed types and the detail layer routes poll items into the poll renderer.

#### Render comments with the proper post parent
File name: src/core/use-cases/get-item.js, src/features/comments/comments-tree.view.js

Critical snippet:
```js
parent: parentId
data-testid: COMMENTS_TREE_TEST_IDS.parent
```

How it is done: each nested comment receives the parent ID during tree building and the view displays that parent value.

#### Order posts and comments newest to oldest
File name: src/core/use-cases/list-items.js, src/features/comments/comments-tree.view.js

Critical snippet:
```js
toSorted(compareNewestFirst)
sortNodesNewestFirst(nodes)
```

How it is done: feed items are sorted before render and comment children are sorted at every depth.

#### Load posts only when users need them
File name: src/features/feed/feed.view.js, src/features/feed/feed.controller.js

Critical snippet:
```js
IntersectionObserver
if (isLoading || !hasMorePages)
```

How it is done: the sentinel triggers page loading lazily and the controller blocks duplicate overlap.

#### Present the newest information and update every 5 seconds or more
File name: src/core/use-cases/poll-updates.js, src/features/live-banner/live-banner.view.js

Critical snippet:
```js
MIN_POLL_INTERVAL_MS = 5000
role: 'status'
aria-live: 'polite'
```

How it is done: the poller has a hard 5-second floor and the banner announces updates as a live region.

#### Avoid overloading the API
File name: src/infra/hn-api-adapter.js, src/infra/cache-adapter.js, src/core/use-cases/list-items.js, src/core/use-cases/get-item.js

Critical snippet:
```js
FETCH_OPTIONS
validateItemPayload(value)
createCacheAdapter()
ITEM_FETCH_CONCURRENCY_CAP = 6
```

How it is done: requests use strict fetch settings, cached reads avoid repeat fetches, and batch fetches are capped at six concurrent calls.

#### Keep routing stable on static hosting
File name: src/shared/router.js, src/main.js

Critical snippet:
```js
createHashRouter()
handleRoute(route)
clearCurrentRoute()
```

How it is done: the app uses hash routing and always unmounts the previous route before mounting the next one.

#### Sanitize API HTML before rendering
File name: src/features/post-detail/post-detail.view.js, src/features/comments/comments-tree.view.js, src/features/polls/poll.view.js, src/shared/dom-helpers.js

Critical snippet:
```js
DOMPurify.sanitize(...)
createElement(...)
textContent
```

How it is done: all rich text goes through DOMPurify, and plain text uses safe DOM sinks.

#### Keep relative time labels consistent
File name: src/shared/time-format.js, src/features/feed/feed.view.js, src/features/post-detail/post-detail.view.js

Critical snippet:
```js
formatRelativeTime(value)
```

How it is done: one shared Temporal-based formatter is used everywhere the UI shows a relative timestamp.

#### Keep the live UI reactive without a framework
File name: src/shared/signals.js, src/features/live-banner/live-banner.controller.js, src/main.js

Critical snippet:
```js
createSignal(initialValue)
subscribe(callback)
pollUpdates
```

How it is done: the update banner listens to a small signal primitive instead of a framework store.

## Short Read

- Use the file name and snippet rows when you want the shortest proof for a question.
- Core and infra files provide the main guarantees: ordering, recursion limits, throttling, caching, and request caps.
- Feature files mostly mount those guarantees and keep rendering safe.
