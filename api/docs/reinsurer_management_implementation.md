# Reinsurer Management System Implementation

## Overview

This document describes the implementation of a comprehensive reinsurer management system in AART, similar to the existing broker management functionality. The system allows users to create, read, update, and delete reinsurer records with the following fields:

- Reinsurer Name
- Reinsurer Code (unique identifier)
- Contact Email
- Contact Person

## Implementation Details

### 1. Backend Implementation

#### 1.1 Database Model (`backend/models/group_pricing.go`)

Added `Reinsurer` struct:
```go
type Reinsurer struct {
    ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
    Name          string    `json:"name" gorm:"size:255;not null"`
    Code          string    `json:"code" gorm:"size:100;uniqueIndex;not null"`
    ContactEmail  string    `json:"contact_email" gorm:"size:255"`
    ContactPerson string    `json:"contact_person" gorm:"size:255"`
    CreationDate  time.Time `json:"creation_date" gorm:"autoCreateTime"`
    CreatedBy     string    `json:"created_by" gorm:"size:100"`
}
```

**Key Features:**
- Auto-incrementing ID
- Unique constraint on `code` field
- Auto-generated creation timestamp
- Tracks who created the record

#### 1.2 Service Layer (`backend/services/group_pricing.go`)

Implemented CRUD operations:

1. **CreateReinsurer**: Creates new reinsurer with duplicate code check
2. **GetReinsurers**: Retrieves all reinsurers ordered by name
3. **GetReinsurer**: Retrieves single reinsurer by ID
4. **EditReinsurer**: Updates reinsurer details (code cannot be changed)
5. **DeleteReinsurer**: Deletes reinsurer by ID

**Error Handling:**
- Returns `gorm.ErrDuplicatedKey` if code already exists
- Returns `gorm.ErrRecordNotFound` if reinsurer not found

#### 1.3 Controller Layer (`backend/controllers/group_pricing.go`)

Created REST API endpoints with proper HTTP status codes:

- **POST** - Returns 201 Created on success, 409 Conflict for duplicates
- **GET** - Returns 200 OK with data, 404 Not Found if missing
- **PUT** - Returns 200 OK on success, 404 if not found
- **DELETE** - Returns 200 OK on success

#### 1.4 Routes (`backend/routes/routes.go`)

Added routes under `/group-pricing/reinsurers`:
```go
groupPricing.POST("reinsurers", controllers.CreateReinsurer)
groupPricing.GET("reinsurers", controllers.GetReinsurers)
groupPricing.GET("reinsurers/:id", controllers.GetReinsurer)
groupPricing.PUT("reinsurers/:id", controllers.EditReinsurer)
groupPricing.DELETE("reinsurers/:id", controllers.DeleteReinsurer)
```

### 2. Database Migration

#### Migration File: `backend/migrations/create_reinsurers_table.sql`

```sql
CREATE TABLE IF NOT EXISTS reinsurers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    contact_email VARCHAR(255),
    contact_person VARCHAR(255),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    INDEX idx_reinsurer_code (code),
    INDEX idx_reinsurer_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**How to Apply:**
```bash
# For MySQL
mysql -u root -p new_aart < backend/migrations/create_reinsurers_table.sql

# For PostgreSQL (adjust syntax as needed)
psql -U postgres -d new_aart -f backend/migrations/create_reinsurers_table.sql
```

### 3. Frontend Implementation

#### 3.1 API Service (`frontend/src/renderer/api/GroupPricingService.js`)

Added reinsurer API methods:
```javascript
createReinsurer(reinsurer)
getReinsurers()
getReinsurer(id)
updateReinsurer(id, reinsurerDetails)
deleteReinsurer(id)
```

#### 3.2 Reinsurer Management Component

**File:** `frontend/src/renderer/components/grouppricing/ReinsurerManagement.vue`

**Features:**
- ✅ Add new reinsurer form with 4 fields
- ✅ AG Grid data table showing all reinsurers
- ✅ Click row to edit functionality
- ✅ Edit dialog with update and delete options
- ✅ Validation: Name and Code are required
- ✅ Code field is disabled during edit (immutable)
- ✅ Success/error notifications
- ✅ Confirmation dialog for deletions
- ✅ Empty state UI when no reinsurers exist
- ✅ Responsive design with Vuetify components

**Grid Columns:**
1. Reinsurer Name
2. Code
3. Contact Email
4. Contact Person
5. Created By

#### 3.3 Integration with Metadata Screen

Updated `frontend/src/renderer/screens/group_pricing/MetaData.vue` to include the ReinsurerManagement component.

**Component Order:**
1. Insurer Data Form
2. Broker Management
3. **Reinsurer Management** ← NEW
4. Scheme Category Management
5. Benefits Customization

### 4. Integration with Treaty Management

#### Enhanced Treaty Form

Updated `frontend/src/renderer/screens/group_pricing/bordereaux_management/components/RITreatyManagement.vue`:

**Changes:**
- Converted "Reinsurer Name" field to autocomplete dropdown
- Populated dropdown with reinsurers from database
- Auto-fills "Reinsurer Code" when reinsurer is selected
- Loads reinsurers on component mount
- Reinsurer code field becomes read-only when name is selected

**Benefits:**
- Consistent reinsurer data across treaties
- No manual entry errors (typos)
- Quick selection from pre-defined list
- Automatic code population reduces data entry

## API Endpoints Summary

### Reinsurer Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/group-pricing/reinsurers` | Create new reinsurer |
| GET | `/group-pricing/reinsurers` | Get all reinsurers |
| GET | `/group-pricing/reinsurers/:id` | Get single reinsurer |
| PUT | `/group-pricing/reinsurers/:id` | Update reinsurer |
| DELETE | `/group-pricing/reinsurers/:id` | Delete reinsurer |

### Request/Response Examples

#### Create Reinsurer
```json
POST /group-pricing/reinsurers
{
  "name": "Munich Re",
  "code": "MURE",
  "contact_email": "contact@munichre.com",
  "contact_person": "John Smith"
}

Response: 201 Created
```

#### Get All Reinsurers
```json
GET /group-pricing/reinsurers

Response: 200 OK
[
  {
    "id": 1,
    "name": "Munich Re",
    "code": "MURE",
    "contact_email": "contact@munichre.com",
    "contact_person": "John Smith",
    "creation_date": "2026-03-14T10:30:00Z",
    "created_by": "admin"
  }
]
```

#### Update Reinsurer
```json
PUT /group-pricing/reinsurers/1
{
  "name": "Munich Re Group",
  "code": "MURE",
  "contact_email": "info@munichre.com",
  "contact_person": "Jane Doe"
}

Response: 200 OK
```

#### Delete Reinsurer
```json
DELETE /group-pricing/reinsurers/1

Response: 200 OK
```

## User Workflow

### Adding a New Reinsurer

1. Navigate to: **Group Pricing → Metadata**
2. Scroll to "Reinsurer Management" section
3. Fill in the form:
   - Reinsurer Name (required)
   - Reinsurer Code (required, unique)
   - Contact Email (optional)
   - Contact Person (optional)
4. Click the **+** button
5. Success notification appears
6. Reinsurer appears in grid below

### Editing a Reinsurer

1. Click on any row in the reinsurer grid
2. Edit dialog opens with pre-filled data
3. Modify fields (except Code, which is read-only)
4. Click **Update Reinsurer**
5. Success notification appears
6. Grid refreshes with updated data

### Deleting a Reinsurer

1. Click on reinsurer row to open edit dialog
2. Click red **Delete** button
3. Confirm deletion in dialog
4. Success notification appears
5. Reinsurer removed from grid

### Using Reinsurers in Treaties

1. Navigate to: **Group Pricing → Bordereaux Management → Treaty Management**
2. Create or edit a treaty
3. In "Reinsurer Name" field, start typing or click dropdown
4. Select reinsurer from list
5. "Reinsurer Code" auto-fills
6. Continue with rest of treaty details

## Validation Rules

### Backend Validation
- ✅ Name: Required, max 255 characters
- ✅ Code: Required, max 100 characters, unique
- ✅ Email: Optional, max 255 characters
- ✅ Contact Person: Optional, max 255 characters

### Frontend Validation
- ✅ Name: Required field
- ✅ Code: Required field
- ✅ Email: Valid email format (if provided)
- ✅ Duplicate code: Server-side check returns error

## Error Handling

### Backend Errors
- **409 Conflict**: Reinsurer code already exists
- **404 Not Found**: Reinsurer ID doesn't exist
- **400 Bad Request**: Invalid input data
- **500 Internal Server Error**: Database or server error

### Frontend Error Messages
- "Reinsurer with this code already exists"
- "Please fill in reinsurer name and code"
- "Reinsurer not found"
- "Failed to add/update/delete reinsurer"

## Testing Checklist

### Backend Testing
- [ ] Create reinsurer with valid data → Success
- [ ] Create reinsurer with duplicate code → 409 error
- [ ] Get all reinsurers → Returns list ordered by name
- [ ] Get reinsurer by ID → Returns correct record
- [ ] Get non-existent reinsurer → 404 error
- [ ] Update reinsurer → Success
- [ ] Update non-existent reinsurer → 404 error
- [ ] Delete reinsurer → Success
- [ ] Delete non-existent reinsurer → Error handling

### Frontend Testing
- [ ] Add reinsurer form appears correctly
- [ ] Grid displays reinsurers in table
- [ ] Empty state shows when no reinsurers
- [ ] Click row opens edit dialog
- [ ] Edit dialog pre-fills data correctly
- [ ] Code field is disabled in edit mode
- [ ] Update button saves changes
- [ ] Delete button shows confirmation
- [ ] Notifications show success/error messages
- [ ] Reinsurer dropdown in treaty form works
- [ ] Auto-fill of code works correctly

## Future Enhancements

Potential improvements for future versions:

1. **Enhanced Validation**
   - Phone number format validation
   - Email domain verification
   - Address fields

2. **Advanced Features**
   - Reinsurer logo upload
   - Multiple contacts per reinsurer
   - Reinsurer rating information
   - Historical treaty relationship tracking

3. **Reporting**
   - Treaties per reinsurer report
   - Contact directory export
   - Reinsurer performance metrics

4. **Integration**
   - Link to external reinsurer databases
   - Auto-populate from industry directories
   - Integration with bordereaux submission

5. **Security**
   - Role-based access control for reinsurer management
   - Audit trail for changes
   - Approval workflow for deletions

## Support

For questions or issues:
- Review this documentation
- Check error messages in browser console
- Verify database migration was applied
- Test API endpoints with Swagger UI at `/api/v1/swagger/`

## Version History

- **v5.5.0** (2026-03-14) - Initial implementation of reinsurer management system
