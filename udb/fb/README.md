## Rules of Firestore database
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

## Rules of Storage
```
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
		// Anyone can upload a public image if the file is less than 1Mb
    match /{allPaths=**} {
      allow create: if request.auth != null && request.resource.size < 1024 * 1024;
    }
    match /{allPaths=**} {
      allow read: if request.auth != null;
    }
  }
}
```
