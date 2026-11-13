# Vapi API CLI Reference

Full command reference for `vapi`.

## Commands

- [`vapi analytics`](#vapi-analytics)
- [`vapi assistants`](#vapi-assistants)
- [`vapi board`](#vapi-board)
- [`vapi calls`](#vapi-calls)
- [`vapi campaigns`](#vapi-campaigns)
- [`vapi chats`](#vapi-chats)
- [`vapi eval`](#vapi-eval)
- [`vapi files`](#vapi-files)
- [`vapi insight`](#vapi-insight)
- [`vapi knowledge-bases-v2`](#vapi-knowledge-bases-v2)
- [`vapi observability-scorecard`](#vapi-observability-scorecard)
- [`vapi phone-numbers`](#vapi-phone-numbers)
- [`vapi provider-resources`](#vapi-provider-resources)
- [`vapi sessions`](#vapi-sessions)
- [`vapi simulation-personalities`](#vapi-simulation-personalities)
- [`vapi simulation-runs`](#vapi-simulation-runs)
- [`vapi simulation-scenarios`](#vapi-simulation-scenarios)
- [`vapi simulation-suites`](#vapi-simulation-suites)
- [`vapi simulations`](#vapi-simulations)
- [`vapi squads`](#vapi-squads)
- [`vapi structured-outputs`](#vapi-structured-outputs)
- [`vapi tools`](#vapi-tools)
- [`vapi traffic-allocations`](#vapi-traffic-allocations)

---

### `vapi analytics`

#### `vapi analytics get`

Runs one or more metric queries against call or subscription data using the requested time range, groupings, and aggregate operations.

`POST /analytics`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi assistants`

#### `vapi assistants create`

Creates a reusable assistant configuration containing the model, voice, transcriber, tools, prompts, and call behavior.

`POST /assistant`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi assistants delete`

Deletes the assistant identified by its ID.

`DELETE /assistant/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the assistant. |

#### `vapi assistants get`

Returns the assistant identified by its ID.

`GET /assistant/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the assistant. |

#### `vapi assistants list`

Returns assistants for the authenticated organization. Filter results by creation or update timestamps and limit the number returned.

`GET /assistant`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi assistants update`

Updates the specified fields of the assistant identified by its ID.

`PATCH /assistant/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the assistant. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi board`

#### `vapi board controller-create`

Create Board

`POST /reporting/board`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi board controller-find-all`

Get Boards

`GET /reporting/board`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi board controller-find-one`

Get Board

`GET /reporting/board/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi board controller-metrics-overview-ensure`

Get Default Metrics Overview Board

`GET /reporting/board/default/metrics-overview`

#### `vapi board controller-remove`

Delete Board

`DELETE /reporting/board/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi board controller-update`

Update Board

`PATCH /reporting/board/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi calls`

#### `vapi calls call-artifact-controller-assistant-recording-download`

Download Call Assistant Recording

`GET /call/{id}/assistant-recording`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls call-artifact-controller-call-logs-download`

Download Call Logs

`GET /call/{id}/call-logs`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls call-artifact-controller-customer-recording-download`

Download Call Customer Recording

`GET /call/{id}/customer-recording`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls call-artifact-controller-mono-recording-download`

Download Call Mono Recording

`GET /call/{id}/mono-recording`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls call-artifact-controller-pcap-download`

Download Call Packet Capture (pcap)

`GET /call/{id}/pcap`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls call-artifact-controller-stereo-recording-download`

Download Call Stereo Recording

`GET /call/{id}/stereo-recording`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls call-artifact-controller-video-recording-download`

Download Call Video Recording

`GET /call/{id}/video-recording`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | Call ID |

#### `vapi calls create`

Creates a call using an assistant or squad. The request can reference saved resources or include transient configurations.

`POST /call`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi calls delete`

Deletes the call identified by its ID.

`DELETE /call/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi calls get`

Returns the call identified by its ID, including its status, configuration, and available call data.

`GET /call/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the call. |

#### `vapi calls list`

Returns calls for the authenticated organization. Filter results by call ID, assistant ID, phone number ID, or creation and update timestamps.

`GET /call`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | This is the unique identifier for the call. |
| `--assistant-id` | `string` | No | This will return calls with the specified assistantId. |
| `--phone-number-id` | `string` | No | This is the phone number that will be used for the call. To use a transient number, use `phoneNumber` instead.

Only relevant for `outboundPhoneCall` and `inboundPhoneCall` type. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi calls update`

Updates the call identified by its ID.

`PATCH /call/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the call. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi campaigns`

#### `vapi campaigns campaign-controller-create`

Creates an outbound calling campaign that calls a set of customers.

`POST /campaign`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi campaigns campaign-controller-create-v2`

Create Campaign V2

`POST /v2/campaign`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi campaigns campaign-controller-find-all`

Returns outbound calling campaigns for the authenticated organization. Filter results by campaign ID, status, or creation and update timestamps.

`GET /campaign`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | Filters campaigns by ID. |
| `--status` | `scheduled | in-progress | ended | cancelled | archived` | No | Filters campaigns by status. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi campaigns campaign-controller-find-all-v2`

List Campaigns V2

`GET /v2/campaign`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No |  |
| `--status` | `scheduled | in-progress | ended | cancelled | archived` | No |  |
| `--include-counters` | `boolean` | No | When true, every campaign in the response includes `contactCounters` and
`callMetrics`. These are aggregate queries over contacts and events —
batched across the page, so the cost is three queries per request rather
than three per campaign, but still opt-in rather than paid for on every
read. Defaults to false. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi campaigns campaign-controller-find-one`

Returns the outbound calling campaign identified by its ID.

`GET /campaign/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the campaign. |

#### `vapi campaigns campaign-controller-find-one-v2`

Get Campaign V2

`GET /v2/campaign/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--include-counters` | `boolean` | No | When true, the response includes `contactCounters` and `callMetrics`.
These are aggregate queries over the campaign's contacts and events, so
they are opt-in rather than paid for on every read. Defaults to false. |

#### `vapi campaigns campaign-controller-get-campaign-v2-contacts`

Get Campaign V2 Contacts

`GET /v2/campaign/{id}/contacts`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--status` | `contact.pending | contact.dispatched | contact.completed | contact.failed | contact.skipped | contact.predial-failed[]` | No | This is the status to filter contacts by. Pass once or multiple times to
filter on any of the provided statuses. |
| `--limit` | `number` | No | This is the maximum number of contacts to return. Defaults to 50. |
| `--sort-by` | `position | status` | No | This is the column to sort by. Defaults to `position` — the order contacts
were uploaded, which is also dial order.

`status` sorts by the enum's declaration order rather than alphabetically,
which means it reads as a lifecycle: pending, dispatched, completed,
failed, skipped, predial-failed.

Only columns on `campaign_contact` are sortable. Call-level values such as
cost or duration live on the call and are attached after this query, so
sorting by them here would only reorder the current page. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |

#### `vapi campaigns campaign-controller-remove`

Deletes the outbound calling campaign identified by its ID.

`DELETE /campaign/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the campaign. |

#### `vapi campaigns campaign-controller-remove-v2`

Delete Campaign

`DELETE /v2/campaign/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi campaigns campaign-controller-update`

Updates the outbound calling campaign identified by its ID. Campaigns can be ended by updating their status to `ended`.

`PATCH /campaign/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the campaign. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi campaigns campaign-controller-update-v2`

Update Campaign

`PATCH /v2/campaign/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi chats`

#### `vapi chats create` `[BETA]`

Creates a new chat with optional SMS delivery via transport field. Requires at least one of: assistantId/assistant, sessionId, or previousChatId. Note: sessionId and previousChatId are mutually exclusive. Transport field enables SMS delivery with two modes: (1) New conversation - provide transport.phoneNumberId and transport.customer to create a new session, (2) Existing conversation - provide sessionId to use existing session data. Cannot specify both sessionId and transport fields together. The transport.useLLMGeneratedMessageForOutbound flag controls whether input is processed by LLM (true, default) or forwarded directly as SMS (false).

`POST /chat`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi chats create-response` `[BETA]`

Create Chat (OpenAI Compatible)

`POST /chat/responses`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi chats delete` `[BETA]`

Delete Chat

`DELETE /chat/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes |  |

#### `vapi chats get` `[BETA]`

Get Chat

`GET /chat/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes |  |

#### `vapi chats list` `[BETA]`

List Chats

`GET /chat`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | This is the unique identifier for the chat to filter by. |
| `--assistant-id` | `string` | No | This is the unique identifier for the assistant that will be used for the chat. |
| `--assistant-id-any` | `string` | No | Filter by multiple assistant IDs. Provide as comma-separated values. |
| `--squad-id` | `string` | No | This is the unique identifier for the squad that will be used for the chat. |
| `--session-id` | `string` | No | This is the unique identifier for the session that will be used for the chat. |
| `--previous-chat-id` | `string` | No | This is the unique identifier for the previous chat to filter by. |
| `--id-any` | `string` | No | Filter by multiple chat IDs. Provide as comma-separated values. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

---

### `vapi eval`

#### `vapi eval controller-create`

Creates a reusable eval that defines a mock conversation and checkpoints for evaluating assistant responses and tool calls.

`POST /eval`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi eval controller-get`

Returns the eval definition identified by its ID.

`GET /eval/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the eval definition. |

#### `vapi eval controller-get-paginated`

Returns eval definitions for the authenticated organization. Filter results by ID or creation and update timestamps.

`GET /eval`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | Filters eval definitions by ID. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi eval controller-get-run`

Returns the eval run identified by its ID.

`GET /eval/run/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the eval run. |

#### `vapi eval controller-get-runs-paginated`

Returns eval runs for the authenticated organization. Filter results by ID or creation and update timestamps.

`GET /eval/run`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--sort-by` | `createdAt` | No |  |
| `--search` | `string` | No | Literal, case-insensitive search across eval and assistant names. |
| `--id` | `string` | No | Filters eval runs by ID. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |

#### `vapi eval controller-remove`

Deletes the eval definition identified by its ID.

`DELETE /eval/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the eval definition. |

#### `vapi eval controller-remove-run`

Deletes the eval run identified by its ID.

`DELETE /eval/run/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the eval run. |

#### `vapi eval controller-run`

Runs a saved or transient eval against an assistant or squad and creates an eval-run record containing the results.

`POST /eval/run`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi eval controller-update`

Updates the eval definition identified by its ID.

`PATCH /eval/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the eval definition. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi files`

#### `vapi files create`

Uploads a file for use with a Vapi knowledge base.

`POST /file`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi files delete`

Deletes the uploaded file identified by its ID.

`DELETE /file/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the file. |

#### `vapi files get`

Returns the uploaded file identified by its ID.

`GET /file/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the file. |

#### `vapi files list`

Returns files uploaded to the authenticated organization.

`GET /file`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--purpose` | `assistant | composer-attachment | knowledge-base-v2` | No | Only return files with this purpose. When omitted, files of every purpose except composer attachments are returned. |

#### `vapi files update`

Updates the name of the uploaded file identified by its ID.

`PATCH /file/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the file. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi insight`

#### `vapi insight controller-create`

Creates a saved reporting insight that queries call data and presents the results as a bar chart, pie chart, line chart, or text value.

`POST /reporting/insight`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi insight controller-find-all`

Returns saved reporting insights for the authenticated organization. Filter results by ID or creation and update timestamps.

`GET /reporting/insight`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | Filters reporting insights by ID. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi insight controller-find-one`

Returns the reporting insight identified by its ID.

`GET /reporting/insight/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the reporting insight. |

#### `vapi insight controller-preview`

Runs an insight definition without first saving it, returning a preview of the resulting chart or text value.

`POST /reporting/insight/preview`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi insight controller-remove`

Deletes the reporting insight identified by its ID.

`DELETE /reporting/insight/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the reporting insight. |

#### `vapi insight controller-run`

Runs a saved reporting insight, optionally overriding its time range and response format.

`POST /reporting/insight/{id}/run`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the reporting insight. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi insight controller-update`

Updates the reporting insight identified by its ID.

`PATCH /reporting/insight/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the reporting insight. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi knowledge-bases-v2`

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-create`

Create a Knowledge Base V2

`POST /v2/knowledge-base`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-file-attach`

Attach and ingest a file

`POST /v2/knowledge-base/{id}/file`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-file-detach`

Detach a file

`DELETE /v2/knowledge-base/{id}/file/{fileId}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--file-id` | `string (uuid)` | Yes |  |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-file-retry`

Retry indexing a failed file

`POST /v2/knowledge-base/{id}/file/{fileId}/retry`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--file-id` | `string (uuid)` | Yes |  |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-files-get`

List files in a Knowledge Base V2

`GET /v2/knowledge-base/{id}/file`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-find-all`

List Knowledge Bases V2

`GET /v2/knowledge-base`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--limit` | `number` | No |  |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-find-one`

Get a Knowledge Base V2

`GET /v2/knowledge-base/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-remove`

Delete a Knowledge Base V2

`DELETE /v2/knowledge-base/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi knowledge-bases-v2 knowledge-base-v2-controller-update`

Update a Knowledge Base V2

`PATCH /v2/knowledge-base/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi observability-scorecard`

#### `vapi observability-scorecard scorecard-controller-create`

Creates a scorecard containing metrics, scoring conditions, and optional links to assistants whose calls should be evaluated.

`POST /observability/scorecard`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi observability-scorecard scorecard-controller-get`

Returns the scorecard identified by its ID.

`GET /observability/scorecard/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes | The unique identifier of the scorecard. |

#### `vapi observability-scorecard scorecard-controller-get-paginated`

Returns scorecards for the authenticated organization. Filter results by ID or creation and update timestamps.

`GET /observability/scorecard`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | Filters scorecards by ID. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi observability-scorecard scorecard-controller-remove`

Deletes the scorecard identified by its ID.

`DELETE /observability/scorecard/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes | The unique identifier of the scorecard. |

#### `vapi observability-scorecard scorecard-controller-update`

Updates the scorecard identified by its ID.

`PATCH /observability/scorecard/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes | The unique identifier of the scorecard. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi phone-numbers`

#### `vapi phone-numbers create`

Creates a Vapi phone number or imports a phone number from a supported provider, including Twilio, Vonage, Telnyx, or a bring-your-own provider.

`POST /phone-number`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi phone-numbers delete`

Deletes the phone number resource identified by its ID.

`DELETE /phone-number/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the phone number. |

#### `vapi phone-numbers get`

Returns the phone number resource identified by its ID.

`GET /phone-number/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the phone number. |

#### `vapi phone-numbers list`

Returns phone numbers for the authenticated organization. Filter results by creation or update timestamps and limit the number returned.

`GET /phone-number`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi phone-numbers phone-number-controller-find-all-paginated`

Returns a paginated list of phone numbers for the authenticated organization. Search by name, number, or SIP URI using a partial, case-insensitive match, and filter by creation or update timestamps.

`GET /v2/phone-number`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--search` | `string` | No | This will search phone numbers by name, number, or SIP URI (partial match, case-insensitive). |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi phone-numbers update`

Updates the specified fields of the phone number resource identified by its ID.

`PATCH /phone-number/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the phone number. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi provider-resources`

#### `vapi provider-resources provider-resource-controller-create-provider-resource`

Creates a pronunciation-dictionary resource for a supported provider, currently Cartesia or ElevenLabs.

`POST /provider/{provider}/{resourceName}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--content-type` | `string` | Yes |  |
| `--provider` | `cartesia | 11labs` | Yes | The provider (e.g., 11labs) |
| `--resource-name` | `pronunciation-dictionary` | Yes | The resource name (e.g., pronunciation-dictionary) |

#### `vapi provider-resources provider-resource-controller-delete-provider-resource`

Deletes the provider resource identified by its Vapi resource ID.

`DELETE /provider/{provider}/{resourceName}/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--provider` | `cartesia | 11labs` | Yes | The provider (e.g., 11labs) |
| `--resource-name` | `pronunciation-dictionary` | Yes | The resource name (e.g., pronunciation-dictionary) |
| `--id` | `string (uuid)` | Yes | The unique identifier of the provider resource. |

#### `vapi provider-resources provider-resource-controller-get-provider-resource`

Returns the provider resource identified by its Vapi resource ID.

`GET /provider/{provider}/{resourceName}/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--provider` | `cartesia | 11labs` | Yes | The provider (e.g., 11labs) |
| `--resource-name` | `pronunciation-dictionary` | Yes | The resource name (e.g., pronunciation-dictionary) |
| `--id` | `string (uuid)` | Yes | The unique identifier of the provider resource. |

#### `vapi provider-resources provider-resource-controller-get-provider-resources-paginated`

Returns a paginated list of provider resources for the authenticated organization. Filter pronunciation dictionaries by provider, resource ID, or creation and update timestamps.

`GET /provider/{provider}/{resourceName}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--provider` | `cartesia | 11labs` | Yes | The provider (e.g., 11labs) |
| `--resource-name` | `pronunciation-dictionary` | Yes | The resource name (e.g., pronunciation-dictionary) |
| `--id` | `string` | No | Filters provider resources by their resource ID. |
| `--resource-id` | `string` | No | Filters provider resources by their provider-specific resource ID. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi provider-resources provider-resource-controller-update-provider-resource`

Updates the provider resource identified by its Vapi resource ID.

`PATCH /provider/{provider}/{resourceName}/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--provider` | `cartesia | 11labs` | Yes | The provider (e.g., 11labs) |
| `--resource-name` | `pronunciation-dictionary` | Yes | The resource name (e.g., pronunciation-dictionary) |
| `--id` | `string (uuid)` | Yes | The unique identifier of the provider resource. |

---

### `vapi sessions`

#### `vapi sessions create` `[BETA]`

Create Session

`POST /session`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi sessions delete` `[BETA]`

Delete Session

`DELETE /session/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi sessions get` `[BETA]`

Get Session

`GET /session/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi sessions list` `[BETA]`

List Sessions

`GET /session`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | This is the unique identifier for the session to filter by. |
| `--name` | `string` | No | This is the name of the customer. This is just for your own reference.

For SIP inbound calls, this is extracted from the `From` SIP header with format `"Display Name" <sip:username@domain>`. |
| `--assistant-id` | `string` | No | This is the ID of the assistant to filter sessions by. |
| `--assistant-id-any` | `string` | No | Filter by multiple assistant IDs. Provide as comma-separated values. |
| `--squad-id` | `string` | No | This is the ID of the squad to filter sessions by. |
| `--workflow-id` | `string` | No | This is the ID of the workflow to filter sessions by. |
| `--number-e164-check-enabled` | `boolean` | No | This is the flag to toggle the E164 check for the `number` field. This is an advanced property which should be used if you know your use case requires it.

Use cases:
- `false`: To allow non-E164 numbers like `+001234567890`, `1234`, or `abc`. This is useful for dialing out to non-E164 numbers on your SIP trunks.
- `true` (default): To allow only E164 numbers like `+14155551234`. This is standard for PSTN calls.

If `false`, the `number` is still required to only contain alphanumeric characters (regex: `/^\+?[a-zA-Z0-9]+$/`).

@default true (E164 check is enabled) |
| `--extension` | `string` | No | This is the extension that will be dialed after the call is answered. |
| `--assistant-overrides` | `string` | No | These are the overrides for the assistant's settings and template variables specific to this customer.
This allows customization of the assistant's behavior for individual customers in batch calls. |
| `--squad-overrides` | `string` | No | These are the overrides applied when the call targets a `squadId`. Mirrors
the call-level `squadOverrides` — use this instead of `assistantOverrides`
when the campaign or call is squad-based. |
| `--number` | `string` | No | This is the number of the customer. |
| `--sip-uri` | `string` | No | This is the SIP URI of the customer. |
| `--email` | `string` | No | This is the email of the customer. |
| `--external-id` | `string` | No | This is the external ID of the customer. |
| `--customer-number-any` | `string` | No | Filter by any of the specified customer phone numbers (comma-separated). |
| `--id-any` | `string` | No | Filter by multiple session IDs. Provide as comma-separated values. |
| `--phone-number-id` | `string` | No | This will return sessions with the specified phoneNumberId. |
| `--phone-number-id-any` | `string[]` | No | This will return sessions with any of the specified phoneNumberIds. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi sessions update` `[BETA]`

Update Session

`PATCH /session/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi simulation-personalities`

#### `vapi simulation-personalities personality-controller-create`

Creates a personality, the AI tester's configuration used in simulations.

`POST /eval/simulation/personality`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-personalities personality-controller-find-all`

Returns the personalities for the authenticated organization.

`GET /eval/simulation/personality`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi simulation-personalities personality-controller-find-one`

Returns the specified personality.

`GET /eval/simulation/personality/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the personality. |

#### `vapi simulation-personalities personality-controller-remove`

Deletes the specified personality.

`DELETE /eval/simulation/personality/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the personality. |

#### `vapi simulation-personalities personality-controller-update`

Updates the specified personality. Changes apply to future runs.

`PATCH /eval/simulation/personality/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the personality. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi simulation-runs`

#### `vapi simulation-runs simulation-run-controller-cancel-group`

Cancels the specified simulation run.

`PATCH /eval/simulation/run/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation run. |

#### `vapi simulation-runs simulation-run-controller-cancel-item`

Cancels the specified run item.

`PATCH /eval/simulation/run/{id}/item/{itemId}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes | The unique identifier of the simulation run. |
| `--item-id` | `string` | Yes | The unique identifier of the run item. |

#### `vapi simulation-runs simulation-run-controller-create`

Starts a simulation run against a target assistant or squad.

`POST /eval/simulation/run`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--idempotency-key` | `string` | No | Makes the request safe to retry. A retry with the same key and the same body returns the run the first request created, with `reused: true`, instead of starting another. Reusing a key with a different body returns 409. Up to 255 printable ASCII characters, scoped to your organization; a key stays bound to its run while the run exists. |
| `--user-agent` | `string` | No | Identifies the client starting the simulation run |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-runs simulation-run-controller-find-all`

Returns the simulation runs for the authenticated organization.

`GET /eval/simulation/run`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--suite-id` | `string (uuid)[]` | No | Filter by simulation suite ids |
| `--status` | `queued | running | ended` | No | Filter by status |
| `--filter-status` | `passed | failed | running` | No | Filter by aggregate run result status |
| `--target-type` | `assistant | squad` | No | Filter by target type |
| `--target-id` | `string (uuid)` | No | Filter by target id |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi simulation-runs simulation-run-controller-find-item`

Returns the specified run item, including its evaluation results and the ID of the call that ran it.

`GET /eval/simulation/run/{id}/item/{itemId}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes | The unique identifier of the simulation run. |
| `--item-id` | `string` | Yes | The unique identifier of the run item. |

#### `vapi simulation-runs simulation-run-controller-find-items`

Returns the run items for the specified simulation run.

`GET /eval/simulation/run/{id}/item`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation run. |
| `--simulation-id` | `string (uuid)` | No | Filters run items to a specific simulation. |
| `--run-id` | `string (uuid)` | No | Filters run items to a specific run. |
| `--status` | `queued | running | evaluating | passed | failed | canceled` | No | Filters run items by status. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi simulation-runs simulation-run-controller-find-one`

Returns the specified simulation run, including its status and item counts.

`GET /eval/simulation/run/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation run. |

#### `vapi simulation-runs simulation-run-controller-generate-suggestions`

Generates AI suggestions for improving the assistant or squad's system prompt, tools, and scenarios, based on the specified run item.

`POST /eval/simulation/run/{id}/item/{itemId}/generate`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes | The unique identifier of the simulation run. |
| `--item-id` | `string` | Yes | The unique identifier of the run item. |
| `--force` | `string` | Yes | Set to the string `true` to regenerate improvement suggestions even if they already exist. |
| `--persist` | `string` | No |  |

---

### `vapi simulation-scenarios`

#### `vapi simulation-scenarios scenario-controller-create`

Creates a scenario, the AI tester's intent plus the success criteria that score a run.

`POST /eval/simulation/scenario`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-scenarios scenario-controller-find-all`

Returns the scenarios for the authenticated organization.

`GET /eval/simulation/scenario`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id-any` | `string[]` | No | Return only scenarios matching the provided ids |
| `--name` | `string` | No | Search by scenario name |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi simulation-scenarios scenario-controller-find-one`

Returns the specified scenario.

`GET /eval/simulation/scenario/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the scenario. |

#### `vapi simulation-scenarios scenario-controller-remove`

Deletes the specified scenario.

`DELETE /eval/simulation/scenario/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the scenario. |

#### `vapi simulation-scenarios scenario-controller-update`

Updates the specified scenario.

`PATCH /eval/simulation/scenario/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the scenario. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi simulation-suites`

#### `vapi simulation-suites simulation-onboarding-controller-requirements-draft`

Generate draft test requirements from an assistant's prompt, tools, and handoff configuration. The response includes numbered requirements describing expected behavior and open questions about behavior that needs clarification. Use this endpoint before generating a simulation suite.

Review the requirements and resolve any open questions. To revise a draft, send it as `previousDraft` together with natural-language `feedback` describing the changes or answers. Once the requirements are ready, pass them to `POST /eval/simulation/onboarding/suite`.

Requests can take up to 2 minutes. Each organization can draft requirements 60 times per UTC day, shared with the dashboard. A 429 `SIMULATION_ONBOARDING_QUOTA_EXCEEDED` response means the daily limit has been reached. Wait until the limit resets at 00:00 UTC before retrying.

`POST /eval/simulation/onboarding/requirements`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-suites simulation-onboarding-controller-suite-create`

Generate and save a simulation suite from reviewed test requirements, and assign it to an assistant. First draft requirements with `POST /eval/simulation/onboarding/requirements`. The response includes the saved suite, its simulations, and requirement coverage, including any gaps.

Requests can take up to 3 minutes. Set a client timeout of at least 200 seconds.

Use a unique `idempotencyKey` for each suite generation request. Retrying with the same key does not create a duplicate suite:
- On a timeout or a 5xx (502, 504, 524), retry with the same key. You get the saved suite, a 409 while it is still generating, or a fresh generation if the first one never saved.
- On a 409 `SIMULATION_ONBOARDING_GENERATION_IN_PROGRESS`, the first request is still generating. Wait `Retry-After` seconds (15) and retry with the same key.
- On a 409 `SIMULATION_ONBOARDING_IDEMPOTENCY_KEY_REUSED`, the key was used for a different assistant or requirements. Send a new key.

Review the generated simulations and coverage gaps. Revise the suite with `PATCH /eval/simulation/onboarding/suite/{suiteId}`, then run it with `POST /eval/simulation/run`.

Each organization can generate 20 suites per UTC day, shared with the dashboard, and 2 at a time. A retry that returns a saved suite does not count toward these limits. On a 429 `SIMULATION_ONBOARDING_CONCURRENCY_LIMITED`, wait `Retry-After` seconds (30) before retrying. On a 429 `SIMULATION_ONBOARDING_QUOTA_EXCEEDED`, wait until the daily limit resets at 00:00 UTC before retrying.

`POST /eval/simulation/onboarding/suite`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-suites simulation-onboarding-controller-suite-update`

Revise a suite created by `POST /eval/simulation/onboarding/suite`. Remove simulations, rewrite scenarios using natural-language instructions, or set `retryFailed` to regenerate scenarios that failed or exceeded generation limits. The response includes the updated suite and requirement coverage.

Requests can take up to 2 minutes.

Send a new `patchKey` for each edit. Replaying an applied key changes nothing and returns the current suite, with `failedEditSimulationIds` listing the simulations that edit could not rewrite, so retry a timeout or a 5xx with the same key. A suite remembers its last 50 keys, so never reuse one.

Each organization can make 60 edits per UTC day, shared with the dashboard. Requests that only remove simulations and replays of applied edits do not count toward this limit. On a 429 `SIMULATION_ONBOARDING_QUOTA_EXCEEDED`, wait until the daily limit resets at 00:00 UTC before retrying.

`PATCH /eval/simulation/onboarding/suite/{suiteId}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--suite-id` | `string` | Yes |  |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-suites simulation-suite-controller-create`

Creates a simulation suite, a group of simulations that run together.

`POST /eval/simulation/suite`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulation-suites simulation-suite-controller-duplicate`

Duplicate simulation suite

`POST /eval/simulation/suite/{id}/duplicate`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier for the resource. |

#### `vapi simulation-suites simulation-suite-controller-find-all`

Returns the simulation suites for the authenticated organization.

`GET /eval/simulation/suite`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--assistant-id` | `string (uuid)[]` | No | Filter by any assigned assistant ID. Supplying both assistant and squad filters requires both to match. |
| `--squad-id` | `string (uuid)[]` | No | Filter by any assigned squad ID. Supplying both assistant and squad filters requires both to match. |
| `--name` | `string` | No | Search by simulation suite name |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi simulation-suites simulation-suite-controller-find-one`

Returns the specified simulation suite.

`GET /eval/simulation/suite/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation suite. |

#### `vapi simulation-suites simulation-suite-controller-remove`

Deletes the specified simulation suite.

`DELETE /eval/simulation/suite/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation suite. |

#### `vapi simulation-suites simulation-suite-controller-update`

Updates the specified simulation suite.

`PATCH /eval/simulation/suite/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation suite. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi simulations`

#### `vapi simulations simulation-controller-create`

Creates a simulation by pairing a scenario with a personality.

`POST /eval/simulation`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi simulations simulation-controller-find-all`

Returns the simulations for the authenticated organization.

`GET /eval/simulation`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id-any` | `string[]` | No | Return only simulations matching the provided ids |
| `--standalone-only` | `boolean` | No | Only include simulations that are not part of a suite |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi simulations simulation-controller-find-one`

Returns the specified simulation.

`GET /eval/simulation/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation. |

#### `vapi simulations simulation-controller-get-concurrency`

Returns the organization's simulation concurrency limit, the number of active simulations, and how many more can start.

`GET /eval/simulation/concurrency`

#### `vapi simulations simulation-controller-remove`

Deletes the specified simulation.

`DELETE /eval/simulation/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation. |

#### `vapi simulations simulation-controller-update`

Updates the specified simulation.

`PATCH /eval/simulation/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the simulation. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi squads`

#### `vapi squads create`

Creates a squad that coordinates multiple assistants and their handoffs during a conversation.

`POST /squad`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi squads delete`

Deletes the squad identified by its ID.

`DELETE /squad/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the squad. |

#### `vapi squads get`

Returns the squad identified by its ID.

`GET /squad/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the squad. |

#### `vapi squads list`

Returns squads for the authenticated organization. Filter results by creation or update timestamps and limit the number returned.

`GET /squad`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id-any` | `string[]` | No | Return only squads matching the provided ids |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi squads update`

Updates the specified fields of the squad identified by its ID.

`PATCH /squad/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the squad. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi structured-outputs`

#### `vapi structured-outputs structured-output-controller-create`

Creates a reusable definition for extracting validated data from conversations using an AI model or regular expression.

`POST /structured-output`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi structured-outputs structured-output-controller-find-all`

Returns structured-output definitions for the authenticated organization. Filter results by ID, name, or creation and update timestamps.

`GET /structured-output`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | No | This will return structured outputs where the id matches the specified value. |
| `--name` | `string` | No | This will return structured outputs where the name matches the specified value. |
| `--page` | `number` | No | This is the page number to return. Defaults to 1. |
| `--sort-order` | `ASC | DESC` | No | This is the sort order for pagination. Defaults to 'DESC'. |
| `--sort-by` | `createdAt | duration | cost` | No | This is the column to sort by. Defaults to 'createdAt'. |
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi structured-outputs structured-output-controller-find-one`

Returns the structured-output definition identified by its ID.

`GET /structured-output/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the structured output. |

#### `vapi structured-outputs structured-output-controller-remove`

Deletes the structured-output definition identified by its ID.

`DELETE /structured-output/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the structured output. |

#### `vapi structured-outputs structured-output-controller-run`

Runs a saved or transient structured-output definition against one or more calls, optionally returning a preview without updating call artifacts.

`POST /structured-output/run`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi structured-outputs structured-output-controller-update`

Updates the structured-output definition identified by its ID.

`PATCH /structured-output/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the structured output. |
| `--schema-override` | `string` | Yes | Set to the string `true` to allow changing the schema's top-level type. Other values do not enable schema type changes. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi tools`

#### `vapi tools create`

Creates a reusable tool that assistants can invoke during conversations.

`POST /tool`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi tools delete`

Deletes the tool identified by its ID.

`DELETE /tool/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the tool. |

#### `vapi tools get`

Returns the tool identified by its ID.

`GET /tool/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the tool. |

#### `vapi tools list`

Returns reusable tools for the authenticated organization. Filter results by creation or update timestamps and limit the number returned.

`GET /tool`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--limit` | `number` | No | This is the maximum number of items to return. Defaults to 100. |
| `--created-at-gt` | `string (date-time)` | No | This will return items where the createdAt is greater than the specified value. |
| `--created-at-lt` | `string (date-time)` | No | This will return items where the createdAt is less than the specified value. |
| `--created-at-ge` | `string (date-time)` | No | This will return items where the createdAt is greater than or equal to the specified value. |
| `--created-at-le` | `string (date-time)` | No | This will return items where the createdAt is less than or equal to the specified value. |
| `--updated-at-gt` | `string (date-time)` | No | This will return items where the updatedAt is greater than the specified value. |
| `--updated-at-lt` | `string (date-time)` | No | This will return items where the updatedAt is less than the specified value. |
| `--updated-at-ge` | `string (date-time)` | No | This will return items where the updatedAt is greater than or equal to the specified value. |
| `--updated-at-le` | `string (date-time)` | No | This will return items where the updatedAt is less than or equal to the specified value. |

#### `vapi tools update`

Updates the specified fields of the tool identified by its ID.

`PATCH /tool/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string (uuid)` | Yes | The unique identifier of the tool. |
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

---

### `vapi traffic-allocations`

#### `vapi traffic-allocations traffic-allocation-controller-create`

Creates a new traffic allocation for an assistant, replacing the one currently in effect. To start or adjust a split, send targets naming published versions (such as "v7") with percentages totaling 100; allocationIntent is inferred as 'explicit'. To stop splitting and send every call to the newest published version, send allocationIntent 'follow-latest' with no targets field; stopping always names its intent, so a dropped targets field can never end a split by accident.

`POST /traffic-allocations`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--json` | `JSON` | Yes | Request body as JSON (or use individual body-field flags) |

#### `vapi traffic-allocations traffic-allocation-controller-find-all-paginated`

The append-only history of an assistant's allocations, newest first.

`GET /traffic-allocations`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--assistant-id` | `string` | No | Filter to allocations for this assistant. |
| `--page` | `integer` | No | The page number to return. Defaults to 1. |
| `--limit` | `integer` | No | The maximum number of items to return. Defaults to 100. |
| `--sort-order` | `ASC | DESC` | No | The sort order for pagination. Defaults to 'DESC'. |

#### `vapi traffic-allocations traffic-allocation-controller-find-one`

Returns a single allocation by id, including its targets and actor attribution.

`GET /traffic-allocations/{id}`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--id` | `string` | Yes |  |

#### `vapi traffic-allocations traffic-allocation-controller-latest-get`

The allocation currently in effect for the assistant, with its targets. The response carries no allocation field when traffic splitting has never been configured.

`GET /traffic-allocations/latest`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--assistant-id` | `string` | Yes | The assistant whose latest allocation to return. |

#### `vapi traffic-allocations traffic-allocation-controller-metrics-get`

Per-assistant-version call metrics over a time window: call counts, configured vs actual traffic share, duration, cost per minute, latency, ended reasons, a bucketed timeseries, and optional structured output aggregations. Calls with no assistant version are reported as an 'unversioned' row.

`GET /traffic-allocations/metrics`

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--assistant-id` | `string` | Yes | The assistant whose per-version traffic metrics to return. |
| `--from` | `string (date-time)` | Yes | Start of the metrics window (inclusive), ISO 8601 datetime. |
| `--to` | `string (date-time)` | Yes | End of the metrics window (exclusive), ISO 8601 datetime. Must be later than from, and the window may hold at most 2000 buckets of the requested width. |
| `--bucket` | `hour | day` | No | Bucket size for the timeseries. Defaults to 'day'. |
| `--structured-output-ids` | `string[]` | No | Structured output ids to aggregate per version, as an array or a comma-separated string. At most 5. |

---

## Global flags

These flags are available on every command:

| Flag | Description |
|------|-------------|
| `--dry-run` | Print the HTTP request without sending it |
| `--json <JSON\|->` | Supply the request body as JSON (or `-` for stdin) |
| `--params <JSON>` | Merge extra parameters as JSON |
| `--format <json\|table\|yaml\|csv>` | Output format (default: `json`) |
| `--output <PATH>` | Write binary responses to a file |
| `--base-url <URL>` | Override the API base URL |
| `--no-extract` | Print the full response body instead of the `x-fern-sdk-return-value` extraction |
| `--no-retry` | Disable retries declared by `x-fern-retries`, including network errors |
| `-q, --quiet` | Suppress stdout on success |
| `-h, --help` | Print help |
| `-V, --version` | Print version |

Operations the spec describes how to page (via `x-fern-pagination` or a root `page_token` parameter) also accept:

| Flag | Description |
|------|-------------|
| `--page-all` | Auto-paginate and stream all results |
| `--page-limit <N>` | Max pages to fetch (default: `10`) |
| `--page-delay <MS>` | Delay between page fetches in milliseconds (default: `100`) |
| `--no-pager` | Disable the pager even on interactive terminals |

Operations the spec marks as streaming (via `x-fern-streaming`) also accept:

| Flag | Description |
|------|-------------|
| `--no-stream` | Buffer the streaming response and print it as a single value once complete |
