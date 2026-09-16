# Logdesk — Log Management Demo

ระบบต้นแบบสำหรับบททดสอบ Full-Stack Intern: **React + TypeScript + Tailwind CSS**, **Go + Fiber v3**, **PostgreSQL 17 + JSONB**

## เริ่มใช้งานบน Windows

ต้องเปิด Docker Desktop ให้ Engine พร้อมก่อน (โหมด Linux containers)

```powershell
cd "D:\My Works\LogManagement"
.\run.ps1
```

เปิด http://localhost:8080 แล้วเข้าสู่ระบบด้วย `admin.a@demo.local` และค่า `ADMIN_PASSWORD` ใน `.env` ซึ่งสคริปต์สร้างให้แบบสุ่ม อย่าเผยแพร่ไฟล์ `.env`

บน Ubuntu/Linux ใช้ `sh run.sh` (ต้องมี Docker Compose และ openssl)

## บัญชีเดโม

| Email | Role | Tenant | Password from .env |
|---|---|---|---|
| admin.a@demo.local | Admin | demo-a | ADMIN_PASSWORD |
| admin.b@demo.local | Admin | demo-b | ADMIN_PASSWORD |
| viewer.a@demo.local | Viewer | demo-a | VIEWER_PASSWORD |
| viewer.b@demo.local | Viewer | demo-b | VIEWER_PASSWORD |

Admin จัดการได้เฉพาะ tenant ของตน Viewer ดูข้อมูลได้แต่ ingest/แก้กฎ/acknowledge ไม่ได้ API key สำหรับส่งข้อมูลผูกกับ tenant โดยตรง และไม่มีสิทธิ์อ่านข้อมูล

Tenant ไม่ได้ถูกเลือกจาก UI หรือรับจาก header ที่ผู้ใช้กำหนดเองสำหรับการอ่านข้อมูล Backend enforce tenant จากตัวตนที่ยืนยันแล้วเท่านั้น: session cookie หลัง login สำหรับ UI/API แบบผู้ใช้ และ `X-API-Key` สำหรับ ingestion แบบ machine-to-machine ถ้า request ส่ง `tenant` มาไม่ตรงกับ session/API key ระบบจะ reject หรือไม่ให้ query ข้าม tenant

## ลองระบบ

1. ไป **Data sources → Send demo batch** เพื่อส่ง 12 events จริงเข้าฐานข้อมูล (รวม failed login 5 ครั้ง)
2. ไป **Overview** ดูกราฟ/Top IP/User/Event และ **Log explorer** เพื่อค้นหาและเปิดรายละเอียด
3. ไป **Alerts** ดูกฎ failed login (ค่าเริ่มต้น 5 ครั้งใน 5 นาที) และ acknowledge
4. Import `samples/aws.json`, `samples/m365.json`, `samples/ad.json` หรือ `samples/events.json`
5. ทดสอบ Syslog ด้วย `python samples/send_syslog.py` และ `python samples/send_syslog.py --tcp`; Overview และ Log explorer refresh อัตโนมัติทุก 10 วินาที และยังมีปุ่ม Refresh สำหรับกดเอง
6. ออกจากระบบแล้วเข้า Viewer A/B เพื่อตรวจ tenant isolation

ทุกไฟล์ sample เป็นข้อมูลสังเคราะห์ ไม่มี timestamp จึงใช้เวลารับเข้า หากส่ง timestamp เก่าต้องเลือก Custom range ใน UI กราฟใช้ UTC; ตารางใช้ timezone ของเบราว์เซอร์

## API หลัก

| Method / path | หน้าที่ | สิทธิ์ |
|---|---|---|
| GET /api/health | ตรวจ API + DB | ไม่ต้อง login |
| POST /api/auth/login | email/password → HttpOnly session cookie | ไม่ต้อง login |
| GET /api/auth/me | ผู้ใช้ปัจจุบัน | Session |
| POST /api/auth/logout | ยกเลิก session | Session |
| POST /ingest หรือ /api/ingest | JSON object/array/AWS Records หรือ multipart field `file` | Admin session / X-API-Key |
| GET /api/logs | ค้นหาและแบ่งหน้า | Session |
| GET /api/stats | จำนวน/กราฟ/Top N | Session |
| GET /api/alerts | ล่าสุด 100 alerts ของ tenant | Session |
| POST /api/alerts/:id/acknowledge | รับทราบ alert | Admin |
| GET /api/rule | ดูกฎ | Session |
| PUT /api/rule | ตั้ง enabled/threshold/window_minutes | Admin |

Query: `q`, `source`, `from`, `to` (RFC3339), `limit` (1–100), `offset` (0–100000) ช่วงเวลามากสุด 31 วัน ค่าเริ่มต้น 24 ชั่วโมง `tenant` ถ้าส่งต้องตรงกับบัญชี ข้าม tenant ไม่ได้

Cookie-authenticated writes ต้องส่ง `Origin` ตรงกับ `APP_ORIGIN`; browser ทำให้อัตโนมัติ API scripts ใช้ `X-API-Key` โดยไม่ต้องใช้ cookie

## โครงโปรเจกต์

```text
backend/cmd/server/       entry point หลักของ backend
backend/internal/config/  โหลดและ validate environment configuration
backend/internal/bootstrap/ init ระบบและ graceful shutdown
backend/internal/collector/ ตัวรับ Syslog UDP/TCP
backend/internal/retention/ worker ลบข้อมูลเก่าตาม retention
backend/internal/router/  ลงทะเบียน HTTP endpoints
backend/internal/handler/ handlers แยกตาม feature
backend/internal/middleware/ auth, origin checks และ error responses
backend/internal/model/   data structures ที่ใช้ร่วมกัน
backend/internal/normalize/ parse และ normalize log
backend/internal/repository/ PostgreSQL และ schema
frontend/src/             React UI
samples/                  JSON samples และ Python senders
tests/                    smoke/integration tests สำหรับ stack ที่รันอยู่
scripts/                  สคริปต์สร้าง environment แบบสุ่ม
deploy/                   ตัวอย่าง Nginx HTTPS
docs/                     สถาปัตยกรรม, คู่มือติดตั้ง และ Postman collection
```

## ทดสอบ

```powershell
cd backend
go test ./...
go vet ./...
cd ../frontend
npm ci
npm run build
cd ..
python tests/smoke.py
```

Smoke test ต้องมี stack รันอยู่และ `.env` จริง จะเพิ่มข้อมูลทดสอบใน demo tenants และทดสอบ alert โดยคืนค่ากฎเดิมหลังจบ

เอกสารส่งมอบหลักอยู่ที่ [สถาปัตยกรรม](docs/architecture.md), [ติดตั้ง Appliance](docs/setup_appliance.md) และ [ติดตั้ง SaaS](docs/setup_saas.md)

Postman/Insomnia collection สำหรับทดสอบ API อยู่ที่ [docs/postman_collection.json](docs/postman_collection.json) และวิธีใช้อยู่ที่ [docs/postman.md](docs/postman.md) หลัง import ให้ตั้งค่า `admin_password`, `api_key_a` และ `api_key_b` จาก `.env`

## ขอบเขตเวอร์ชันนี้

เป็น local MVP ที่มี ingestion, normalization, search, dashboard, RBAC, tenant isolation, UI alerts และ retention พร้อมโค้ดติดตั้ง Cloud แต่ **ยังไม่มี URL Cloud หรือวิดีโอเดโม** ไม่มีการเชื่อมบัญชี CrowdStrike/AWS/M365 ของจริง

Syslog รองรับ RFC3164 / key=value และ newline-framed TCP โดยใช้ tenant ที่ตั้งบน collector; ยังไม่มี TLS Syslog, octet-counted framing, durable queue, retry/dedup หรือ full-text index สำหรับค้นคำใน raw จึงต้องวัดประสิทธิภาพตามปริมาณงานก่อนใช้จริง

## ที่มาของตัวเลือก

- Fiber v3 ต้องใช้ Go 1.25+: https://docs.gofiber.io/
- Vite: https://vite.dev/guide/
- PostgreSQL JSONB/indexing: https://www.postgresql.org/docs/17/datatype-json.html
- Docker Compose บนเซิร์ฟเวอร์เดียว: https://docs.docker.com/compose/how-tos/production/
