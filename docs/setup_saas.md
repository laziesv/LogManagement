# การ deploy แบบ SaaS / Cloud

สถานะปัจจุบันของ demo: ระบบ deploy บน Azure VM แล้ว และเข้าใช้งานผ่าน HTTPS ได้ที่

```text
https://logdisk.malaysiawest.cloudapp.azure.com
```

แนวทางนี้ใช้ VM เครื่องเดียวรัน Docker Compose โดยให้ frontend เปิดที่ `127.0.0.1:8080` แล้วใช้ host Nginx reverse proxy ออก public HTTPS ส่วน backend และ PostgreSQL ไม่เปิด public โดยตรง

## สเปกแนะนำ

- Ubuntu 22.04+ หรือใกล้เคียง
- 2 ถึง 4 vCPU
- RAM 4 ถึง 8 GB
- Disk 30 GB ขึ้นไป
- Docker Engine + Docker Compose plugin
- Nginx บน host
- Certbot / Let's Encrypt สำหรับ HTTPS

## Environment บน VM

สร้าง `.env` ที่ root project บน VM และตั้งค่าหลักดังนี้

```env
APP_ORIGIN=https://logdisk.malaysiawest.cloudapp.azure.com
COOKIE_SECURE=true
HTTP_BIND=127.0.0.1
SYSLOG_BIND=0.0.0.0
SYSLOG_TENANT=demo-a
RETENTION_DAYS=7
```

ค่าจริงของ password และ API key ต้องเก็บเฉพาะใน `.env` บน VM ห้าม commit เข้า Git

## ขั้นตอน deploy บน VM

1. ติดตั้ง Docker และ Docker Compose plugin

   ```sh
   docker compose version
   docker ps
   ```

2. Clone repository ไปที่ VM

   ```sh
   git clone https://github.com/laziesv/LogManagement.git
   cd LogManagement
   ```

3. สร้าง `.env`

   ```sh
   sh scripts/init-env.sh
   ```

4. แก้ `.env` ให้เป็นค่า production/demo เช่น `APP_ORIGIN`, `COOKIE_SECURE`, `HTTP_BIND`, `SYSLOG_BIND`

5. Start container

   ```sh
   docker compose up --build -d
   docker compose ps
   ```

6. ตรวจ health จากเครื่อง VM

   ```sh
   curl http://127.0.0.1:8080/api/health
   ```

7. ตั้งค่า Nginx reverse proxy โดยใช้ `deploy/nginx-tls.conf.example` เป็น template แล้วแก้ domain ให้ตรงกับ VM

8. ขอ TLS certificate ด้วย Certbot

   ```sh
   sudo certbot --nginx -d logdisk.malaysiawest.cloudapp.azure.com
   ```

9. ตรวจ Nginx และ reload

   ```sh
   sudo nginx -t
   sudo systemctl reload nginx
   ```

10. เปิด Azure Network Security Group เฉพาะ port ที่ต้องใช้

    | Port | Protocol | ใช้สำหรับ |
    |---|---|---|
    | 22 | TCP | SSH เฉพาะ trusted source |
    | 80 | TCP | HTTP / Certbot challenge |
    | 443 | TCP | HTTPS web app |
    | 5514 | TCP | Syslog TCP demo |
    | 5514 | UDP | Syslog UDP demo |

    ไม่ควรเปิด PostgreSQL port `5432` ออก public

## ทดสอบหลัง deploy

เปิดเว็บ

```text
https://logdisk.malaysiawest.cloudapp.azure.com
```

ตรวจสิ่งเหล่านี้

- Login ด้วย admin demo ได้
- Session cookie มี `Secure` และ `HttpOnly`
- Viewer เห็นเฉพาะ tenant ของตัวเอง
- ส่ง HTTP JSON ผ่าน `POST /ingest` แล้วค้นหาได้
- ส่ง Syslog TCP/UDP ไป port `5514` แล้วเห็นใน UI
- Import sample AWS/M365/AD แล้ว normalize ได้
- ส่ง abnormal log ครบ 5 ครั้งภายใน 5 นาทีแล้วเห็น alert ใน UI

## ตัวอย่างส่ง Syslog ไป SaaS

UDP

```powershell
py samples/send_syslog.py --host logdisk.malaysiawest.cloudapp.azure.com --port 5514
```

TCP

```powershell
py samples/send_syslog.py --tcp --host logdisk.malaysiawest.cloudapp.azure.com --port 5514
```

ถ้า UDP ไม่ถึง แต่ TCP ถึง ให้ตรวจ network/firewall ของฝั่งผู้ส่งด้วย เพราะบางเครือข่าย block UDP ขาออก

## CI/CD กับ Azure VM

หลังตั้ง VM ครั้งแรกแล้ว การ deploy ต่อไปใช้ GitHub Actions ได้ โดย push เข้า branch `deploy`

Workflow จะ SSH เข้า VM, pull code ล่าสุด และรัน

```sh
docker compose up -d --build
```

รายละเอียดอยู่ที่ [cicd.md](cicd.md)

## ข้อควรระวัง

- Backup PostgreSQL volume เป็นระยะ
- จำกัด SSH ให้ trusted source เท่านั้น
- Monitor disk usage เพราะ log เพิ่มขึ้นเรื่อย ๆ
- Syslog ไม่มี authentication ควรใช้ firewall allowlist หรือ private network ในงานจริง
- Demo นี้ใช้ shared table พร้อม tenant column และ enforce tenant ที่ application layer
