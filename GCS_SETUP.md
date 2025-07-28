# Google Cloud Storage Setup Guide

This document provides instructions for setting up Google Cloud Storage integration for image uploads in the CT Backend Course application.

## Prerequisites

- Google Cloud Platform account
- Google Cloud SDK installed (optional but recommended)
- A Google Cloud Storage bucket created and configured

## Step 1: Create a Google Cloud Storage Bucket

1. **Create a new bucket:**
   ```bash
   gsutil mb gs://your-unique-bucket-name
   ```

2. **Set bucket permissions for public read access:**
   ```bash
   gsutil iam ch allUsers:objectViewer gs://your-unique-bucket-name
   ```

3. **Enable uniform bucket-level access (recommended):**
   ```bash
   gsutil uniformbucketlevelaccess set on gs://your-unique-bucket-name
   ```

## Step 2: Create Service Account and Credentials

1. **Create a service account:**
   ```bash
   gcloud iam service-accounts create gcs-image-uploader \
       --description="Service account for image uploads to GCS" \
       --display-name="GCS Image Uploader"
   ```

2. **Grant necessary permissions:**
   ```bash
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
       --member="serviceAccount:gcs-image-uploader@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
       --role="roles/storage.objectCreator"
   
   gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
       --member="serviceAccount:gcs-image-uploader@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
       --role="roles/storage.legacyBucketReader"
   ```

3. **Create and download the service account key:**
   ```bash
   gcloud iam service-accounts keys create ./gcs-credentials.json \
       --iam-account=gcs-image-uploader@YOUR_PROJECT_ID.iam.gserviceaccount.com
   ```

## Step 3: Configure Environment Variables

Create a `.env` file in the project root or set environment variables:

```bash
# Required GCS Configuration
GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/gcs-credentials.json
GOOGLE_APPLICATION_BUCKET=your-unique-bucket-name

# PostgreSQL Configuration (existing)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=ct_backend_course
POSTGRES_SSLMODE=disable

# Other configurations...
PORT=8090
```

## Step 4: Application Configuration

The application will automatically detect and use Google Cloud Storage if the credentials and bucket name are properly configured. If not, it will fall back to a fake implementation for development.

### Configuration Validation

The application includes validation to ensure:
- Credentials file exists and is accessible
- Bucket name is specified
- Service account has proper permissions

### File Upload Features

The implemented GCS integration includes:

1. **File Type Validation**: Only allows image files (jpg, jpeg, png, gif, webp)
2. **File Size Limits**: Maximum 10MB per file
3. **Unique File Naming**: Uses UUID with date-based folder structure
4. **Content Type Detection**: Automatically sets appropriate MIME types
5. **Public URL Generation**: Returns publicly accessible URLs
6. **Error Handling**: Comprehensive error reporting and logging

### File Organization

Uploaded files are organized in the bucket with the following structure:
```
your-bucket/
├── 2024/
│   ├── 01/
│   │   ├── 15/
│   │   │   ├── 550e8400-e29b-41d4-a716-446655440000.jpg
│   │   │   └── 6ba7b810-9dad-11d1-80b4-00c04fd430c8.png
│   │   └── 16/
│   │       └── 6ba7b811-9dad-11d1-80b4-00c04fd430c8.gif
│   └── 02/
│       └── ...
```

## Step 5: Testing the Setup

1. **Run the application:**
   ```bash
   go run main.go
   ```

2. **Check logs for GCS initialization:**
   ```
   INFO Using Google Cloud Storage bucket: your-unique-bucket-name
   ```

3. **Test image upload via API:**
   ```bash
   # First register and login to get a token
   curl -X POST http://localhost:8090/api/public/register \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"testpass","fullName":"Test User","address":"Test Address"}'

   curl -X POST http://localhost:8090/api/public/login \
     -H "Content-Type: application/json" \
     -d '{"username":"testuser","password":"testpass"}'

   # Upload an image (replace TOKEN with actual token)
   curl -X POST http://localhost:8090/api/private/upload \
     -H "Authorization: Bearer TOKEN" \
     -F "file=@/path/to/your/test-image.jpg"
   ```

## Troubleshooting

### Common Issues

1. **"Service account key not found"**
   - Ensure the credentials file path is correct
   - Check file permissions (readable by the application)

2. **"Access denied" errors**
   - Verify service account has `storage.objectCreator` role
   - Ensure bucket exists and is accessible

3. **"Bucket not found"**
   - Double-check bucket name spelling
   - Verify bucket exists in the same project as the service account

4. **"File too large" errors**
   - Default limit is 10MB, configurable in `pkg/bucket/google.go`
   - Check client-side file size before upload

### Logs and Monitoring

The application provides detailed logging for GCS operations:
- Upload start/completion
- File size and content type information
- Error details with context
- Public URL generation

Monitor these logs to troubleshoot issues and track usage patterns.

## Security Considerations

1. **Service Account Permissions**: Use minimal required permissions
2. **Credential Storage**: Never commit credentials to version control
3. **Public Access**: Configure bucket permissions carefully
4. **File Validation**: The app validates file types and sizes
5. **Rate Limiting**: Consider implementing upload rate limits for production

## Production Deployment

For production deployment:

1. Use Google Cloud's Application Default Credentials when possible
2. Store credentials securely (e.g., Kubernetes secrets, cloud provider secret managers)
3. Configure monitoring and alerting for GCS operations
4. Set up proper backup and disaster recovery procedures
5. Monitor storage costs and implement lifecycle policies if needed

## Cost Optimization

- Implement lifecycle policies to automatically delete old files
- Use regional storage classes for better performance
- Monitor usage patterns and optimize bucket configuration
- Consider CDN integration for frequently accessed images