# Ingestion

เอกสารนี้อธิบายช่องทางรับ log ของ Logdesk

Backend ใช้ normalization pipeline เดียวกันสำหรับ HTTP JSON, File batch และ Syslog เพื่อให้ log จากหลาย source ถูกแปลงเป็น field กลางก่อนบันทึกลง PostgreSQL

## ช่องทางรับ log

| ช่องทาง | Protocol / Format | Tenant |
|---|---|---|
| HTTP JSON | `POST /ingest` หรือ `/api/ingest` | จาก `X-API-Key` หรือ Admin session |
| File batch | JSON file upload/import | จาก Admin session หรือ API key |
| Syslog TCP | newline framed TCP port `5514` | จาก `SYSLOG_TENANT` |
| Syslog UDP | UDP port `5514` | จาก `SYSLOG_TENANT` |

Syslog ไม่มี authentication ในตัว จึงไม่เชื่อ tenant ที่มากับข้อความจากอุปกรณ์ และใช้ค่า `SYSLOG_TENANT` จาก config แทน

## Source ที่ normalize ได้

ตัวอย่าง source ที่รองรับใน demo

- firewall
- network
- api
- crowdstrike
- aws
- m365
- ad
- nginx
- syslog

## Field กลาง

Normalizer พยายาม map ข้อมูลเข้า field กลาง เช่น

- `tenant`
- `source`
- `event_type`
- `severity`
- `src_ip`
- `username`
- `host`
- `action`
- `timestamp`
- `raw`
- `fields`

ถ้า log ไม่มี timestamp ระบบใช้เวลา ingest ปัจจุบันแทน

## HTTP JSON example

```powershell
py samples/post_logs.py --file samples/tenant-a/events.json
```

หรือใช้ Postman collection ที่ `docs/postman_collection.json`

## Syslog example

UDP

```powershell
py samples/send_syslog.py --host localhost --port 5514
```

TCP

```powershell
py samples/send_syslog.py --tcp --host localhost --port 5514
```

บน SaaS demo ให้เปลี่ยน host เป็น

```text
logdisk.malaysiawest.cloudapp.azure.com
```

## Alert test

การ trigger alert ใช้ abnormal sample ที่มี log field ตรง rule ครบ **5 ครั้งภายใน 5 นาที** ภายใน tenant เดียวกัน

```powershell
py samples/post_logs.py --file samples/tenant-a/abnormal.json
```

จากนั้นเปิดหน้า Alerts หรือเรียก `GET /api/alerts`

## ข้อจำกัด

- UDP เป็น best effort และไม่มี acknowledgement
- Syslog ยังไม่มี TLS และไม่มี authentication
- ยังไม่มี durable queue, retry หรือ deduplication
- ถ้าส่ง log ซ้ำ ระบบจะนับเป็น event ใหม่
