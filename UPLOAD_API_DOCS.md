# Upload API Documentation

## Overview
Backend now supports image uploads with persistent storage for:
- **User Profile Pictures** - User profile images
- **Project Images** - Project logos/images

## Endpoints

### 1. Upload Profile Picture
```
POST /api/upload/profile-picture
Authorization: Bearer {token}
Content-Type: multipart/form-data

Body:
- image: {file} - Image file (JPEG, PNG, GIF, WebP)
```

**Constraints:**
- Max file size: 5MB
- Allowed formats: JPEG, PNG, GIF, WebP
- Field name must be: `image`

**Response:**
```json
{
  "message": "Profile picture uploaded successfully",
  "image_url": "/uploads/profiles/user_9_1778163111.png"
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8000/api/upload/profile-picture \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@/path/to/image.png"
```

### 2. Upload Project Image
```
POST /api/project/:projectId/upload-image
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Constraints:**
- Max file size: 10MB
- Allowed formats: JPEG, PNG, GIF, WebP
- Only Admin or Project Manager can upload
- Field name must be: `image`

**Response:**
```json
{
  "message": "Project image uploaded successfully",
  "image_url": "/uploads/projects/project_4_1778162776.png"
}
```

**Example cURL:**
```bash
curl -X POST http://localhost:8000/api/project/4/upload-image \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "image=@/path/to/image.png"
```

### 3. Retrieve Uploaded Image
```
GET /uploads/:type/:filename

Available types:
- profiles
- projects
```

**Example:**
```bash
curl http://localhost:8000/uploads/profiles/user_9_1778163111.png -o downloaded_image.png
```

## Storage
- All uploaded files are stored in: `./uploads/`
- Directory structure:
  ```
  uploads/
  ├── profiles/     - User profile pictures
  └── projects/     - Project images
  ```
- Files are persistent using Docker volumes
- Files survive container restarts

## Error Handling

### Common Errors

**400 - No file provided:**
```json
{"error": "No file provided"}
```

**400 - File too large:**
```json
{"error": "File size must be less than 5MB"}
```

**400 - Invalid file type:**
```json
{"error": "Only JPEG, PNG, GIF, and WebP images are allowed"}
```

**401 - Unauthorized:**
```json
{"error": "Invalid or expired token"}
```

**403 - Forbidden (not Admin/PM for project upload):**
```json
{"error": "Unauthorized action"}
```

**404 - User/Project not found:**
```json
{"error": "User not found"}
```

**500 - Server error:**
```json
{"error": "Failed to save file"}
```

## Integration Example

### React/Frontend Example
```javascript
async function uploadProfilePicture(file, token) {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch('http://localhost:8000/api/upload/profile-picture', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  });

  const data = await response.json();
  if (response.ok) {
    console.log('Uploaded:', data.image_url);
    // Store image_url in UI or update profile
    return data.image_url;
  } else {
    console.error('Upload failed:', data.error);
    throw new Error(data.error);
  }
}
```

## Development Setup

1. **Container Volume Mapping:**
   ```yaml
   volumes:
     - ./uploads:/root/uploads
   ```

2. **Create Directories:**
   ```bash
   mkdir -p uploads/profiles
   mkdir -p uploads/projects
   ```

3. **Start Backend:**
   ```bash
   docker-compose up --build -d
   ```

4. **Verify Upload Working:**
   ```bash
   bash /tmp/test_persistent_upload.sh
   ```

## Testing

Run complete upload test:
```bash
bash /tmp/test_persistent_upload.sh
```

Test specific features:
```bash
# Test profile upload
curl -X POST http://localhost:8000/api/upload/profile-picture \
  -H "Authorization: Bearer TOKEN" \
  -F "image=@test.png"

# Test project upload
curl -X POST http://localhost:8000/api/project/1/upload-image \
  -H "Authorization: Bearer TOKEN" \
  -F "image=@logo.png"

# Access uploaded image
curl http://localhost:8000/uploads/profiles/user_9_1778163111.png
```

## Important Notes

1. **Image URLs**: Returned URLs are relative paths. Use them with your API domain:
   - API: `http://localhost:8000/uploads/profiles/user_9_1778163111.png`
   - Frontend: `<img src="/uploads/profiles/user_9_1778163111.png" />`

2. **Permissions**: 
   - Any authenticated user can upload their own profile picture
   - Only Admin/PM can upload project images

3. **Storage**: Files are stored permanently on host filesystem until manually deleted

4. **Access**: Uploaded files are accessible via HTTP without authentication

## Troubleshooting

**Issue: "File not found" after restart**
- Solution: Ensure `docker-compose.yml` has volumes mapping for uploads

**Issue: File uploaded but can't access**
- Check file permissions: `ls -la uploads/`
- Check volume mapping: `docker-compose config`

**Issue: Error "Failed to create upload directory"**
- Ensure `/uploads/profiles` and `/uploads/projects` directories exist
- Check disk space and permissions
