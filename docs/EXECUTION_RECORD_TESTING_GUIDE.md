# Execution Record Feature - Testing Guide

## Overview
This document provides instructions for testing the newly implemented Execution Record Management feature.

## Prerequisites
- Docker and Docker Compose installed
- PostgreSQL database running (via docker-compose)
- Backend built successfully
- Frontend dependencies installed

## Quick Start

### 1. Start the Database
```bash
cd /home/runner/work/ops-core/ops-core
docker-compose up -d db
```

### 2. Start the Backend
```bash
cd backend
export DATABASE_URL="postgres://opscore_user:opscore_password@localhost:5432/opscore_db?sslmode=disable"
export ENCRYPTION_KEY="dev-key-123456789012345678901234"
go run cmd/server/main.go
```

The backend will start on http://localhost:8080

### 3. Start the Frontend
```bash
cd frontend
npm run dev
```

The frontend will start on http://localhost:5173

## Testing the Feature

### Create an Execution Record (via UI)

1. **Navigate to Documents**
   - Go to http://localhost:5173/documents
   - Select a document from the list

2. **Start Execution**
   - Click on a document to view it
   - Navigate to `/documents/{docId}/execute`
   - You should see the 3-pane layout:
     - Left: Variable input form
     - Center: Document content
     - Right: Execution panel

3. **Fill Variables**
   - Enter values for all required variables in the left pane
   - Values will automatically substitute in the document

4. **Create Execution**
   - Review the execution title (editable)
   - Click "Start Execution"
   - The execution record is created with status "in_progress"

5. **Add Steps**
   - In the right pane, add execution steps
   - Each step has:
     - Step number
     - Description
     - Notes (optional)
   - Click "Add Step" to create

6. **Update Notes**
   - Click "Edit Notes" on any step
   - Add detailed notes about what was done
   - Click "Save Notes"

7. **Complete Execution**
   - Add overall notes in the "Overall Notes" field
   - Click "Complete" to mark execution as successful
   - OR click "Mark as Failed" if something went wrong

8. **View History**
   - Navigate back to the document
   - (Future: Add ExecutionRecordList component to document page)
   - Navigate to `/documents/{docId}/execute/{recordId}` to view past executions

### Test via API (cURL)

```bash
# Set the API URL
API_URL="http://localhost:8080/api/v1"

# Create an execution record
curl -X POST "$API_URL/execution-records" \
  -H "Content-Type: application/json" \
  -H "user_id: test-user-123" \
  -d '{
    "document_id": "doc-123",
    "document_version_id": "ver-456",
    "title": "Test Execution - $(date)",
    "variable_values": [
      {"name": "server_name", "value": "prod-server-01"},
      {"name": "backup_path", "value": "/data/backups"}
    ]
  }'

# Get the record (use ID from response above)
curl -X GET "$API_URL/execution-records/{id}"

# Add a step
curl -X POST "$API_URL/execution-records/{id}/steps" \
  -H "Content-Type: application/json" \
  -d '{
    "step_number": 1,
    "description": "Verify server connection"
  }'

# Update step notes
curl -X PUT "$API_URL/execution-records/{id}/steps/1/notes" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Successfully connected to server. Latency: 15ms"
  }'

# Complete the execution
curl -X POST "$API_URL/execution-records/{id}/complete" \
  -H "Content-Type: application/json"

# Search executions
curl -X GET "$API_URL/execution-records?status=completed"
```

## Test Scenarios

### Scenario 1: Complete Execution Flow
1. Start a new execution
2. Fill in all variables
3. Add 3-5 execution steps
4. Add notes to each step
5. Add overall notes
6. Mark as completed
7. Verify completion timestamp is set
8. Verify status badge shows "completed"

### Scenario 2: Failed Execution
1. Start a new execution
2. Add 2 steps
3. Add note explaining the failure
4. Mark as failed
5. Verify status shows "failed"
6. Verify completion timestamp is set

### Scenario 3: In-Progress Execution
1. Start a new execution
2. Add 1 step
3. Leave it in progress
4. Navigate away and back
5. Verify state is preserved (in-memory only)

### Scenario 4: Variable Substitution
1. Create document with variables
2. Start execution
3. Enter variable values
4. Verify document content updates in real-time
5. Verify variables are saved with execution

### Scenario 5: Search Executions
1. Create multiple executions with different statuses
2. Use the search endpoint with filters
3. Filter by status, document ID, executor ID
4. Verify correct results returned

## Expected Behavior

### 3-Pane Layout
- ✅ Left pane shows variable input form
- ✅ Center pane shows document with substituted variables
- ✅ Right pane shows execution tracking
- ✅ All panes are responsive and scrollable

### Execution States
- ✅ `in_progress` - Can edit, add steps, update notes
- ✅ `completed` - Read-only, shows completion time
- ✅ `failed` - Read-only, shows completion time

### Step Management
- ✅ Can add steps with any step number
- ✅ Can update notes on existing steps
- ✅ Steps show execution timestamp
- ✅ Cannot add steps to completed/failed executions

### Data Persistence
- ⚠️ **In-Memory Only** - Data lost on server restart
- ⚠️ **No Database** - Need to implement PostgreSQL repository

## Known Issues

1. **In-Memory Storage**: Data is not persisted to database
2. **No Authentication**: User ID is mocked in context
3. **No File Attachments**: Feature referenced but not implemented
4. **Page Refresh**: May lose state (depends on routing)

## Troubleshooting

### Backend Issues
- **Connection refused**: Check if backend is running on port 8080
- **Database errors**: Verify PostgreSQL is running and migrations applied
- **404 errors**: Check routes are registered in main.go

### Frontend Issues
- **Blank page**: Check console for errors
- **API errors**: Verify backend is running and accessible
- **Type errors**: Rebuild frontend (`npm run build`)

### Common Problems
- **CORS errors**: Backend has CORS middleware configured
- **401 Unauthorized**: Mock user_id in header (backend needs update)
- **500 errors**: Check backend logs for details

## Next Steps

1. **Database Persistence**: Implement PostgreSQL repository
2. **Authentication**: Add proper auth middleware
3. **Integration Tests**: Add handler and integration tests
4. **Attachment Upload**: Implement file attachment feature
5. **UI Polish**: Add loading states, better error handling

## Success Criteria

- ✅ Can create execution records via UI
- ✅ Can add and update execution steps
- ✅ Can mark executions as complete or failed
- ✅ Variable substitution works correctly
- ✅ 3-pane layout displays properly
- ✅ API endpoints respond correctly
- ✅ No security vulnerabilities
- ✅ TypeScript types are correct
- ✅ Backend tests pass
