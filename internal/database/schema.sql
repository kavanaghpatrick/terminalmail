-- Emails table
CREATE TABLE IF NOT EXISTS emails (
    id TEXT PRIMARY KEY,
    subject TEXT,
    body_text TEXT,
    body_html TEXT,
    from_email TEXT,
    to_email TEXT,
    cc_email TEXT,
    date INTEGER,
    labels TEXT,
    ai_summary TEXT,
    status TEXT,
    thread_id TEXT,
    ai_category TEXT
);

-- FTS5 virtual table for full-text search
CREATE VIRTUAL TABLE IF NOT EXISTS email_search USING fts5(
    id,
    subject,
    body_text,
    from_email,
    to_email,
    cc_email,
    labels,
    ai_summary,
    ai_category
);

-- Triggers to keep FTS5 table in sync with emails table
CREATE TRIGGER IF NOT EXISTS email_insert_trigger
AFTER INSERT ON emails
BEGIN
    INSERT INTO email_search (
        id, subject, body_text, from_email, to_email, cc_email, 
        labels, ai_summary, ai_category
    ) VALUES (
        new.id, new.subject, new.body_text, new.from_email, new.to_email, 
        new.cc_email, new.labels, new.ai_summary, new.ai_category
    );
END;

CREATE TRIGGER IF NOT EXISTS email_update_trigger
AFTER UPDATE ON emails
BEGIN
    DELETE FROM email_search WHERE id = old.id;
    INSERT INTO email_search (
        id, subject, body_text, from_email, to_email, cc_email, 
        labels, ai_summary, ai_category
    ) VALUES (
        new.id, new.subject, new.body_text, new.from_email, new.to_email, 
        new.cc_email, new.labels, new.ai_summary, new.ai_category
    );
END;

CREATE TRIGGER IF NOT EXISTS email_delete_trigger
AFTER DELETE ON emails
BEGIN
    DELETE FROM email_search WHERE id = old.id;
END;

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_emails_labels_date ON emails(labels, date);
CREATE INDEX IF NOT EXISTS idx_emails_ai_category ON emails(ai_category);

-- Terminal preferences singleton table
CREATE TABLE IF NOT EXISTS terminal_preferences (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    preferences TEXT
);