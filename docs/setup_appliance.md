# การติดตั้งแบบ Appliance

เป้าหมาย: รันบนเครื่องหรือ VM เดียว เช่น Ubuntu 22.04+, 4 vCPU, 8 GB RAM, 40 GB disk พร้อม Docker Engine + Compose plugin  
สำหรับ Windows development ให้ใช้ Docker Desktop โหมด Linux containers

## เริ่มระบบ

Windows PowerShell จาก root project:

```powershell
.\run.ps1
```

Ubuntu/Linux:

```sh
sh run.sh
```

สคริปต์จะสร้าง `.env` แบบสุ่มให้เฉพาะตอนที่ไฟล์ยังไม่มีอยู่ จากนั้น build/start PostgreSQL, backend และ frontend การ start ครั้งแรกต้องใช้อินเทอร์เน็ตเพื่อ pull image และติดตั้ง dependencies

เปิดเว็บ:

```text
http://localhost:8080
```

อ่านค่า `ADMIN_PASSWORD` ใน `.env` เพื่อ login ด้วย `admin.a@demo.local` หรือ `admin.b@demo.local` ส่วน Viewer ใช้ `VIEWER_PASSWORD` ห้ามส่งไฟล์ `.env` จริงขึ้น Git

คำสั่งตรวจสถานะ:

```sh
docker compose ps
docker compose logs --tail=100 backend
curl http://localhost:8080/api/health
```

services จะ start ตาม health-check order ข้อมูล DB อยู่ใน named volume `postgres_data` และยังอยู่หลัง `docker compose down` ถ้าไม่ใส่ `-v` อย่าลบ volume เพื่อแก้ปัญหา เว้นแต่ยอมทิ้งข้อมูลได้

## ตัวอย่างการส่ง log

ผ่าน UI:

- Data sources → Send demo batch
- Import logs → `samples/events.json`

Syslog:

```sh
python samples/send_syslog.py
python samples/send_syslog.py --tcp
```

HTTP simulator:

```sh
python samples/post_logs.py
python samples/post_logs.py --tenant b
python samples/post_logs.py --alert
python samples/post_logs.py --file samples/aws.json
```

`post_logs.py` จะอ่าน `API_KEY_A` หรือ `API_KEY_B` จาก `.env` ให้อัตโนมัติ ถ้าต้อง override เองให้ set `LOG_API_KEY`

ไฟล์ sample ตั้งใจ omit tenant เพื่อให้ credentials เป็นตัวกำหนด tenant:

- UI ใช้ tenant จาก session ของ user ที่ login
- HTTP API ใช้ tenant จาก `X-API-Key`
- ถ้า payload ใส่ `tenant` ไม่ตรงกับ credential ระบบจะ reject

ถ้า record ไม่มี timestamp จะใช้ ingestion time ถ้าส่ง timestamp เก่าต้องเลือก Custom range ใน UI ให้ครอบคลุมช่วงเวลา สูงสุด 31 วัน

## การเปิด port และ network

ค่าเริ่มต้น:

- HTTP UI bind ที่ `127.0.0.1:8080`
- Syslog UDP/TCP bind ที่ `127.0.0.1:5514`
- PostgreSQL และ backend HTTP port ไม่ publish ออก host

Syslog events จะเข้า tenant ตาม `SYSLOG_TENANT` ใน `.env` ค่า default คือ `demo-a` ถ้าจะทดสอบ tenant B ให้เปลี่ยนเป็น:

```env
SYSLOG_TENANT=demo-b
```

แล้ว restart backend:

```sh
docker compose up -d --build backend
```

ถ้าทดสอบใน LAN ที่ไว้ใจได้ ให้ตั้ง `HTTP_BIND` หรือ `SYSLOG_BIND` เป็น private interface ที่ต้องการ และตั้ง `APP_ORIGIN` ให้ตรงกับ URL ที่เปิดใน browser ควร firewall Syslog ให้รับเฉพาะอุปกรณ์ที่ไว้ใจได้ อย่าเปิด unauthenticated Syslog สู่ Internet

host Syslog port ใช้ `5514` เพื่อเลี่ยง privileged port ถ้าต้องใช้ `514` ให้เปลี่ยนเฉพาะ host port ใน Compose และส่งด้วย `--port 514`

## พัฒนาแบบ local โดยไม่ rebuild container

ใช้ PostgreSQL local หรือ start เฉพาะ Compose database พร้อม loopback port mapping สำหรับ dev โดยยึดค่า DB จาก `.env` เป็นหลัก:

- `DB_HOST`
- `DB_PORT`
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`
- `ADMIN_PASSWORD`
- `VIEWER_PASSWORD`
- `API_KEY_A`
- `API_KEY_B`
- `APP_ORIGIN`
- `COOKIE_SECURE`

เมื่อรันผ่าน Docker Compose ระบบจะประกอบ `DATABASE_URL` ให้ backend จากค่า `DB_*` เหล่านี้เอง จึงไม่ต้องแก้ `docker-compose.yml`

สำหรับ Vite default:

```text
APP_ORIGIN=http://127.0.0.1:5173
COOKIE_SECURE=false
```

รัน backend:

```sh
cd backend
go run ./cmd/server
```

อีก terminal รัน frontend:

```sh
cd frontend
npm ci
npm run dev
```

## วิธีแก้ปัญหาเบื้องต้น

- Docker pipe/daemon unavailable: เปิด Docker Desktop และรอ Engine running ถ้า Docker ขอ WSL setup, license acceptance หรือ reboot ให้ทำบน host ให้เสร็จก่อน
- Login invalid: อ่าน password จาก `.env` เดิม seed ไม่ update user/key ที่มีอยู่แล้วหลัง restart
- 403 ตอนเขียนข้อมูล: ตรวจ `APP_ORIGIN` ให้ตรงกับ browser URL ทั้ง scheme และ port; cloud ต้องตั้ง `COOKIE_SECURE=true`
- ไม่เห็น logs: ตรวจ source/time/search filters และ tenant account; collector tenant มาจาก `SYSLOG_TENANT`
- Overview และ Log explorer auto-refresh ทุก 10 วินาที และยังมีปุ่ม Refresh สำหรับโหลดทันที
- API unavailable: ดู `docker compose logs backend postgres`; UI ไม่มี mock fallback
