# File Upload Endpoints

Hệ thống cung cấp 3 endpoints upload khác nhau tùy theo level authentication:

## 1. Public Upload (Không cần authentication)
```
POST /api/v1/public/upload/image
```
- **Mục đích**: Upload ảnh cho guest users (reviews, comments, etc.)
- **Authentication**: Không cần
- **Thư mục lưu**: `uploads/public/images/`
- **Use cases**: Review images, comment attachments

## 2. User Upload (Cần authentication)
```
POST /api/v1/upload/image
```
- **Mục đích**: Upload ảnh cho logged-in users
- **Authentication**: Bearer token required
- **Thư mục lưu**: `uploads/images/`
- **Use cases**: Profile pictures, user-generated content

## 3. Admin Upload (Cần admin role)
```
POST /api/v1/admin/upload/image
```
- **Mục đích**: Upload ảnh cho admin functions
- **Authentication**: Bearer token + admin role required  
- **Thư mục lưu**: `uploads/admin/images/`
- **Use cases**: Product images, category images, CMS content

## Request Format
Tất cả endpoints đều sử dụng `multipart/form-data`:
```
Content-Type: multipart/form-data
Field name: file
```

## Response Format
```json
{
  "url": "http://localhost:8080/uploads/images/uuid_timestamp.jpg",
  "message": "File uploaded successfully"
}
```

## File Constraints
- **Max size**: 5MB
- **Allowed types**: JPG, JPEG, PNG, GIF, WebP
- **Filename**: Auto-generated UUID + timestamp

## Frontend Usage

### Using the utility function:
```typescript
import { uploadImageFile } from '@/lib/utils/image-upload'

// For admin/user context (auto-detects based on current route/auth)
const imageUrl = await uploadImageFile(file)

// For specific endpoint (if needed)
const imageUrl = await uploadImageFile(file, '/admin/upload/image')
```

### Manual API call:
```typescript
import api from '@/lib/api'

const formData = new FormData()
formData.append('file', file)

const response = await api.post('/upload/image', formData, {
  headers: {
    'Content-Type': 'multipart/form-data',
  },
})

const imageUrl = response.data.url
```

## Static File Access
Tất cả uploaded files đều accessible via:
```
GET /uploads/{folder}/{filename}
```

Ví dụ:
- Public: `http://localhost:8080/uploads/public/images/uuid_timestamp.jpg`
- User: `http://localhost:8080/uploads/images/uuid_timestamp.jpg`  
- Admin: `http://localhost:8080/uploads/admin/images/uuid_timestamp.jpg`

## Security Notes
1. **Public uploads**: Cần implement rate limiting để tránh abuse
2. **File validation**: Backend validate file type và size
3. **Malware scanning**: Nên implement trong production
4. **CDN**: Trong production nên dùng CDN cho static files
