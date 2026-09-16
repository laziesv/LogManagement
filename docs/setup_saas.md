# การ deploy แบบ SaaS / Cloud

สถานะปัจจุบัน: repository มีขั้นตอน deploy และตัวอย่าง HTTPS Nginx config แล้ว แต่ยังไม่ได้สร้าง cloud resource, DNS record หรือ public URL จริง

ใช้ VM เช่น Ubuntu 22.04+, 4 vCPU, 8 GB RAM, 40 GB disk พร้อม Docker/Compose, host Nginx และ domain ที่ชี้มาที่ VM ควรให้ database และ backend อยู่หลัง reverse proxy ไม่เปิดตรงสู่ public

## ขั้นตอน

1. Copy repository ไปที่ VM แล้วรัน:

   ```sh
   sh scripts/init-env.sh
   ```

2. แก้ `.env`:

   ```env
   APP_ORIGIN=https://YOUR_DOMAIN
   COOKIE_SECURE=true
   HTTP_BIND=127.0.0.1
   SYSLOG_BIND=127.0.0.1
   ```

3. Start services:

   ```sh
   docker compose up --build -d
   curl http://localhost:8080/api/health
   ```

4. ขอ TLS certificate สำหรับ domain ด้วย ACME client/issuer ที่เลือก แล้วติดตั้ง host Nginx โดยใช้ `deploy/nginx-tls.conf.example` เป็น template แก้ domain และ certificate paths ให้ตรงเครื่องจริง

5. ตรวจ config ก่อน reload:

   ```sh
   sudo nginx -t
   sudo systemctl reload nginx
   ```

6. เปิด cloud firewall เฉพาะ 80/443 และจำกัด SSH ให้ trusted source อย่าเปิด PostgreSQL port 5432 หรือ public Syslog

7. เปิด `https://YOUR_DOMAIN`, login และส่ง sample

8. ตรวจว่า session cookie มี Secure และ HttpOnly attributes

9. ตรวจว่า Viewer B ไม่เห็น logs ของ `demo-a`

10. Tenant scope มาจาก session หรือ API key จึงไม่สามารถ switch tenant จาก browser query string ได้

11. รัน smoke test ด้วย `.env` เดียวกันและ `APP_ORIGIN` แบบ HTTPS:

    ```sh
    python tests/smoke.py
    ```

12. ส่งเฉพาะ public URL และ demo credentials ให้ผู้ประเมินโดยตรง อย่าใส่ secrets ใน Git หรือเอกสาร

## ใบรับรองแบบ self-signed

โจทย์ยอมรับ self-signed certificate ถ้าอธิบายขั้นตอนชัดเจน แต่ควรเตรียมวิธี trust/import certificate ให้กรรมการ และไม่ควร disable TLS validation ใน application หรือ scripts

ถ้าเป็นไปได้ ใช้ certificate ที่ browser trust อยู่แล้วจะทดสอบง่ายกว่า

## ข้อควรระวัง

- backup PostgreSQL volume
- จำกัดสิทธิ์ SSH/server access
- monitor disk capacity
- demo นี้ใช้ shared-table tenant isolation ที่ application layer ยังไม่ได้ load test สำหรับ production scale
