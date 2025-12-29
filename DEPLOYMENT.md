# Deployment Guide

Guide untuk deploy Go Gin Backend Starter Kit ke production.

## Pre-Deployment Checklist

### Security
- [ ] Change `JWT_SECRET` to strong value (generate with `openssl rand -base64 32`)
- [ ] Use strong database password
- [ ] Set `APP_ENV=production`
- [ ] Configure HTTPS/SSL certificate
- [ ] Update CORS origins to production domain
- [ ] Review rate limiting settings
- [ ] Enable database backups

### Performance
- [ ] Optimize database connection pool settings
- [ ] Setup Redis for caching
- [ ] Enable compression
- [ ] Test with production-like load
- [ ] Setup monitoring & logging

### Configuration
- [ ] Update all .env values for production
- [ ] Configure proper CORS origins
- [ ] Set appropriate rate limits
- [ ] Configure logging to file/service
- [ ] Setup error tracking (Sentry, etc.)

---

## Local Production Build

### Build executable
```bash
make build
# Or manually:
go build -o bin/main main.go
```

### Run production build
```bash
APP_ENV=production ./bin/main
```

### Test
```bash
curl http://localhost:8080/health
```

---

## Docker Deployment

### Prerequisites
- Docker installed
- Docker Compose installed

### Build Docker image
```bash
docker-compose build
```

### Run with Docker
```bash
docker-compose up -d
```

### View logs
```bash
docker-compose logs -f app
```

### Stop containers
```bash
docker-compose down
```

### Production docker-compose
Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      APP_ENV: production
      DB_HOST: db
      REDIS_HOST: redis
    depends_on:
      - db
      - redis
    restart: unless-stopped

  db:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: go_backend_db
      MYSQL_ROOT_PASSWORD: your-strong-password
      MYSQL_PASSWORD: your-strong-password
    volumes:
      - db_data:/var/lib/mysql
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    restart: unless-stopped

volumes:
  db_data:
```

Run:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

---

## Cloud Deployment

### Heroku

1. **Setup**
```bash
heroku create your-app-name
heroku addons:create cleardb:ignite  # MySQL
heroku addons:create heroku-redis:premium-0
```

2. **Configure environment**
```bash
heroku config:set APP_ENV=production
heroku config:set JWT_SECRET=$(openssl rand -base64 32)
heroku config:set DB_HOST=<cleardb-host>
heroku config:set DB_USERNAME=<cleardb-user>
heroku config:set DB_PASSWORD=<cleardb-password>
```

3. **Create Procfile**
```
web: ./bin/main
```

4. **Build & Deploy**
```bash
git push heroku main
```

### AWS (EC2)

1. **Launch EC2 instance**
   - Ubuntu 20.04 LTS
   - Security group: Allow ports 22 (SSH), 80 (HTTP), 443 (HTTPS)

2. **SSH into instance**
```bash
ssh -i your-key.pem ubuntu@your-instance-ip
```

3. **Install Go**
```bash
wget https://go.dev/dl/go1.23.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.23.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

4. **Clone repository**
```bash
git clone your-repo.git
cd your-repo
```

5. **Build & Run**
```bash
make install
make build
./bin/main
```

6. **Setup systemd service**
Create `/etc/systemd/system/go-backend.service`:

```ini
[Unit]
Description=Go Backend API
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/your-repo
Environment="PATH=/usr/local/go/bin:/home/ubuntu/go/bin:$PATH"
EnvironmentFile=/home/ubuntu/your-repo/.env
ExecStart=/home/ubuntu/your-repo/bin/main
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable & start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable go-backend
sudo systemctl start go-backend
sudo systemctl status go-backend
```

### DigitalOcean

Similar to AWS EC2:
1. Create Droplet (Ubuntu 20.04)
2. SSH into droplet
3. Install Go & dependencies
4. Clone & setup application
5. Configure systemd service
6. Setup Nginx as reverse proxy

### Railway

1. **Connect repository**
   - Push to GitHub
   - Login to Railway.app
   - Create new project
   - Connect GitHub repo

2. **Add MySQL plugin**
   - Add service → MySQL
   - Configure environment variables

3. **Add Redis plugin**
   - Add service → Redis
   - Configure environment variables

4. **Deploy**
   - Configure environment variables
   - Deploy button

---

## Nginx Reverse Proxy

Setup Nginx to proxy requests to Go backend:

```nginx
upstream go_backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://go_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 90;
    }
}
```

Reload Nginx:
```bash
sudo nginx -t
sudo systemctl reload nginx
```

---

## SSL/TLS Setup

### Let's Encrypt with Certbot

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot certonly --nginx -d yourdomain.com
```

Update Nginx config:
```nginx
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://go_backend;
        # ... proxy settings
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

Setup auto-renewal:
```bash
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

---

## Monitoring & Logging

### Application Logs
```bash
# Follow logs
docker-compose logs -f app

# Or for systemd
sudo journalctl -u go-backend -f
```

### Database Backups
MySQL backup:
```bash
mysqldump -u root -p go_backend_db > backup.sql
```

Restore:
```bash
mysql -u root -p go_backend_db < backup.sql
```

### Health Checks
```bash
curl https://yourdomain.com/health
```

Monitor endpoint:
```bash
watch -n 5 'curl -s https://yourdomain.com/health | jq'
```

---

## Environment Variables - Production

```env
APP_NAME=My API Production
APP_ENV=production
APP_PORT=8080

DB_CONNECTION=mysql
DB_HOST=production-db-host
DB_PORT=3306
DB_DATABASE=production_database
DB_USERNAME=prod_user
DB_PASSWORD=strong-password-here

DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=300

JWT_SECRET=generate-with-openssl-rand-base64-32
JWT_ACCESS_TOKEN_EXPIRATION=15
JWT_REFRESH_TOKEN_EXPIRATION=10080

RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_REQUESTS=1000
RATE_LIMIT_WINDOW_SECONDS=3600

CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS

REDIS_HOST=production-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=redis-password
REDIS_DB=0
```

---

## Troubleshooting

### Application won't start
- Check environment variables are set
- Verify database connection
- Check logs for errors: `docker-compose logs app`

### High latency
- Check database queries
- Monitor Redis usage
- Adjust rate limiting if needed
- Scale horizontally if needed

### Database connection errors
- Verify database is running
- Check credentials in .env
- Verify database user permissions
- Check firewall rules

### Out of memory
- Reduce `DB_MAX_OPEN_CONNS`
- Enable Redis caching
- Scale to larger instance

---

## Performance Optimization

1. **Database**
   - Use indexes on frequently queried columns
   - Monitor slow queries
   - Optimize complex queries

2. **Caching**
   - Cache frequently accessed data in Redis
   - Implement cache invalidation strategy
   - Use appropriate TTLs

3. **Application**
   - Monitor goroutines for leaks
   - Profile with `pprof`
   - Optimize hot paths

4. **Infrastructure**
   - Use CDN for static assets
   - Load balancer for multiple instances
   - Database replication/clustering

---

## Rollback Plan

If deployment fails:

1. **Docker**
```bash
docker-compose down
docker-compose up -d # Previous version
```

2. **Systemd**
```bash
sudo systemctl stop go-backend
# Restore previous binary
sudo systemctl start go-backend
```

3. **Database**
```bash
# Restore from backup
mysql -u root -p go_backend_db < backup.sql
```

---

## Post-Deployment

1. Test all endpoints
2. Monitor error logs
3. Check performance metrics
4. Verify backups are running
5. Document any issues
6. Setup alerting & monitoring
7. Schedule regular reviews

---

**Remember:** Always test in staging before deploying to production!
