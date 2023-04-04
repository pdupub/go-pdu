## rules of firestore database
```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
  	// All data is readable for users
    match /{document=**} {
      allow read: if request.auth != null;
    }
    
    // user can only create new quantum, not update or delete
    match /quantum/{document=**} {
    	allow create: if request.auth != null;
    }
  }
}
```