package db

const createQuantumTable = `
CREATE TABLE IF NOT EXISTS quantum (
  signature   TEXT PRIMARY KEY,
  last        TEXT,
  nonce       INTEGER,
  type        INTEGER,
  signer      TEXT,
  timestamp   INTEGER
);`

const createContentTable = `
CREATE TABLE IF NOT EXISTS content (
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  quantum_signature  TEXT NOT NULL,
  data               BLOB,
  format             TEXT,
  FOREIGN KEY (quantum_signature) REFERENCES quantum(signature)
);`

const createReferenceTable = `
CREATE TABLE IF NOT EXISTS reference (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  ref_text  TEXT NOT NULL
);`

const createQuantumReferenceTable = `
CREATE TABLE IF NOT EXISTS quantum_reference (
  quantum_signature TEXT NOT NULL,
  reference_id      INTEGER NOT NULL,
  PRIMARY KEY (quantum_signature, reference_id),
  FOREIGN KEY (quantum_signature) REFERENCES quantum(signature),
  FOREIGN KEY (reference_id) REFERENCES reference(id)
);`
