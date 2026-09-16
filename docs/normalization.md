# Normalized fields

Newly ingested events expose canonical fields in `fields` (PostgreSQL JSONB). Core fields also remain in their existing indexed columns and API properties. Original names and additional provider fields remain available; `raw` preserves the input. No database migration or historical rewrite is performed.

| Canonical field | Accepted aliases / conversion |
|---|---|
| vendor, product | Vendor, Product; defaults for AWS, M365, AD, CrowdStrike when absent |
| dst_ip | dst, destinationIPAddress, DestinationIp; validate and format IP |
| src_port | spt, sport, sourcePort, IpPort; integer 0–65535 |
| dst_port | dpt, dport, destinationPort, DestinationPort; integer 0–65535 |
| protocol | proto, Protocol; lowercase; 6→tcp, 17→udp, 1→icmp, 58→icmpv6; HTTP/1.1→http/1.1 |
| url | request_uri, RequestUri; preserve text |
| http_method | method, HttpMethod; uppercase |
| status_code | http_status, statusCode; integer 100–599 |
| outcome | success, failure, unknown; lowercase |
| cloud.account_id | nested cloud value, dotted key, recipientAccountId (AWS) |
| cloud.region | nested cloud value, dotted key, awsRegion (AWS) |
| cloud.service | nested cloud value, dotted key, eventSource (AWS); AWS suffix .amazonaws.com removed |

Canonical values take precedence over aliases. Optional null/empty values fall back to aliases or are omitted. Invalid provided values reject the event; HTTP batches are validated before insertion so a malformed event rejects the whole batch. Account IDs must be strings to preserve leading zeroes. Additional keys inside `cloud` are preserved.

Outcome uses an explicit valid `outcome` first. Otherwise it reads textual `ResultStatus`/`status`; failed-login event types and provider `errorCode` imply failure. Known successful login types imply success when no provider result is available. Finally HTTP 200–399 means success, 400–599 failure, and informational or unclassified events remain unknown. A firewall deny is not automatically a failed action, and a malware detection is not automatically a failed login. Missing login action becomes `login`; provided actions are lowercased.

Search now includes normalized JSONB fields in addition to the original payload and existing core values. The substring query does not use the JSONB GIN index; this remains suitable for a small demo, not a claim of full-text indexing. Details display available canonical fields above the complete normalized JSON.

Historical events retain their original fields; send new HTTP, Syslog, AWS, M365 or AD samples to see the expanded schema.
