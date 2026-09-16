# Logdesk — Log Management Demo

Logdesk เป็นระบบจัดการ log แบบ multi tenant สำหรับงาน Full Stack Developer Intern Assignment ระบบนี้รับ log ได้หลายช่องทาง, normalize field กลาง, ค้นหา log, แสดง dashboard, แจ้งเตือน และทดสอบ RBAC ได้

Tech stack หลัก

- Frontend: React, TypeScript, Vite, Tailwind CSS, Lucide icons
- Backend: Go, Fiber
- Database: PostgreSQL 17 + JSONB
- Deployment: Docker Compose, Nginx, Azure VM, Certbot HTTPS
- CI/CD: GitHub Actions + SSH deploy

## เริ่มใช้งานแบบ Appliance

บน Windows ให้เปิด Docker Desktop ก่อน แล้วรันจาก root project

```powershell
.\run.ps1
```

ถ้า PowerShell block script ให้ใช้คำสั่งนี้แทน

```powershell
powershell -ExecutionPolicy Bypass -File .\run.ps1
```

หรือใช้ Docker Compose โดยตรง

```powershell
docker compose up --build -d
```

เปิดเว็บที่

```text
http://localhost:8080
```

บน Linux หรือ Azure VM ใช้

```sh
sh run.sh
```

## SaaS demo

ระบบ demo deploy บน Azure VM และเข้าใช้งานผ่าน HTTPS ได้ที่

```text
https://logdisk.malaysiawest.cloudapp.azure.com
```

Syslog public demo ใช้ port `5514` ทั้ง TCP และ UDP ถ้าเปิด firewall / Azure NSG แล้ว

## บัญชี demo

| Email | Role | Tenant | Password |
|---|---|---|---|
| `admin.a@demo.local` | Admin | `demo-a` | `ADMIN_PASSWORD` จาก `.env` |
| `admin.b@demo.local` | Admin | `demo-b` | `ADMIN_PASSWORD` จาก `.env` |
| `viewer.a@demo.local` | Viewer | `demo-a` | `VIEWER_PASSWORD` จาก `.env` |
| `viewer.b@demo.local` | Viewer | `demo-b` | `VIEWER_PASSWORD` จาก `.env` |

Admin เป็น admin เฉพาะ tenant ของตัวเอง ไม่ใช่ super admin ข้าม tenant ส่วน Viewer อ่านข้อมูลได้เฉพาะ tenant ของตัวเอง

## Tenant enforcement

ระบบไม่ให้ผู้ใช้เลือก tenant เองจาก UI เพื่อป้องกันการข้าม tenant

| ช่องทาง | วิธีกำหนด tenant |
|---|---|
| Browser UI | tenant จาก session หลัง login |
| HTTP ingest | tenant จาก `X-API-Key` |
| Syslog | tenant จาก `SYSLOG_TENANT` |

ถ้า payload ส่ง `tenant` มาไม่ตรงกับ credential ระบบจะ reject หรือไม่ให้ query ข้าม tenant

## วิธี demo ระบบ

1. เปิดเว็บและ login เป็น `admin.a@demo.local`
2. ไปหน้า Data Sources แล้วส่ง demo batch หรือ import sample
3. เปิด Overview เพื่อดู Top N, Timeline และ filter ตาม source/time
4. เปิด Log Explorer เพื่อค้นหา log
5. ส่ง abnormal logs 5 รายการเพื่อ trigger alert
6. เปิด Alerts เพื่อดู alert และ acknowledge
7. login เป็น Viewer A/B เพื่อตรวจว่าเห็นเฉพาะ tenant ของตัวเอง
8. ยิง Syslog ด้วย TCP/UDP แล้วดู log ใน UI ภายในประมาณ 1 นาที

Alert ปัจจุบันมีเงื่อนไขหลักเดียว คือพบ log field ที่ rule กำหนดครบ **5 ครั้งภายใน 5 นาที** ภายใน tenant เดียวกัน เช่น failed login จาก IP เดียวกัน 5 ครั้ง

## ช่องทาง ingest

- HTTP JSON: `POST /ingest` หรือ `POST /api/ingest` พร้อม `X-API-Key`
- File batch: import JSON sample ผ่าน UI หรือ API
- Syslog TCP/UDP: port `5514`

ตัวอย่างส่ง Syslog

```powershell
py samples/send_syslog.py --host localhost --port 5514
py samples/send_syslog.py --tcp --host localhost --port 5514
```

ตัวอย่างส่ง HTTP logs

```powershell
py samples/post_logs.py --file samples/tenant-a/events.json
```

## API หลัก

| Method / path | หน้าที่ | สิทธิ์ |
|---|---|---|
| `GET /api/health` | ตรวจ backend และ database | public |
| `POST /api/auth/login` | login และสร้าง session cookie | public |
| `GET /api/auth/me` | ดู session ปัจจุบัน | session |
| `POST /api/auth/logout` | logout | session |
| `POST /ingest` หรือ `/api/ingest` | ingest JSON/file batch | Admin session หรือ API key |
| `GET /api/logs` | ค้นหา log | session |
| `GET /api/stats` | dashboard stats | session |
| `GET /api/alerts` | ดู alerts ล่าสุดของ tenant | session |
| `POST /api/alerts/:id/acknowledge` | acknowledge alert | Admin |
| `GET /api/rule` | ดู alert rule | session |
| `PUT /api/rule` | แก้ alert rule | Admin |

Cookie-authenticated writes ต้องมี `Origin` ตรงกับ `APP_ORIGIN` ส่วน script/API automation ควรใช้ `X-API-Key`

## โครงสร้างโปรเจกต์

```text
backend/                  Go Fiber backend
frontend/                 React frontend
samples/                  sample logs และ sender scripts
ingest/                   คำอธิบายช่องทาง ingest
tests/                    live/smoke test scripts
docs/                     คู่มือ setup, architecture, Postman, CI/CD
deploy/                   ตัวอย่าง Nginx HTTPS config
scripts/                  script สร้าง .env
.github/workflows/        GitHub Actions CI/CD
```

## ทดสอบ

Backend

```powershell
cd backend
go test ./...
go vet ./...
```

Frontend

```powershell
cd frontend
npm ci
npm test
npm run build
```

Live tests ต้องเปิด stack ก่อน

```powershell
py tests/smoke.py
py tests/normalization_live.py
py tests/retention_live.py
```

ถ้าเครื่องไม่มี `python` command ให้ใช้ Python Launcher `py` หรือ path Python ที่ติดตั้งไว้

## เอกสารสำคัญ

- [Architecture](docs/architecture.md)
- [Setup Appliance](docs/setup_appliance.md)
- [Setup SaaS](docs/setup_saas.md)
- [Postman guide](docs/postman.md)
- [CI/CD](docs/cicd.md)
- [Postman collection](docs/postman_collection.json)

## ขอบเขตปัจจุบัน

ระบบนี้เป็น demo/MVP สำหรับ assignment มี ingestion, normalization, search, dashboard, RBAC, tenant isolation, UI alert, retention, SaaS HTTPS และ CI/CD แล้ว

ข้อจำกัดที่ยังมี

- Alert แจ้งเตือนผ่าน UI เท่านั้น ยังไม่มี email/webhook
- Syslog ยังไม่มี authentication และยังไม่มี TLS Syslog
- UDP เป็น best effort ไม่มี acknowledgement
- ยังไม่มี durable queue, retry, deduplication หรือ benchmark production load
- Sample logs เป็น synthetic data ไม่ได้เชื่อม AWS/M365/AD จริง
