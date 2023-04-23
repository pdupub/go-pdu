# Setup test env

## Steps

1. Apply Firebase Service
2. Setup Authentication (Anonymouse)
3. Firestore Database and set rules
4. Store and set rules
5. download adminsdk.json from Project Settings
6. change TestFirebaseProjectID to your app
7. cd udb/fb and go test


## Rules of Firestore Database
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

## Indexes of Firestore Database

| Collection ID |	Fields indexed | Query scope | Status |
|---------------|----------------|-------------|--------|
| quantum |	address Ascending createTime Ascending __name__ Ascending	| Collection | Enabled	|
| quantum	| address Ascending seq Ascending __name__ Ascending	| Collection | Enabled	|
| quantum	| address Ascending createTime Descending __name__ Descending	| Collection | Enabled	|
| quantum	| address Ascending seq Descending __name__ Descending	| Collection	|	Enabled	|
| quantum	| ios.action Ascending ios.param Ascending seq Descending __name__ Descending	| Collection	|	Enabled	|
| individua | attitude.level Ascending updateTime Descending __name__ Descending | Collection | Enabled |


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
