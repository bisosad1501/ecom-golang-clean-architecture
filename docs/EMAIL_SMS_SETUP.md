# Email & SMS Service Setup Guide

This guide will help you configure real email and SMS providers for your e-commerce application.

## 📧 Email Service Configuration

### Supported Providers

1. **SMTP** - Generic SMTP server (Gmail, Outlook, etc.)
2. **SendGrid** - Cloud-based email service
3. **AWS SES** - Amazon Simple Email Service
4. **Mailgun** - Email API service
5. **Mock** - Development/testing mode

### Configuration Steps

#### 1. Choose Your Provider

Set the `EMAIL_PROVIDER` environment variable:

```bash
EMAIL_PROVIDER=smtp          # For SMTP
EMAIL_PROVIDER=sendgrid      # For SendGrid
EMAIL_PROVIDER=aws_ses       # For AWS SES
EMAIL_PROVIDER=mailgun       # For Mailgun
EMAIL_PROVIDER=mock          # For development/testing
```

#### 2. Provider-Specific Setup

##### SMTP Configuration (Gmail Example)

```bash
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_USE_TLS=true
SMTP_USE_SSL=false
FROM_EMAIL=your-email@gmail.com
FROM_NAME=Your Store Name
```

**Gmail Setup:**
1. Enable 2-factor authentication
2. Generate an App Password: https://myaccount.google.com/apppasswords
3. Use the App Password as `SMTP_PASSWORD`

##### SendGrid Configuration

```bash
EMAIL_PROVIDER=sendgrid
SENDGRID_API_KEY=SG.your-sendgrid-api-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your Store Name
```

**SendGrid Setup:**
1. Sign up at https://sendgrid.com
2. Create an API key with "Mail Send" permissions
3. Verify your sender domain/email

##### AWS SES Configuration

```bash
EMAIL_PROVIDER=aws_ses
AWS_SES_REGION=us-east-1
AWS_SES_ACCESS_KEY_ID=your-aws-access-key
AWS_SES_SECRET_ACCESS_KEY=your-aws-secret-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your Store Name
```

**AWS SES Setup:**
1. Create AWS account and enable SES
2. Verify your email/domain
3. Create IAM user with SES permissions
4. Move out of sandbox mode for production

##### Mailgun Configuration

```bash
EMAIL_PROVIDER=mailgun
MAILGUN_DOMAIN=mg.yourdomain.com
MAILGUN_API_KEY=your-mailgun-api-key
MAILGUN_PUBLIC_KEY=your-mailgun-public-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Your Store Name
```

**Mailgun Setup:**
1. Sign up at https://mailgun.com
2. Add and verify your domain
3. Get API keys from dashboard

#### 3. Common Email Settings

```bash
# Email settings
REPLY_TO_EMAIL=support@yourdomain.com
MAX_EMAILS_PER_HOUR=1000
MAX_EMAILS_PER_MINUTE=50
EMAIL_MAX_RETRIES=3
EMAIL_RETRY_INTERVAL=300
EMAIL_TEMPLATE_ENGINE=html/template
```

## 📱 SMS Service Configuration

### Supported Providers

1. **Twilio** - Popular SMS API service
2. **AWS SNS** - Amazon Simple Notification Service
3. **Nexmo/Vonage** - SMS API service
4. **Mock** - Development/testing mode

### Configuration Steps

#### 1. Choose Your Provider

Set the `SMS_PROVIDER` environment variable:

```bash
SMS_PROVIDER=twilio          # For Twilio
SMS_PROVIDER=aws_sns         # For AWS SNS
SMS_PROVIDER=nexmo           # For Nexmo/Vonage
SMS_PROVIDER=mock            # For development/testing
```

#### 2. Provider-Specific Setup

##### Twilio Configuration

```bash
SMS_PROVIDER=twilio
TWILIO_ACCOUNT_SID=your-twilio-account-sid
TWILIO_AUTH_TOKEN=your-twilio-auth-token
TWILIO_FROM_NUMBER=+1234567890
```

**Twilio Setup:**
1. Sign up at https://twilio.com
2. Get Account SID and Auth Token from console
3. Purchase a phone number for sending SMS
4. Verify your account for production use

##### AWS SNS Configuration

```bash
SMS_PROVIDER=aws_sns
AWS_SNS_REGION=us-east-1
AWS_SNS_ACCESS_KEY_ID=your-aws-access-key
AWS_SNS_SECRET_ACCESS_KEY=your-aws-secret-key
```

**AWS SNS Setup:**
1. Create AWS account and enable SNS
2. Create IAM user with SNS permissions
3. Configure SMS preferences in SNS console

##### Nexmo/Vonage Configuration

```bash
SMS_PROVIDER=nexmo
NEXMO_API_KEY=your-nexmo-api-key
NEXMO_API_SECRET=your-nexmo-api-secret
NEXMO_FROM_NAME=YourStore
```

**Nexmo Setup:**
1. Sign up at https://vonage.com
2. Get API key and secret from dashboard
3. Configure sender ID

#### 3. Common SMS Settings

```bash
# SMS settings
MAX_SMS_PER_HOUR=100
MAX_SMS_PER_MINUTE=10
SMS_MAX_RETRIES=3
SMS_RETRY_INTERVAL=60
```

## 🧪 Testing Your Configuration

### 1. Start the Application

```bash
# Copy environment variables
cp .env.example .env

# Edit .env with your provider settings
nano .env

# Start the application
docker-compose up -d --build
```

### 2. Test Email Functionality

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

**Check logs for email sending:**
```bash
docker-compose logs api | grep -i email
```

**Test password reset:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

### 3. Test SMS Functionality

**Send phone verification:**
```bash
curl -X POST http://localhost:8080/api/v1/users/send-phone-verification \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "phone": "+1234567890"
  }'
```

**Check logs for SMS sending:**
```bash
docker-compose logs api | grep -i sms
```

## 🔧 Troubleshooting

### Common Issues

#### Email Issues

1. **SMTP Authentication Failed**
   - Check username/password
   - Enable "Less secure app access" for Gmail
   - Use App Password for Gmail with 2FA

2. **SendGrid API Error**
   - Verify API key permissions
   - Check sender verification status
   - Ensure domain is verified

3. **AWS SES Access Denied**
   - Check IAM permissions
   - Verify email/domain in SES
   - Check if still in sandbox mode

#### SMS Issues

1. **Twilio Authentication Failed**
   - Verify Account SID and Auth Token
   - Check account status and balance

2. **AWS SNS Permission Denied**
   - Check IAM permissions for SNS
   - Verify region settings

3. **Phone Number Format**
   - Use international format (+1234567890)
   - Verify number is valid for SMS

### Debug Mode

Enable debug logging by setting:

```bash
LOG_LEVEL=debug
```

This will show detailed logs for email/SMS operations.

### Fallback Behavior

If a provider fails to initialize or send, the system will:

1. Log the error
2. Fall back to console logging (mock mode)
3. Continue operation without failing

This ensures your application remains functional even with misconfigured providers.

## 📝 Production Recommendations

### Email

1. **Use a dedicated email service** (SendGrid, AWS SES, Mailgun)
2. **Verify your domain** for better deliverability
3. **Set up SPF, DKIM, and DMARC** records
4. **Monitor bounce rates** and reputation
5. **Implement rate limiting** to avoid being flagged as spam

### SMS

1. **Use a reliable SMS provider** (Twilio, AWS SNS)
2. **Register your business** with carriers
3. **Use short codes** for high-volume sending
4. **Implement opt-out mechanisms**
5. **Monitor delivery rates** and costs

### Security

1. **Use environment variables** for sensitive credentials
2. **Rotate API keys** regularly
3. **Implement rate limiting** to prevent abuse
4. **Monitor usage** for unusual patterns
5. **Use HTTPS** for all webhook endpoints

## 🆘 Support

If you encounter issues:

1. Check the application logs
2. Verify your provider credentials
3. Test with mock providers first
4. Check provider status pages
5. Review provider documentation

For provider-specific support:
- **SendGrid**: https://support.sendgrid.com
- **Twilio**: https://support.twilio.com
- **AWS**: https://aws.amazon.com/support
- **Mailgun**: https://help.mailgun.com
- **Nexmo**: https://help.nexmo.com
