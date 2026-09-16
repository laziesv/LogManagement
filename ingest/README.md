# Ingestion

collector สำหรับ Syslog เขียนด้วย Go อยู่ที่ `backend/internal/collector/syslog.go` และถูก wire จาก `backend/internal/bootstrap/run.go` โดยเรียกใช้ normalizer ชุดเดียวกับ Fiber API (`backend/internal/normalize/normalize.go`)

collector รันอยู่ใน backend container เพื่อให้ package แบบ appliance ง่ายขึ้น

## ช่องทางรับ log

ระบบรองรับช่องทางหลักเหล่านี้:

- HTTP JSON POST
- multipart JSON file
- Syslog UDP/TCP
- Python simulator ใน `samples/`

## แหล่งข้อมูลที่รองรับ

ตัวอย่าง source ที่ normalize ได้:

- firewall
- network
- api
- crowdstrike
- aws
- m365
- ad

## Tenant ของแต่ละช่องทาง

- Syslog ใช้ tenant จาก `SYSLOG_TENANT` ไม่เชื่อ tenant ที่มากับข้อความจากอุปกรณ์
- JSON ingestion ผ่าน HTTP ใช้ tenant จาก API key ที่ authenticate แล้ว หรือจาก admin session ถ้า import ผ่าน UI
- ถ้า payload ระบุ tenant ไม่ตรงกับ credential ระบบจะ reject

ดูตัวอย่าง sender และ payload ได้ที่ `samples/` และอ่านรูปแบบ/ข้อจำกัดเพิ่มเติมที่ `docs/architecture.md`
